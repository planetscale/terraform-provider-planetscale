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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                 = &branchResource{}
	_ resource.ResourceWithImportState  = &branchResource{}
	_ resource.ResourceWithUpgradeState = &branchResource{}
)

func newBranchResource() resource.Resource {
	return &branchResource{}
}

// branchResource defines the resource implementation.
type branchResource struct {
	client *planetscale.Client
}

// branchResourceModelV* and branchSchemaV* are defined in branch_resource_migration.go

func branchResourceFromClient(ctx context.Context, branch *planetscale.Branch, organization, database types.String, diags diag.Diagnostics) *branchResourceModelV1 {
	if branch == nil {
		return nil
	}
	actor, diags := types.ObjectValueFrom(ctx, actorResourceAttrTypes, branch.Actor)
	diags.Append(diags...)
	region, diags := types.ObjectValueFrom(ctx, regionResourceAttrTypes, branch.Region)
	diags.Append(diags...)
	restoredFromBranch, diags := types.ObjectValueFrom(ctx, restoredFromBranchResourceAttrTypes, branch.RestoredFromBranch)
	diags.Append(diags...)
	return &branchResourceModelV1{
		Organization: organization,
		Database:     database,

		Actor:                       actor,
		Region:                      region,
		RestoredFromBranch:          restoredFromBranch,
		Name:                        types.StringValue(branch.Name),
		ParentBranch:                types.StringPointerValue(branch.ParentBranch),
		CreatedAt:                   types.StringValue(branch.CreatedAt),
		HtmlUrl:                     types.StringValue(branch.HtmlUrl),
		Id:                          types.StringValue(branch.Id),
		MysqlAddress:                types.StringValue(branch.MysqlAddress),
		MysqlEdgeAddress:            types.StringValue(branch.MysqlEdgeAddress),
		Production:                  types.BoolValue(branch.Production),
		Ready:                       types.BoolValue(branch.Ready),
		RestoreChecklistCompletedAt: types.StringPointerValue(branch.RestoreChecklistCompletedAt),
		SchemaLastUpdatedAt:         types.StringValue(branch.SchemaLastUpdatedAt),
		ShardCount:                  types.Float64PointerValue(branch.ShardCount),
		Sharded:                     types.BoolValue(branch.Sharded),
		UpdatedAt:                   types.StringValue(branch.UpdatedAt),
	}
}

func (r *branchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

func (r *branchResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Schema definition is sourced from the migration file
	resp.Schema = *branchSchemaV1()
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
	var data *branchResourceModelV1
	tflog.Info(ctx, "getting current branch resource from plan")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization
	database := data.Database
	name := data.Name
	seedData := data.SeedData

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if database.IsNull() || database.IsUnknown() || database.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("database"), "database is required", "a database must be provided and cannot be empty")
		return
	}
	if name.IsNull() || name.IsUnknown() || name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "name is required", "a name must be provided and cannot be empty")
		return
	}

	parentBranch := stringValueIfKnown(data.ParentBranch)
	if parentBranch == nil {
		resp.Diagnostics.AddAttributeError(path.Root("parent_branch"), "Missing parent branch", "All newly created branches require a parent branch.")
		return
	}

	createReq := planetscale.CreateBranchReq{
		Name:         name.ValueString(),
		ParentBranch: *parentBranch,
		SeedData:     seedData.ValueString(),
	}
	if !data.RestoredFromBranch.IsNull() && !data.RestoredFromBranch.IsUnknown() {
		var rfb restoredFromBranchResource
		resp.Diagnostics.Append(data.RestoredFromBranch.As(ctx, &rfb, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		backupID := rfb.Id.String()
		createReq.BackupId = &backupID
	}
	res, err := r.client.CreateBranch(ctx, org.ValueString(), database.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create branch, got error: %s", err))
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to create branchs", "no data")
		return
	}

	// wait for branch to enter ready state
	createState := &retry.StateChangeConf{
		Delay:      5 * time.Second, // initial delay before the first check
		Timeout:    10 * time.Minute,
		MinTimeout: 5 * time.Second,

		Pending: []string{"not-ready"},
		Target:  []string{"ready"},

		Refresh: func() (interface{}, string, error) {
			res, err := r.client.GetBranch(ctx, org.ValueString(), database.ValueString(), name.ValueString())
			if err != nil {
				return nil, "", err
			}
			if res.Branch.Ready {
				return res.Branch, "ready", nil
			}
			return res.Branch, "not-ready", nil
		},
	}
	branchRaw, err := createState.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create database",
			fmt.Sprintf("Branch %s never became ready; got error: %s", name.ValueString(), err),
		)
		return
	}

	branch, ok := branchRaw.(planetscale.Branch)
	if !ok {
		resp.Diagnostics.AddError("Unable to create branch", "no data")
		return
	}

	// After branch is ready, check if we need to promote it to a production branch
	if !data.Production.IsNull() && data.Production.ValueBool() {
		res, err := r.client.PromoteBranch(ctx, org.ValueString(), database.ValueString(), name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to promote branch during creation",
				fmt.Sprintf("Branch %s could not be promoted to production: %s", name.ValueString(), err),
			)
			return
		}
		branch = res.Branch
	}

	data = branchResourceFromClient(ctx, &branch, data.Organization, data.Database, resp.Diagnostics)
	data.SeedData = seedData

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *branchResourceModelV1

	tflog.Info(ctx, "getting current branch resource from state")
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization
	database := data.Database
	name := data.Name

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if database.IsNull() || database.IsUnknown() || database.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("database"), "database is required", "a database must be provided and cannot be empty")
		return
	}
	if name.IsNull() || name.IsUnknown() || name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "name is required", "a name must be provided and cannot be empty")
		return
	}

	res, err := r.client.GetBranch(ctx, org.ValueString(), database.ValueString(), name.ValueString())
	if err != nil {
		if notFoundErr, ok := err.(*planetscale.GetBranchRes404); ok {
			tflog.Warn(ctx, fmt.Sprintf("Branch not found, removing from state: %s", notFoundErr.Message))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read branch, got error: %s", err))
		return
	}

	data = branchResourceFromClient(ctx, &res.Branch, data.Organization, data.Database, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		old  *branchResourceModelV1
		data *branchResourceModelV1
	)
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization
	database := data.Database
	name := data.Name

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if database.IsNull() || database.IsUnknown() || database.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("database"), "database is required", "a database must be provided and cannot be empty")
		return
	}
	if name.IsNull() || name.IsUnknown() || name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "name is required", "a name must be provided and cannot be empty")
		return
	}

	productionWasChanged := false
	isProduction := boolIfDifferent(old.Production, data.Production, &productionWasChanged)
	var branch planetscale.Branch
	if productionWasChanged {
		if *isProduction {
			res, err := r.client.PromoteBranch(ctx, org.ValueString(), database.ValueString(), name.ValueString())
			if err != nil {
				resp.Diagnostics.AddAttributeError(path.Root("production"), "Failed to promote branch", "Unable to promote branch to production: "+err.Error())
				if resp.Diagnostics.HasError() {
					return
				}
			}
			branch = res.Branch
		} else {
			res, err := r.client.DemoteBranch(ctx, org.ValueString(), database.ValueString(), name.ValueString())
			if err != nil {
				resp.Diagnostics.AddAttributeError(path.Root("production"), "Failed to demote branch", "Unable to demote branch from production: "+err.Error())
				if resp.Diagnostics.HasError() {
					return
				}
			}
			branch = res.Branch
		}
	}
	data = branchResourceFromClient(ctx, &branch, data.Organization, data.Database, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *branchResourceModelV1

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	org := data.Organization
	database := data.Database
	name := data.Name

	if org.IsNull() || org.IsUnknown() || org.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("organization"), "organization is required", "an organization must be provided and cannot be empty")
		return
	}
	if database.IsNull() || database.IsUnknown() || database.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("database"), "database is required", "a database must be provided and cannot be empty")
		return
	}
	if name.IsNull() || name.IsUnknown() || name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("name"), "name is required", "a name must be provided and cannot be empty")
		return
	}

	res, err := r.client.DeleteBranch(ctx, org.ValueString(), database.ValueString(), name.ValueString())
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
