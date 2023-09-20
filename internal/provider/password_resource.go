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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &passwordResource{}
var _ resource.ResourceWithImportState = &passwordResource{}

func newPasswordResource() resource.Resource {
	return &passwordResource{}
}

// passwordResource defines the resource implementation.
type passwordResource struct {
	client *planetscale.Client
}

type passwordResourceModel struct {
	Organization   types.String  `tfsdk:"organization"`
	Database       types.String  `tfsdk:"database"`
	Branch         types.String  `tfsdk:"branch"`
	Id             types.String  `tfsdk:"id"`
	Name           types.String  `tfsdk:"name"`
	AccessHostUrl  types.String  `tfsdk:"access_host_url"`
	Actor          types.Object  `tfsdk:"actor"`
	CreatedAt      types.String  `tfsdk:"created_at"`
	DatabaseBranch types.Object  `tfsdk:"database_branch"`
	DeletedAt      types.String  `tfsdk:"deleted_at"`
	ExpiresAt      types.String  `tfsdk:"expires_at"`
	Integrations   types.List    `tfsdk:"integrations"`
	Region         types.Object  `tfsdk:"region"`
	Renewable      types.Bool    `tfsdk:"renewable"`
	Role           types.String  `tfsdk:"role"`
	TtlSeconds     types.Float64 `tfsdk:"ttl_seconds"`
	Username       types.String  `tfsdk:"username"`

	PlainText types.String `tfsdk:"plaintext"`
}

func passwordResourceFromClient(ctx context.Context, password *planetscale.Password, organization, database, branch types.String, diags diag.Diagnostics) *passwordResourceModel {
	if password == nil {
		return nil
	}
	actor, diags := types.ObjectValueFrom(ctx, actorResourceAttrTypes, password.Actor)
	diags.Append(diags...)
	databaseBranch, diags := types.ObjectValueFrom(ctx, databaseBranchResourceAttrTypes, password.DatabaseBranch)
	diags.Append(diags...)
	region, diags := types.ObjectValueFrom(ctx, regionResourceAttrTypes, password.Region)
	diags.Append(diags...)
	return &passwordResourceModel{
		Organization: organization,
		Database:     database,
		Branch:       branch,

		Name:           types.StringValue(password.Name),
		AccessHostUrl:  types.StringValue(password.AccessHostUrl),
		Actor:          actor,
		CreatedAt:      types.StringValue(password.CreatedAt),
		DatabaseBranch: databaseBranch,
		DeletedAt:      types.StringPointerValue(password.DeletedAt),
		ExpiresAt:      types.StringPointerValue(password.ExpiresAt),
		Id:             types.StringValue(password.Id),
		Integrations:   stringsToListValue(password.Integrations, diags),
		Region:         region,
		Renewable:      types.BoolValue(password.Renewable),
		Role:           types.StringValue(password.Role),
		TtlSeconds:     types.Float64Value(password.TtlSeconds),
		Username:       types.StringPointerValue(password.Username),
	}
}

func passwordWithPlaintextResourceFromClient(ctx context.Context, password *planetscale.PasswordWithPlaintext, organization, database, branch types.String, diags diag.Diagnostics) *passwordResourceModel {
	if password == nil {
		return nil
	}
	actor, diags := types.ObjectValueFrom(ctx, actorResourceAttrTypes, password.Actor)
	diags.Append(diags...)
	databaseBranch, diags := types.ObjectValueFrom(ctx, databaseBranchResourceAttrTypes, password.DatabaseBranch)
	diags.Append(diags...)
	region, diags := types.ObjectValueFrom(ctx, regionResourceAttrTypes, password.Region)
	diags.Append(diags...)
	return &passwordResourceModel{
		Organization: organization,
		Database:     database,
		Branch:       branch,

		Name:           types.StringValue(password.Name),
		AccessHostUrl:  types.StringValue(password.AccessHostUrl),
		Actor:          actor,
		CreatedAt:      types.StringValue(password.CreatedAt),
		DatabaseBranch: databaseBranch,
		DeletedAt:      types.StringPointerValue(password.DeletedAt),
		ExpiresAt:      types.StringPointerValue(password.ExpiresAt),
		Id:             types.StringValue(password.Id),
		Integrations:   stringsToListValue(password.Integrations, diags),
		Region:         region,
		Renewable:      types.BoolValue(password.Renewable),
		Role:           types.StringValue(password.Role),
		TtlSeconds:     types.Float64Value(password.TtlSeconds),
		Username:       types.StringPointerValue(password.Username),

		PlainText: types.StringValue(password.PlainText),
	}
}

func (r *passwordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

func (r *passwordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A PlanetScale password",
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

			"role": schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"ttl_seconds": schema.Float64Attribute{Optional: true, PlanModifiers: []planmodifier.Float64{
				float64planmodifier.RequiresReplace(),
			}},
			// updatable
			"name": schema.StringAttribute{Optional: true},

			// read-only
			"id":              schema.StringAttribute{Computed: true},
			"access_host_url": schema.StringAttribute{Computed: true},
			"actor":           schema.ObjectAttribute{Computed: true},
			"created_at":      schema.StringAttribute{Computed: true},
			"database_branch": schema.ObjectAttribute{Computed: true},
			"deleted_at":      schema.StringAttribute{Computed: true},
			"expires_at":      schema.StringAttribute{Computed: true},
			"integrations":    schema.ListAttribute{Computed: true},
			"region":          schema.ObjectAttribute{Computed: true},
			"renewable":       schema.BoolAttribute{Computed: true},
			"username":        schema.StringAttribute{Computed: true},
		},
	}
}

func (r *passwordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *passwordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *passwordResourceModel
	tflog.Info(ctx, "getting current password resource from plan")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org := data.Organization
	database := data.Database
	branch := data.Branch

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

	name := data.Name
	role := data.Role
	ttl := data.TtlSeconds

	createReq := planetscale.CreatePasswordReq{
		Name: name.ValueStringPointer(),
		Role: role.ValueStringPointer(),
		Ttl:  ttl.ValueFloat64Pointer(),
	}
	res, err := r.client.CreatePassword(ctx, org.ValueString(), database.ValueString(), branch.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create password, got error: %s", err))
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to create passwords", "no data")
		return
	}
	data = passwordWithPlaintextResourceFromClient(ctx, &res.PasswordWithPlaintext, data.Organization, data.Database, data.Branch, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *passwordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *passwordResourceModel

	tflog.Info(ctx, "getting current password resource from state")
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

	res, err := r.client.GetPassword(ctx,
		org.ValueString(),
		database.ValueString(),
		branch.ValueString(),
		id.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read password, got error: %s", err))
		return
	}

	data = passwordResourceFromClient(
		ctx,
		&res.Password,
		data.Organization,
		data.Database,
		data.Branch,
		resp.Diagnostics,
	)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *passwordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		old  *passwordResourceModel
		data *passwordResourceModel
	)
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	changedUpdatableSettings := false
	name := stringIfDifferent(old.Name, data.Name, &changedUpdatableSettings)

	if changedUpdatableSettings && name != nil {
		updateReq := planetscale.UpdatePasswordReq{
			Name: *name,
		}
		res, err := r.client.UpdatePassword(
			ctx,
			org.ValueString(),
			database.ValueString(),
			branch.ValueString(),
			id.ValueString(),
			updateReq,
		)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update password settings, got error: %s", err))
			return
		}
		data = passwordResourceFromClient(ctx, &res.Password, data.Organization, data.Database, data.Branch, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *passwordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *passwordResourceModel

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

	res, err := r.client.DeletePassword(ctx, org.ValueString(), database.ValueString(), branch.ValueString(), id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete password, got error: %s", err))
		return
	}
	_ = res
}

func (r *passwordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
