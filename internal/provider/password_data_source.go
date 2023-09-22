package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &passwordDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordDataSource{}
)

func newPasswordDataSource() datasource.DataSource {
	return &passwordDataSource{}
}

type passwordDataSource struct {
	client *planetscale.Client
}

func (d *passwordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

func (d *passwordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A PlanetScale database password.",
		MarkdownDescription: "A PlanetScale database password.",
		Attributes:          passwordDataSourceSchemaAttribute(false),
	}
}

func (d *passwordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *passwordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *passwordDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.GetPassword(
		ctx,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
		data.Id.ValueString(),
		data.ReadOnlyRegionId.ValueStringPointer(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database password", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read database password", "no data")
		return
	}
	state := passwordFromClient(
		&res.Password,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
		data.ReadOnlyRegionId.ValueStringPointer(),
		resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
