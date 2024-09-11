// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &databaseResource{}
	_ resource.ResourceWithImportState = &databaseResource{}
)

func newDatabaseResource() resource.Resource {
	return &databaseResource{}
}

// databaseResource defines the resource implementation.
type databaseResource struct {
	client *planetscale.Client
}

type databaseResourceModel struct {
	Organization                      types.String  `tfsdk:"organization"`
	Id                                types.String  `tfsdk:"id"`
	AllowDataBranching                types.Bool    `tfsdk:"allow_data_branching"`
	AtBackupRestoreBranchesLimit      types.Bool    `tfsdk:"at_backup_restore_branches_limit"`
	AtDevelopmentBranchLimit          types.Bool    `tfsdk:"at_development_branch_limit"`
	AutomaticMigrations               types.Bool    `tfsdk:"automatic_migrations"`
	BranchesCount                     types.Float64 `tfsdk:"branches_count"`
	BranchesUrl                       types.String  `tfsdk:"branches_url"`
	CreatedAt                         types.String  `tfsdk:"created_at"`
	DataImport                        types.Object  `tfsdk:"data_import"`
	DefaultBranch                     types.String  `tfsdk:"default_branch"`
	DefaultBranchReadOnlyRegionsCount types.Float64 `tfsdk:"default_branch_read_only_regions_count"`
	DefaultBranchShardCount           types.Float64 `tfsdk:"default_branch_shard_count"`
	DefaultBranchTableCount           types.Float64 `tfsdk:"default_branch_table_count"`
	DevelopmentBranchesCount          types.Float64 `tfsdk:"development_branches_count"`
	HtmlUrl                           types.String  `tfsdk:"html_url"`
	InsightsRawQueries                types.Bool    `tfsdk:"insights_raw_queries"`
	IssuesCount                       types.Float64 `tfsdk:"issues_count"`
	MigrationFramework                types.String  `tfsdk:"migration_framework"`
	MigrationTableName                types.String  `tfsdk:"migration_table_name"`
	MultipleAdminsRequiredForDeletion types.Bool    `tfsdk:"multiple_admins_required_for_deletion"`
	Name                              types.String  `tfsdk:"name"`
	Plan                              types.String  `tfsdk:"plan"`
	ClusterSize                       types.String  `tfsdk:"cluster_size"`
	ProductionBranchWebConsole        types.Bool    `tfsdk:"production_branch_web_console"`
	ProductionBranchesCount           types.Float64 `tfsdk:"production_branches_count"`
	Ready                             types.Bool    `tfsdk:"ready"`
	Region                            types.String  `tfsdk:"region"`
	RequireApprovalForDeploy          types.Bool    `tfsdk:"require_approval_for_deploy"`
	RestrictBranchRegion              types.Bool    `tfsdk:"restrict_branch_region"`
	SchemaLastUpdatedAt               types.String  `tfsdk:"schema_last_updated_at"`
	Sharded                           types.Bool    `tfsdk:"sharded"`
	State                             types.String  `tfsdk:"state"`
	UpdatedAt                         types.String  `tfsdk:"updated_at"`
	Url                               types.String  `tfsdk:"url"`
}

func databaseResourcefromClient(ctx context.Context, database *planetscale.Database, organization, clusterSize types.String, diags diag.Diagnostics) *databaseResourceModel {
	if database == nil {
		return nil
	}
	if clusterSize.IsUnknown() {
		clusterSize = types.StringNull()
	}
	dataImport, diags := types.ObjectValueFrom(ctx, importResourceAttrTypes, database.DataImport)
	diags.Append(diags...)
	return &databaseResourceModel{
		Organization:                      organization,
		DataImport:                        dataImport,
		Id:                                types.StringValue(database.Id),
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
		InsightsRawQueries:                types.BoolValue(database.InsightsRawQueries),
		IssuesCount:                       types.Float64Value(database.IssuesCount),
		MigrationFramework:                types.StringPointerValue(database.MigrationFramework),
		MigrationTableName:                types.StringPointerValue(database.MigrationTableName),
		MultipleAdminsRequiredForDeletion: types.BoolValue(database.MultipleAdminsRequiredForDeletion),
		Name:                              types.StringValue(database.Name),
		Plan:                              types.StringValue(database.Plan),
		ClusterSize:                       clusterSize,
		ProductionBranchWebConsole:        types.BoolValue(database.ProductionBranchWebConsole),
		ProductionBranchesCount:           types.Float64Value(database.ProductionBranchesCount),
		Ready:                             types.BoolValue(database.Ready),
		Region:                            types.StringValue(database.Region.Slug),
		RequireApprovalForDeploy:          types.BoolValue(database.RequireApprovalForDeploy),
		RestrictBranchRegion:              types.BoolValue(database.RestrictBranchRegion),
		SchemaLastUpdatedAt:               types.StringPointerValue(database.SchemaLastUpdatedAt),
		Sharded:                           types.BoolValue(database.Sharded),
		State:                             types.StringValue(database.State),
		UpdatedAt:                         types.StringValue(database.UpdatedAt),
		Url:                               types.StringValue(database.Url),
	}
}

func (r *databaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *databaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A PlanetScale database.",
		MarkdownDescription: `A PlanetScale database.

Known limitations:
- When the provider is configured with a service token, the service token needs to manually be granted permission on this database resource. This can be done in the UI or via the CLI (` + "`pscale service-token add-access`" + `).`,
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description: "The organization this database belongs to.",
				Required:    true, PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of this database.",
				Required:    true, PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"id": schema.StringAttribute{
				Description: "The ID of the database.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Computed:    true,
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
			"cluster_size": schema.StringAttribute{
				Description: "The size of the database cluster plan.",
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
			"region": schema.StringAttribute{
				Description: "The region the database lives in.",
				Computed:    true, Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
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
	tflog.Info(ctx, "getting current database resource from plan")
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

	createDbReq := planetscale.CreateDatabaseReq{
		Name:        name.ValueString(),
		Plan:        stringValueIfKnown(data.Plan),
		ClusterSize: stringValueIfKnown(data.ClusterSize),
		Region:      stringValueIfKnown(data.Region),
	}
	res, err := r.client.CreateDatabase(ctx, org.ValueString(), createDbReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create database, got error: %s", err))
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to create databases", "no data")
		return
	}

	// wait for db to enter ready state
	createState := &retry.StateChangeConf{
		Delay:      5 * time.Second, // initial delay before the first check
		Timeout:    10 * time.Minute,
		MinTimeout: 5 * time.Second,

		Pending: []string{"not-ready"},
		Target:  []string{"ready"},

		Refresh: func() (interface{}, string, error) {
			res, err := r.client.GetDatabase(ctx, org.ValueString(), name.ValueString())
			if err != nil {
				return nil, "", err
			}
			if res.Database.Ready {
				return res, "ready", nil
			}
			return res, "not-ready", nil
		},
	}
	_, err = createState.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create database",
			fmt.Sprintf("Database %s never became ready; got error: %s", name.ValueString(), err),
		)
		return
	}

	data = databaseResourcefromClient(ctx, &res.Database, data.Organization, data.ClusterSize, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *databaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *databaseResourceModel

	tflog.Info(ctx, "getting current database resource from state")
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

	res, err := r.client.GetDatabase(ctx, org.ValueString(), name.ValueString())
	if err != nil {
		if notFoundErr, ok := err.(*planetscale.GetDatabaseRes404); ok {
			tflog.Warn(ctx, fmt.Sprintf("Database not found, removing from state: %s", notFoundErr.Message))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err))
		return
	}

	data = databaseResourcefromClient(ctx, &res.Database, data.Organization, data.ClusterSize, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
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
		ProductionBranchWebConsole: boolIfDifferent(old.ProductionBranchWebConsole, data.ProductionBranchWebConsole, &changedUpdatableSettings),
		RequireApprovalForDeploy:   boolIfDifferent(old.RequireApprovalForDeploy, data.RequireApprovalForDeploy, &changedUpdatableSettings),
		RestrictBranchRegion:       boolIfDifferent(old.RestrictBranchRegion, data.RestrictBranchRegion, &changedUpdatableSettings),
	}

	if changedUpdatableSettings {
		res, err := r.client.UpdateDatabaseSettings(ctx, org.ValueString(), name.ValueString(), updateReq)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update database settings, got error: %s", err))
			return
		}
		data = databaseResourcefromClient(ctx, &res.Database, data.Organization, data.ClusterSize, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

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

	_, err := r.client.DeleteDatabase(ctx, org.ValueString(), name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete database, got error: %s", err))
		return
	}
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
