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
	_ datasource.DataSource              = &databasesDataSource{}
	_ datasource.DataSourceWithConfigure = &databasesDataSource{}
)

func newDatabasesDataSource() datasource.DataSource {
	return &databasesDataSource{}
}

type databasesDataSource struct {
	client *planetscale.Client
}

type databasesDataSourceModel struct {
	Organization string                    `tfsdk:"organization"`
	Databases    []databaseDataSourceModel `tfsdk:"databases"`
}

func (d *databasesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_databases"
}

func (d *databasesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"databases": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"organization":                     schema.StringAttribute{Computed: true},
					"name":                             schema.StringAttribute{Computed: true},
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
				},
			},
		},
	}}
}

func (d *databasesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *databasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data *databasesDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgName := data.Organization

	res200, err := d.client.ListDatabases(ctx, orgName, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read databases", err.Error())
		return
	}
	if res200 == nil {
		resp.Diagnostics.AddError("Unable to read databases", "no data")
		return
	}
	dbs := make([]databaseDataSourceModel, 0, len(res200.Data))
	for _, item := range res200.Data {

		out := databaseDataSourceModel{
			Organization:                      data.Organization,
			Id:                                item.Id,
			AllowDataBranching:                item.AllowDataBranching,
			AtBackupRestoreBranchesLimit:      item.AtBackupRestoreBranchesLimit,
			AtDevelopmentBranchLimit:          item.AtDevelopmentBranchLimit,
			AutomaticMigrations:               item.AutomaticMigrations,
			BranchesCount:                     item.BranchesCount,
			BranchesUrl:                       item.BranchesUrl,
			CreatedAt:                         item.CreatedAt,
			DefaultBranch:                     item.DefaultBranch,
			DefaultBranchReadOnlyRegionsCount: item.DefaultBranchReadOnlyRegionsCount,
			DefaultBranchShardCount:           item.DefaultBranchShardCount,
			DefaultBranchTableCount:           item.DefaultBranchTableCount,
			DevelopmentBranchesCount:          item.DevelopmentBranchesCount,
			HtmlUrl:                           item.HtmlUrl,
			InsightsRawQueries:                item.InsightsRawQueries,
			IssuesCount:                       item.IssuesCount,
			MigrationFramework:                item.MigrationFramework,
			MigrationTableName:                item.MigrationTableName,
			MultipleAdminsRequiredForDeletion: item.MultipleAdminsRequiredForDeletion,
			Name:                              item.Name,
			Notes:                             item.Notes,
			Plan:                              item.Plan,
			ProductionBranchWebConsole:        item.ProductionBranchWebConsole,
			ProductionBranchesCount:           item.ProductionBranchesCount,
			Ready:                             item.Ready,
			RequireApprovalForDeploy:          item.RequireApprovalForDeploy,
			RestrictBranchRegion:              item.RestrictBranchRegion,
			SchemaLastUpdatedAt:               item.SchemaLastUpdatedAt,
			Sharded:                           item.Sharded,
			State:                             item.State,
			Type:                              item.Type,
			UpdatedAt:                         item.UpdatedAt,
			Url:                               item.Url,
		}

		if item.DataImport != nil {
			out.DataImport = &importDataSourceModel{
				DataSource: importDataSourceDataSourceModel{
					Database: item.DataImport.DataSource.Database,
					Hostname: item.DataImport.DataSource.Hostname,
					Port:     item.DataImport.DataSource.Port,
				},
				FinishedAt:        item.DataImport.FinishedAt,
				ImportCheckErrors: item.DataImport.ImportCheckErrors,
				StartedAt:         item.DataImport.StartedAt,
				State:             item.DataImport.State,
			}
		}
		out.Region = regionDataSourceModel{
			DisplayName:       item.Region.DisplayName,
			Enabled:           item.Region.Enabled,
			Id:                item.Region.Id,
			Location:          item.Region.Location,
			Provider:          item.Region.Provider,
			PublicIpAddresses: item.Region.PublicIpAddresses,
			Slug:              item.Region.Slug,
		}

		dbs = append(dbs, out)
	}
	state := databasesDataSourceModel{Organization: data.Organization, Databases: dbs}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
