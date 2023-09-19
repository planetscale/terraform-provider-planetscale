package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &organizationRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationRegionsDataSource{}
)

func newOrganizationRegionsDataSource() datasource.DataSource {
	return &organizationRegionsDataSource{}
}

type organizationRegionsDataSource struct {
	client *planetscale.Client
}

type organizationRegionsDataSourceModel struct {
	Organization string                  `tfsdk:"organization"`
	Regions      []regionDataSourceModel `tfsdk:"regions"`
}

func (d *organizationRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_regions"
}

func (d *organizationRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"regions": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: regionDataSourceSchemaAttribute,
			},
		},
	}}
}

func (d *organizationRegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *organizationRegionsDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgName := data.Organization

	res, err := d.client.ListRegionsForOrganization(ctx, orgName, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read organization regions", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read organization regions", "no data")
		return
	}
	state := organizationRegionsDataSourceModel{
		Organization: data.Organization,
		Regions:      make([]regionDataSourceModel, 0, len(res.Data)),
	}
	for _, item := range res.Data {
		item := item
		state.Regions = append(state.Regions, *regionFromClient(&item, resp.Diagnostics))
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
