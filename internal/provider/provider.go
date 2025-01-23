package provider

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkg/errors"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
	"golang.org/x/oauth2"
)

var _ provider.ProviderWithConfigValidators = &PlanetScaleProvider{}

type PlanetScaleProvider struct {
	version string
	debug   bool
}

type PlanetScaleProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`

	AccessToken types.String `tfsdk:"access_token"`

	ServiceTokenID    types.String `tfsdk:"service_token_id"`   // new preferred field
	ServiceTokenName  types.String `tfsdk:"service_token_name"` // deprecated
	ServiceTokenValue types.String `tfsdk:"service_token"`
}

func (p *PlanetScaleProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "planetscale"
	resp.Version = p.version
}

func (p *PlanetScaleProvider) ConfigValidators(context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{
		providervalidator.Conflicting(path.MatchRoot("access_token"), path.MatchRoot("service_token")),
		providervalidator.Conflicting(path.MatchRoot("access_token"), path.MatchRoot("service_token_name")),
	}
}

func (p *PlanetScaleProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The PlanetScale provider allows using the OpenAPI surface of our public API. To use this provider, one of the following are required:

- access token credentials, configured or stored in the environment variable ` + "`PLANETSCALE_ACCESS_TOKEN`" + `
- service token credentials, configured or stored in the environment variables ` + "`PLANETSCALE_SERVICE_TOKEN_NAME`" + ` and ` + "`PLANETSCALE_SERVICE_TOKEN`" + `

Note that the provider is not production ready and only for early testing at this time.

Known limitations:
- Support for deployments, deploy queues, deploy requests and reverts is not implemented at this time. If you have a use case for it, please let us know in the repository issues.
- When using service tokens (recommended), ensure the token has the ` + "`create_databases`" + ` organization-level permission. This allows terraform to create new databases and automatically grants the token all other permissions on the databases created by the token.`,
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "If set, points the API client to a different endpoint than `https:://api.planetscale.com/v1`.",
				Optional:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Name of the service token to use. Alternatively, use `PLANETSCALE_SERVICE_TOKEN_NAME`. Mutually exclusive with `service_token_name` and `service_token`.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("service_token_name")),
					stringvalidator.ConflictsWith(path.MatchRoot("service_token_id")),
					stringvalidator.ConflictsWith(path.MatchRoot("service_token")),
				},
			},
			"service_token_id": schema.StringAttribute{
				MarkdownDescription: "ID of the service token to use. Alternatively, use `PLANETSCALE_SERVICE_TOKEN_ID`. Mutually exclusive with `access_token`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("access_token")),
					stringvalidator.ConflictsWith(path.MatchRoot("service_token_name")),
				},
			},
			"service_token_name": schema.StringAttribute{
				MarkdownDescription: "Name of the service token to use. Alternatively, use `PLANETSCALE_SERVICE_TOKEN_NAME`. Mutually exclusive with `access_token`. (Deprecated, use `service_token_id` instead)",
				Optional:            true,
				DeprecationMessage:  "Use service_token_id instead. This field will be removed in a future version.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("access_token")),
					stringvalidator.ConflictsWith(path.MatchRoot("service_token_id")),
				},
			},
			"service_token": schema.StringAttribute{
				MarkdownDescription: "Value of the service token to use. Alternatively, use `PLANETSCALE_SERVICE_TOKEN`. Mutually exclusive with `access_token`.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("access_token")),
				},
			},
		},
	}
}

func (p *PlanetScaleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PlanetScaleProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		initrt  http.RoundTripper
		rt      http.RoundTripper
		baseURL *url.URL
	)
	if !p.debug {
		initrt = http.DefaultTransport
	} else {
		initrt = debugRoundTripper(func(req, res []byte) {
			tflog.Debug(ctx, "roundtripper", map[string]interface{}{
				"req": string(req),
				"res": string(res),
			})
		}, http.DefaultTransport)
	}
	if !data.Endpoint.IsNull() {
		u, err := url.Parse(data.Endpoint.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("endpoint"), "invalid URL", err.Error())
			return
		}
		baseURL = u
	}
	var (
		accessToken       = stringValueOrDefault(data.AccessToken, os.Getenv("PLANETSCALE_ACCESS_TOKEN"))
		serviceTokenID    = stringValueOrDefault(data.ServiceTokenID, os.Getenv("PLANETSCALE_SERVICE_TOKEN_ID"))
		serviceTokenName  = stringValueOrDefault(data.ServiceTokenName, os.Getenv("PLANETSCALE_SERVICE_TOKEN_NAME"))
		serviceTokenValue = stringValueOrDefault(data.ServiceTokenValue, os.Getenv("PLANETSCALE_SERVICE_TOKEN"))
	)

	// Warn if the deprecated PLANETSCALE_SERVICE_TOKEN_NAME env var is used.
	// Adding this to `resp.Diagnostics` ensures it will be printed during typical
	// terraform operations, whereas logging with `tflog.Warn()` will only show if the
	// users specifies the `TF_LOG` env var.
	if serviceTokenName != "" && serviceTokenID == "" {
		resp.Diagnostics.AddWarning(
			"Deprecated Configuration",
			"PLANETSCALE_SERVICE_TOKEN_NAME is deprecated. Please use PLANETSCALE_SERVICE_TOKEN_ID instead.",
		)
	}

	// Use serviceTokenID if available, fall back to serviceTokenName
	effectiveTokenID := serviceTokenID
	if effectiveTokenID == "" {
		effectiveTokenID = serviceTokenName
	}

	switch {
	case accessToken != "" && effectiveTokenID == "" && serviceTokenValue == "":
		tok := &oauth2.Token{AccessToken: accessToken}
		rt = &oauth2.Transport{Base: initrt, Source: oauth2.StaticTokenSource(tok)}
	case accessToken == "" && effectiveTokenID != "" && serviceTokenValue != "":
		rt = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Set("Authorization", effectiveTokenID+":"+serviceTokenValue)
			return initrt.RoundTrip(r)
		})
	case accessToken == "" && effectiveTokenID == "" && serviceTokenValue == "":
		resp.Diagnostics.AddError("Missing PlanetScale credentials.",
			"You must set either of:\n"+
				"- `PLANETSCALE_ACCESS_TOKEN`\n"+
				"- `PLANETSCALE_SERVICE_TOKEN_ID` and `PLANETSCALE_SERVICE_TOKEN`")
	case accessToken == "" && effectiveTokenID != "" && serviceTokenValue == "",
		accessToken == "" && effectiveTokenID == "" && serviceTokenValue != "":
		resp.Diagnostics.AddError("Incomplete PlanetScale service token credentials.",
			"Both of `PLANETSCALE_SERVICE_TOKEN_ID` and `PLANETSCALE_SERVICE_TOKEN` must be set.")
	default:
		resp.Diagnostics.AddError("Ambiguous PlanetScale credentials.", "You must set only an access token or a service token, but not both:\n"+
			"- `PLANETSCALE_ACCESS_TOKEN`\n"+
			"- `PLANETSCALE_SERVICE_TOKEN_ID` and `PLANETSCALE_SERVICE_TOKEN`")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := planetscale.NewClient(
		&http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				r.Header.Set("User-Agent", "PlanetScale_Terraform_Provider/"+p.version+" (Terraform "+req.TerraformVersion+")")
				return rt.RoundTrip(r)
			}),
		}, baseURL,
	)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PlanetScaleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newDatabaseResource,
		newBranchResource,
		newBranchSafeMigrationsResource,
		newBackupResource,
		newPasswordResource,
	}
}

func (p *PlanetScaleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newOrganizationsDataSource,
		newOrganizationDataSource,
		newOrganizationRegionsDataSource,
		newDatabasesDataSource,
		newDatabaseDataSource,
		newDatabaseRegionsDataSource,
		newDatabaseReadOnlyRegionsDataSource,
		newBranchesDataSource,
		newBranchDataSource,
		newBranchSchemaDataSource,
		newBranchSchemaLintDataSource,
		newbranchSafeMigrationsDataSource,
		newBackupDataSource,
		newBackupsDataSource,
		newPasswordDataSource,
		newPasswordsDataSource,
		newOAuthApplicationsDataSource,
		newUserDataSource,
	}
}

func New(version string, debug bool) func() provider.Provider {
	return func() provider.Provider {
		return &PlanetScaleProvider{
			version: version,
			debug:   debug,
		}
	}
}

func debugRoundTripper(log func(req, res []byte), tpt http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		debugReq, err := httputil.DumpRequestOut(r, true)
		if err != nil {
			return nil, errors.Wrap(err, "dumping request output")
		}
		res, err := tpt.RoundTrip(r)
		if res == nil {
			return res, err
		}
		debugRes, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, errors.Wrap(err, "dumping response output")
		}
		log(debugReq, debugRes)
		return res, err
	})
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func boolIfDifferent(oldBool, newBool types.Bool, wasChanged *bool) *bool {
	if oldBool.Equal(newBool) {
		return nil
	}
	*wasChanged = true
	return boolValueIfKnown(newBool)
}

func stringIfDifferent(oldString, newString types.String, wasChanged *bool) *string {
	if oldString == newString {
		return nil
	}
	*wasChanged = true
	return stringValueIfKnown(newString)
}

func boolValueIfKnown(v basetypes.BoolValue) *bool {
	if v.IsUnknown() || v.IsNull() {
		return nil
	}
	return v.ValueBoolPointer()
}

func stringValueIfKnown(v basetypes.StringValue) *string {
	if v.IsUnknown() || v.IsNull() {
		return nil
	}
	return v.ValueStringPointer()
}

//nolint:unused
func float64ValueIfKnown(v basetypes.Float64Value) *float64 {
	if v.IsUnknown() || v.IsNull() {
		return nil
	}
	return v.ValueFloat64Pointer()
}

func stringValueOrDefault(v basetypes.StringValue, def string) string {
	if v.IsUnknown() || v.IsNull() {
		return def
	}
	return v.ValueString()
}
