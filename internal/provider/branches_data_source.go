package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &branchesDataSource{}
	_ datasource.DataSourceWithConfigure = &branchesDataSource{}
)

func newBranchesDataSource() datasource.DataSource {
	return &branchesDataSource{}
}

type branchesDataSource struct {
	client *planetscale.Client
}

type branchesDataSourceModel struct {
	Organization string                  `tfsdk:"organization"`
	Database     string                  `tfsdk:"database"`
	Branches     []branchDataSourceModel `tfsdk:"branches"`
}

func (d *branchesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branches"
}

func (d *branchesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A list of PlanetScale branches.",
		MarkdownDescription: "A list of PlanetScale branches.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true},
			"database":     schema.StringAttribute{Required: true},
			"branches": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: branchDataSourceSchemaAttribute(true),
				},
			},
		},
	}
}

func (d *branchesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *branchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data *branchesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.ListBranches(ctx, data.Organization, data.Database, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database branches", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read database branches", "no data")
		return
	}
	state := branchesDataSourceModel{
		Organization: data.Organization,
		Database:     data.Database,
		Branches:     make([]branchDataSourceModel, 0, len(res.Data)),
	}

	for _, item := range res.Data {
		item := item
		state.Branches = append(state.Branches, *branchFromClient(&item, data.Organization, data.Database, resp.Diagnostics))
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
