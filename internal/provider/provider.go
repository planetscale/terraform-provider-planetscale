// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/planetscale-go/planetscale"
)

// Ensure PlanetScaleProvider satisfies various provider interfaces.
var _ provider.Provider = &PlanetScaleProvider{}

// PlanetScaleProvider defines the provider implementation.
type PlanetScaleProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string

	client *planetscale.Client
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

	var clientOpts []planetscale.ClientOption
	if !data.Endpoint.IsNull() {
		clientOpts = append(clientOpts, planetscale.WithBaseURL(data.Endpoint.ValueString()))
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
		clientOpts = append(clientOpts, planetscale.WithAccessToken(accessToken))
	case accessToken == "" && serviceTokenName != "" && serviceTokenValue != "":
		clientOpts = append(clientOpts, planetscale.WithServiceToken(serviceTokenName, serviceTokenValue))
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

	client, err := planetscale.NewClient(clientOpts...)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create PlanetScale client",
			"An unexpected error occurred when creating the PlanetScale API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"PlanetScale Client Error: "+err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PlanetScaleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *PlanetScaleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PlanetScaleProvider{
			version: version,
		}
	}
}
