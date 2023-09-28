package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &oauthApplicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthApplicationsDataSource{}
)

func newOAuthApplicationsDataSource() datasource.DataSource {
	return &oauthApplicationsDataSource{}
}

type oauthApplicationsDataSource struct {
	client *planetscale.Client
}

type oauthApplicationsDataSourceModel struct {
	Organization string                            `tfsdk:"organization"`
	Applications []oauthApplicationDataSourceModel `tfsdk:"applications"`
}

func (d *oauthApplicationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_applications"
}

func (d *oauthApplicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A list of PlanetScale OAuth applications. (requires feature flag)",
		MarkdownDescription: "A list of PlanetScale OAuth applications. (requires feature flag)",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true},
			"applications": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: oauthApplicationAttribute,
				},
			},
		},
	}
}

func (d *oauthApplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*planetscale.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *planetscale.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *oauthApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *oauthApplicationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.ListOauthApplications(ctx, data.Organization, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list oauth applications", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to list oauth applications", "no data")
		return
	}

	state := oauthApplicationsDataSourceModel{
		Organization: data.Organization,
		Applications: make([]oauthApplicationDataSourceModel, 0, len(res.Data)),
	}
	for _, item := range res.Data {
		item := item
		state.Applications = append(state.Applications, *oauthApplicationFromClient(&item, resp.Diagnostics))
		if resp.Diagnostics.HasError() {
			return
		}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
