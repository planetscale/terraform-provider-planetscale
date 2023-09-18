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
	_ datasource.DataSource              = &databaseDataSource{}
	_ datasource.DataSourceWithConfigure = &databaseDataSource{}
)

func newDatabaseDataSource() datasource.DataSource {
	return &databaseDataSource{}
}

type databaseDataSource struct {
	client *planetscale.Client
}

type importDataSourceDataSourceModel struct {
	Database string `tfsdk:"database"`
	Hostname string `tfsdk:"hostname"`
	Port     string `tfsdk:"port"`
}

type importDataSourceModel struct {
	DataSource        importDataSourceDataSourceModel `tfsdk:"data_source"`
	FinishedAt        string                          `tfsdk:"finished_at"`
	ImportCheckErrors string                          `tfsdk:"import_check_errors"`
	StartedAt         string                          `tfsdk:"started_at"`
	State             string                          `tfsdk:"state"`
}

type regionDataSourceModel struct {
	DisplayName       string   `tfsdk:"display_name"`
	Enabled           bool     `tfsdk:"enabled"`
	Id                string   `tfsdk:"id"`
	Location          string   `tfsdk:"location"`
	Provider          string   `tfsdk:"provider"`
	PublicIpAddresses []string `tfsdk:"public_ip_addresses"`
	Slug              string   `tfsdk:"slug"`
}

type databaseDataSourceModel struct {
	Organization                      string                 `tfsdk:"organization"`
	Name                              string                 `tfsdk:"name"`
	Id                                types.String           `tfsdk:"id"`
	AllowDataBranching                types.Bool             `tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      types.Bool             `tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          types.Bool             `tfsdk:"at_development_branch_limit"`
	AutomaticMigrations               types.Bool             `tfsdk:"automatic_migrations"`
	BranchesCount                     types.Float64          `tfsdk:"branches_count"`
	BranchesUrl                       types.String           `tfsdk:"branches_url"`
	CreatedAt                         types.String           `tfsdk:"created_at"`
	DataImport                        *importDataSourceModel `tfsdk:"data_import"`
	DefaultBranch                     types.String           `tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount types.Float64          `tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           types.Float64          `tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           types.Float64          `tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          types.Float64          `tfsdk:"development_branches_count"`
	HtmlUrl                           types.String           `tfsdk:"html_url"`
	InsightsRawQueries                types.Bool             `tfsdk:"insights_raw_queries"`
	IssuesCount                       types.Float64          `tfsdk:"issues_count"`
	MigrationFramework                types.String           `tfsdk:"migration_framework"`
	MigrationTableName                types.String           `tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion types.Bool             `tfsdk:"multiple_admins_required_for_deletion"`
	Plan                              types.String           `tfsdk:"plan"`
	ProductionBranchWebConsole        types.Bool             `tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           types.Float64          `tfsdk:"production_branches_count"`
	Ready                             types.Bool             `tfsdk:"ready"`
	Region                            *regionDataSourceModel `tfsdk:"region"`
	RequireApprovalForDeploy          types.Bool             `tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              types.Bool             `tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               types.String           `tfsdk:"schema_last_updated_at"`
	Sharded                           types.Bool             `tfsdk:"sharded"`
	State                             types.String           `tfsdk:"state"`
	Type                              types.String           `tfsdk:"type"`
	UpdatedAt                         types.String           `tfsdk:"updated_at"`
	Url                               types.String           `tfsdk:"url"`
}

func (d *databaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *databaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization":                     schema.StringAttribute{Required: true},
		"name":                             schema.StringAttribute{Required: true},
		"id":                               schema.StringAttribute{Computed: true},
		"allow_data_branching":             schema.BoolAttribute{Computed: true, Optional: true},
		"at_backup_restore_branches_limit": schema.BoolAttribute{Computed: true},
		"at_development_branch_limit":      schema.BoolAttribute{Computed: true},
		"automatic_migrations":             schema.BoolAttribute{Computed: true, Optional: true},
		"branches_count":                   schema.Float64Attribute{Computed: true},
		"branches_url":                     schema.StringAttribute{Computed: true},
		"created_at":                       schema.StringAttribute{Computed: true},
		"data_import": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"data_source": schema.SingleNestedAttribute{
					Computed: true,
					Attributes: map[string]schema.Attribute{
						"database": schema.StringAttribute{Required: true},
						"hostname": schema.StringAttribute{Required: true},
						"port":     schema.StringAttribute{Required: true},
					},
				},
				"finished_at":         schema.StringAttribute{Computed: true},
				"import_check_errors": schema.StringAttribute{Computed: true},
				"started_at":          schema.StringAttribute{Computed: true},
				"state":               schema.StringAttribute{Computed: true},
			},
		},
		"default_branch":                         schema.StringAttribute{Computed: true, Optional: true},
		"default_branch_read_only_regions_count": schema.Float64Attribute{Computed: true},
		"default_branch_shard_count":             schema.Float64Attribute{Computed: true},
		"default_branch_table_count":             schema.Float64Attribute{Computed: true},
		"development_branches_count":             schema.Float64Attribute{Computed: true},
		"html_url":                               schema.StringAttribute{Computed: true},
		"insights_raw_queries":                   schema.BoolAttribute{Computed: true, Optional: true},
		"issues_count":                           schema.Float64Attribute{Computed: true, Optional: true},
		"migration_framework":                    schema.StringAttribute{Computed: true, Optional: true},
		"migration_table_name":                   schema.StringAttribute{Computed: true, Optional: true},
		"multiple_admins_required_for_deletion":  schema.BoolAttribute{Computed: true, Optional: true},
		"notes":                                  schema.StringAttribute{Computed: true, Optional: true},
		"plan":                                   schema.StringAttribute{Computed: true, Optional: true},
		"production_branch_web_console":          schema.BoolAttribute{Computed: true, Optional: true},
		"production_branches_count":              schema.Float64Attribute{Computed: true},
		"ready":                                  schema.BoolAttribute{Computed: true},
		"region": schema.SingleNestedAttribute{
			Computed: true, Optional: true,
			Attributes: map[string]schema.Attribute{
				"slug": schema.StringAttribute{
					// the actual argument to set the region to use
					Computed: true, Optional: true,
				},
				"display_name":        schema.StringAttribute{Computed: true},
				"enabled":             schema.BoolAttribute{Computed: true},
				"id":                  schema.StringAttribute{Computed: true},
				"location":            schema.StringAttribute{Computed: true},
				"provider":            schema.StringAttribute{Computed: true},
				"public_ip_addresses": schema.ListAttribute{Computed: true, ElementType: types.StringType},
			},
		},
		"require_approval_for_deploy": schema.BoolAttribute{Computed: true, Optional: true},
		"restrict_branch_region":      schema.BoolAttribute{Computed: true, Optional: true},
		"schema_last_updated_at":      schema.StringAttribute{Computed: true},
		"sharded":                     schema.BoolAttribute{Computed: true},
		"state":                       schema.StringAttribute{Computed: true},
		"type":                        schema.StringAttribute{Computed: true},
		"updated_at":                  schema.StringAttribute{Computed: true},
		"url":                         schema.StringAttribute{Computed: true},
	}}
}

func (d *databaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *databaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *databaseDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res200, err := d.client.GetDatabase(ctx, data.Organization, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database", err.Error())
		return
	}
	if res200 == nil {
		resp.Diagnostics.AddError("Received a nil database", "")
		return
	}
	state := databaseDataSourceModel{
		Organization:                      data.Organization,
		Name:                              res200.Name,
		Id:                                types.StringValue(res200.Id),
		AllowDataBranching:                types.BoolValue(res200.AllowDataBranching),
		AtBackupRestoreBranchesLimit:      types.BoolValue(res200.AtBackupRestoreBranchesLimit),
		AtDevelopmentBranchLimit:          types.BoolValue(res200.AtDevelopmentBranchLimit),
		AutomaticMigrations:               types.BoolPointerValue(res200.AutomaticMigrations),
		BranchesCount:                     types.Float64Value(res200.BranchesCount),
		BranchesUrl:                       types.StringValue(res200.BranchesUrl),
		CreatedAt:                         types.StringValue(res200.CreatedAt),
		DefaultBranch:                     types.StringValue(res200.DefaultBranch),
		DefaultBranchReadOnlyRegionsCount: types.Float64Value(res200.DefaultBranchReadOnlyRegionsCount),
		DefaultBranchShardCount:           types.Float64Value(res200.DefaultBranchShardCount),
		DefaultBranchTableCount:           types.Float64Value(res200.DefaultBranchTableCount),
		DevelopmentBranchesCount:          types.Float64Value(res200.DevelopmentBranchesCount),
		HtmlUrl:                           types.StringValue(res200.HtmlUrl),
		InsightsRawQueries:                types.BoolValue(res200.InsightsRawQueries),
		IssuesCount:                       types.Float64Value(res200.IssuesCount),
		MigrationFramework:                types.StringPointerValue(res200.MigrationFramework),
		MigrationTableName:                types.StringPointerValue(res200.MigrationTableName),
		MultipleAdminsRequiredForDeletion: types.BoolValue(res200.MultipleAdminsRequiredForDeletion),
		Plan:                              types.StringValue(res200.Plan),
		ProductionBranchWebConsole:        types.BoolValue(res200.ProductionBranchWebConsole),
		ProductionBranchesCount:           types.Float64Value(res200.ProductionBranchesCount),
		Ready:                             types.BoolValue(res200.Ready),
		RequireApprovalForDeploy:          types.BoolValue(res200.RequireApprovalForDeploy),
		RestrictBranchRegion:              types.BoolValue(res200.RestrictBranchRegion),
		SchemaLastUpdatedAt:               types.StringPointerValue(res200.SchemaLastUpdatedAt),
		Sharded:                           types.BoolValue(res200.Sharded),
		State:                             types.StringValue(res200.State),
		Type:                              types.StringValue(res200.Type),
		UpdatedAt:                         types.StringValue(res200.UpdatedAt),
		Url:                               types.StringValue(res200.Url),
	}

	if res200.DataImport != nil {
		data.DataImport = &importDataSourceModel{
			DataSource: importDataSourceDataSourceModel{
				Database: res200.DataImport.DataSource.Database,
				Hostname: res200.DataImport.DataSource.Hostname,
				Port:     res200.DataImport.DataSource.Port,
			},
			FinishedAt:        res200.DataImport.FinishedAt,
			ImportCheckErrors: res200.DataImport.ImportCheckErrors,
			StartedAt:         res200.DataImport.StartedAt,
			State:             res200.DataImport.State,
		}
	}
	data.Region = &regionDataSourceModel{
		DisplayName:       res200.Region.DisplayName,
		Enabled:           res200.Region.Enabled,
		Id:                res200.Region.Id,
		Location:          res200.Region.Location,
		Provider:          res200.Region.Provider,
		PublicIpAddresses: res200.Region.PublicIpAddresses,
		Slug:              res200.Region.Slug,
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
