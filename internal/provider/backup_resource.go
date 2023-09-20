// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &backupResource{}
var _ resource.ResourceWithImportState = &backupResource{}

func newBackupResource() resource.Resource {
	return &backupResource{}
}

// backupResource defines the resource implementation.
type backupResource struct {
	client *planetscale.Client
}

type backupResourceModel struct {
	Organization         types.String  `tfsdk:"organization"`
	Database             types.String  `tfsdk:"database"`
	Branch               types.String  `tfsdk:"branch"`
	Name                 types.String  `tfsdk:"name"`
	Actor                types.Object  `tfsdk:"actor"`
	BackupPolicy         types.Object  `tfsdk:"backup_policy"`
	CreatedAt            types.String  `tfsdk:"created_at"`
	EstimatedStorageCost types.String  `tfsdk:"estimated_storage_cost"`
	Id                   types.String  `tfsdk:"id"`
	Required             types.Bool    `tfsdk:"required"`
	RestoredBranches     types.List    `tfsdk:"restored_branches"`
	SchemaSnapshot       types.Object  `tfsdk:"schema_snapshot"`
	Size                 types.Float64 `tfsdk:"size"`
	State                types.String  `tfsdk:"state"`
	UpdatedAt            types.String  `tfsdk:"updated_at"`
}

func backupResourceFromClient(ctx context.Context, backup *planetscale.Backup, organization, database, branch types.String, diags diag.Diagnostics) *backupResourceModel {
	if backup == nil {
		return nil
	}
	actor, diags := types.ObjectValueFrom(ctx, actorResourceAttrTypes, backup.Actor)
	diags.Append(diags...)
	backupPolicy, diags := types.ObjectValueFrom(ctx, backupPolicyResourceAttrTypes, backup.BackupPolicy)
	diags.Append(diags...)
	schemaSnapshot, diags := types.ObjectValueFrom(ctx, schemaSnapshotResourceAttrTypes, backup.SchemaSnapshot)
	diags.Append(diags...)

	restoredBranch := types.ListNull(types.StringType)
	if backup.RestoredBranches != nil {
		restoredBranch = stringsToListValue(*backup.RestoredBranches, diags)
	}
	return &backupResourceModel{
		Organization: organization,
		Database:     database,
		Branch:       branch,

		// partially required
		BackupPolicy: backupPolicy,

		Name:                 types.StringValue(backup.Name),
		Actor:                actor,
		SchemaSnapshot:       schemaSnapshot,
		CreatedAt:            types.StringValue(backup.CreatedAt),
		EstimatedStorageCost: types.StringValue(backup.EstimatedStorageCost),
		Id:                   types.StringValue(backup.Id),
		Required:             types.BoolValue(backup.Required),
		RestoredBranches:     restoredBranch,
		Size:                 types.Float64Value(backup.Size),
		State:                types.StringValue(backup.State),
		UpdatedAt:            types.StringValue(backup.UpdatedAt),
	}
}

func (r *backupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup"
}

func (r *backupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A PlanetScale backup",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"database": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"branch": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"backup_policy": schema.SingleNestedAttribute{
				Required:   true,
				Attributes: backupPolicyResourceAttribute,
			},

			// read only
			"actor": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: actorResourceSchemaAttribute,
			},
			"schema_snapshot": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: schemaSnapshotResourceAttribute,
			},
			"id":                     schema.StringAttribute{Computed: true},
			"created_at":             schema.StringAttribute{Computed: true},
			"estimated_storage_cost": schema.StringAttribute{Computed: true},
			"required":               schema.BoolAttribute{Computed: true},
			"restored_branches":      schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"size":                   schema.Float64Attribute{Computed: true},
			"state":                  schema.StringAttribute{Computed: true},
			"updated_at":             schema.StringAttribute{Computed: true},
		},
	}
}

func (r *backupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *backupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *backupResourceModel
	tflog.Info(ctx, "getting current backup resource from plan")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization
	database := data.Database
	branch := data.Branch
	name := data.Name
	backupPolicy := data.BackupPolicy

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if database.IsNull() || database.IsUnknown() || database.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("database"), "database is required", "a database must be provided and cannot be empty")
		return
	}
	if branch.IsNull() || branch.IsUnknown() || branch.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("branch"), "branch is required", "a branch must be provided and cannot be empty")
		return
	}
	if name.IsNull() || name.IsUnknown() || name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "name is required", "a name must be provided and cannot be empty")
		return
	}
	if backupPolicy.IsNull() || backupPolicy.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("backup_policy"), "backup_policy is required", "a backup_policy must be provided and cannot be empty")
		return
	}
	var bp backupPolicyDataSourceModel
	resp.Diagnostics.Append(backupPolicy.As(ctx, &bp, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := planetscale.CreateBackupReq{
		Name:           name.ValueStringPointer(),
		RetentionUnit:  bp.RetentionUnit.ValueStringPointer(),
		RetentionValue: bp.RetentionValue.ValueFloat64Pointer(),
	}
	res, err := r.client.CreateBackup(ctx, org.ValueString(), database.ValueString(), branch.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create backup, got error: %s", err))
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to create backups", "no data")
		return
	}

	data = backupResourceFromClient(ctx, &res.Backup, data.Organization, data.Database, data.Branch, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *backupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *backupResourceModel

	tflog.Info(ctx, "getting current backup resource from state")
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization
	database := data.Database
	branch := data.Branch
	id := data.Id

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if database.IsNull() || database.IsUnknown() || database.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("database"), "database is required", "a database must be provided and cannot be empty")
		return
	}
	if branch.IsNull() || branch.IsUnknown() || branch.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("branch"), "branch is required", "a branch must be provided and cannot be empty")
		return
	}
	if id.IsNull() || id.IsUnknown() || id.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "id is required", "an ID must be provided and cannot be empty")
		return
	}

	res, err := r.client.GetBackup(ctx, org.ValueString(), database.ValueString(), branch.ValueString(), id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read backup, got error: %s", err))
		return
	}

	data = backupResourceFromClient(ctx, &res.Backup, data.Organization, data.Database, data.Branch, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *backupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// nothing to do, backups have no updatable settings
}

func (r *backupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *backupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	org := data.Organization
	database := data.Database
	branch := data.Branch
	id := data.Id

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if database.IsNull() || database.IsUnknown() || database.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("database"), "database is required", "a database must be provided and cannot be empty")
		return
	}
	if branch.IsNull() || branch.IsUnknown() || branch.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("branch"), "branch is required", "a branch must be provided and cannot be empty")
		return
	}
	if id.IsNull() || id.IsUnknown() || id.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "id is required", "an ID must be provided and cannot be empty")
		return
	}

	res, err := r.client.DeleteBackup(ctx, org.ValueString(), database.ValueString(), branch.ValueString(), id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete backup, got error: %s", err))
		return
	}
	_ = res
}

func (r *backupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization,database,name,id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("branch"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[3])...)
}
