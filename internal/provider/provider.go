// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pkg/errors"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
	"golang.org/x/oauth2"
)

// Ensure PlanetScaleProvider satisfies various provider interfaces.
var _ provider.Provider = &PlanetScaleProvider{}

// PlanetScaleProvider defines the provider implementation.
type PlanetScaleProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PlanetScaleProviderModel describes the provider data model.
type PlanetScaleProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`

	ServiceTokenName types.String `tfsdk:"service_token_name"`
}

func (p *PlanetScaleProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "planetscale"
	resp.Version = p.version
}

func (p *PlanetScaleProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
			"service_token_name": schema.StringAttribute{
				MarkdownDescription: "Name of the service token to use",
				Optional:            true,
			},
		},
	}
}

func (p *PlanetScaleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PlanetScaleProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	var (
		initrt = debugRoundTripper(func(req, res []byte) {
			tflog.Debug(ctx, "roundtripper", map[string]interface{}{
				"req": string(req),
				"res": string(res),
			})
		}, http.DefaultTransport)
		rt      http.RoundTripper
		baseURL *url.URL
	)
	if !data.Endpoint.IsNull() {
		u, err := url.Parse(data.Endpoint.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("endpoint"), "invalid URL", err.Error())
			return
		}
		baseURL = u
	}
	var (
		accessToken       = os.Getenv("PLANETSCALE_ACCESS_TOKEN")
		serviceTokenName  = os.Getenv("PLANETSCALE_SERVICE_TOKEN_NAME")
		serviceTokenValue = os.Getenv("PLANETSCALE_SERVICE_TOKEN")
	)
	if !data.ServiceTokenName.IsNull() {
		serviceTokenName = data.ServiceTokenName.ValueString()
	}
	switch {
	case accessToken != "" && serviceTokenName == "" && serviceTokenValue == "":
		tok := &oauth2.Token{AccessToken: accessToken}
		rt = &oauth2.Transport{Base: initrt, Source: oauth2.StaticTokenSource(tok)}
	case accessToken == "" && serviceTokenName != "" && serviceTokenValue != "":
		rt = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			r.Header.Set("Authorization", serviceTokenName+":"+serviceTokenValue)
			return initrt.RoundTrip(r)
		})
	case accessToken == "" && serviceTokenName == "" && serviceTokenValue == "":
		resp.Diagnostics.AddError("Missing PlanetScale credentials.",
			"You must set either of:\n"+
				"- `PLANETSCALE_ACCESS_TOKEN`\n"+
				"- `PLANETSCALE_SERVICE_TOKEN_NAME` and `PLANETSCALE_SERVICE_TOKEN`")
	case accessToken == "" && serviceTokenName != "" && serviceTokenValue == "",
		accessToken == "" && serviceTokenName == "" && serviceTokenValue != "":
		resp.Diagnostics.AddError("Incomplete PlanetScale service token credentials.",
			"Both of `PLANETSCALE_SERVICE_TOKEN_NAME` and `PLANETSCALE_SERVICE_TOKEN` must be set.")
	default:
		resp.Diagnostics.AddError("Ambiguous PlanetScale credentials.", "You must set only either of an access token or a service token, but not both:\n"+
			"- `PLANETSCALE_ACCESS_TOKEN`\n"+
			"- `PLANETSCALE_SERVICE_TOKEN_NAME` and `PLANETSCALE_SERVICE_TOKEN`")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := planetscale.NewClient(&http.Client{Transport: rt}, baseURL)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PlanetScaleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newDatabaseResource,
		newBranchResource,
		newBackupResource,
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
		newBackupDataSource,
		newBackupsDataSource,
		newOAuthApplicationsDataSource,
		newUserDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PlanetScaleProvider{
			version: version,
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
