package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &branchSafeMigrationsDataSource{}
	_ datasource.DataSourceWithConfigure = &branchSafeMigrationsDataSource{}
)

func newbranchSafeMigrationsDataSource() datasource.DataSource {
	return &branchSafeMigrationsDataSource{}
}

type branchSafeMigrationsDataSource struct {
	client *planetscale.Client
}

func (d *branchSafeMigrationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch_safe_migrations"
}

func (d *branchSafeMigrationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Safe migration state on a PlanetScale branch.",
		MarkdownDescription: "Safe migration state on a PlanetScale branch.",
		Attributes:          branchSafeMigrationsDataSourceSchemaAttribute,
	}
}

func (d *branchSafeMigrationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *branchSafeMigrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *branchSafeMigrationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetBranch(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read branch safe migrations status", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read branch safe migrations status", "no data")
		return
	}

	state := branchSafeMigrationsFromClient(&res.Branch, data.Organization.ValueString(), data.Database.ValueString())
	if resp.Diagnostics.HasError() {
		return
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
