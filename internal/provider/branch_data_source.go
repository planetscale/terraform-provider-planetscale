package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &branchDataSource{}
	_ datasource.DataSourceWithConfigure = &branchDataSource{}
)

func newBranchDataSource() datasource.DataSource {
	return &branchDataSource{}
}

type branchDataSource struct {
	client *planetscale.Client
}

func (d *branchDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

func (d *branchDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A PlanetScale branch.",
		MarkdownDescription: "A PlanetScale branch.",
		Attributes:          branchDataSourceSchemaAttribute(false),
	}
}

func (d *branchDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *branchDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *branchDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.GetBranch(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database branch", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read database branch", "no data")
		return
	}
	state := branchFromClient(&res.Branch, data.Organization.ValueString(), data.Database.ValueString(), resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
