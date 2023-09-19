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
		"name": schema.StringAttribute{Required: !computedName, Computed: computedName},

		"admin_only_production_access": schema.BoolAttribute{Computed: true},
		"billing_email":                schema.StringAttribute{Computed: true},
		"can_create_databases":         schema.BoolAttribute{Computed: true},
		"created_at":                   schema.StringAttribute{Computed: true},
		"database_count":               schema.Float64Attribute{Computed: true},
		"features": schema.SingleNestedAttribute{
			Computed:   true,
			Attributes: featuresDataSourceSchemaAttribute,
		},
		"flags": schema.SingleNestedAttribute{
			Computed:   true,
			Attributes: flagsDataSourceSchemaAttribute,
		},
		"free_databases_remaining": schema.Float64Attribute{Computed: true},
		"has_past_due_invoices":    schema.BoolAttribute{Computed: true},
		"id":                       schema.StringAttribute{Computed: true},
		"plan":                     schema.StringAttribute{Computed: true},
		"single_tenancy":           schema.BoolAttribute{Computed: true},
		"sleeping_database_count":  schema.Float64Attribute{Computed: true},
		"sso":                      schema.BoolAttribute{Computed: true},
		"sso_directory":            schema.BoolAttribute{Computed: true},
		"sso_portal_url":           schema.StringAttribute{Computed: true},
		"updated_at":               schema.StringAttribute{Computed: true},
		"valid_billing_info":       schema.BoolAttribute{Computed: true},
	}
}

type organizationDataSourceModel struct {
	AdminOnlyProductionAccess types.Bool               `tfsdk:"admin_only_production_access"`
	BillingEmail              types.String             `tfsdk:"billing_email"`
	CanCreateDatabases        types.Bool               `tfsdk:"can_create_databases"`
	CreatedAt                 types.String             `tfsdk:"created_at"`
	DatabaseCount             types.Float64            `tfsdk:"database_count"`
	Features                  *featuresDataSourceModel `tfsdk:"features"`
	Flags                     *flagsDataSourceModel    `tfsdk:"flags"`
	FreeDatabasesRemaining    types.Float64            `tfsdk:"free_databases_remaining"`
	HasPastDueInvoices        types.Bool               `tfsdk:"has_past_due_invoices"`
	Id                        types.String             `tfsdk:"id"`
	Name                      types.String             `tfsdk:"name"`
	Plan                      types.String             `tfsdk:"plan"`
	SingleTenancy             types.Bool               `tfsdk:"single_tenancy"`
	SleepingDatabaseCount     types.Float64            `tfsdk:"sleeping_database_count"`
	Sso                       types.Bool               `tfsdk:"sso"`
	SsoDirectory              types.Bool               `tfsdk:"sso_directory"`
	SsoPortalUrl              types.String             `tfsdk:"sso_portal_url"`
	UpdatedAt                 types.String             `tfsdk:"updated_at"`
	ValidBillingInfo          types.Bool               `tfsdk:"valid_billing_info"`
}

func (mdl *organizationDataSourceModel) fromClient(org *planetscale.Organization) (diags diag.Diagnostics) {
	if org == nil {
		return diags
	}
	if org.Features != nil {
		mdl.Features = &featuresDataSourceModel{}
	}
	if org.Flags != nil {
		mdl.Flags = &flagsDataSourceModel{}
	}
	diags.Append(mdl.Features.fromClient(org.Features)...)
	diags.Append(mdl.Flags.fromClient(org.Flags)...)

	mdl.AdminOnlyProductionAccess = types.BoolValue(org.AdminOnlyProductionAccess)
	mdl.BillingEmail = types.StringPointerValue(org.BillingEmail)
	mdl.CanCreateDatabases = types.BoolValue(org.CanCreateDatabases)
	mdl.CreatedAt = types.StringValue(org.CreatedAt)
	mdl.DatabaseCount = types.Float64Value(org.DatabaseCount)
	mdl.FreeDatabasesRemaining = types.Float64Value(org.FreeDatabasesRemaining)
	mdl.HasPastDueInvoices = types.BoolValue(org.HasPastDueInvoices)
	mdl.Id = types.StringValue(org.Id)
	mdl.Name = types.StringValue(org.Name)
	mdl.Plan = types.StringValue(org.Plan)
	mdl.SingleTenancy = types.BoolValue(org.SingleTenancy)
	mdl.SleepingDatabaseCount = types.Float64Value(org.SleepingDatabaseCount)
	mdl.Sso = types.BoolValue(org.Sso)
	mdl.SsoDirectory = types.BoolValue(org.SsoDirectory)
	mdl.SsoPortalUrl = types.StringPointerValue(org.SsoPortalUrl)
	mdl.UpdatedAt = types.StringValue(org.UpdatedAt)
	mdl.ValidBillingInfo = types.BoolValue(org.ValidBillingInfo)
	return diags
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

func (mdl *featuresDataSourceModel) fromClient(features *planetscale.Features) (diags diag.Diagnostics) {
	if features == nil {
		return diags
	}
	mdl.Insights = types.BoolPointerValue(features.Insights)
	mdl.SingleTenancy = types.BoolPointerValue(features.SingleTenancy)
	mdl.Sso = types.BoolPointerValue(features.Sso)
	return diags
}

var flagsDataSourceSchemaAttribute = map[string]schema.Attribute{
	"example_flag": schema.StringAttribute{Computed: true},
}

type flagsDataSourceModel struct {
	ExampleFlag types.String `tfsdk:"example_flag"`
}

func (mdl *flagsDataSourceModel) fromClient(flags *planetscale.Flags) (diags diag.Diagnostics) {
	if flags == nil {
		return diags
	}
	mdl.ExampleFlag = types.StringPointerValue(flags.ExampleFlag)
	return diags
}

type dataSourceDataSourceModel struct {
	Database types.String `tfsdk:"database"`
	Hostname types.String `tfsdk:"hostname"`
	Port     types.String `tfsdk:"port"`
}

func (mdl *dataSourceDataSourceModel) fromClient(dataSource *planetscale.DataSource) (diags diag.Diagnostics) {
	if dataSource == nil {
		return diags
	}
	mdl.Database = types.StringValue(dataSource.Database)
	mdl.Hostname = types.StringValue(dataSource.Hostname)
	mdl.Port = types.StringValue(dataSource.Port)
	return diags
}

type dataImportDataSourceModel struct {
	DataSource        dataSourceDataSourceModel `tfsdk:"data_source"`
	FinishedAt        types.String              `tfsdk:"finished_at"`
	ImportCheckErrors types.String              `tfsdk:"import_check_errors"`
	StartedAt         types.String              `tfsdk:"started_at"`
	State             types.String              `tfsdk:"state"`
}

func (mdl *dataImportDataSourceModel) fromClient(dataImport *planetscale.DataImport) (diags diag.Diagnostics) {
	if dataImport == nil {
		return diags
	}
	diags.Append(mdl.DataSource.fromClient(&dataImport.DataSource)...)
	mdl.FinishedAt = types.StringValue(dataImport.FinishedAt)
	mdl.ImportCheckErrors = types.StringValue(dataImport.ImportCheckErrors)
	mdl.StartedAt = types.StringValue(dataImport.StartedAt)
	mdl.State = types.StringValue(dataImport.State)
	return diags
}

func databaseDataSourceSchemaAttribute(computedName bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"organization":                     schema.StringAttribute{Required: true},
		"name":                             schema.StringAttribute{Required: !computedName, Computed: computedName},
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
			Attributes: regionDataSourceSchemaAttribute,
		},
		"require_approval_for_deploy": schema.BoolAttribute{Computed: true, Optional: true},
		"restrict_branch_region":      schema.BoolAttribute{Computed: true, Optional: true},
		"schema_last_updated_at":      schema.StringAttribute{Computed: true},
		"sharded":                     schema.BoolAttribute{Computed: true},
		"state":                       schema.StringAttribute{Computed: true},
		"type":                        schema.StringAttribute{Computed: true},
		"updated_at":                  schema.StringAttribute{Computed: true},
		"url":                         schema.StringAttribute{Computed: true},
	}
}

type databaseDataSourceModel struct {
	AllowDataBranching                types.Bool                 `tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      types.Bool                 `tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          types.Bool                 `tfsdk:"at_development_branch_limit"`
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
	Region                            regionDataSourceModel      `tfsdk:"region"`
	RequireApprovalForDeploy          types.Bool                 `tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              types.Bool                 `tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               types.String               `tfsdk:"schema_last_updated_at"`
	Sharded                           types.Bool                 `tfsdk:"sharded"`
	State                             types.String               `tfsdk:"state"`
	Type                              types.String               `tfsdk:"type"`
	UpdatedAt                         types.String               `tfsdk:"updated_at"`
	Url                               types.String               `tfsdk:"url"`
}

func (mdl *databaseDataSourceModel) fromClient(database *planetscale.Database) (diags diag.Diagnostics) {
	if database == nil {
		return diags
	}
	if database.DataImport != nil {
		mdl.DataImport = &dataImportDataSourceModel{}
	}
	diags.Append(mdl.DataImport.fromClient(database.DataImport)...)
	diags.Append(mdl.Region.fromClient(&database.Region)...)

	mdl.AllowDataBranching = types.BoolValue(database.AllowDataBranching)
	mdl.AtBackupRestoreBranchesLimit = types.BoolValue(database.AtBackupRestoreBranchesLimit)
	mdl.AtDevelopmentBranchLimit = types.BoolValue(database.AtDevelopmentBranchLimit)
	mdl.AutomaticMigrations = types.BoolPointerValue(database.AutomaticMigrations)
	mdl.BranchesCount = types.Float64Value(database.BranchesCount)
	mdl.BranchesUrl = types.StringValue(database.BranchesUrl)
	mdl.CreatedAt = types.StringValue(database.CreatedAt)
	mdl.DefaultBranch = types.StringValue(database.DefaultBranch)
	mdl.DefaultBranchReadOnlyRegionsCount = types.Float64Value(database.DefaultBranchReadOnlyRegionsCount)
	mdl.DefaultBranchShardCount = types.Float64Value(database.DefaultBranchShardCount)
	mdl.DefaultBranchTableCount = types.Float64Value(database.DefaultBranchTableCount)
	mdl.DevelopmentBranchesCount = types.Float64Value(database.DevelopmentBranchesCount)
	mdl.HtmlUrl = types.StringValue(database.HtmlUrl)
	mdl.Id = types.StringValue(database.Id)
	mdl.InsightsRawQueries = types.BoolValue(database.InsightsRawQueries)
	mdl.IssuesCount = types.Float64Value(database.IssuesCount)
	mdl.MigrationFramework = types.StringPointerValue(database.MigrationFramework)
	mdl.MigrationTableName = types.StringPointerValue(database.MigrationTableName)
	mdl.MultipleAdminsRequiredForDeletion = types.BoolValue(database.MultipleAdminsRequiredForDeletion)
	mdl.Name = types.StringValue(database.Name)
	mdl.Plan = types.StringValue(database.Plan)
	mdl.ProductionBranchWebConsole = types.BoolValue(database.ProductionBranchWebConsole)
	mdl.ProductionBranchesCount = types.Float64Value(database.ProductionBranchesCount)
	mdl.Ready = types.BoolValue(database.Ready)
	mdl.RequireApprovalForDeploy = types.BoolValue(database.RequireApprovalForDeploy)
	mdl.RestrictBranchRegion = types.BoolValue(database.RestrictBranchRegion)
	mdl.SchemaLastUpdatedAt = types.StringPointerValue(database.SchemaLastUpdatedAt)
	mdl.Sharded = types.BoolValue(database.Sharded)
	mdl.State = types.StringValue(database.State)
	mdl.Type = types.StringValue(database.Type)
	mdl.UpdatedAt = types.StringValue(database.UpdatedAt)
	mdl.Url = types.StringValue(database.Url)
	return diags
}

var branchDataSourceSchemaAttribute = map[string]schema.Attribute{
	"organization": schema.StringAttribute{Required: true},
	"database":     schema.StringAttribute{Required: true},
	"name":         schema.StringAttribute{Required: true},

	"access_host_url":                schema.StringAttribute{Computed: true},
	"cluster_rate_name":              schema.StringAttribute{Computed: true},
	"created_at":                     schema.StringAttribute{Computed: true},
	"html_url":                       schema.StringAttribute{Computed: true},
	"id":                             schema.StringAttribute{Computed: true},
	"initial_restore_id":             schema.StringAttribute{Computed: true},
	"mysql_address":                  schema.StringAttribute{Computed: true},
	"mysql_edge_address":             schema.StringAttribute{Computed: true},
	"parent_branch":                  schema.StringAttribute{Computed: true},
	"production":                     schema.BoolAttribute{Computed: true},
	"ready":                          schema.BoolAttribute{Computed: true},
	"restore_checklist_completed_at": schema.StringAttribute{Computed: true},
	"schema_last_updated_at":         schema.StringAttribute{Computed: true},
	"shard_count":                    schema.Float64Attribute{Computed: true},
	"sharded":                        schema.BoolAttribute{Computed: true},
	"updated_at":                     schema.StringAttribute{Computed: true},

	"actor": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: actorDataSourceSchemaAttribute,
	},
	"region": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: regionDataSourceSchemaAttribute,
	},
	"restored_from_branch": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: restoredFromBranchDataSourceSchemaAttribute,
	},
}

type branchDataSourceModel struct {
	AccessHostUrl               types.String                       `tfsdk:"access_host_url"`
	Actor                       *actorDataSourceModel              `tfsdk:"actor"`
	ClusterRateName             types.String                       `tfsdk:"cluster_rate_name"`
	CreatedAt                   types.String                       `tfsdk:"created_at"`
	HtmlUrl                     types.String                       `tfsdk:"html_url"`
	Id                          types.String                       `tfsdk:"id"`
	InitialRestoreId            types.String                       `tfsdk:"initial_restore_id"`
	MysqlAddress                types.String                       `tfsdk:"mysql_address"`
	MysqlEdgeAddress            types.String                       `tfsdk:"mysql_edge_address"`
	Name                        types.String                       `tfsdk:"name"`
	ParentBranch                types.String                       `tfsdk:"parent_branch"`
	Production                  types.Bool                         `tfsdk:"production"`
	Ready                       types.Bool                         `tfsdk:"ready"`
	Region                      *regionDataSourceModel             `tfsdk:"region"`
	RestoreChecklistCompletedAt types.String                       `tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *restoredFromBranchDataSourceModel `tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         types.String                       `tfsdk:"schema_last_updated_at"`
	ShardCount                  types.Float64                      `tfsdk:"shard_count"`
	Sharded                     types.Bool                         `tfsdk:"sharded"`
	UpdatedAt                   types.String                       `tfsdk:"updated_at"`
}

func (mdl *branchDataSourceModel) fromClient(branch *planetscale.Branch) (diags diag.Diagnostics) {
	if branch == nil {
		return diags
	}
	if branch.Actor != nil {
		mdl.Actor = &actorDataSourceModel{}
	}
	diags.Append(mdl.Actor.fromClient(branch.Actor)...)

	if branch.Region != nil {
		mdl.Region = &regionDataSourceModel{}
	}
	diags.Append(mdl.Region.fromClient(branch.Region)...)

	if branch.RestoredFromBranch != nil {
		mdl.RestoredFromBranch = &restoredFromBranchDataSourceModel{}
	}
	diags.Append(mdl.RestoredFromBranch.fromClient(branch.RestoredFromBranch)...)

	mdl.Name = types.StringValue(branch.Name)
	mdl.AccessHostUrl = types.StringPointerValue(branch.AccessHostUrl)
	mdl.ClusterRateName = types.StringValue(branch.ClusterRateName)
	mdl.CreatedAt = types.StringValue(branch.CreatedAt)
	mdl.HtmlUrl = types.StringValue(branch.HtmlUrl)
	mdl.Id = types.StringValue(branch.Id)
	mdl.InitialRestoreId = types.StringPointerValue(branch.InitialRestoreId)
	mdl.MysqlAddress = types.StringValue(branch.MysqlAddress)
	mdl.MysqlEdgeAddress = types.StringValue(branch.MysqlEdgeAddress)
	mdl.ParentBranch = types.StringPointerValue(branch.ParentBranch)
	mdl.Production = types.BoolValue(branch.Production)
	mdl.Ready = types.BoolValue(branch.Ready)
	mdl.RestoreChecklistCompletedAt = types.StringPointerValue(branch.RestoreChecklistCompletedAt)
	mdl.SchemaLastUpdatedAt = types.StringValue(branch.SchemaLastUpdatedAt)
	mdl.ShardCount = types.Float64PointerValue(branch.ShardCount)
	mdl.Sharded = types.BoolValue(branch.Sharded)
	mdl.UpdatedAt = types.StringValue(branch.UpdatedAt)
	return diags
}

var actorDataSourceSchemaAttribute = map[string]schema.Attribute{
	"avatar_url":   schema.StringAttribute{Computed: true},
	"display_name": schema.StringAttribute{Computed: true},
	"id":           schema.StringAttribute{Computed: true},
}

type actorDataSourceModel struct {
	AvatarUrl   types.String `tfsdk:"avatar_url"`
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

func (mdl *actorDataSourceModel) fromClient(actor *planetscale.Actor) (diags diag.Diagnostics) {
	if actor == nil {
		return diags
	}
	mdl.AvatarUrl = types.StringValue(actor.AvatarUrl)
	mdl.DisplayName = types.StringValue(actor.DisplayName)
	mdl.Id = types.StringValue(actor.Id)
	return diags
}

var regionDataSourceSchemaAttribute = map[string]schema.Attribute{
	"display_name":        schema.StringAttribute{Computed: true},
	"enabled":             schema.BoolAttribute{Computed: true},
	"id":                  schema.StringAttribute{Computed: true},
	"location":            schema.StringAttribute{Computed: true},
	"provider":            schema.StringAttribute{Computed: true},
	"public_ip_addresses": schema.ListAttribute{Computed: true, ElementType: types.StringType},
	"slug":                schema.StringAttribute{Computed: true},
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

func (mdl *regionDataSourceModel) fromClient(region *planetscale.Region) (diags diag.Diagnostics) {
	if region == nil {
		return diags
	}
	mdl.DisplayName = types.StringValue(region.DisplayName)
	mdl.Enabled = types.BoolValue(region.Enabled)
	mdl.Id = types.StringValue(region.Id)
	mdl.Location = types.StringValue(region.Location)
	mdl.Provider = types.StringValue(region.Provider)
	mdl.PublicIpAddresses = stringsToListValue(region.PublicIpAddresses, diags)
	mdl.Slug = types.StringValue(region.Slug)
	return diags
}

var readOnlyRegionDataSourceSchemaAttribute = map[string]schema.Attribute{
	"organization": schema.StringAttribute{Required: true},
	"name":         schema.StringAttribute{Required: true},
	"regions": schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"actor": schema.SingleNestedAttribute{
					Computed:   true,
					Attributes: actorDataSourceSchemaAttribute,
				},
				"created_at":   schema.StringAttribute{Computed: true},
				"display_name": schema.StringAttribute{Computed: true},
				"id":           schema.StringAttribute{Computed: true},
				"ready":        schema.BoolAttribute{Computed: true},
				"ready_at":     schema.StringAttribute{Computed: true},
				"updated_at":   schema.StringAttribute{Computed: true},
				"region": schema.SingleNestedAttribute{
					Computed:   true,
					Attributes: regionDataSourceSchemaAttribute,
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

func (mdl *readOnlyRegionDataSourceModel) fromClient(readOnlyRegion *planetscale.ReadOnlyRegion) (diags diag.Diagnostics) {
	if readOnlyRegion == nil {
		return diags
	}
	diags.Append(mdl.Actor.fromClient(&readOnlyRegion.Actor)...)
	diags.Append(mdl.Region.fromClient(&readOnlyRegion.Region)...)

	mdl.CreatedAt = types.StringValue(readOnlyRegion.CreatedAt)
	mdl.DisplayName = types.StringValue(readOnlyRegion.DisplayName)
	mdl.Id = types.StringValue(readOnlyRegion.Id)
	mdl.Ready = types.BoolValue(readOnlyRegion.Ready)
	mdl.ReadyAt = types.StringValue(readOnlyRegion.ReadyAt)
	mdl.UpdatedAt = types.StringValue(readOnlyRegion.UpdatedAt)

	return diags
}

var restoredFromBranchDataSourceSchemaAttribute = map[string]schema.Attribute{
	"created_at": schema.StringAttribute{Computed: true},
	"deleted_at": schema.StringAttribute{Computed: true},
	"id":         schema.StringAttribute{Computed: true},
	"name":       schema.StringAttribute{Computed: true},
	"updated_at": schema.StringAttribute{Computed: true},
}

type restoredFromBranchDataSourceModel struct {
	CreatedAt types.String `tfsdk:"created_at"`
	DeletedAt types.String `tfsdk:"deleted_at"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (mdl *restoredFromBranchDataSourceModel) fromClient(rfb *planetscale.RestoredFromBranch) (diags diag.Diagnostics) {
	if rfb == nil {
		return diags
	}
	mdl.CreatedAt = types.StringValue(rfb.CreatedAt)
	mdl.DeletedAt = types.StringValue(rfb.DeletedAt)
	mdl.Id = types.StringValue(rfb.Id)
	mdl.Name = types.StringValue(rfb.Name)
	mdl.UpdatedAt = types.StringValue(rfb.UpdatedAt)
	return diags
}

var tableSchemaDataSourceSchemaAttribute = map[string]schema.Attribute{
	"html": schema.StringAttribute{Computed: true},
	"name": schema.StringAttribute{Computed: true},
	"raw":  schema.StringAttribute{Computed: true},
}

type tableSchemaDataSourceModel struct {
	Html types.String `tfsdk:"html"`
	Name types.String `tfsdk:"name"`
	Raw  types.String `tfsdk:"raw"`
}

func (mdl *tableSchemaDataSourceModel) fromClient(ts *planetscale.TableSchema) (diags diag.Diagnostics) {
	if ts == nil {
		return diags
	}
	mdl.Html = types.StringValue(ts.Html)
	mdl.Name = types.StringValue(ts.Name)
	mdl.Raw = types.StringValue(ts.Raw)
	return diags
}

var lintErrorDataSourceSchemaAttribute = map[string]schema.Attribute{
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

func (mdl *lintErrorDataSourceModel) fromClient(le *planetscale.LintError) (diags diag.Diagnostics) {
	if le == nil {
		return diags
	}
	mdl.AutoIncrementColumnNames = stringsToListValue(le.AutoIncrementColumnNames, diags)
	mdl.CharsetName = types.StringValue(le.CharsetName)
	mdl.CheckConstraintName = types.StringValue(le.CheckConstraintName)
	mdl.ColumnName = types.StringValue(le.ColumnName)
	mdl.DocsUrl = types.StringValue(le.DocsUrl)
	mdl.EngineName = types.StringValue(le.EngineName)
	mdl.EnumValue = types.StringValue(le.EnumValue)
	mdl.ErrorDescription = types.StringValue(le.ErrorDescription)
	mdl.ForeignKeyColumnNames = stringsToListValue(le.ForeignKeyColumnNames, diags)
	mdl.JsonPath = types.StringValue(le.JsonPath)
	mdl.KeyspaceName = types.StringValue(le.KeyspaceName)
	mdl.LintError = types.StringValue(le.LintError)
	mdl.PartitionName = types.StringValue(le.PartitionName)
	mdl.PartitioningType = types.StringValue(le.PartitioningType)
	mdl.SubjectType = types.StringValue(le.SubjectType)
	mdl.TableName = types.StringValue(le.TableName)
	mdl.VindexName = types.StringValue(le.VindexName)
	return diags
}

var oauthApplicationAttribute = map[string]schema.Attribute{
	"avatar":       schema.StringAttribute{Computed: true},
	"client_id":    schema.StringAttribute{Computed: true},
	"created_at":   schema.StringAttribute{Computed: true},
	"domain":       schema.StringAttribute{Computed: true},
	"id":           schema.StringAttribute{Computed: true},
	"name":         schema.StringAttribute{Computed: true},
	"redirect_uri": schema.StringAttribute{Computed: true},
	"scopes":       schema.ListAttribute{Computed: true, ElementType: types.StringType},
	"tokens":       schema.Float64Attribute{Computed: true},
	"updated_at":   schema.StringAttribute{Computed: true},
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

func (mdl *oauthApplicationDataSourceModel) fromClient(oa *planetscale.OauthApplication) (diags diag.Diagnostics) {
	if oa == nil {
		return diags
	}
	mdl.Avatar = types.StringPointerValue(oa.Avatar)
	mdl.ClientId = types.StringValue(oa.ClientId)
	mdl.CreatedAt = types.StringValue(oa.CreatedAt)
	mdl.Domain = types.StringValue(oa.Domain)
	mdl.Id = types.StringValue(oa.Id)
	mdl.Name = types.StringValue(oa.Name)
	mdl.RedirectUri = types.StringValue(oa.RedirectUri)
	mdl.Scopes = stringsToListValue(oa.Scopes, diags)
	mdl.Tokens = types.Float64Value(oa.Tokens)
	mdl.UpdatedAt = types.StringValue(oa.UpdatedAt)
	return diags
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
