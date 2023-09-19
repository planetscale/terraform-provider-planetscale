package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &databaseReadOnlyRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &databaseReadOnlyRegionsDataSource{}
)

func newDatabaseReadOnlyRegionsDataSource() datasource.DataSource {
	return &databaseReadOnlyRegionsDataSource{}
}

type databaseReadOnlyRegionsDataSource struct {
	client *planetscale.Client
}

type TTreadOnlyRegionsDataSourceModel struct {
	Organization string                          `tfsdk:"organization"`
	Name         string                          `tfsdk:"name"`
	Regions      []readOnlyRegionDataSourceModel `tfsdk:"regions"`
}

func (d *databaseReadOnlyRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_read_only_regions"
}

func (d *databaseReadOnlyRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: readOnlyRegionDataSourceSchemaAttribute}
}

func (d *databaseReadOnlyRegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *databaseReadOnlyRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *TTreadOnlyRegionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.ListReadOnlyRegions(ctx, data.Organization, data.Name, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list database read only regions", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Received a nil database read only regions list", "")
		return
	}
	state := TTreadOnlyRegionsDataSourceModel{
		Organization: data.Organization,
		Name:         data.Name,
		Regions:      make([]readOnlyRegionDataSourceModel, 0, len(res.Data)),
	}
	for _, item := range res.Data {
		item := item
		state.Regions = append(state.Regions, *readOnlyRegionFromClient(&item, resp.Diagnostics))
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
