package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &databasesDataSource{}
	_ datasource.DataSourceWithConfigure = &databasesDataSource{}
)

func newDatabasesDataSource() datasource.DataSource {
	return &databasesDataSource{}
}

type databasesDataSource struct {
	client *planetscale.Client
}

type databasesDataSourceModel struct {
	Organization string                    `tfsdk:"organization"`
	Databases    []databaseDataSourceModel `tfsdk:"databases"`
}

func (d *databasesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_databases"
}

func (d *databasesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A list of PlanetScale databases.",
		MarkdownDescription: "A list of PlanetScale databases.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true},
			"databases": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: databaseDataSourceSchemaAttribute(true),
				},
			},
		},
	}
}

func (d *databasesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *databasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *databasesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgName := data.Organization
	res, err := d.client.ListDatabases(ctx, orgName, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read databases", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read databases", "no data")
		return
	}
	state := databasesDataSourceModel{
		Organization: data.Organization,
		Databases:    make([]databaseDataSourceModel, 0, len(res.Data)),
	}
	for _, item := range res.Data {
		item := item
		state.Databases = append(state.Databases, *databaseFromClient(&item, orgName, resp.Diagnostics))
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
