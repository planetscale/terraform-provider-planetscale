// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &databaseResource{}
var _ resource.ResourceWithImportState = &databaseResource{}

func newDatabaseResource() resource.Resource {
	return &databaseResource{}
}

// databaseResource defines the resource implementation.
type databaseResource struct {
	client *planetscale.Client
}

type importDataSourceResourceModel struct {
	Database types.String `tfsdk:"database"`
	Hostname types.String `tfsdk:"hostname"`
	Port     types.String `tfsdk:"port"`
}

type importResourceModel struct {
	DataSource        importDataSourceResourceModel `tfsdk:"data_source"`
	FinishedAt        types.String                  `tfsdk:"finished_at"`
	ImportCheckErrors types.String                  `tfsdk:"import_check_errors"`
	StartedAt         types.String                  `tfsdk:"started_at"`
	State             types.String                  `tfsdk:"state"`
}

type regionResourceModel struct {
	DisplayName       types.String `tfsdk:"display_name"`
	Enabled           types.String `tfsdk:"enabled"`
	Id                types.String `tfsdk:"id"`
	Location          types.String `tfsdk:"location"`
	Provider          types.String `tfsdk:"provider"`
	PublicIpAddresses types.List   `tfsdk:"public_ip_addresses"`
	Slug              types.String `tfsdk:"slug"`
}

type databaseResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Id           types.String `tfsdk:"id"`

	AllowDataBranching                types.Bool           `tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      types.Bool           `tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          types.Bool           `tfsdk:"at_development_branch_limit"`
	AutomaticMigrations               types.Bool           `tfsdk:"automatic_migrations"`
	BranchesCount                     types.Float64        `tfsdk:"branches_count"`
	BranchesUrl                       types.String         `tfsdk:"branches_url"`
	CreatedAt                         types.String         `tfsdk:"created_at"`
	DataImport                        *importResourceModel `tfsdk:"data_import"`
	DefaultBranch                     types.String         `tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount types.Float64        `tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           types.Float64        `tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           types.Float64        `tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          types.Float64        `tfsdk:"development_branches_count"`
	HtmlUrl                           types.String         `tfsdk:"html_url"`
	InsightsRawQueries                types.Bool           `tfsdk:"insights_raw_queries"`
	IssuesCount                       types.Float64        `tfsdk:"issues_count"`
	MigrationFramework                types.String         `tfsdk:"migration_framework"`
	MigrationTableName                types.String         `tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion types.Bool           `tfsdk:"multiple_admins_required_for_deletion"`
	Name                              types.String         `tfsdk:"name"`
	Notes                             types.String         `tfsdk:"notes"`
	Plan                              types.String         `tfsdk:"plan"`
	ProductionBranchWebConsole        types.Bool           `tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           types.Float64        `tfsdk:"production_branches_count"`
	Ready                             types.Bool           `tfsdk:"ready"`
	Region                            regionResourceModel  `tfsdk:"region"`
	RequireApprovalForDeploy          types.Bool           `tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              types.Bool           `tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               types.String         `tfsdk:"schema_last_updated_at"`
	Sharded                           types.Bool           `tfsdk:"sharded"`
	State                             types.String         `tfsdk:"state"`
	Type                              types.String         `tfsdk:"type"`
	UpdatedAt                         types.String         `tfsdk:"updated_at"`
	Url                               types.String         `tfsdk:"url"`
}

func (r *databaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *databaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},

			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_data_branching":             schema.BoolAttribute{Computed: true, Optional: true},
			"at_backup_restore_branches_limit": schema.BoolAttribute{Computed: true},
			"at_development_branch_limit":      schema.BoolAttribute{Computed: true},
			"automatic_migrations":             schema.BoolAttribute{Computed: true, Optional: true},
			"branches_count":                   schema.Float64Attribute{Computed: true},
			"branches_url":                     schema.StringAttribute{Computed: true},
			"created_at":                       schema.StringAttribute{Computed: true},
			"data_import": schema.SingleNestedAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"data_source": schema.SingleNestedAttribute{
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
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"slug": schema.StringAttribute{
						// the actual argument to set the region to use
						Computed: true, Optional: true,
					},
					"display_name":        schema.StringAttribute{Computed: true},
					"enabled":             schema.StringAttribute{Computed: true},
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
	}
}

func (r *databaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*planetscale.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *planetscale.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *databaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *databaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res201, _, _, _, err := r.client.CreateDatabase(ctx, data.Organization.String(), planetscale.CreateDatabaseReq{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create database, got error: %s", err))
		return
	}

	ipAddrValue, diags := types.ListValue(types.StringType, stringListAttr(res201.Region.PublicIpAddresses))
	resp.Diagnostics.Append(diags...)

	data.Id = types.StringValue(res201.Id)

	data.AllowDataBranching = types.BoolValue(res201.AllowDataBranching)
	data.AtBackupRestoreBranchesLimit = types.BoolValue(res201.AtBackupRestoreBranchesLimit)
	data.AtDevelopmentBranchLimit = types.BoolValue(res201.AtDevelopmentBranchLimit)
	data.AutomaticMigrations = types.BoolValue(res201.AutomaticMigrations)
	data.BranchesCount = types.Float64Value(res201.BranchesCount)
	data.BranchesUrl = types.StringValue(res201.BranchesUrl)
	data.CreatedAt = types.StringValue(res201.CreatedAt)
	data.DefaultBranch = types.StringValue(res201.DefaultBranch)
	data.DefaultBranchReadOnlyRegionsCount = types.Float64Value(res201.DefaultBranchReadOnlyRegionsCount)
	data.DefaultBranchShardCount = types.Float64Value(res201.DefaultBranchShardCount)
	data.DefaultBranchTableCount = types.Float64Value(res201.DefaultBranchTableCount)
	data.DevelopmentBranchesCount = types.Float64Value(res201.DevelopmentBranchesCount)
	data.HtmlUrl = types.StringValue(res201.HtmlUrl)
	data.InsightsRawQueries = types.BoolValue(res201.InsightsRawQueries)
	data.IssuesCount = types.Float64Value(res201.IssuesCount)
	data.MigrationFramework = types.StringPointerValue(res201.MigrationFramework)
	data.MigrationTableName = types.StringPointerValue(res201.MigrationTableName)
	data.MultipleAdminsRequiredForDeletion = types.BoolValue(res201.MultipleAdminsRequiredForDeletion)
	data.Name = types.StringValue(res201.Name)
	data.Notes = types.StringPointerValue(res201.Notes)
	data.Plan = types.StringValue(res201.Plan)
	data.ProductionBranchWebConsole = types.BoolValue(res201.ProductionBranchWebConsole)
	data.ProductionBranchesCount = types.Float64Value(res201.ProductionBranchesCount)
	data.Ready = types.BoolValue(res201.Ready)
	data.RequireApprovalForDeploy = types.BoolValue(res201.RequireApprovalForDeploy)
	data.RestrictBranchRegion = types.BoolValue(res201.RestrictBranchRegion)
	data.SchemaLastUpdatedAt = types.StringPointerValue(res201.SchemaLastUpdatedAt)
	data.Sharded = types.BoolValue(res201.Sharded)
	data.State = types.StringValue(res201.State)
	data.Type = types.StringValue(res201.Type)
	data.UpdatedAt = types.StringValue(res201.UpdatedAt)
	data.Url = types.StringValue(res201.Url)
	if res201.DataImport != nil {
		data.DataImport = &importResourceModel{
			DataSource: importDataSourceResourceModel{
				Database: types.StringValue(res201.DataImport.DataSource.Database),
				Hostname: types.StringValue(res201.DataImport.DataSource.Hostname),
				Port:     types.StringValue(res201.DataImport.DataSource.Port),
			},
			FinishedAt:        types.StringValue(res201.DataImport.FinishedAt),
			ImportCheckErrors: types.StringValue(res201.DataImport.ImportCheckErrors),
			StartedAt:         types.StringValue(res201.DataImport.StartedAt),
			State:             types.StringValue(res201.DataImport.State),
		}
	}
	data.Region = regionResourceModel{
		DisplayName:       types.StringValue(res201.Region.DisplayName),
		Enabled:           types.StringValue(res201.Region.Enabled),
		Id:                types.StringValue(res201.Region.Id),
		Location:          types.StringValue(res201.Region.Location),
		Provider:          types.StringValue(res201.Region.Provider),
		PublicIpAddresses: ipAddrValue,
		Slug:              types.StringValue(res201.Region.Slug),
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a database resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *databaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *databaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization
	name := data.Name

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if name.IsNull() || name.IsUnknown() || name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "name is required", "a name must be provided and cannot be empty")
		return
	}

	res200, _, _, _, err := r.client.GetDatabase(ctx, org.ValueString(), name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err))
		return
	}

	ipAddrValue, diags := types.ListValue(types.StringType, stringListAttr(res200.Region.PublicIpAddresses))
	resp.Diagnostics.Append(diags...)

	data.Id = types.StringValue(res200.Id)
	data.AllowDataBranching = types.BoolValue(res200.AllowDataBranching)
	data.AtBackupRestoreBranchesLimit = types.BoolValue(res200.AtBackupRestoreBranchesLimit)
	data.AtDevelopmentBranchLimit = types.BoolValue(res200.AtDevelopmentBranchLimit)
	data.AutomaticMigrations = types.BoolValue(res200.AutomaticMigrations)
	data.BranchesCount = types.Float64Value(res200.BranchesCount)
	data.BranchesUrl = types.StringValue(res200.BranchesUrl)
	data.CreatedAt = types.StringValue(res200.CreatedAt)
	data.DefaultBranch = types.StringValue(res200.DefaultBranch)
	data.DefaultBranchReadOnlyRegionsCount = types.Float64Value(res200.DefaultBranchReadOnlyRegionsCount)
	data.DefaultBranchShardCount = types.Float64Value(res200.DefaultBranchShardCount)
	data.DefaultBranchTableCount = types.Float64Value(res200.DefaultBranchTableCount)
	data.DevelopmentBranchesCount = types.Float64Value(res200.DevelopmentBranchesCount)
	data.HtmlUrl = types.StringValue(res200.HtmlUrl)
	data.InsightsRawQueries = types.BoolValue(res200.InsightsRawQueries)
	data.IssuesCount = types.Float64Value(res200.IssuesCount)
	data.MigrationFramework = types.StringPointerValue(res200.MigrationFramework)
	data.MigrationTableName = types.StringPointerValue(res200.MigrationTableName)
	data.MultipleAdminsRequiredForDeletion = types.BoolValue(res200.MultipleAdminsRequiredForDeletion)
	data.Name = types.StringValue(res200.Name)
	data.Notes = types.StringPointerValue(res200.Notes)
	data.Plan = types.StringValue(res200.Plan)
	data.ProductionBranchWebConsole = types.BoolValue(res200.ProductionBranchWebConsole)
	data.ProductionBranchesCount = types.Float64Value(res200.ProductionBranchesCount)
	data.Ready = types.BoolValue(res200.Ready)
	data.RequireApprovalForDeploy = types.BoolValue(res200.RequireApprovalForDeploy)
	data.RestrictBranchRegion = types.BoolValue(res200.RestrictBranchRegion)
	data.SchemaLastUpdatedAt = types.StringPointerValue(res200.SchemaLastUpdatedAt)
	data.Sharded = types.BoolValue(res200.Sharded)
	data.State = types.StringValue(res200.State)
	data.Type = types.StringValue(res200.Type)
	data.UpdatedAt = types.StringValue(res200.UpdatedAt)
	data.Url = types.StringValue(res200.Url)
	if res200.DataImport != nil {
		data.DataImport = &importResourceModel{
			DataSource: importDataSourceResourceModel{
				Database: types.StringValue(res200.DataImport.DataSource.Database),
				Hostname: types.StringValue(res200.DataImport.DataSource.Hostname),
				Port:     types.StringValue(res200.DataImport.DataSource.Port),
			},
			FinishedAt:        types.StringValue(res200.DataImport.FinishedAt),
			ImportCheckErrors: types.StringValue(res200.DataImport.ImportCheckErrors),
			StartedAt:         types.StringValue(res200.DataImport.StartedAt),
			State:             types.StringValue(res200.DataImport.State),
		}
	}
	data.Region = regionResourceModel{
		DisplayName:       types.StringValue(res200.Region.DisplayName),
		Enabled:           types.StringValue(res200.Region.Enabled),
		Id:                types.StringValue(res200.Region.Id),
		Location:          types.StringValue(res200.Region.Location),
		Provider:          types.StringValue(res200.Region.Provider),
		PublicIpAddresses: ipAddrValue,
		Slug:              types.StringValue(res200.Region.Slug),
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *databaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		old  *databaseResourceModel
		data *databaseResourceModel
	)
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization
	name := data.Name

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if name.IsNull() || name.IsUnknown() || name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "name is required", "a name must be provided and cannot be empty")
		return
	}

	changedUpdatableSettings := false
	updateReq := planetscale.UpdateDatabaseSettingsReq{
		AllowDataBranching:         boolIfDifferent(old.AllowDataBranching, data.AllowDataBranching, &changedUpdatableSettings),
		AutomaticMigrations:        boolIfDifferent(old.AutomaticMigrations, data.AutomaticMigrations, &changedUpdatableSettings),
		DefaultBranch:              stringIfDifferent(old.DefaultBranch, data.DefaultBranch, &changedUpdatableSettings),
		InsightsRawQueries:         boolIfDifferent(old.InsightsRawQueries, data.InsightsRawQueries, &changedUpdatableSettings),
		MigrationFramework:         stringIfDifferent(old.MigrationFramework, data.MigrationFramework, &changedUpdatableSettings),
		MigrationTableName:         stringIfDifferent(old.MigrationTableName, data.MigrationTableName, &changedUpdatableSettings),
		Notes:                      stringIfDifferent(old.Notes, data.Notes, &changedUpdatableSettings),
		ProductionBranchWebConsole: boolIfDifferent(old.ProductionBranchWebConsole, data.ProductionBranchWebConsole, &changedUpdatableSettings),
		RequireApprovalForDeploy:   boolIfDifferent(old.RequireApprovalForDeploy, data.RequireApprovalForDeploy, &changedUpdatableSettings),
		RestrictBranchRegion:       boolIfDifferent(old.RestrictBranchRegion, data.RestrictBranchRegion, &changedUpdatableSettings),
	}

	if changedUpdatableSettings {
		res200, _, _, _, err := r.client.UpdateDatabaseSettings(ctx, org.ValueString(), name.ValueString(), updateReq)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update database settings, got error: %s", err))
			return
		}
		ipAddrValue, diags := types.ListValue(types.StringType, stringListAttr(res200.Region.PublicIpAddresses))
		resp.Diagnostics.Append(diags...)
		data.Id = types.StringValue(res200.Id)
		data.AllowDataBranching = types.BoolValue(res200.AllowDataBranching)
		data.AtBackupRestoreBranchesLimit = types.BoolValue(res200.AtBackupRestoreBranchesLimit)
		data.AtDevelopmentBranchLimit = types.BoolValue(res200.AtDevelopmentBranchLimit)
		data.AutomaticMigrations = types.BoolValue(res200.AutomaticMigrations)
		data.BranchesCount = types.Float64Value(res200.BranchesCount)
		data.BranchesUrl = types.StringValue(res200.BranchesUrl)
		data.CreatedAt = types.StringValue(res200.CreatedAt)
		data.DefaultBranch = types.StringValue(res200.DefaultBranch)
		data.DefaultBranchReadOnlyRegionsCount = types.Float64Value(res200.DefaultBranchReadOnlyRegionsCount)
		data.DefaultBranchShardCount = types.Float64Value(res200.DefaultBranchShardCount)
		data.DefaultBranchTableCount = types.Float64Value(res200.DefaultBranchTableCount)
		data.DevelopmentBranchesCount = types.Float64Value(res200.DevelopmentBranchesCount)
		data.HtmlUrl = types.StringValue(res200.HtmlUrl)
		data.InsightsRawQueries = types.BoolValue(res200.InsightsRawQueries)
		data.IssuesCount = types.Float64Value(res200.IssuesCount)
		data.MigrationFramework = types.StringPointerValue(res200.MigrationFramework)
		data.MigrationTableName = types.StringPointerValue(res200.MigrationTableName)
		data.MultipleAdminsRequiredForDeletion = types.BoolValue(res200.MultipleAdminsRequiredForDeletion)
		data.Name = types.StringValue(res200.Name)
		data.Notes = types.StringPointerValue(res200.Notes)
		data.Plan = types.StringValue(res200.Plan)
		data.ProductionBranchWebConsole = types.BoolValue(res200.ProductionBranchWebConsole)
		data.ProductionBranchesCount = types.Float64Value(res200.ProductionBranchesCount)
		data.Ready = types.BoolValue(res200.Ready)
		data.RequireApprovalForDeploy = types.BoolValue(res200.RequireApprovalForDeploy)
		data.RestrictBranchRegion = types.BoolValue(res200.RestrictBranchRegion)
		data.SchemaLastUpdatedAt = types.StringPointerValue(res200.SchemaLastUpdatedAt)
		data.Sharded = types.BoolValue(res200.Sharded)
		data.State = types.StringValue(res200.State)
		data.Type = types.StringValue(res200.Type)
		data.UpdatedAt = types.StringValue(res200.UpdatedAt)
		data.Url = types.StringValue(res200.Url)
		if res200.DataImport != nil {
			data.DataImport = &importResourceModel{
				DataSource: importDataSourceResourceModel{
					Database: types.StringValue(res200.DataImport.DataSource.Database),
					Hostname: types.StringValue(res200.DataImport.DataSource.Hostname),
					Port:     types.StringValue(res200.DataImport.DataSource.Port),
				},
				FinishedAt:        types.StringValue(res200.DataImport.FinishedAt),
				ImportCheckErrors: types.StringValue(res200.DataImport.ImportCheckErrors),
				StartedAt:         types.StringValue(res200.DataImport.StartedAt),
				State:             types.StringValue(res200.DataImport.State),
			}
		}
		data.Region = regionResourceModel{
			DisplayName:       types.StringValue(res200.Region.DisplayName),
			Enabled:           types.StringValue(res200.Region.Enabled),
			Id:                types.StringValue(res200.Region.Id),
			Location:          types.StringValue(res200.Region.Location),
			Provider:          types.StringValue(res200.Region.Provider),
			PublicIpAddresses: ipAddrValue,
			Slug:              types.StringValue(res200.Region.Slug),
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *databaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *databaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	org := data.Organization
	name := data.Name

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if name.IsNull() || name.IsUnknown() || name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "name is required", "a name must be provided and cannot be empty")
		return
	}

	res204, _, _, _, err := r.client.DeleteDatabase(ctx, org.ValueString(), name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete database, got error: %s", err))
		return
	}
	_ = res204
}

func (r *databaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization,name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)

}
