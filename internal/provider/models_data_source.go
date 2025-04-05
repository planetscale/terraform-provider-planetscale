package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

func organizationDataSourceSchemaAttribute(computedName bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "The name of the organization.",
			Required:    !computedName, Computed: computedName,
		},
		"billing_email": schema.StringAttribute{
			Description: "The billing email of the organization.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "When the organization was created.",
			Computed:    true,
		},
		"database_count": schema.Float64Attribute{
			Description: "The number of databases in the organization.",
			Computed:    true,
		},
		"features": schema.SingleNestedAttribute{
			Description: "Features that are enabled on the organization.",
			Computed:    true,
			Attributes:  featuresDataSourceSchemaAttribute,
		},
		"flags": schema.SingleNestedAttribute{
			Description: ".",
			Computed:    true,
			Attributes:  flagsDataSourceSchemaAttribute,
		},
		"has_past_due_invoices": schema.BoolAttribute{
			Description: "Whether or not the organization has past due billing invoices.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The ID for the organization.",
			Computed:    true,
		},
		"idp_managed_roles": schema.BoolAttribute{
			Description: "Whether or not the IdP provider is be responsible for managing roles in PlanetScale.",
			Computed:    true,
		},
		"plan": schema.StringAttribute{
			Description: "The billing plan of the organization.",
			Computed:    true,
		},
		"single_tenancy": schema.BoolAttribute{
			Description: "Whether or not the organization has single tenancy enabled.",
			Computed:    true,
		},
		"sso": schema.BoolAttribute{
			Description: "Whether or not SSO is enabled on the organization.",
			Computed:    true,
		},
		"sso_directory": schema.BoolAttribute{
			Description: "Whether or not the organization uses a WorkOS directory.",
			Computed:    true,
		},
		"sso_portal_url": schema.StringAttribute{
			Description: "The URL of the organization's SSO portal.",
			Computed:    true,
		},
		"updated_at": schema.StringAttribute{
			Description: "When the organization was last updated.",
			Computed:    true,
		},
		"valid_billing_info": schema.BoolAttribute{
			Description: "Whether or not the organization's billing information is valid.",
			Computed:    true,
		},
	}
}

type organizationDataSourceModel struct {
	BillingEmail       types.String             `tfsdk:"billing_email"`
	CreatedAt          types.String             `tfsdk:"created_at"`
	DatabaseCount      types.Float64            `tfsdk:"database_count"`
	Features           *featuresDataSourceModel `tfsdk:"features"`
	Flags              *flagsDataSourceModel    `tfsdk:"flags"`
	HasPastDueInvoices types.Bool               `tfsdk:"has_past_due_invoices"`
	Id                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	Plan               types.String             `tfsdk:"plan"`
	SingleTenancy      types.Bool               `tfsdk:"single_tenancy"`
	Sso                types.Bool               `tfsdk:"sso"`
	SsoDirectory       types.Bool               `tfsdk:"sso_directory"`
	SsoPortalUrl       types.String             `tfsdk:"sso_portal_url"`
	UpdatedAt          types.String             `tfsdk:"updated_at"`
	ValidBillingInfo   types.Bool               `tfsdk:"valid_billing_info"`
	IdpManagedRoles    types.Bool               `tfsdk:"idp_managed_roles"`
}

func organizationFromClient(org *planetscale.Organization) *organizationDataSourceModel {
	if org == nil {
		return nil
	}
	return &organizationDataSourceModel{
		Features:           featuresFromClient(org.Features),
		Flags:              flagsFromClient(org.Flags),
		BillingEmail:       types.StringPointerValue(org.BillingEmail),
		CreatedAt:          types.StringValue(org.CreatedAt),
		DatabaseCount:      types.Float64Value(org.DatabaseCount),
		HasPastDueInvoices: types.BoolValue(org.HasPastDueInvoices),
		Id:                 types.StringValue(org.Id),
		Name:               types.StringValue(org.Name),
		Plan:               types.StringValue(org.Plan),
		SingleTenancy:      types.BoolValue(org.SingleTenancy),
		Sso:                types.BoolValue(org.Sso),
		SsoDirectory:       types.BoolValue(org.SsoDirectory),
		SsoPortalUrl:       types.StringPointerValue(org.SsoPortalUrl),
		UpdatedAt:          types.StringValue(org.UpdatedAt),
		ValidBillingInfo:   types.BoolValue(org.ValidBillingInfo),
		IdpManagedRoles:    types.BoolValue(org.IdpManagedRoles),
	}
}

var featuresDataSourceSchemaAttribute = map[string]schema.Attribute{
	"insights":       schema.BoolAttribute{Computed: true},
	"single_tenancy": schema.BoolAttribute{Computed: true},
	"sso":            schema.BoolAttribute{Computed: true},
}

type featuresDataSourceModel struct {
	Insights      types.Bool `tfsdk:"insights"`
	SingleTenancy types.Bool `tfsdk:"single_tenancy"`
	Sso           types.Bool `tfsdk:"sso"`
}

func featuresFromClient(features *planetscale.Features) *featuresDataSourceModel {
	if features == nil {
		return nil
	}
	return &featuresDataSourceModel{
		Insights:      types.BoolPointerValue(features.Insights),
		SingleTenancy: types.BoolPointerValue(features.SingleTenancy),
		Sso:           types.BoolPointerValue(features.Sso),
	}
}

var flagsDataSourceSchemaAttribute = map[string]schema.Attribute{
	"example_flag": schema.StringAttribute{Computed: true},
}

type flagsDataSourceModel struct {
	ExampleFlag types.String `tfsdk:"example_flag"`
}

func flagsFromClient(flags *planetscale.Flags) *flagsDataSourceModel {
	if flags == nil {
		return nil
	}
	return &flagsDataSourceModel{
		ExampleFlag: types.StringPointerValue(flags.ExampleFlag),
	}
}

type dataSourceDataSourceModel struct {
	Database types.String  `tfsdk:"database"`
	Hostname types.String  `tfsdk:"hostname"`
	Port     types.Float64 `tfsdk:"port"`
}

func dataSourceFromClient(dataSource planetscale.DataSource) dataSourceDataSourceModel {
	return dataSourceDataSourceModel{
		Database: types.StringValue(dataSource.Database),
		Hostname: types.StringValue(dataSource.Hostname),
		Port:     types.Float64Value(dataSource.Port),
	}
}

type dataImportDataSourceModel struct {
	DataSource        dataSourceDataSourceModel `tfsdk:"data_source"`
	FinishedAt        types.String              `tfsdk:"finished_at"`
	ImportCheckErrors types.String              `tfsdk:"import_check_errors"`
	StartedAt         types.String              `tfsdk:"started_at"`
	State             types.String              `tfsdk:"state"`
}

func dataImportFromClient(dataImport *planetscale.DataImport) *dataImportDataSourceModel {
	if dataImport == nil {
		return nil
	}
	return &dataImportDataSourceModel{
		DataSource:        dataSourceFromClient(dataImport.DataSource),
		FinishedAt:        types.StringValue(dataImport.FinishedAt),
		ImportCheckErrors: types.StringValue(dataImport.ImportCheckErrors),
		StartedAt:         types.StringValue(dataImport.StartedAt),
		State:             types.StringValue(dataImport.State),
	}
}

func databaseDataSourceSchemaAttribute(computedName bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"organization": schema.StringAttribute{
			Description: "The organization this database belongs to.",
			Required:    !computedName, Computed: computedName,
		},
		"name": schema.StringAttribute{
			Description: "The name of this database.",
			Required:    !computedName, Computed: computedName,
		},
		"id": schema.StringAttribute{
			Description: "The ID of the database.",
			Computed:    true,
		},
		"allow_data_branching": schema.BoolAttribute{
			Description: "Whether seeding branches with data is enabled for all branches.",
			Computed:    true, Optional: true,
		},
		"at_backup_restore_branches_limit": schema.BoolAttribute{
			Description: "If the database has reached its backup restored branch limit.",
			Computed:    true,
		},
		"at_development_branch_limit": schema.BoolAttribute{
			Description: "If the database has reached its development branch limit.",
			Computed:    true,
		},
		"automatic_migrations": schema.BoolAttribute{
			Description: "Whether to automatically manage Rails migrations during deploy requests.",
			Computed:    true, Optional: true,
		},
		"branches_count": schema.Float64Attribute{
			Description: "The total number of database branches.",
			Computed:    true,
		},
		"branches_url": schema.StringAttribute{
			Description: "The URL to retrieve this database's branches via the API.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "When the database was created.",
			Computed:    true,
		},
		"data_import": schema.SingleNestedAttribute{
			Description: "If the database was created from an import, describes the import process.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"data_source": schema.SingleNestedAttribute{
					Description: "Connection information for the source of the data for the import.",
					Computed:    true,
					Attributes: map[string]schema.Attribute{
						"database": schema.StringAttribute{
							Description: "The name of the database imported from.",
							Required:    true,
						},
						"hostname": schema.StringAttribute{
							Description: "The hostname where the database lives.",
							Required:    true,
						},
						"port": schema.StringAttribute{
							Description: "The port on which the database listens on the host.",
							Required:    true,
						},
					},
				},
				"finished_at": schema.StringAttribute{
					Description: "When the import finished.",
					Computed:    true,
				},
				"import_check_errors": schema.StringAttribute{
					Description: "Errors encountered while preparing the import.",
					Computed:    true,
				},
				"started_at": schema.StringAttribute{
					Description: "When the import started.",
					Computed:    true,
				},
				"state": schema.StringAttribute{
					Description: "The state of the import, one of: pending, queued, in_progress, complete, cancelled, error.",
					Computed:    true,
				},
			},
		},
		"default_branch": schema.StringAttribute{
			Description: "The default branch for the database.",
			Computed:    true, Optional: true,
		},
		"default_branch_read_only_regions_count": schema.Float64Attribute{
			Description: "Number of read only regions in the default branch.",
			Computed:    true,
		},
		"default_branch_shard_count": schema.Float64Attribute{
			Description: "Number of shards in the default branch.",
			Computed:    true,
		},
		"default_branch_table_count": schema.Float64Attribute{
			Description: "Number of tables in the default branch schema.",
			Computed:    true,
		},
		"development_branches_count": schema.Float64Attribute{
			Description: "The total number of database development branches.",
			Computed:    true,
		},
		"html_url": schema.StringAttribute{
			Description: "The total number of database development branches.",
			Computed:    true,
		},
		"insights_raw_queries": schema.BoolAttribute{
			Description: "The URL to see this database's branches in the web UI.",
			Computed:    true, Optional: true,
		},
		"issues_count": schema.Float64Attribute{
			Description: "The total number of ongoing issues within a database.",
			Computed:    true, Optional: true,
		},
		"migration_framework": schema.StringAttribute{
			Description: "Framework used for applying migrations.",
			Computed:    true, Optional: true,
		},
		"migration_table_name": schema.StringAttribute{
			Description: "Table name to use for copying schema migration data.",
			Computed:    true, Optional: true,
		},
		"multiple_admins_required_for_deletion": schema.BoolAttribute{
			Description: "If the database requires multiple admins for deletion.",
			Computed:    true, Optional: true,
		},
		"plan": schema.StringAttribute{
			Description: "The database plan.",
			Computed:    true, Optional: true,
		},
		"production_branch_web_console": schema.BoolAttribute{
			Description: "Whether web console is enabled for production branches.",
			Computed:    true, Optional: true,
		},
		"production_branches_count": schema.Float64Attribute{
			Description: "The total number of database production branches.",
			Computed:    true,
		},
		"ready": schema.BoolAttribute{
			Description: "If the database is ready to be used.",
			Computed:    true,
		},
		"region": schema.SingleNestedAttribute{
			Description: "The region the database lives in.",
			Computed:    true, Optional: true,
			Attributes: regionDataSourceSchemaAttribute,
		},
		"require_approval_for_deploy": schema.BoolAttribute{
			Description: "Whether an approval is required to deploy schema changes to this database.",
			Computed:    true, Optional: true,
		},
		"restrict_branch_region": schema.BoolAttribute{
			Description: "Whether to restrict branch creation to one region.",
			Computed:    true, Optional: true,
		},
		"schema_last_updated_at": schema.StringAttribute{
			Description: "When the default branch schema was last changed.",
			Computed:    true,
		},
		"sharded": schema.BoolAttribute{
			Description: "If the database is sharded.",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "State of the database.",
			Computed:    true,
		},
		"updated_at": schema.StringAttribute{
			Description: "When the database was last updated.",
			Computed:    true,
		},
		"url": schema.StringAttribute{
			Description: "The URL to the database API endpoint.",
			Computed:    true,
		},
	}
}

type databaseDataSourceModel struct {
	Organization                      string                     `tfsdk:"organization"`
	AllowDataBranching                types.Bool                 `tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      types.Bool                 `tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          types.Bool                 `tfsdk:"at_development_branch_limit"` // XXX: deprecated. no longer exists in api
	AutomaticMigrations               types.Bool                 `tfsdk:"automatic_migrations"`
	BranchesCount                     types.Float64              `tfsdk:"branches_count"`
	BranchesUrl                       types.String               `tfsdk:"branches_url"`
	CreatedAt                         types.String               `tfsdk:"created_at"`
	DataImport                        *dataImportDataSourceModel `tfsdk:"data_import"`
	DefaultBranch                     types.String               `tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount types.Float64              `tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           types.Float64              `tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           types.Float64              `tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          types.Float64              `tfsdk:"development_branches_count"`
	HtmlUrl                           types.String               `tfsdk:"html_url"`
	Id                                types.String               `tfsdk:"id"`
	InsightsRawQueries                types.Bool                 `tfsdk:"insights_raw_queries"`
	IssuesCount                       types.Float64              `tfsdk:"issues_count"`
	MigrationFramework                types.String               `tfsdk:"migration_framework"`
	MigrationTableName                types.String               `tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion types.Bool                 `tfsdk:"multiple_admins_required_for_deletion"`
	Name                              types.String               `tfsdk:"name"`
	Plan                              types.String               `tfsdk:"plan"`
	ProductionBranchWebConsole        types.Bool                 `tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           types.Float64              `tfsdk:"production_branches_count"`
	Ready                             types.Bool                 `tfsdk:"ready"`
	Region                            *regionDataSourceModel     `tfsdk:"region"`
	RequireApprovalForDeploy          types.Bool                 `tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              types.Bool                 `tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               types.String               `tfsdk:"schema_last_updated_at"`
	Sharded                           types.Bool                 `tfsdk:"sharded"`
	State                             types.String               `tfsdk:"state"`
	UpdatedAt                         types.String               `tfsdk:"updated_at"`
	Url                               types.String               `tfsdk:"url"`
}

func databaseFromClient(database *planetscale.Database, orgName string, diags diag.Diagnostics) *databaseDataSourceModel {
	if database == nil {
		return nil
	}
	return &databaseDataSourceModel{
		Organization:                      orgName,
		DataImport:                        dataImportFromClient(database.DataImport),
		Region:                            regionFromClient(&database.Region, diags),
		AllowDataBranching:                types.BoolValue(database.AllowDataBranching),
		AtBackupRestoreBranchesLimit:      types.BoolValue(database.AtBackupRestoreBranchesLimit),
		AutomaticMigrations:               types.BoolPointerValue(database.AutomaticMigrations),
		BranchesCount:                     types.Float64Value(database.BranchesCount),
		BranchesUrl:                       types.StringValue(database.BranchesUrl),
		CreatedAt:                         types.StringValue(database.CreatedAt),
		DefaultBranch:                     types.StringValue(database.DefaultBranch),
		DefaultBranchReadOnlyRegionsCount: types.Float64Value(database.DefaultBranchReadOnlyRegionsCount),
		DefaultBranchShardCount:           types.Float64Value(database.DefaultBranchShardCount),
		DefaultBranchTableCount:           types.Float64Value(database.DefaultBranchTableCount),
		DevelopmentBranchesCount:          types.Float64Value(database.DevelopmentBranchesCount),
		HtmlUrl:                           types.StringValue(database.HtmlUrl),
		Id:                                types.StringValue(database.Id),
		InsightsRawQueries:                types.BoolValue(database.InsightsRawQueries),
		IssuesCount:                       types.Float64Value(database.IssuesCount),
		MigrationFramework:                types.StringPointerValue(database.MigrationFramework),
		MigrationTableName:                types.StringPointerValue(database.MigrationTableName),
		MultipleAdminsRequiredForDeletion: types.BoolValue(database.MultipleAdminsRequiredForDeletion),
		Name:                              types.StringValue(database.Name),
		Plan:                              types.StringValue(database.Plan),
		ProductionBranchWebConsole:        types.BoolValue(database.ProductionBranchWebConsole),
		ProductionBranchesCount:           types.Float64Value(database.ProductionBranchesCount),
		Ready:                             types.BoolValue(database.Ready),
		RequireApprovalForDeploy:          types.BoolValue(database.RequireApprovalForDeploy),
		RestrictBranchRegion:              types.BoolValue(database.RestrictBranchRegion),
		SchemaLastUpdatedAt:               types.StringPointerValue(database.SchemaLastUpdatedAt),
		Sharded:                           types.BoolValue(database.Sharded),
		State:                             types.StringValue(database.State),
		UpdatedAt:                         types.StringValue(database.UpdatedAt),
		Url:                               types.StringValue(database.Url),
	}
}

func branchDataSourceSchemaAttribute(computedName bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"organization": schema.StringAttribute{
			Description: "The organization this branch belongs to.",
			Required:    !computedName, Computed: computedName,
		},
		"database": schema.StringAttribute{
			Description: "The database this branch belongs to.",
			Required:    !computedName, Computed: computedName,
		},
		"name": schema.StringAttribute{
			Description: "The name of the branch.",
			Required:    !computedName, Computed: computedName,
		},

		"access_host_url": schema.StringAttribute{
			Description: "The access host URL for the branch. This is a legacy field, use `mysql_edge_address`.",
			Computed:    true,
		},
		"cluster_rate_name": schema.StringAttribute{
			Description: "The SKU representing the branch's cluster size.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "When the branch was created.",
			Computed:    true,
		},
		"html_url": schema.StringAttribute{
			Description: "Planetscale app URL for the branch.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The ID of the branch.",
			Computed:    true,
		},
		"initial_restore_id": schema.StringAttribute{
			Description: "The ID of the backup from which the branch was restored.",
			Computed:    true,
		},
		"mysql_address": schema.StringAttribute{
			Description: "The MySQL address for the branch.",
			Computed:    true,
		},
		"mysql_edge_address": schema.StringAttribute{
			Description: "The address of the MySQL provider for the branch.",
			Computed:    true,
		},
		"parent_branch": schema.StringAttribute{
			Description: "The name of the parent branch from which the branch was created.",
			Computed:    true,
		},
		"production": schema.BoolAttribute{
			Description: "Whether or not the branch is a production branch.",
			Computed:    true,
		},
		"ready": schema.BoolAttribute{
			Description: "Whether or not the branch is ready to serve queries.",
			Computed:    true,
		},
		"restore_checklist_completed_at": schema.StringAttribute{
			Description: "When a user last marked a backup restore checklist as completed.",
			Computed:    true,
		},
		"safe_migrations": schema.BoolAttribute{
			Description: "Whether safe migrations are enabled for this branch.",
			Computed:    true,
		},
		"schema_last_updated_at": schema.StringAttribute{
			Description: "When the schema for the branch was last updated.",
			Computed:    true,
		},
		"shard_count": schema.Float64Attribute{
			Description: "The number of shards in the branch.",
			Computed:    true,
		},
		"sharded": schema.BoolAttribute{
			Description: "Whether or not the branch is sharded.",
			Computed:    true,
		},
		"updated_at": schema.StringAttribute{
			Description: "When the branch was last updated.",
			Computed:    true,
		},

		"actor": schema.SingleNestedAttribute{
			Description: "The actor who created this branch.",
			Computed:    true,
			Attributes:  actorDataSourceSchemaAttribute,
		},
		"region": schema.SingleNestedAttribute{
			Description: "The region in which this branch lives.",
			Computed:    true,
			Attributes:  regionDataSourceSchemaAttribute,
		},
		"restored_from_branch": schema.SingleNestedAttribute{
			Description: "",
			Computed:    true,
			Attributes:  restoredFromBranchDataSourceSchemaAttribute,
		},
	}
}

type branchDataSourceModel struct {
	Organization                types.String                       `tfsdk:"organization"`
	Database                    types.String                       `tfsdk:"database"`
	Name                        types.String                       `tfsdk:"name"`
	AccessHostUrl               types.String                       `tfsdk:"access_host_url"`
	Actor                       *actorDataSourceModel              `tfsdk:"actor"`
	ClusterRateName             types.String                       `tfsdk:"cluster_rate_name"`
	CreatedAt                   types.String                       `tfsdk:"created_at"`
	HtmlUrl                     types.String                       `tfsdk:"html_url"`
	Id                          types.String                       `tfsdk:"id"`
	InitialRestoreId            types.String                       `tfsdk:"initial_restore_id"`
	MysqlAddress                types.String                       `tfsdk:"mysql_address"`
	MysqlEdgeAddress            types.String                       `tfsdk:"mysql_edge_address"`
	ParentBranch                types.String                       `tfsdk:"parent_branch"`
	Production                  types.Bool                         `tfsdk:"production"`
	Ready                       types.Bool                         `tfsdk:"ready"`
	Region                      *regionDataSourceModel             `tfsdk:"region"`
	RestoreChecklistCompletedAt types.String                       `tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *restoredFromBranchDataSourceModel `tfsdk:"restored_from_branch"`
	SafeMigrations              types.Bool                         `tfsdk:"safe_migrations"`
	SchemaLastUpdatedAt         types.String                       `tfsdk:"schema_last_updated_at"`
	ShardCount                  types.Float64                      `tfsdk:"shard_count"`
	Sharded                     types.Bool                         `tfsdk:"sharded"`
	UpdatedAt                   types.String                       `tfsdk:"updated_at"`
}

func branchFromClient(branch *planetscale.Branch, organization, database string, diags diag.Diagnostics) *branchDataSourceModel {
	if branch == nil {
		return nil
	}
	return &branchDataSourceModel{
		Organization:                types.StringValue(organization),
		Database:                    types.StringValue(database),
		Actor:                       actorFromClient(branch.Actor),
		Region:                      regionFromClient(branch.Region, diags),
		RestoredFromBranch:          restoredFromBranchFromClient(branch.RestoredFromBranch),
		Name:                        types.StringValue(branch.Name),
		CreatedAt:                   types.StringValue(branch.CreatedAt),
		HtmlUrl:                     types.StringValue(branch.HtmlUrl),
		Id:                          types.StringValue(branch.Id),
		MysqlAddress:                types.StringValue(branch.MysqlAddress),
		MysqlEdgeAddress:            types.StringValue(branch.MysqlEdgeAddress),
		ParentBranch:                types.StringPointerValue(branch.ParentBranch),
		Production:                  types.BoolValue(branch.Production),
		Ready:                       types.BoolValue(branch.Ready),
		RestoreChecklistCompletedAt: types.StringPointerValue(branch.RestoreChecklistCompletedAt),
		SafeMigrations:              types.BoolValue(branch.SafeMigrations),
		SchemaLastUpdatedAt:         types.StringValue(branch.SchemaLastUpdatedAt),
		ShardCount:                  types.Float64PointerValue(branch.ShardCount),
		Sharded:                     types.BoolValue(branch.Sharded),
		UpdatedAt:                   types.StringValue(branch.UpdatedAt),
	}
}

var branchSafeMigrationsDataSourceSchemaAttribute = map[string]schema.Attribute{
	"organization": schema.StringAttribute{
		Description: "The organization this branch belongs to.",
		Required:    true,
	},
	"database": schema.StringAttribute{
		Description: "The database this branch belongs to.",
		Required:    true,
	},
	"branch": schema.StringAttribute{
		Description: "The name of the branch this safe migrations configuration belongs to.",
		Required:    true,
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether safe migrations are enabled for this branch.",
		Computed:    true,
	},
}

type branchSafeMigrationsDataSourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Database     types.String `tfsdk:"database"`
	Branch       types.String `tfsdk:"branch"`
	Enabled      types.Bool   `tfsdk:"enabled"`
}

func branchSafeMigrationsFromClient(branch *planetscale.Branch, organization, database string) *branchSafeMigrationsDataSourceModel {
	if branch == nil {
		return nil
	}
	return &branchSafeMigrationsDataSourceModel{
		Organization: types.StringValue(organization),
		Database:     types.StringValue(database),
		Branch:       types.StringValue(branch.Name),
		Enabled:      types.BoolValue(branch.SafeMigrations),
	}
}

var actorDataSourceSchemaAttribute = map[string]schema.Attribute{
	"avatar_url": schema.StringAttribute{
		Computed: true, Description: "The URL of the actor's avatar",
	},
	"display_name": schema.StringAttribute{
		Computed: true, Description: "The name of the actor",
	},
	"id": schema.StringAttribute{
		Computed: true, Description: "The ID of the actor",
	},
}

type actorDataSourceModel struct {
	AvatarUrl   types.String `tfsdk:"avatar_url"`
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

func actorFromClient(actor *planetscale.Actor) *actorDataSourceModel {
	if actor == nil {
		return nil
	}
	return &actorDataSourceModel{
		AvatarUrl:   types.StringValue(actor.AvatarUrl),
		DisplayName: types.StringValue(actor.DisplayName),
		Id:          types.StringValue(actor.Id),
	}
}

var regionDataSourceSchemaAttribute = map[string]schema.Attribute{
	"display_name": schema.StringAttribute{
		Description: "Name of the region.",
		Computed:    true,
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether or not the region is currently active.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID of the region.",
		Computed:    true,
	},
	"location": schema.StringAttribute{
		Description: "Location of the region.",
		Computed:    true,
	},
	"provider": schema.StringAttribute{
		Description: "Provider for the region (ex. AWS).",
		Computed:    true,
	},
	"public_ip_addresses": schema.ListAttribute{
		Description: "Public IP addresses for the region.",
		Computed:    true, ElementType: types.StringType,
	},
	"slug": schema.StringAttribute{
		Description: "The slug of the region.",
		Computed:    true,
	},
}

type regionDataSourceModel struct {
	DisplayName       types.String `tfsdk:"display_name"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Id                types.String `tfsdk:"id"`
	Location          types.String `tfsdk:"location"`
	Provider          types.String `tfsdk:"provider"`
	PublicIpAddresses types.List   `tfsdk:"public_ip_addresses"`
	Slug              types.String `tfsdk:"slug"`
}

func regionFromClient(region *planetscale.Region, diags diag.Diagnostics) *regionDataSourceModel {
	if region == nil {
		return nil
	}
	return &regionDataSourceModel{
		DisplayName:       types.StringValue(region.DisplayName),
		Enabled:           types.BoolValue(region.Enabled),
		Id:                types.StringValue(region.Id),
		Location:          types.StringValue(region.Location),
		Provider:          types.StringValue(region.Provider),
		PublicIpAddresses: stringsToListValue(region.PublicIpAddresses, diags),
		Slug:              types.StringValue(region.Slug),
	}
}

var readOnlyRegionsDataSourceSchemaAttribute = map[string]schema.Attribute{
	"organization": schema.StringAttribute{
		Description: "The organization for which the read-only regions are available.",
		Required:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the database for which the read-only regions are available.",
		Required:    true,
	},
	"regions": schema.ListNestedAttribute{
		Description: "The list of read-only regions available for the database.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"actor": schema.SingleNestedAttribute{
					Description: "The actor that created the read-only region.",
					Computed:    true,
					Attributes:  actorDataSourceSchemaAttribute,
				},
				"created_at": schema.StringAttribute{
					Description: "When the read-only region was created.",
					Computed:    true,
				},
				"display_name": schema.StringAttribute{
					Description: "The name of the read-only region.",
					Computed:    true,
				},
				"id": schema.StringAttribute{
					Description: "The ID of the read-only region.",
					Computed:    true,
				},
				"ready": schema.BoolAttribute{
					Description: "Whether or not the read-only region is ready to serve queries.",
					Computed:    true,
				},
				"ready_at": schema.StringAttribute{
					Description: "When the read-only region was ready to serve queries.",
					Computed:    true,
				},
				"updated_at": schema.StringAttribute{
					Description: "When the read-only region was last updated.",
					Computed:    true,
				},
				"region": schema.SingleNestedAttribute{
					Description: "The details of the read-only region.",
					Computed:    true,
					Attributes:  regionDataSourceSchemaAttribute,
				},
			},
		},
	},
}

type readOnlyRegionDataSourceModel struct {
	Actor       actorDataSourceModel  `tfsdk:"actor"`
	CreatedAt   types.String          `tfsdk:"created_at"`
	DisplayName types.String          `tfsdk:"display_name"`
	Id          types.String          `tfsdk:"id"`
	Ready       types.Bool            `tfsdk:"ready"`
	ReadyAt     types.String          `tfsdk:"ready_at"`
	Region      regionDataSourceModel `tfsdk:"region"`
	UpdatedAt   types.String          `tfsdk:"updated_at"`
}

func readOnlyRegionFromClient(readOnlyRegion *planetscale.ReadOnlyRegion, diags diag.Diagnostics) *readOnlyRegionDataSourceModel {
	if readOnlyRegion == nil {
		return nil
	}

	return &readOnlyRegionDataSourceModel{
		Actor:       *actorFromClient(&readOnlyRegion.Actor),
		Region:      *regionFromClient(&readOnlyRegion.Region, diags),
		CreatedAt:   types.StringValue(readOnlyRegion.CreatedAt),
		DisplayName: types.StringValue(readOnlyRegion.DisplayName),
		Id:          types.StringValue(readOnlyRegion.Id),
		Ready:       types.BoolValue(readOnlyRegion.Ready),
		ReadyAt:     types.StringValue(readOnlyRegion.ReadyAt),
		UpdatedAt:   types.StringValue(readOnlyRegion.UpdatedAt),
	}
}

var restoredFromBranchDataSourceSchemaAttribute = map[string]schema.Attribute{
	"created_at": schema.StringAttribute{
		Description: "When the resource was created.",
		Computed:    true,
	},
	"deleted_at": schema.StringAttribute{
		Description: "When the resource was deleted, if deleted.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID for the resource.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name for the resource.",
		Computed:    true,
	},
	"updated_at": schema.StringAttribute{
		Description: "When the resource was last updated.",
		Computed:    true,
	},
}

type restoredFromBranchDataSourceModel struct {
	CreatedAt types.String `tfsdk:"created_at"`
	DeletedAt types.String `tfsdk:"deleted_at"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func restoredFromBranchFromClient(rfb *planetscale.RestoredFromBranch) *restoredFromBranchDataSourceModel {
	if rfb == nil {
		return nil
	}
	return &restoredFromBranchDataSourceModel{
		CreatedAt: types.StringValue(rfb.CreatedAt),
		DeletedAt: types.StringValue(rfb.DeletedAt),
		Id:        types.StringValue(rfb.Id),
		Name:      types.StringValue(rfb.Name),
		UpdatedAt: types.StringValue(rfb.UpdatedAt),
	}
}

var tableSchemaDataSourceSchemaAttribute = map[string]schema.Attribute{
	"html": schema.StringAttribute{
		Description: "Syntax highlighted HTML for the table's schema.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "Name of the table.",
		Computed:    true,
	},
	"raw": schema.StringAttribute{
		Description: "The table's schema.",
		Computed:    true,
	},
}

type tableSchemaDataSourceModel struct {
	Html types.String `tfsdk:"html"`
	Name types.String `tfsdk:"name"`
	Raw  types.String `tfsdk:"raw"`
}

func tableSchemaFromClient(ts *planetscale.TableSchema) *tableSchemaDataSourceModel {
	if ts == nil {
		return nil
	}
	return &tableSchemaDataSourceModel{
		Html: types.StringValue(ts.Html),
		Name: types.StringValue(ts.Name),
		Raw:  types.StringValue(ts.Raw),
	}
}

var lintErrorDataSourceSchemaAttribute = map[string]schema.Attribute{
	"auto_increment_column_names": schema.ListAttribute{
		Description: "A list of invalid auto-incremented columns.",
		Computed:    true, ElementType: types.StringType,
	},
	"charset_name": schema.StringAttribute{
		Description: "The charset of the schema.",
		Computed:    true,
	},
	"check_constraint_name": schema.StringAttribute{
		Description: "The name of the invalid check constraint.",
		Computed:    true,
	},
	"column_name": schema.StringAttribute{
		Description: "The column in a table relevant to the error.",
		Computed:    true,
	},
	"docs_url": schema.StringAttribute{
		Description: "A link to the documentation related to the error.",
		Computed:    true,
	},
	"engine_name": schema.StringAttribute{
		Description: "The engine of the schema.",
		Computed:    true,
	},
	"enum_value": schema.StringAttribute{
		Description: "The name of the invalid enum value.",
		Computed:    true,
	},
	"error_description": schema.StringAttribute{
		Description: "A description for the error that occurred.",
		Computed:    true,
	},
	"foreign_key_column_names": schema.ListAttribute{
		Description: "A list of invalid foreign key columns in a table.",
		Computed:    true, ElementType: types.StringType,
	},
	"json_path": schema.StringAttribute{
		Description: "The path for an invalid JSON column.",
		Computed:    true,
	},
	"keyspace_name": schema.StringAttribute{
		Description: "The keyspace of the schema with the error.",
		Computed:    true,
	},
	"lint_error": schema.StringAttribute{
		Description: "Code representing.",
		Computed:    true,
	},
	"partition_name": schema.StringAttribute{
		Description: "The name of the invalid partition in the schema.",
		Computed:    true,
	},
	"partitioning_type": schema.StringAttribute{
		Description: "The name of the invalid partitioning type.",
		Computed:    true,
	},
	"subject_type": schema.StringAttribute{
		Description: "The subject for the errors.",
		Computed:    true,
	},
	"table_name": schema.StringAttribute{
		Description: "The table with the error.",
		Computed:    true,
	},
	"vindex_name": schema.StringAttribute{
		Description: "The name of the vindex for the schema.",
		Computed:    true,
	},
}

type lintErrorDataSourceModel struct {
	AutoIncrementColumnNames types.List   `tfsdk:"auto_increment_column_names"`
	CharsetName              types.String `tfsdk:"charset_name"`
	CheckConstraintName      types.String `tfsdk:"check_constraint_name"`
	ColumnName               types.String `tfsdk:"column_name"`
	DocsUrl                  types.String `tfsdk:"docs_url"`
	EngineName               types.String `tfsdk:"engine_name"`
	EnumValue                types.String `tfsdk:"enum_value"`
	ErrorDescription         types.String `tfsdk:"error_description"`
	ForeignKeyColumnNames    types.List   `tfsdk:"foreign_key_column_names"`
	JsonPath                 types.String `tfsdk:"json_path"`
	KeyspaceName             types.String `tfsdk:"keyspace_name"`
	LintError                types.String `tfsdk:"lint_error"`
	PartitionName            types.String `tfsdk:"partition_name"`
	PartitioningType         types.String `tfsdk:"partitioning_type"`
	SubjectType              types.String `tfsdk:"subject_type"`
	TableName                types.String `tfsdk:"table_name"`
	VindexName               types.String `tfsdk:"vindex_name"`
}

func lintErrorFromClient(le *planetscale.LintError, diags diag.Diagnostics) *lintErrorDataSourceModel {
	if le == nil {
		return nil
	}
	return &lintErrorDataSourceModel{
		AutoIncrementColumnNames: stringsToListValue(le.AutoIncrementColumnNames, diags),
		CharsetName:              types.StringValue(le.CharsetName),
		CheckConstraintName:      types.StringValue(le.CheckConstraintName),
		ColumnName:               types.StringValue(le.ColumnName),
		DocsUrl:                  types.StringValue(le.DocsUrl),
		EngineName:               types.StringValue(le.EngineName),
		EnumValue:                types.StringValue(le.EnumValue),
		ErrorDescription:         types.StringValue(le.ErrorDescription),
		ForeignKeyColumnNames:    stringsToListValue(le.ForeignKeyColumnNames, diags),
		JsonPath:                 types.StringValue(le.JsonPath),
		KeyspaceName:             types.StringValue(le.KeyspaceName),
		LintError:                types.StringValue(le.LintError),
		PartitionName:            types.StringValue(le.PartitionName),
		PartitioningType:         types.StringValue(le.PartitioningType),
		SubjectType:              types.StringValue(le.SubjectType),
		TableName:                types.StringValue(le.TableName),
		VindexName:               types.StringValue(le.VindexName),
	}
}

var oauthApplicationAttribute = map[string]schema.Attribute{
	"avatar": schema.StringAttribute{
		Description: "The image source for the OAuth application's avatar.",
		Computed:    true,
	},
	"client_id": schema.StringAttribute{
		Description: "The OAuth application's unique client id.",
		Computed:    true,
	},
	"created_at": schema.StringAttribute{
		Description: "When the OAuth application was created.",
		Computed:    true,
	},
	"domain": schema.StringAttribute{
		Description: "The domain of the OAuth application. Used for verification of a valid redirect uri.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID of the OAuth application.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the OAuth application.",
		Computed:    true,
	},
	"redirect_uri": schema.StringAttribute{
		Description: "The redirect URI of the OAuth application.",
		Computed:    true,
	},
	"scopes": schema.ListAttribute{
		Description: "The scopes that the OAuth application requires on a user's accout.",
		Computed:    true, ElementType: types.StringType,
	},
	"tokens": schema.Float64Attribute{
		Description: "The number of tokens issued by the OAuth application.",
		Computed:    true,
	},
	"updated_at": schema.StringAttribute{
		Description: "When the OAuth application was last updated.",
		Computed:    true,
	},
}

type oauthApplicationDataSourceModel struct {
	Avatar      types.String  `tfsdk:"avatar"`
	ClientId    types.String  `tfsdk:"client_id"`
	CreatedAt   types.String  `tfsdk:"created_at"`
	Domain      types.String  `tfsdk:"domain"`
	Id          types.String  `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	RedirectUri types.String  `tfsdk:"redirect_uri"`
	Scopes      types.List    `tfsdk:"scopes"`
	Tokens      types.Float64 `tfsdk:"tokens"`
	UpdatedAt   types.String  `tfsdk:"updated_at"`
}

func oauthApplicationFromClient(oa *planetscale.OauthApplication, diags diag.Diagnostics) *oauthApplicationDataSourceModel {
	if oa == nil {
		return nil
	}
	return &oauthApplicationDataSourceModel{
		Avatar:      types.StringPointerValue(oa.Avatar),
		ClientId:    types.StringValue(oa.ClientId),
		CreatedAt:   types.StringValue(oa.CreatedAt),
		Domain:      types.StringValue(oa.Domain),
		Id:          types.StringValue(oa.Id),
		Name:        types.StringValue(oa.Name),
		RedirectUri: types.StringValue(oa.RedirectUri),
		Scopes:      stringsToListValue(oa.Scopes, diags),
		Tokens:      types.Float64Value(oa.Tokens),
		UpdatedAt:   types.StringValue(oa.UpdatedAt),
	}
}

var backupPolicyDataSourceSchemaAttribute = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The ID of the backup policy.",
		Computed:    true,
	},
	"created_at": schema.StringAttribute{
		Description: "When the backup policy was created.",
		Computed:    true,
	},
	"frequency_unit": schema.StringAttribute{
		Description: "The unit for the frequency of the backup policy.",
		Computed:    true,
	},
	"frequency_value": schema.Float64Attribute{
		Description: "A number value for the frequency of the backup policy.",
		Computed:    true,
	},
	"last_ran_at": schema.StringAttribute{
		Description: "When the backup was last run.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the backup policy.",
		Computed:    true,
	},
	"next_run_at": schema.StringAttribute{
		Description: "When the backup will next run.",
		Computed:    true,
	},
	"retention_unit": schema.StringAttribute{
		Description: "The unit for the retention period of the backup policy.",
		Computed:    true,
	},
	"retention_value": schema.Float64Attribute{
		Description: "A number value for the retention period of the backup policy.",
		Computed:    true,
	},
	"schedule_day": schema.StringAttribute{
		Description: "Day of the week that the backup is scheduled.",
		Computed:    true,
	},
	"schedule_week": schema.StringAttribute{
		Description: "Week of the month that the backup is scheduled.",
		Computed:    true,
	},
	"target": schema.StringAttribute{
		Description: "Whether the backup policy is for a production or development database, or for a database branch.",
		Computed:    true,
	},
	"updated_at": schema.StringAttribute{
		Description: "When the backup policy was last updated.",
		Computed:    true,
	},
}

type backupPolicyDataSourceModel struct {
	CreatedAt      types.String  `tfsdk:"created_at"`
	FrequencyUnit  types.String  `tfsdk:"frequency_unit"`
	FrequencyValue types.Float64 `tfsdk:"frequency_value"`
	Id             types.String  `tfsdk:"id"`
	LastRanAt      types.String  `tfsdk:"last_ran_at"`
	Name           types.String  `tfsdk:"name"`
	NextRunAt      types.String  `tfsdk:"next_run_at"`
	RetentionUnit  types.String  `tfsdk:"retention_unit"`
	RetentionValue types.Float64 `tfsdk:"retention_value"`
	ScheduleDay    types.String  `tfsdk:"schedule_day"`
	ScheduleWeek   types.String  `tfsdk:"schedule_week"`
	Target         types.String  `tfsdk:"target"`
	UpdatedAt      types.String  `tfsdk:"updated_at"`
}

func backupPolicyFromClient(backupPolicy *planetscale.BackupPolicy) *backupPolicyDataSourceModel {
	if backupPolicy == nil {
		return nil
	}
	return &backupPolicyDataSourceModel{
		CreatedAt:      types.StringValue(backupPolicy.CreatedAt),
		FrequencyUnit:  types.StringValue(backupPolicy.FrequencyUnit),
		FrequencyValue: types.Float64Value(backupPolicy.FrequencyValue),
		Id:             types.StringValue(backupPolicy.Id),
		LastRanAt:      types.StringValue(backupPolicy.LastRanAt),
		Name:           types.StringValue(backupPolicy.Name),
		NextRunAt:      types.StringValue(backupPolicy.NextRunAt),
		RetentionUnit:  types.StringValue(backupPolicy.RetentionUnit),
		RetentionValue: types.Float64Value(backupPolicy.RetentionValue),
		ScheduleDay:    types.StringValue(backupPolicy.ScheduleDay),
		ScheduleWeek:   types.StringValue(backupPolicy.ScheduleWeek),
		Target:         types.StringValue(backupPolicy.Target),
		UpdatedAt:      types.StringValue(backupPolicy.UpdatedAt),
	}
}

var branchForPasswordDataSourceSchemaAttribute = map[string]schema.Attribute{
	"access_host_url":    schema.StringAttribute{Computed: true},
	"id":                 schema.StringAttribute{Computed: true},
	"mysql_edge_address": schema.StringAttribute{Computed: true},
	"name":               schema.StringAttribute{Computed: true},
	"production":         schema.BoolAttribute{Computed: true},
}

type branchForPasswordDataSourceModel struct {
	AccessHostUrl    types.String `tfsdk:"access_host_url"`
	Id               types.String `tfsdk:"id"`
	MysqlEdgeAddress types.String `tfsdk:"mysql_edge_address"`
	Name             types.String `tfsdk:"name"`
	Production       types.Bool   `tfsdk:"production"`
}

func branchForPasswordFromClient(branchForPassword *planetscale.BranchForPassword) *branchForPasswordDataSourceModel {
	if branchForPassword == nil {
		return nil
	}
	return &branchForPasswordDataSourceModel{
		AccessHostUrl:    types.StringValue(branchForPassword.AccessHostUrl),
		Id:               types.StringValue(branchForPassword.Id),
		MysqlEdgeAddress: types.StringValue(branchForPassword.MysqlEdgeAddress),
		Name:             types.StringValue(branchForPassword.Name),
		Production:       types.BoolValue(branchForPassword.Production),
	}
}

func passwordDataSourceSchemaAttribute(computedName bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"organization": schema.StringAttribute{
			Description: "The organization this database branch password belongs to.",
			Required:    !computedName, Computed: computedName,
		},
		"database": schema.StringAttribute{
			Description: "The database this branch password belongs to.",
			Required:    !computedName, Computed: computedName,
		},
		"branch": schema.StringAttribute{
			Description: "The branch this password belongs to..",
			Required:    !computedName, Computed: computedName,
		},
		"id": schema.StringAttribute{
			Description: "The ID for the password.",
			Required:    !computedName, Computed: computedName,
		},
		"read_only_region_id": schema.StringAttribute{
			Description: "If the password is for a read-only region, the ID of the region.",
			Optional:    !computedName, Computed: computedName,
		},
		"access_host_url": schema.StringAttribute{
			Description: "The host URL for the password.",
			Computed:    true,
		},
		"actor": schema.SingleNestedAttribute{
			Description: "The actor that created this branch.",
			Computed:    true,
			Attributes:  actorDataSourceSchemaAttribute,
		},
		"created_at": schema.StringAttribute{
			Description: "When the password was created.",
			Computed:    true,
		},
		"database_branch": schema.SingleNestedAttribute{
			Description: "The branch this password is allowed to access.",
			Computed:    true,
			Attributes:  branchForPasswordDataSourceSchemaAttribute,
		},
		"deleted_at": schema.StringAttribute{
			Description: "When the password was deleted.",
			Computed:    true,
		},
		"expires_at": schema.StringAttribute{
			Description: "When the password will expire.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The display name for the password.",
			Computed:    true,
		},
		"region": schema.SingleNestedAttribute{
			Description: "The region in which this password can be used.",
			Computed:    true,
			Attributes:  regionDataSourceSchemaAttribute,
		},
		"renewable": schema.BoolAttribute{
			Description: "Whether or not the password can be renewed.",
			Computed:    true,
		},
		"role": schema.StringAttribute{
			Description: "The role for the password.",
			Computed:    true,
		},
		"ttl_seconds": schema.Float64Attribute{
			Description: "Time to live (in seconds) for the password. The password will be invalid and unrenewable when TTL has passed.",
			Computed:    true,
		},
		"username": schema.StringAttribute{
			Description: "The username for the password.",
			Computed:    true,
		},

		// manually removed from spec because currently buggy
		// "integrations": schema.ListAttribute{Computed: true, ElementType: types.StringType},
	}
}

type passwordDataSourceModel struct {
	Organization     types.String                      `tfsdk:"organization"`
	Database         types.String                      `tfsdk:"database"`
	Branch           types.String                      `tfsdk:"branch"`
	ReadOnlyRegionId types.String                      `tfsdk:"read_only_region_id"`
	Id               types.String                      `tfsdk:"id"`
	AccessHostUrl    types.String                      `tfsdk:"access_host_url"`
	Actor            *actorDataSourceModel             `tfsdk:"actor"`
	CreatedAt        types.String                      `tfsdk:"created_at"`
	DatabaseBranch   *branchForPasswordDataSourceModel `tfsdk:"database_branch"`
	DeletedAt        types.String                      `tfsdk:"deleted_at"`
	ExpiresAt        types.String                      `tfsdk:"expires_at"`
	Name             types.String                      `tfsdk:"name"`
	Region           *regionDataSourceModel            `tfsdk:"region"`
	Renewable        types.Bool                        `tfsdk:"renewable"`
	Role             types.String                      `tfsdk:"role"`
	TtlSeconds       types.Float64                     `tfsdk:"ttl_seconds"`
	Username         types.String                      `tfsdk:"username"`

	// manually removed from spec because currently buggy
	// Integrations     types.List                        `tfsdk:"integrations"`
}

func passwordFromClient(password *planetscale.Password, organization, database, branch string, readOnlyRegionID *string, diags diag.Diagnostics) *passwordDataSourceModel {
	if password == nil {
		return nil
	}
	return &passwordDataSourceModel{
		Organization:     types.StringValue(organization),
		Database:         types.StringValue(database),
		Branch:           types.StringValue(branch),
		ReadOnlyRegionId: types.StringPointerValue(readOnlyRegionID),
		AccessHostUrl:    types.StringValue(password.AccessHostUrl),
		Actor:            actorFromClient(password.Actor),
		CreatedAt:        types.StringValue(password.CreatedAt),
		DatabaseBranch:   branchForPasswordFromClient(&password.DatabaseBranch),
		DeletedAt:        types.StringPointerValue(password.DeletedAt),
		ExpiresAt:        types.StringPointerValue(password.ExpiresAt),
		Id:               types.StringValue(password.Id),
		Name:             types.StringValue(password.Name),
		Region:           regionFromClient(password.Region, diags),
		Renewable:        types.BoolValue(password.Renewable),
		Role:             types.StringValue(password.Role),
		TtlSeconds:       types.Float64Value(password.TtlSeconds),
		Username:         types.StringPointerValue(password.Username),
		// manually removed from spec because currently buggy
		// Integrations:     stringsToListValue(password.Integrations, diags),
	}
}

func backupDataSourceSchemaAttribute(computedID bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"organization": schema.StringAttribute{
			Description: "The organization this backup belongs to.",
			Required:    !computedID, Computed: computedID,
		},
		"database": schema.StringAttribute{
			Description: "The database this backup belongs to.",
			Required:    !computedID, Computed: computedID,
		},
		"branch": schema.StringAttribute{
			Description: "The branch this backup belongs to.",
			Required:    !computedID, Computed: computedID,
		},
		"id": schema.StringAttribute{
			Description: "The ID of the backup.",
			Required:    !computedID, Computed: computedID,
		},
		"name": schema.StringAttribute{
			Description: "The name of the backup.",
			Computed:    true,
		},
		"actor": schema.SingleNestedAttribute{
			Description: "The actor that created the backup.",
			Computed:    true, Attributes: actorDataSourceSchemaAttribute,
		},
		"backup_policy": schema.SingleNestedAttribute{
			Description: "The backup policy being followed.",
			Computed:    true,
			Attributes:  backupPolicyDataSourceSchemaAttribute,
		},
		"created_at": schema.StringAttribute{
			Description: "When the backup was created.",
			Computed:    true,
		},
		"estimated_storage_cost": schema.Float64Attribute{
			Description: "The estimated storage cost of the backup.",
			Computed:    true,
		},
		"required": schema.BoolAttribute{
			Description: "Whether or not the backup policy is required.",
			Computed:    true,
		},
		"restored_branches": schema.ListAttribute{
			Description: "Branches that have been restored with this backup.",
			Computed:    true, ElementType: types.StringType,
		},
		"size": schema.Float64Attribute{
			Description: "The size of the backup.",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "The current state of the backup.",
			Computed:    true,
		},
		"updated_at": schema.StringAttribute{
			Description: "When the backup was last updated.",
			Computed:    true,
		},
	}
}

type backupDataSourceModel struct {
	Organization         types.String                 `tfsdk:"organization"`
	Database             types.String                 `tfsdk:"database"`
	Branch               types.String                 `tfsdk:"branch"`
	Name                 types.String                 `tfsdk:"name"`
	Id                   types.String                 `tfsdk:"id"`
	Actor                *actorDataSourceModel        `tfsdk:"actor"`
	BackupPolicy         *backupPolicyDataSourceModel `tfsdk:"backup_policy"`
	CreatedAt            types.String                 `tfsdk:"created_at"`
	EstimatedStorageCost types.Float64                `tfsdk:"estimated_storage_cost"`
	Required             types.Bool                   `tfsdk:"required"`
	RestoredBranches     types.List                   `tfsdk:"restored_branches"`
	Size                 types.Float64                `tfsdk:"size"`
	State                types.String                 `tfsdk:"state"`
	UpdatedAt            types.String                 `tfsdk:"updated_at"`
}

func backupFromClient(backup *planetscale.Backup, organization, database, branch string, diags diag.Diagnostics) *backupDataSourceModel {
	if backup == nil {
		return nil
	}
	restoredBranches := types.ListNull(types.StringType)
	if backup.RestoredBranches != nil {
		restoredBranches = stringsToListValue(*backup.RestoredBranches, diags)
	}
	return &backupDataSourceModel{
		Organization:         types.StringValue(organization),
		Database:             types.StringValue(database),
		Branch:               types.StringValue(branch),
		Name:                 types.StringValue(backup.Name),
		Actor:                actorFromClient(backup.Actor),
		BackupPolicy:         backupPolicyFromClient(backup.BackupPolicy),
		CreatedAt:            types.StringValue(backup.CreatedAt),
		EstimatedStorageCost: types.Float64Value(backup.EstimatedStorageCost),
		Id:                   types.StringValue(backup.Id),
		Required:             types.BoolValue(backup.Required),
		RestoredBranches:     restoredBranches,
		Size:                 types.Float64Value(backup.Size),
		State:                types.StringValue(backup.State),
		UpdatedAt:            types.StringValue(backup.UpdatedAt),
	}
}

func stringsToListValue(in []string, diags diag.Diagnostics) types.List {
	out := make([]attr.Value, 0, len(in))
	for _, el := range in {
		out = append(out, types.StringValue(el))
	}
	lv, diag := types.ListValue(types.StringType, out)
	if diag.HasError() {
		diags.Append(diag...)
	}
	return lv
}
