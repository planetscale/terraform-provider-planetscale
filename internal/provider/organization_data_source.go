package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/planetscale-go/planetscale"
)

var (
	_ datasource.DataSource              = &organizationDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationDataSource{}
)

func newOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

type organizationDataSource struct {
	client *planetscale.Client
}

type organizationDataSourceModel struct {
	Name                   string `tfsdk:"name"`
	CreatedAt              string `tfsdk:"created_at"`
	UpdatedAt              string `tfsdk:"updated_at"`
	FreeDatabasesRemaining int64  `tfsdk:"free_databases_remaining"`
}

func (d *organizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *organizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: orgSchemaAttributes(true)}
}

func (d *organizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *organizationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := d.client.Organizations.Get(ctx, &planetscale.GetOrganizationRequest{
		Organization: data.Name,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to read organization", err.Error())
		return
	}
	if org == nil {
		resp.Diagnostics.AddError("Received a nil organization", "")
		return
	}
	state := orgToModel(*org)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func orgToModel(in planetscale.Organization) organizationDataSourceModel {
	return organizationDataSourceModel{
		Name:                   in.Name,
		CreatedAt:              in.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:              in.UpdatedAt.Format(time.RFC3339Nano),
		FreeDatabasesRemaining: int64(in.RemainingFreeDatabases),
	}
}
