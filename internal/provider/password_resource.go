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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &passwordResource{}
	_ resource.ResourceWithImportState = &passwordResource{}
)

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
	Region         types.Object  `tfsdk:"region"`
	Renewable      types.Bool    `tfsdk:"renewable"`
	Role           types.String  `tfsdk:"role"`
	Cidrs          types.List    `tfsdk:"cidrs"`
	TtlSeconds     types.Float64 `tfsdk:"ttl_seconds"`
	Username       types.String  `tfsdk:"username"`

	PlainText types.String `tfsdk:"plaintext"`

	// manually removed from spec because currently buggy
	// Integrations   types.List    `tfsdk:"integrations"`
}

func passwordResourceFromClient(ctx context.Context, password *planetscale.Password, organization, database, branch, plainText types.String, diags diag.Diagnostics) *passwordResourceModel {
	if password == nil {
		return nil
	}
	actor, diags := types.ObjectValueFrom(ctx, actorResourceAttrTypes, password.Actor)
	diags.Append(diags...)
	databaseBranch, diags := types.ObjectValueFrom(ctx, databaseBranchResourceAttrTypes, password.DatabaseBranch)
	diags.Append(diags...)
	region, diags := types.ObjectValueFrom(ctx, regionResourceAttrTypes, password.Region)
	diags.Append(diags...)

	var cidrs types.List
	if password.Cidrs != nil {
		cidrs = stringsToListValue(password.Cidrs, diags)
	} else {
		cidrs = types.ListNull(types.StringType)
	}

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
		Region:         region,
		Renewable:      types.BoolValue(password.Renewable),
		Role:           types.StringValue(password.Role),
		TtlSeconds:     types.Float64Value(password.TtlSeconds),
		Username:       types.StringPointerValue(password.Username),
		Cidrs:          cidrs,
		PlainText:      plainText,

		// manually removed from spec because currently buggy
		// Integrations:   stringsToListValue(password.Integrations, diags),
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

	var cidrs types.List
	if password.Cidrs != nil {
		cidrs = stringsToListValue(password.Cidrs, diags)
	} else {
		cidrs = types.ListNull(types.StringType)
	}

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
		Region:         region,
		Renewable:      types.BoolValue(password.Renewable),
		Role:           types.StringValue(password.Role),
		TtlSeconds:     types.Float64Value(password.TtlSeconds),
		Username:       types.StringPointerValue(password.Username),
		Cidrs:          cidrs,
		PlainText:      types.StringValue(password.PlainText),

		// manually removed from spec because currently buggy
		// Integrations:   stringsToListValue(password.Integrations, diags),
	}
}

func (r *passwordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

func (r *passwordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A PlanetScale database password.",
		MarkdownDescription: "A PlanetScale database password.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description: "The organization this database branch password belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"database": schema.StringAttribute{
				Description: "The database this branch password belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"branch": schema.StringAttribute{
				Description: "The branch this password belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"role": schema.StringAttribute{
				Description: "The role for the password.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"ttl_seconds": schema.Float64Attribute{
				Description: "Time to live (in seconds) for the password. The password will be invalid and unrenewable when TTL has passed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.RequiresReplaceIfConfigured(),
				},
			},
			// updatable
			"name": schema.StringAttribute{
				Description: "The display name for the password.",
				Optional:    true,
			},
			"cidrs": schema.ListAttribute{
				Description: "List of IP addresses or CIDR ranges that can use this password.",
				Optional:    true,
				Computed:    false,
				ElementType: types.StringType,
			},

			// read-only
			"id": schema.StringAttribute{
				Description: "The ID for the password.",
				Computed:    true,
			},
			"actor": schema.SingleNestedAttribute{
				Description: "The actor that created this branch.",
				Computed:    true,
				Attributes:  actorResourceSchemaAttribute,
			},
			"database_branch": schema.SingleNestedAttribute{
				Description: "The branch this password is allowed to access.",
				Computed:    true,
				Attributes:  databaseBranchResourceAttribute,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.SingleNestedAttribute{
				Description: "The region in which this password can be used.",
				Computed:    true,
				Attributes:  regionResourceSchemaAttribute,
			},
			"access_host_url": schema.StringAttribute{
				Description: "The host URL for the password.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When the password was created.",
				Computed:    true,
			},
			"deleted_at": schema.StringAttribute{
				Description: "When the password was deleted.",
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "When the password will expire.",
				Computed:    true,
			},
			"renewable": schema.BoolAttribute{
				Description: "Whether or not the password can be renewed.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username for the password.",
				Computed:    true,
			},

			// read-only, sensitive
			"plaintext": schema.StringAttribute{
				Description: "The plaintext password, only available if the password was created by this provider.",
				Sensitive:   true, Computed: true,
			},

			// manually removed from spec because currently buggy
			// "integrations":    schema.ListAttribute{Computed: true, ElementType: types.StringType},
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
	data = passwordWithPlaintextResourceFromClient(
		ctx,
		&res.PasswordWithPlaintext,
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
	)
	if err != nil {
		if notFoundErr, ok := err.(*planetscale.GetPasswordRes404); ok {
			tflog.Warn(ctx, fmt.Sprintf("Password not found, removing from state: %s", notFoundErr.Message))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read password, got error: %s", err))
		return
	}
	data = passwordResourceFromClient(
		ctx,
		&res.Password,
		data.Organization,
		data.Database,
		data.Branch,
		data.PlainText,
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

	var state *passwordResourceModel
	if changedUpdatableSettings && name != nil {
		updateReq := planetscale.UpdatePasswordReq{
			Name: name,
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
		state = passwordResourceFromClient(
			ctx,
			&res.Password,
			data.Organization,
			data.Database,
			data.Branch,
			data.PlainText,
			resp.Diagnostics,
		)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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
			fmt.Sprintf("Expected import identifier with format: organization,database,branch,id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("branch"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[3])...)
}
