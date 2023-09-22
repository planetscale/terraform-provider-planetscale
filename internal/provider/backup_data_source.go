package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &backupDataSource{}
	_ datasource.DataSourceWithConfigure = &backupDataSource{}
)

func newBackupDataSource() datasource.DataSource {
	return &backupDataSource{}
}

type backupDataSource struct {
	client *planetscale.Client
}

func (d *backupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

func (d *backupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A PlanetScale backup.",
		MarkdownDescription: "A PlanetScale backup.",
		Attributes:          backupDataSourceSchemaAttribute(false),
	}
}

func (d *backupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *backupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *backupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.GetBackup(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database branch backup", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read database branch backup", "no data")
		return
	}
	state := backupFromClient(&res.Backup, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString(), resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
