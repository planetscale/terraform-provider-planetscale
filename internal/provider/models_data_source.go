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

func organizationFromClient(org *planetscale.Organization, diags diag.Diagnostics) *organizationDataSourceModel {
	if org == nil {
		return nil
	}
	return &organizationDataSourceModel{
		Features:                  featuresFromClient(org.Features, diags),
		Flags:                     flagsFromClient(org.Flags, diags),
		AdminOnlyProductionAccess: types.BoolValue(org.AdminOnlyProductionAccess),
		BillingEmail:              types.StringPointerValue(org.BillingEmail),
		CanCreateDatabases:        types.BoolValue(org.CanCreateDatabases),
		CreatedAt:                 types.StringValue(org.CreatedAt),
		DatabaseCount:             types.Float64Value(org.DatabaseCount),
		FreeDatabasesRemaining:    types.Float64Value(org.FreeDatabasesRemaining),
		HasPastDueInvoices:        types.BoolValue(org.HasPastDueInvoices),
		Id:                        types.StringValue(org.Id),
		Name:                      types.StringValue(org.Name),
		Plan:                      types.StringValue(org.Plan),
		SingleTenancy:             types.BoolValue(org.SingleTenancy),
		SleepingDatabaseCount:     types.Float64Value(org.SleepingDatabaseCount),
		Sso:                       types.BoolValue(org.Sso),
		SsoDirectory:              types.BoolValue(org.SsoDirectory),
		SsoPortalUrl:              types.StringPointerValue(org.SsoPortalUrl),
		UpdatedAt:                 types.StringValue(org.UpdatedAt),
		ValidBillingInfo:          types.BoolValue(org.ValidBillingInfo),
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

func featuresFromClient(features *planetscale.Features, diags diag.Diagnostics) *featuresDataSourceModel {
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

func flagsFromClient(flags *planetscale.Flags, diags diag.Diagnostics) *flagsDataSourceModel {
	if flags == nil {
		return nil
	}
	return &flagsDataSourceModel{
		ExampleFlag: types.StringPointerValue(flags.ExampleFlag),
	}
}

type dataSourceDataSourceModel struct {
	Database types.String `tfsdk:"database"`
	Hostname types.String `tfsdk:"hostname"`
	Port     types.String `tfsdk:"port"`
}

func dataSourceFromClient(dataSource planetscale.DataSource, diags diag.Diagnostics) dataSourceDataSourceModel {
	return dataSourceDataSourceModel{
		Database: types.StringValue(dataSource.Database),
		Hostname: types.StringValue(dataSource.Hostname),
		Port:     types.StringValue(dataSource.Port),
	}
}

type dataImportDataSourceModel struct {
	DataSource        dataSourceDataSourceModel `tfsdk:"data_source"`
	FinishedAt        types.String              `tfsdk:"finished_at"`
	ImportCheckErrors types.String              `tfsdk:"import_check_errors"`
	StartedAt         types.String              `tfsdk:"started_at"`
	State             types.String              `tfsdk:"state"`
}

func dataImportFromClient(dataImport *planetscale.DataImport, diags diag.Diagnostics) *dataImportDataSourceModel {
	if dataImport == nil {
		return nil
	}
	return &dataImportDataSourceModel{
		DataSource:        dataSourceFromClient(dataImport.DataSource, diags),
		FinishedAt:        types.StringValue(dataImport.FinishedAt),
		ImportCheckErrors: types.StringValue(dataImport.ImportCheckErrors),
		StartedAt:         types.StringValue(dataImport.StartedAt),
		State:             types.StringValue(dataImport.State),
	}
}

func databaseDataSourceSchemaAttribute(computedName bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"organization":                     schema.StringAttribute{Required: !computedName, Computed: computedName},
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
	Organization                      string                     `tfsdk:"organization"`
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
	Region                            *regionDataSourceModel     `tfsdk:"region"`
	RequireApprovalForDeploy          types.Bool                 `tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              types.Bool                 `tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               types.String               `tfsdk:"schema_last_updated_at"`
	Sharded                           types.Bool                 `tfsdk:"sharded"`
	State                             types.String               `tfsdk:"state"`
	Type                              types.String               `tfsdk:"type"`
	UpdatedAt                         types.String               `tfsdk:"updated_at"`
	Url                               types.String               `tfsdk:"url"`
}

func databaseFromClient(database *planetscale.Database, orgName string, diags diag.Diagnostics) *databaseDataSourceModel {
	if database == nil {
		return nil
	}
	return &databaseDataSourceModel{
		Organization:                      orgName,
		DataImport:                        dataImportFromClient(database.DataImport, diags),
		Region:                            regionFromClient(&database.Region, diags),
		AllowDataBranching:                types.BoolValue(database.AllowDataBranching),
		AtBackupRestoreBranchesLimit:      types.BoolValue(database.AtBackupRestoreBranchesLimit),
		AtDevelopmentBranchLimit:          types.BoolValue(database.AtDevelopmentBranchLimit),
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
		Type:                              types.StringValue(database.Type),
		UpdatedAt:                         types.StringValue(database.UpdatedAt),
		Url:                               types.StringValue(database.Url),
	}
}

func branchDataSourceSchemaAttribute(computedName bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: !computedName, Computed: computedName},
		"database":     schema.StringAttribute{Required: !computedName, Computed: computedName},
		"name":         schema.StringAttribute{Required: !computedName, Computed: computedName},

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
		Actor:                       actorFromClient(branch.Actor, diags),
		Region:                      regionFromClient(branch.Region, diags),
		RestoredFromBranch:          restoredFromBranchFromClient(branch.RestoredFromBranch, diags),
		Name:                        types.StringValue(branch.Name),
		AccessHostUrl:               types.StringPointerValue(branch.AccessHostUrl),
		ClusterRateName:             types.StringValue(branch.ClusterRateName),
		CreatedAt:                   types.StringValue(branch.CreatedAt),
		HtmlUrl:                     types.StringValue(branch.HtmlUrl),
		Id:                          types.StringValue(branch.Id),
		InitialRestoreId:            types.StringPointerValue(branch.InitialRestoreId),
		MysqlAddress:                types.StringValue(branch.MysqlAddress),
		MysqlEdgeAddress:            types.StringValue(branch.MysqlEdgeAddress),
		ParentBranch:                types.StringPointerValue(branch.ParentBranch),
		Production:                  types.BoolValue(branch.Production),
		Ready:                       types.BoolValue(branch.Ready),
		RestoreChecklistCompletedAt: types.StringPointerValue(branch.RestoreChecklistCompletedAt),
		SchemaLastUpdatedAt:         types.StringValue(branch.SchemaLastUpdatedAt),
		ShardCount:                  types.Float64PointerValue(branch.ShardCount),
		Sharded:                     types.BoolValue(branch.Sharded),
		UpdatedAt:                   types.StringValue(branch.UpdatedAt),
	}
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

func actorFromClient(actor *planetscale.Actor, diags diag.Diagnostics) *actorDataSourceModel {
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

func readOnlyRegionFromClient(readOnlyRegion *planetscale.ReadOnlyRegion, diags diag.Diagnostics) *readOnlyRegionDataSourceModel {
	if readOnlyRegion == nil {
		return nil
	}

	return &readOnlyRegionDataSourceModel{
		Actor:       *actorFromClient(&readOnlyRegion.Actor, diags),
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

func restoredFromBranchFromClient(rfb *planetscale.RestoredFromBranch, diags diag.Diagnostics) *restoredFromBranchDataSourceModel {
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
	"html": schema.StringAttribute{Computed: true},
	"name": schema.StringAttribute{Computed: true},
	"raw":  schema.StringAttribute{Computed: true},
}

type tableSchemaDataSourceModel struct {
	Html types.String `tfsdk:"html"`
	Name types.String `tfsdk:"name"`
	Raw  types.String `tfsdk:"raw"`
}

func tableSchemaFromClient(ts *planetscale.TableSchema, diags diag.Diagnostics) *tableSchemaDataSourceModel {
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
