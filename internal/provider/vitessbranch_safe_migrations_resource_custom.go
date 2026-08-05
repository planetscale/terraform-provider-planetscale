package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk"
	sdkerrors "github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/errors"
)

var (
	_ resource.Resource                = &VitessBranchSafeMigrationsResource{}
	_ resource.ResourceWithImportState = &VitessBranchSafeMigrationsResource{}
)

// NewVitessBranchSafeMigrationsResource returns a resource that manages the
// safe migrations setting of an existing Vitess branch.
func NewVitessBranchSafeMigrationsResource() resource.Resource {
	return &VitessBranchSafeMigrationsResource{}
}

// VitessBranchSafeMigrationsResource manages safe migrations independently of
// the generated branch resource because the API exposes dedicated enable and
// disable operations.
type VitessBranchSafeMigrationsResource struct {
	client *sdk.PlanetScale
}

type vitessBranchSafeMigrationsResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Organization   types.String `tfsdk:"organization"`
	Database       types.String `tfsdk:"database"`
	Branch         types.String `tfsdk:"branch"`
	SafeMigrations types.Bool   `tfsdk:"safe_migrations"`
}

func (r *VitessBranchSafeMigrationsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vitess_branch_safe_migrations"
}

func (r *VitessBranchSafeMigrationsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages safe migrations for an existing PlanetScale Vitess branch.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The public ID of the branch.",
			},
			"organization": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The organization the branch belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The database the branch belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"branch": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The public ID of the Vitess branch. Prefer `planetscale_vitess_branch.<name>.id` so branch renames do not replace this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"safe_migrations": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether safe migrations (DDL protection) are enabled.",
			},
		},
	}
}

func (r *VitessBranchSafeMigrationsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.PlanetScale)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sdk.PlanetScale, got %T.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *VitessBranchSafeMigrationsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data vitessBranchSafeMigrationsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.apply(ctx, &data, &resp.Diagnostics)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *VitessBranchSafeMigrationsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data vitessBranchSafeMigrationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.DatabaseBranches.GetVitessBranchSettings(
		ctx,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
	)
	if err != nil {
		if isBranchSettingNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read safe migrations", err.Error())
		return
	}

	data.ID = types.StringValue(settings.ID)
	data.SafeMigrations = types.BoolValue(settings.SafeMigrations)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VitessBranchSafeMigrationsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data vitessBranchSafeMigrationsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.apply(ctx, &data, &resp.Diagnostics)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *VitessBranchSafeMigrationsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vitessBranchSafeMigrationsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DatabaseBranches.SetVitessBranchSafeMigrations(
		ctx,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
		false,
	)
	if err != nil && !isBranchSettingNotFound(err) {
		resp.Diagnostics.AddError("Unable to disable safe migrations", err.Error())
	}
}

func (r *VitessBranchSafeMigrationsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var identity struct {
		Organization string `json:"organization"`
		Database     string `json:"database"`
		Branch       string `json:"branch"`
	}
	if err := json.Unmarshal([]byte(req.ID), &identity); err != nil ||
		identity.Organization == "" || identity.Database == "" || identity.Branch == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			`Expected a JSON object with non-empty "organization", "database", and "branch" properties.`,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), identity.Organization)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), identity.Database)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("branch"), identity.Branch)...)
}

func (r *VitessBranchSafeMigrationsResource) apply(ctx context.Context, data *vitessBranchSafeMigrationsResourceModel, diagnostics interface {
	AddError(string, string)
	HasError() bool
}) {
	settings, err := r.client.DatabaseBranches.SetVitessBranchSafeMigrations(
		ctx,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
		data.SafeMigrations.ValueBool(),
	)
	if err != nil {
		diagnostics.AddError("Unable to update safe migrations", err.Error())
		return
	}

	data.ID = types.StringValue(settings.ID)
	data.SafeMigrations = types.BoolValue(settings.SafeMigrations)
}

func isBranchSettingNotFound(err error) bool {
	var apiError *sdkerrors.APIError
	return errors.As(err, &apiError) && apiError.StatusCode == 404
}
