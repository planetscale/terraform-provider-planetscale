package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &branchSchemaDataSource{}
	_ datasource.DataSourceWithConfigure = &branchSchemaDataSource{}
)

func newBranchSchemaDataSource() datasource.DataSource {
	return &branchSchemaDataSource{}
}

type branchSchemaDataSource struct {
	client *planetscale.Client
}

type branchSchemaDataSourceModel struct {
	Organization string       `tfsdk:"organization"`
	Database     string       `tfsdk:"database"`
	Branch       string       `tfsdk:"branch"`
	Keyspace     types.String `tfsdk:"keyspace"`

	Tables []tableSchemaDataSourceModel `tfsdk:"tables"`
}

func (d *branchSchemaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch_schema"
}

func (d *branchSchemaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"database":     schema.StringAttribute{Required: true},
		"branch":       schema.StringAttribute{Required: true},
		"keyspace":     schema.StringAttribute{Optional: true},
		"tables": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: tableSchemaDataSourceSchemaAttribute,
			},
		},
	}}
}

func (d *branchSchemaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *branchSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *branchSchemaDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.GetBranchSchema(ctx, data.Organization, data.Database, data.Branch, stringValueIfKnown(data.Keyspace))
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database branch schema", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read database branch schema", "no data")
		return
	}
	state := branchSchemaDataSourceModel{
		Organization: data.Organization,
		Database:     data.Database,
		Branch:       data.Branch,
		Keyspace:     data.Keyspace,
		Tables:       make([]tableSchemaDataSourceModel, 0, len(res.Data)),
	}
	for _, item := range res.Data {
		item := item
		el := tableSchemaDataSourceModel{}
		resp.Diagnostics.Append(el.fromClient(&item)...)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
