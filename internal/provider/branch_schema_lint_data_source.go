package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &branchSchemaLintDataSource{}
	_ datasource.DataSourceWithConfigure = &branchSchemaLintDataSource{}
)

func newBranchSchemaLintDataSource() datasource.DataSource {
	return &branchSchemaLintDataSource{}
}

type branchSchemaLintDataSource struct {
	client *planetscale.Client
}

type branchSchemaLintDataSourceModel struct {
	Organization string `tfsdk:"organization"`
	Database     string `tfsdk:"database"`
	Branch       string `tfsdk:"branch"`

	Errors []lintErrorDataSourceModel `tfsdk:"errors"`
}

func (d *branchSchemaLintDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch_schema_lint"
}

func (d *branchSchemaLintDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Linting errors for the schema of a PlanetScale branch.",
		MarkdownDescription: "Linting errors for the schema of a PlanetScale branch.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true},
			"database":     schema.StringAttribute{Required: true},
			"branch":       schema.StringAttribute{Required: true},
			"errors": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: lintErrorDataSourceSchemaAttribute,
				},
			},
		},
	}
}

func (d *branchSchemaLintDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *branchSchemaLintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data *branchSchemaLintDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res200, err := d.client.LintBranchSchema(ctx, data.Organization, data.Database, data.Branch, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database branch schema", err.Error())
		return
	}
	if res200 == nil {
		resp.Diagnostics.AddError("Unable to read database branch schema", "no data")
		return
	}
	state := branchSchemaLintDataSourceModel{
		Organization: data.Organization,
		Database:     data.Database,
		Branch:       data.Branch,
		Errors:       make([]lintErrorDataSourceModel, 0, len(res200.Data)),
	}
	for _, item := range res200.Data {
		item := item
		state.Errors = append(state.Errors, *lintErrorFromClient(&item, resp.Diagnostics))
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
