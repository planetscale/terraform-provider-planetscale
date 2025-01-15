package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ resource.Resource                = &branchSafeMigrationsResource{}
	_ resource.ResourceWithImportState = &branchSafeMigrationsResource{}
)

func newBranchSafeMigrationsResource() resource.Resource {
	return &branchSafeMigrationsResource{}
}

type branchSafeMigrationsResource struct {
	client *planetscale.Client
}

type branchSafeMigrationsResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Database     types.String `tfsdk:"database"`
	Branch       types.String `tfsdk:"branch"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Id           types.String `tfsdk:"id"`
}

func (r *branchSafeMigrationsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch_safe_migrations" // planetscale_branch_safe_migrations
}

func (r *branchSafeMigrationsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages safe migrations settings for a PlanetScale branch.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the branch.",
				Computed:    true,
			},
			"organization": schema.StringAttribute{
				Description: "The organization this branch belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"database": schema.StringAttribute{
				Description: "The database this branch belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"branch": schema.StringAttribute{
				Description: "The name of the branch to configure safe migrations on..",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether safe migrations are enabled for this branch.",
				Required:    true,
			},
		},
	}
}

func (r *branchSafeMigrationsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*planetscale.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *planetscale.Client, got: %T.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *branchSafeMigrationsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *branchSafeMigrationsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var branch *planetscale.Branch

	if data.Enabled.ValueBool() {
		res, err := r.client.EnableSafeMigrations(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to enable safe migrations, got error: %s", err))
			return
		}
		branch = &res.Branch
	} else {
		res, err := r.client.DisableSafeMigrations(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to disable safe migrations, got error: %s", err))
			return
		}
		branch = &res.Branch
	}

	data.Id = types.StringValue(branch.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchSafeMigrationsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *branchSafeMigrationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetBranch(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString())
	if err != nil {
		if notFoundErr, ok := err.(*planetscale.GetBranchRes404); ok {
			tflog.Warn(ctx, fmt.Sprintf("Branch not found, removing from state: %s", notFoundErr.Message))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read branch safe migrations, got error: %s", err))
		return
	}

	data.Id = types.StringValue(res.Branch.Id)
	data.Enabled = types.BoolValue(res.Branch.SafeMigrations)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchSafeMigrationsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *branchSafeMigrationsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var branch *planetscale.Branch

	if data.Enabled.ValueBool() {
		res, err := r.client.EnableSafeMigrations(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to enable safe migrations, got error: %s", err))
			return
		}
		branch = &res.Branch
	} else {
		res, err := r.client.DisableSafeMigrations(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to disable safe migrations, got error: %s", err))
			return
		}
		branch = &res.Branch
	}

	data.Id = types.StringValue(branch.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *branchSafeMigrationsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *branchSafeMigrationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DisableSafeMigrations(ctx, data.Organization.ValueString(), data.Database.ValueString(), data.Branch.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to disable safe migrations, got error: %s", err))
		return
	}
}

func (r *branchSafeMigrationsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization,database,branch. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("branch"), idParts[2])...)
}
