package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &databaseRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &databaseRegionsDataSource{}
)

func newDatabaseRegionsDataSource() datasource.DataSource {
	return &databaseRegionsDataSource{}
}

type databaseRegionsDataSource struct {
	client *planetscale.Client
}

type databaseRegionsDataSourceModel struct {
	Organization string                          `tfsdk:"organization"`
	Name         string                          `tfsdk:"name"`
	Regions      []databaseRegionDataSourceModel `tfsdk:"regions"`
}

type databaseRegionDataSourceModel struct {
	DisplayName       string   `tfsdk:"display_name"`
	Enabled           bool     `tfsdk:"enabled"`
	Id                string   `tfsdk:"id"`
	Location          string   `tfsdk:"location"`
	Provider          string   `tfsdk:"provider"`
	PublicIpAddresses []string `tfsdk:"public_ip_addresses"`
	Slug              string   `tfsdk:"slug"`
}

func (d *databaseRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_regions"
}

func (d *databaseRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A list of PlanetScale regions.",
		MarkdownDescription: "A list of PlanetScale regions.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description: "The organization for which the regions are available.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The database for which the regions are available.",
				Required:    true,
			},
			"regions": schema.ListNestedAttribute{
				Description: "The list of regions available for the database.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: regionDataSourceSchemaAttribute,
				},
			},
		},
	}
}

func (d *databaseRegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *databaseRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *databaseRegionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.ListDatabaseRegions(ctx, data.Organization, data.Name, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list database regions", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Received a nil database regions list", "")
		return
	}
	state := databaseRegionsDataSourceModel{
		Organization: data.Organization,
		Name:         data.Name,
	}
	for _, region := range res.Data {
		state.Regions = append(state.Regions, databaseRegionDataSourceModel{
			DisplayName:       region.DisplayName,
			Enabled:           region.Enabled,
			Id:                region.Id,
			Location:          region.Location,
			Provider:          region.Provider,
			PublicIpAddresses: region.PublicIpAddresses,
			Slug:              region.Slug,
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
