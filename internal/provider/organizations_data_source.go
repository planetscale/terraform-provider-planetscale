package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &organizationsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationsDataSource{}
)

func newOrganizationsDataSource() datasource.DataSource {
	return &organizationsDataSource{}
}

type organizationsDataSource struct {
	client *planetscale.Client
}

type organizationsDataSourceModel struct {
	Organizations []organizationDataSourceModel `tfsdk:"organizations"`
}

func (d *organizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *organizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A list of PlanetScale organizations.",
		MarkdownDescription: "A list of PlanetScale organizations.",
		Attributes: map[string]schema.Attribute{
			"organizations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: organizationDataSourceSchemaAttribute(true),
				},
			},
		},
	}
}

func (d *organizationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	res, err := d.client.ListOrganizations(ctx, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read organizations", err.Error())
		return
	}
	state := organizationsDataSourceModel{
		Organizations: make([]organizationDataSourceModel, 0, len(res.Data)),
	}
	for _, item := range res.Data {
		item := item
		state.Organizations = append(state.Organizations, *organizationFromClient(&item, resp.Diagnostics))
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
