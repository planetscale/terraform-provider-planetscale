package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &databaseDataSource{}
	_ datasource.DataSourceWithConfigure = &databaseDataSource{}
)

func newDatabaseDataSource() datasource.DataSource {
	return &databaseDataSource{}
}

type databaseDataSource struct {
	client *planetscale.Client
}

type TTdatabaseDataSourceModel struct {
	Organization string `tfsdk:"organization"`

	databaseDataSourceModel
}

func (d *databaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *databaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: databaseDataSourceSchemaAttribute(false)}
}

func (d *databaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *databaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *TTdatabaseDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.GetDatabase(ctx, data.Organization, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Received a nil database", "")
		return
	}
	state := TTdatabaseDataSourceModel{
		Organization: data.Organization,
	}
	resp.Diagnostics.Append(state.fromClient(&res.Database)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
