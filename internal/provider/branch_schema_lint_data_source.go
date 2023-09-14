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

	Errors []branchSchemaLintErrorDataSourceModel `tfsdk:"errors"`
}

type branchSchemaLintErrorDataSourceModel struct {
	AutoIncrementColumnNames []string `tfsdk:"auto_increment_column_names"`
	CharsetName              string   `tfsdk:"charset_name"`
	CheckConstraintName      string   `tfsdk:"check_constraint_name"`
	ColumnName               string   `tfsdk:"column_name"`
	DocsUrl                  string   `tfsdk:"docs_url"`
	EngineName               string   `tfsdk:"engine_name"`
	EnumValue                string   `tfsdk:"enum_value"`
	ErrorDescription         string   `tfsdk:"error_description"`
	ForeignKeyColumnNames    []string `tfsdk:"foreign_key_column_names"`
	JsonPath                 string   `tfsdk:"json_path"`
	KeyspaceName             string   `tfsdk:"keyspace_name"`
	LintError                string   `tfsdk:"lint_error"`
	PartitionName            string   `tfsdk:"partition_name"`
	PartitioningType         string   `tfsdk:"partitioning_type"`
	SubjectType              string   `tfsdk:"subject_type"`
	TableName                string   `tfsdk:"table_name"`
	VindexName               string   `tfsdk:"vindex_name"`
}

func (d *branchSchemaLintDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch_schema_lint"
}

func (d *branchSchemaLintDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"database":     schema.StringAttribute{Required: true},
		"branch":       schema.StringAttribute{Required: true},
		"errors": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"auto_increment_column_names": schema.ListAttribute{Computed: true, ElementType: types.StringType},
					"charset_name":                schema.StringAttribute{Computed: true},
					"check_constraint_name":       schema.StringAttribute{Computed: true},
					"column_name":                 schema.StringAttribute{Computed: true},
					"docs_url":                    schema.StringAttribute{Computed: true},
					"engine_name":                 schema.StringAttribute{Computed: true},
					"enum_value":                  schema.StringAttribute{Computed: true},
					"error_description":           schema.StringAttribute{Computed: true},
					"foreign_key_column_names":    schema.ListAttribute{Computed: true, ElementType: types.StringType},
					"json_path":                   schema.StringAttribute{Computed: true},
					"keyspace_name":               schema.StringAttribute{Computed: true},
					"lint_error":                  schema.StringAttribute{Computed: true},
					"partition_name":              schema.StringAttribute{Computed: true},
					"partitioning_type":           schema.StringAttribute{Computed: true},
					"subject_type":                schema.StringAttribute{Computed: true},
					"table_name":                  schema.StringAttribute{Computed: true},
					"vindex_name":                 schema.StringAttribute{Computed: true},
				},
			},
		},
	}}
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
		Errors:       make([]branchSchemaLintErrorDataSourceModel, 0, len(res200.Data)),
	}
	for _, item := range res200.Data {
		state.Errors = append(state.Errors, branchSchemaLintErrorDataSourceModel{
			AutoIncrementColumnNames: item.AutoIncrementColumnNames,
			CharsetName:              item.CharsetName,
			CheckConstraintName:      item.CheckConstraintName,
			ColumnName:               item.ColumnName,
			DocsUrl:                  item.DocsUrl,
			EngineName:               item.EngineName,
			EnumValue:                item.EnumValue,
			ErrorDescription:         item.ErrorDescription,
			ForeignKeyColumnNames:    item.ForeignKeyColumnNames,
			JsonPath:                 item.JsonPath,
			KeyspaceName:             item.KeyspaceName,
			LintError:                item.LintError,
			PartitionName:            item.PartitionName,
			PartitioningType:         item.PartitioningType,
			SubjectType:              item.SubjectType,
			TableName:                item.TableName,
			VindexName:               item.VindexName,
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
