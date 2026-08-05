package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk"
)

var (
	_ resource.Resource                = &VitessBranchVTGateAutoscalingResource{}
	_ resource.ResourceWithImportState = &VitessBranchVTGateAutoscalingResource{}
)

// NewVitessBranchVTGateAutoscalingResource returns a resource that manages the
// VTGate autoscaling configuration of an existing Vitess branch.
func NewVitessBranchVTGateAutoscalingResource() resource.Resource {
	return &VitessBranchVTGateAutoscalingResource{}
}

// VitessBranchVTGateAutoscalingResource models autoscaling as a durable branch
// setting rather than exposing transient resize requests as Terraform objects.
type VitessBranchVTGateAutoscalingResource struct {
	client *sdk.PlanetScale
}

type vitessBranchVTGateAutoscalingResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Organization               types.String `tfsdk:"organization"`
	Database                   types.String `tfsdk:"database"`
	Branch                     types.String `tfsdk:"branch"`
	VTGateAutoscaling          types.Bool   `tfsdk:"vtgate_autoscaling"`
	VTGateSize                 types.String `tfsdk:"vtgate_size"`
	VTGateCount                types.Int64  `tfsdk:"vtgate_count"`
	VTGateMaxCount             types.Int64  `tfsdk:"vtgate_max_count"`
	VTGateTargetCPUUtilization types.Int64  `tfsdk:"vtgate_target_cpu_utilization"`
}

type currentVTGateConfiguration struct {
	branchID             string
	autoscaling          bool
	size                 string
	count                int64
	maxCount             *int64
	targetCPUUtilization *int64
}

func (r *VitessBranchVTGateAutoscalingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vitess_branch_vtgate_autoscaling"
}

func (r *VitessBranchVTGateAutoscalingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages VTGate autoscaling for an existing PlanetScale Vitess branch.",
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
			"vtgate_autoscaling": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether VTGate autoscaling is enabled.",
			},
			"vtgate_size": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The VTGate size, such as `VTG_320`. Autoscaling requires a configurable size of at least VTG_320.",
			},
			"vtgate_count": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The minimum number of VTGates per availability zone.",
			},
			"vtgate_max_count": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The maximum number of VTGates per availability zone. This must be set and greater than `vtgate_count` when enabling autoscaling for the first time.",
			},
			"vtgate_target_cpu_utilization": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The target average CPU utilization percentage. Supported values are 40, 50, 60, and 70.",
			},
		},
	}
}

func (r *VitessBranchVTGateAutoscalingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VitessBranchVTGateAutoscalingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data vitessBranchVTGateAutoscalingResourceModel
	var config vitessBranchVTGateAutoscalingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.apply(ctx, &data, &config); err != nil {
		resp.Diagnostics.AddError("Unable to update VTGate autoscaling", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VitessBranchVTGateAutoscalingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data vitessBranchVTGateAutoscalingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	current, err := r.readCurrent(ctx, &data)
	if err != nil {
		if isBranchSettingNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read VTGate autoscaling", err.Error())
		return
	}

	refreshVTGateAutoscalingModel(&data, current)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VitessBranchVTGateAutoscalingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data vitessBranchVTGateAutoscalingResourceModel
	var config vitessBranchVTGateAutoscalingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.apply(ctx, &data, &config); err != nil {
		resp.Diagnostics.AddError("Unable to update VTGate autoscaling", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VitessBranchVTGateAutoscalingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vitessBranchVTGateAutoscalingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	current, err := r.readCurrent(ctx, &data)
	if err != nil {
		if isBranchSettingNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Unable to read VTGate autoscaling", err.Error())
		return
	}
	if !current.autoscaling {
		return
	}

	disabled := false
	_, err = r.client.DatabaseBranches.UpdateVitessBranchVTGateConfiguration(
		ctx,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
		sdk.UpdateVitessBranchVTGateConfigurationRequest{
			VTGateAutoscaling: &disabled,
		},
	)
	if err != nil && !isBranchSettingNotFound(err) {
		resp.Diagnostics.AddError("Unable to disable VTGate autoscaling", err.Error())
	}
}

func (r *VitessBranchVTGateAutoscalingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *VitessBranchVTGateAutoscalingResource) apply(
	ctx context.Context,
	data *vitessBranchVTGateAutoscalingResourceModel,
	config *vitessBranchVTGateAutoscalingResourceModel,
) error {
	current, err := r.readCurrent(ctx, data)
	if err != nil {
		return err
	}

	desiredAutoscaling := config.VTGateAutoscaling.ValueBool()
	update := sdk.UpdateVitessBranchVTGateConfigurationRequest{
		VTGateAutoscaling: &desiredAutoscaling,
	}
	changed := current.autoscaling != desiredAutoscaling

	if !config.VTGateSize.IsNull() && !config.VTGateSize.IsUnknown() {
		value := config.VTGateSize.ValueString()
		update.VTGateSize = &value
		changed = changed || current.size != value
	}
	if !config.VTGateCount.IsNull() && !config.VTGateCount.IsUnknown() {
		value := config.VTGateCount.ValueInt64()
		update.VTGateCount = &value
		changed = changed || current.count != value
	}
	if !config.VTGateMaxCount.IsNull() && !config.VTGateMaxCount.IsUnknown() {
		value := config.VTGateMaxCount.ValueInt64()
		update.VTGateMaxCount = &value
		changed = changed || current.maxCount == nil || *current.maxCount != value
	}
	if !config.VTGateTargetCPUUtilization.IsNull() && !config.VTGateTargetCPUUtilization.IsUnknown() {
		value := config.VTGateTargetCPUUtilization.ValueInt64()
		update.VTGateTargetCPUUtilization = &value
		changed = changed || current.targetCPUUtilization == nil || *current.targetCPUUtilization != value
	}

	if !changed {
		refreshVTGateAutoscalingModel(data, current)
		return nil
	}

	resize, err := r.client.DatabaseBranches.UpdateVitessBranchVTGateConfiguration(
		ctx,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
		update,
	)
	if err != nil {
		return err
	}

	refreshVTGateAutoscalingModel(data, configurationFromResize(current.branchID, *resize))
	return nil
}

func (r *VitessBranchVTGateAutoscalingResource) readCurrent(
	ctx context.Context,
	data *vitessBranchVTGateAutoscalingResourceModel,
) (currentVTGateConfiguration, error) {
	settings, err := r.client.DatabaseBranches.GetVitessBranchSettings(
		ctx,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
	)
	if err != nil {
		return currentVTGateConfiguration{}, err
	}

	size := settings.VTGateName
	if size == "" {
		size = settings.VTGateSize
	}

	return currentVTGateConfiguration{
		branchID:             settings.ID,
		autoscaling:          settings.VTGateAutoscaling,
		size:                 size,
		count:                settings.VTGateCount,
		maxCount:             settings.VTGateMaxCount,
		targetCPUUtilization: settings.VTGateTargetCPUUtilization,
	}, nil
}

func configurationFromResize(branchID string, resize sdk.VitessBranchResizeRequest) currentVTGateConfiguration {
	size := resize.VTGateName
	if size == "" {
		size = resize.VTGateSize
	}
	return currentVTGateConfiguration{
		branchID:             branchID,
		autoscaling:          resize.VTGateAutoscaling,
		size:                 size,
		count:                resize.VTGateCount,
		maxCount:             resize.VTGateMaxCount,
		targetCPUUtilization: resize.VTGateTargetCPUUtilization,
	}
}

func refreshVTGateAutoscalingModel(data *vitessBranchVTGateAutoscalingResourceModel, current currentVTGateConfiguration) {
	data.ID = types.StringValue(current.branchID)
	data.VTGateAutoscaling = types.BoolValue(current.autoscaling)
	data.VTGateSize = types.StringValue(current.size)
	data.VTGateCount = types.Int64Value(current.count)
	if current.maxCount == nil {
		data.VTGateMaxCount = types.Int64Null()
	} else {
		data.VTGateMaxCount = types.Int64Value(*current.maxCount)
	}
	if current.targetCPUUtilization == nil {
		data.VTGateTargetCPUUtilization = types.Int64Null()
	} else {
		data.VTGateTargetCPUUtilization = types.Int64Value(*current.targetCPUUtilization)
	}
}
