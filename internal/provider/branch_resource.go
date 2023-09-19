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
var _ resource.Resource = &branchResource{}
var _ resource.ResourceWithImportState = &branchResource{}

func newBranchResource() resource.Resource {
	return &branchResource{}
}

// branchResource defines the resource implementation.
type branchResource struct {
	client *planetscale.Client
}

type branchResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Database     types.String `tfsdk:"database"`

	Name                        types.String  `tfsdk:"name"`
	ParentBranch                types.String  `tfsdk:"parent_branch"`
	AccessHostUrl               types.String  `tfsdk:"access_host_url"`
	Actor                       types.Object  `tfsdk:"actor"`
	ClusterRateName             types.String  `tfsdk:"cluster_rate_name"`
	CreatedAt                   types.String  `tfsdk:"created_at"`
	HtmlUrl                     types.String  `tfsdk:"html_url"`
	Id                          types.String  `tfsdk:"id"`
	InitialRestoreId            types.String  `tfsdk:"initial_restore_id"`
	MysqlAddress                types.String  `tfsdk:"mysql_address"`
	MysqlEdgeAddress            types.String  `tfsdk:"mysql_edge_address"`
	Region                      types.Object  `tfsdk:"region"`
	Production                  types.Bool    `tfsdk:"production"`
	Ready                       types.Bool    `tfsdk:"ready"`
	RestoreChecklistCompletedAt types.String  `tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          types.Object  `tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         types.String  `tfsdk:"schema_last_updated_at"`
	ShardCount                  types.Float64 `tfsdk:"shard_count"`
	Sharded                     types.Bool    `tfsdk:"sharded"`
	UpdatedAt                   types.String  `tfsdk:"updated_at"`
}

func (mdl *branchResourceModel) fromClient(ctx context.Context, branch *planetscale.Branch) (diags diag.Diagnostics) {
	if branch == nil {
		return diags
	}
	var (
		actorDiags  diag.Diagnostics
		regionDiags diag.Diagnostics
		rfbDiags    diag.Diagnostics
	)
	mdl.Actor, actorDiags = types.ObjectValueFrom(ctx, actorResourceAttrTypes, branch.Actor)
	mdl.Region, regionDiags = types.ObjectValueFrom(ctx, regionResourceAttrTypes, branch.Region)
	mdl.RestoredFromBranch, rfbDiags = types.ObjectValueFrom(ctx, restoredFromBranchResourceAttrTypes, branch.RestoredFromBranch)

	diags.Append(actorDiags...)
	diags.Append(regionDiags...)
	diags.Append(rfbDiags...)

	mdl.Name = types.StringValue(branch.Name)
	mdl.ParentBranch = types.StringPointerValue(branch.ParentBranch)
	mdl.AccessHostUrl = types.StringPointerValue(branch.AccessHostUrl)
	mdl.ClusterRateName = types.StringValue(branch.ClusterRateName)
	mdl.CreatedAt = types.StringValue(branch.CreatedAt)
	mdl.HtmlUrl = types.StringValue(branch.HtmlUrl)
	mdl.Id = types.StringValue(branch.Id)
	mdl.InitialRestoreId = types.StringPointerValue(branch.InitialRestoreId)
	mdl.MysqlAddress = types.StringValue(branch.MysqlAddress)
	mdl.MysqlEdgeAddress = types.StringValue(branch.MysqlEdgeAddress)
	mdl.Production = types.BoolValue(branch.Production)
	mdl.Ready = types.BoolValue(branch.Ready)
	mdl.RestoreChecklistCompletedAt = types.StringPointerValue(branch.RestoreChecklistCompletedAt)
	mdl.SchemaLastUpdatedAt = types.StringValue(branch.SchemaLastUpdatedAt)
	mdl.ShardCount = types.Float64PointerValue(branch.ShardCount)
	mdl.Sharded = types.BoolValue(branch.Sharded)
	mdl.UpdatedAt = types.StringValue(branch.UpdatedAt)
	return diags
}

func (r *branchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

func (r *branchResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A PlanetScale branch",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"database": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"parent_branch": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},

			// updatable
			"production": schema.BoolAttribute{Computed: true, Optional: true},

			// read only
			"id":              schema.StringAttribute{Computed: true},
			"access_host_url": schema.StringAttribute{Computed: true},
			"actor": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: actorResourceSchemaAttribute,
			},
			"cluster_rate_name":  schema.StringAttribute{Computed: true},
			"created_at":         schema.StringAttribute{Computed: true},
			"html_url":           schema.StringAttribute{Computed: true},
			"initial_restore_id": schema.StringAttribute{Computed: true},
			"mysql_address":      schema.StringAttribute{Computed: true},
			"mysql_edge_address": schema.StringAttribute{Computed: true},
			"region": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: regionResourceSchemaAttribute,
			},
			"ready":                          schema.BoolAttribute{Computed: true},
			"restore_checklist_completed_at": schema.StringAttribute{Computed: true},
			"restored_from_branch": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: restoredFromBranchSchemaAttribute,
			},
			"schema_last_updated_at": schema.StringAttribute{Computed: true},
			"shard_count":            schema.Float64Attribute{Computed: true},
			"sharded":                schema.BoolAttribute{Computed: true},
			"updated_at":             schema.StringAttribute{Computed: true},
		},
	}
}

func (r *branchResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *branchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *branchResourceModel
	tflog.Info(ctx, "getting current branch resource from plan")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgName := data.Organization.ValueString()
	dbName := data.Database.ValueString()

	parentBranch := stringValueIfKnown(data.ParentBranch)
	if parentBranch == nil {
		resp.Diagnostics.AddAttributeError(path.Root("parent_branch"), "Missing parent branch", "All newly created branches require a parent branch.")
		return
	}

	createReq := planetscale.CreateBranchReq{
		Name:         data.Name.ValueString(),
		ParentBranch: *parentBranch,
	}
	if !(data.RestoredFromBranch.IsNull() || data.RestoredFromBranch.IsUnknown()) {
		var rfb restoredFromBranchResource
		resp.Diagnostics.Append(data.RestoredFromBranch.As(ctx, &rfb, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		backupID := rfb.Id.String()
		createReq.BackupId = &backupID
	}
	res, err := r.client.CreateBranch(ctx, orgName, dbName, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create branch, got error: %s", err))
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to create branchs", "no data")
		return
	}

	resp.Diagnostics.Append(data.fromClient(ctx, &res.Branch)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *branchResourceModel

	tflog.Info(ctx, "getting current branch resource from state")
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization.ValueString()
	database := data.Database.ValueString()
	name := data.Name.ValueString()

	res, err := r.client.GetBranch(ctx, org, database, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read branch, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(data.fromClient(ctx, &res.Branch)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		old  *branchResourceModel
		data *branchResourceModel
	)
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization.ValueString()
	database := data.Database.ValueString()
	name := data.Name.ValueString()

	productionWasChanged := false
	isProduction := boolIfDifferent(old.Production, data.Production, &productionWasChanged)
	var branch planetscale.Branch
	if productionWasChanged {
		if *isProduction {
			res, err := r.client.PromoteBranch(ctx, org, database, name)
			if err != nil {
				resp.Diagnostics.AddAttributeError(path.Root("production"), "Failed to promote branch", "Unable to promote branch to production: "+err.Error())
			}
			branch = res.Branch
		} else {
			res, err := r.client.DemoteBranch(ctx, org, database, name)
			if err != nil {
				resp.Diagnostics.AddAttributeError(path.Root("production"), "Failed to demote branch", "Unable to demote branch from production: "+err.Error())
			}
			branch = res.Branch
		}
	}
	resp.Diagnostics.Append(data.fromClient(ctx, &branch)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *branchResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	org := data.Organization.ValueString()
	database := data.Database.ValueString()
	name := data.Name.ValueString()

	res, err := r.client.DeleteBranch(ctx, org, database, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete branch, got error: %s", err))
		return
	}
	_ = res
}

func (r *branchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization,database,name. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[2])...)
}
