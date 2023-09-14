// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

type branchApiActorResource struct {
	AvatarUrl   string `tfsdk:"avatar_url"`
	DisplayName string `tfsdk:"display_name"`
	Id          string `tfsdk:"id"`
}

var apiActorResourceAttrTypes = map[string]attr.Type{
	"avatar_url":   types.StringType,
	"display_name": types.StringType,
	"id":           types.StringType,
}

type branchRegionResource struct {
	DisplayName       string   `tfsdk:"display_name"`
	Enabled           bool     `tfsdk:"enabled"`
	Id                string   `tfsdk:"id"`
	Location          string   `tfsdk:"location"`
	Provider          string   `tfsdk:"provider"`
	PublicIpAddresses []string `tfsdk:"public_ip_addresses"`
	Slug              string   `tfsdk:"slug"`
}

var planetscaleRegionResourceAttrTypes = map[string]attr.Type{
	"display_name":        types.StringType,
	"enabled":             types.BoolType,
	"id":                  types.StringType,
	"location":            types.StringType,
	"provider":            types.StringType,
	"public_ip_addresses": types.ListType{ElemType: types.StringType},
	"slug":                types.StringType,
}

type branchRestoredFromBranchResource struct {
	CreatedAt string `tfsdk:"created_at"`
	DeletedAt string `tfsdk:"deleted_at"`
	Id        string `tfsdk:"id"`
	Name      string `tfsdk:"name"`
	UpdatedAt string `tfsdk:"updated_at"`
}

var restoredFromBranchResourceAttrTypes = map[string]attr.Type{
	"created_at": types.StringType,
	"deleted_at": types.StringType,
	"id":         types.StringType,
	"name":       types.StringType,
	"updated_at": types.StringType,
}

type branchResourceModel struct {
	Organization string       `tfsdk:"organization"`
	Database     string       `tfsdk:"database"`
	Name         string       `tfsdk:"name"`
	ParentBranch types.String `tfsdk:"parent_branch"`

	AccessHostUrl               types.String  `tfsdk:"access_host_url"`
	ApiActor                    types.Object  `tfsdk:"api_actor"`
	ClusterRateName             types.String  `tfsdk:"cluster_rate_name"`
	CreatedAt                   types.String  `tfsdk:"created_at"`
	HtmlUrl                     types.String  `tfsdk:"html_url"`
	Id                          types.String  `tfsdk:"id"`
	InitialRestoreId            types.String  `tfsdk:"initial_restore_id"`
	MysqlAddress                types.String  `tfsdk:"mysql_address"`
	MysqlEdgeAddress            types.String  `tfsdk:"mysql_edge_address"`
	PlanetscaleRegion           types.Object  `tfsdk:"region"`
	Production                  types.Bool    `tfsdk:"production"`
	Ready                       types.Bool    `tfsdk:"ready"`
	RestoreChecklistCompletedAt types.String  `tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          types.Object  `tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         types.String  `tfsdk:"schema_last_updated_at"`
	ShardCount                  types.Float64 `tfsdk:"shard_count"`
	Sharded                     types.Bool    `tfsdk:"sharded"`
	UpdatedAt                   types.String  `tfsdk:"updated_at"`
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
			"api_actor": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"avatar_url":   schema.StringAttribute{Computed: true},
					"display_name": schema.StringAttribute{Computed: true},
					"id":           schema.StringAttribute{Computed: true},
				},
			},
			"cluster_rate_name":  schema.StringAttribute{Computed: true},
			"created_at":         schema.StringAttribute{Computed: true},
			"html_url":           schema.StringAttribute{Computed: true},
			"initial_restore_id": schema.StringAttribute{Computed: true},
			"mysql_address":      schema.StringAttribute{Computed: true},
			"mysql_edge_address": schema.StringAttribute{Computed: true},
			"region": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"display_name":        schema.StringAttribute{Computed: true},
					"enabled":             schema.BoolAttribute{Computed: true},
					"id":                  schema.StringAttribute{Computed: true},
					"location":            schema.StringAttribute{Computed: true},
					"provider":            schema.StringAttribute{Computed: true},
					"public_ip_addresses": schema.ListAttribute{Computed: true, ElementType: types.StringType},
					"slug":                schema.StringAttribute{Computed: true},
				},
			},
			"ready":                          schema.BoolAttribute{Computed: true},
			"restore_checklist_completed_at": schema.StringAttribute{Computed: true},
			"restored_from_branch": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{Computed: true},
					"deleted_at": schema.StringAttribute{Computed: true},
					"id":         schema.StringAttribute{Computed: true},
					"name":       schema.StringAttribute{Computed: true},
					"updated_at": schema.StringAttribute{Computed: true},
				},
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

	orgName := data.Organization
	dbName := data.Database

	parentBranch := stringValueIfKnown(data.ParentBranch)
	if parentBranch == nil {
		resp.Diagnostics.AddAttributeError(path.Root("parent_branch"), "Missing parent branch", "All newly created branches require a parent branch.")
		return
	}

	createReq := planetscale.CreateBranchReq{
		Name:         data.Name,
		ParentBranch: *parentBranch,
	}
	if !(data.RestoredFromBranch.IsNull() || data.RestoredFromBranch.IsUnknown()) {
		var rfb branchRestoredFromBranchResource
		resp.Diagnostics.Append(data.RestoredFromBranch.As(ctx, &rfb, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.BackupId = &rfb.Id
	}
	res201, err := r.client.CreateBranch(ctx, orgName, dbName, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create branch, got error: %s", err))
		return
	}
	if res201 == nil {
		resp.Diagnostics.AddError("Unable to create branchs", "no data")
		return
	}

	data.Name = res201.Name
	data.ParentBranch = types.StringValue(res201.ParentBranch)
	data.AccessHostUrl = types.StringPointerValue(res201.AccessHostUrl)
	data.ClusterRateName = types.StringValue(res201.ClusterRateName)
	data.CreatedAt = types.StringValue(res201.CreatedAt)
	data.HtmlUrl = types.StringValue(res201.HtmlUrl)
	data.Id = types.StringValue(res201.Id)
	data.InitialRestoreId = types.StringPointerValue(res201.InitialRestoreId)
	data.MysqlAddress = types.StringValue(res201.MysqlAddress)
	data.MysqlEdgeAddress = types.StringValue(res201.MysqlEdgeAddress)
	data.Production = types.BoolValue(res201.Production)
	data.Ready = types.BoolValue(res201.Ready)
	data.RestoreChecklistCompletedAt = types.StringPointerValue(res201.RestoreChecklistCompletedAt)
	data.SchemaLastUpdatedAt = types.StringValue(res201.SchemaLastUpdatedAt)
	data.ShardCount = types.Float64PointerValue(res201.ShardCount)
	data.Sharded = types.BoolValue(res201.Sharded)
	data.UpdatedAt = types.StringValue(res201.UpdatedAt)

	var diErr diag.Diagnostics
	if res201.ApiActor != nil {
		data.ApiActor, diErr = types.ObjectValueFrom(ctx, apiActorResourceAttrTypes, res201.ApiActor)
		if diErr.HasError() {
			resp.Diagnostics.Append(diErr.Errors()...)
			return
		}
	}
	if res201.PlanetscaleRegion != nil {
		data.PlanetscaleRegion, diErr = types.ObjectValueFrom(ctx, planetscaleRegionResourceAttrTypes, res201.PlanetscaleRegion)
		if diErr.HasError() {
			resp.Diagnostics.Append(diErr.Errors()...)
			return
		}
	}
	if res201.RestoredFromBranch != nil {
		data.RestoredFromBranch, diErr = types.ObjectValueFrom(ctx, restoredFromBranchResourceAttrTypes, res201.RestoredFromBranch)
		if diErr.HasError() {
			resp.Diagnostics.Append(diErr.Errors()...)
			return
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a branch resource")

	// Save data into Terraform state
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

	org := data.Organization
	database := data.Database
	name := data.Name

	res200, err := r.client.GetBranch(ctx, org, database, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read branch, got error: %s", err))
		return
	}

	data.Name = res200.Name
	// TODO: fix to a pointer if the openapi spec gets fixed to mark parent branch as nullable
	if res200.ParentBranch == "" {
		data.ParentBranch = types.StringNull()
	} else {
		data.ParentBranch = types.StringValue(res200.ParentBranch)
	}
	data.AccessHostUrl = types.StringPointerValue(res200.AccessHostUrl)
	data.ClusterRateName = types.StringValue(res200.ClusterRateName)
	data.CreatedAt = types.StringValue(res200.CreatedAt)
	data.HtmlUrl = types.StringValue(res200.HtmlUrl)
	data.Id = types.StringValue(res200.Id)
	data.InitialRestoreId = types.StringPointerValue(res200.InitialRestoreId)
	data.MysqlAddress = types.StringValue(res200.MysqlAddress)
	data.MysqlEdgeAddress = types.StringValue(res200.MysqlEdgeAddress)
	data.Production = types.BoolValue(res200.Production)
	data.Ready = types.BoolValue(res200.Ready)
	data.RestoreChecklistCompletedAt = types.StringPointerValue(res200.RestoreChecklistCompletedAt)
	data.SchemaLastUpdatedAt = types.StringValue(res200.SchemaLastUpdatedAt)
	data.ShardCount = types.Float64PointerValue(res200.ShardCount)
	data.Sharded = types.BoolValue(res200.Sharded)
	data.UpdatedAt = types.StringValue(res200.UpdatedAt)

	var diErr diag.Diagnostics
	if res200.ApiActor != nil {
		data.ApiActor, diErr = types.ObjectValueFrom(ctx, apiActorResourceAttrTypes, res200.ApiActor)
		if diErr.HasError() {
			resp.Diagnostics.Append(diErr.Errors()...)
			return
		}
	}
	if res200.PlanetscaleRegion != nil {
		data.PlanetscaleRegion, diErr = types.ObjectValueFrom(ctx, planetscaleRegionResourceAttrTypes, res200.PlanetscaleRegion)
		if diErr.HasError() {
			resp.Diagnostics.Append(diErr.Errors()...)
			return
		}
	}
	if res200.RestoredFromBranch != nil {
		data.RestoredFromBranch, diErr = types.ObjectValueFrom(ctx, restoredFromBranchResourceAttrTypes, res200.RestoredFromBranch)
		if diErr.HasError() {
			resp.Diagnostics.Append(diErr.Errors()...)
			return
		}
	}

	// Save updated data into Terraform state
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

	// todo

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *branchResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	org := data.Organization
	database := data.Database
	name := data.Name

	res204, err := r.client.DeleteBranch(ctx, org, database, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete branch, got error: %s", err))
		return
	}
	_ = res204
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
