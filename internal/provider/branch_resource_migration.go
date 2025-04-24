package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// This file contains the state migration logic and historical models/schemas
// for the planetscale_branch resource.

// --- Version 0 ---

// branchResourceModelV1 represents the state data for schema version 1.
// v1 differs from v0 in the following ways:
//  1. ClusterRateName (cluster_rate_name) removed in v1.
type branchResourceModelV1 struct {
	Organization types.String `tfsdk:"organization"`
	Database     types.String `tfsdk:"database"`

	Name                        types.String  `tfsdk:"name"`
	ParentBranch                types.String  `tfsdk:"parent_branch"`
	Actor                       types.Object  `tfsdk:"actor"`
	CreatedAt                   types.String  `tfsdk:"created_at"`
	HtmlUrl                     types.String  `tfsdk:"html_url"`
	Id                          types.String  `tfsdk:"id"`
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
	SeedData                    types.String  `tfsdk:"seed_data"`
	UpdatedAt                   types.String  `tfsdk:"updated_at"`
}

// branchSchemaV1 defines the schema for version 1.
// v1 differs from v0 in the following ways:
//  1. cluster_rate_name removed in v1.
func branchSchemaV1() *schema.Schema {
	return &schema.Schema{
		Version:             1, // Explicitly set version
		Description:         "A PlanetScale branch.",
		MarkdownDescription: "A PlanetScale branch.",
		Attributes: map[string]schema.Attribute{
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
			"name": schema.StringAttribute{
				Description: "The name of the branch.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parent_branch": schema.StringAttribute{
				Description: "The name of the parent branch from which the branch was created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"production": schema.BoolAttribute{
				Description: "Whether or not the branch is a production branch.",
				Computed:    true, Optional: true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the branch.",
				Computed:    true,
			},
			"actor": schema.SingleNestedAttribute{
				Description: "The actor who created this branch.",
				Computed:    true,
				Attributes:  actorResourceSchemaAttribute,
			},
			"created_at": schema.StringAttribute{
				Description: "When the branch was created.",
				Computed:    true,
			},
			"html_url": schema.StringAttribute{
				Description: "Planetscale app URL for the branch.",
				Computed:    true,
			},
			"mysql_address": schema.StringAttribute{
				Description: "The MySQL address for the branch.",
				Computed:    true,
			},
			"mysql_edge_address": schema.StringAttribute{
				Description: "The address of the MySQL provider for the branch.",
				Computed:    true,
			},
			"region": schema.SingleNestedAttribute{
				Description: "The region in which this branch lives.",
				Computed:    true,
				Attributes:  regionResourceSchemaAttribute,
			},
			"ready": schema.BoolAttribute{
				Description: "Whether or not the branch is ready to serve queries.",
				Computed:    true,
			},
			"restore_checklist_completed_at": schema.StringAttribute{
				Description: "When a user last marked a backup restore checklist as completed.",
				Computed:    true,
			},
			"restored_from_branch": schema.SingleNestedAttribute{
				Description: "todo",
				Computed:    true,
				Attributes:  restoredFromBranchSchemaAttribute,
			},
			"schema_last_updated_at": schema.StringAttribute{
				Description: "When the schema for the branch was last updated.",
				Computed:    true,
			},
			"shard_count": schema.Float64Attribute{
				Description: "The number of shards in the branch.",
				Computed:    true,
			},
			"sharded": schema.BoolAttribute{
				Description: "Whether or not the branch is sharded.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When the branch was last updated.",
				Computed:    true,
			},
			"seed_data": schema.StringAttribute{
				Description: "Seed data",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// --- Version 0 ---

// branchResourceModelV0 represents the state data for schema version 0.
type branchResourceModelV0 struct {
	Organization types.String `tfsdk:"organization"`
	Database     types.String `tfsdk:"database"`

	Name                        types.String  `tfsdk:"name"`
	ParentBranch                types.String  `tfsdk:"parent_branch"`
	Actor                       types.Object  `tfsdk:"actor"`
	ClusterRateName             types.String  `tfsdk:"cluster_rate_name"` // This field only exists in v0
	CreatedAt                   types.String  `tfsdk:"created_at"`
	HtmlUrl                     types.String  `tfsdk:"html_url"`
	Id                          types.String  `tfsdk:"id"`
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

// branchSchemaV0 defines the schema for version 0.
func branchSchemaV0() *schema.Schema {
	return &schema.Schema{
		Attributes: map[string]schema.Attribute{
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
			"name": schema.StringAttribute{
				Description: "The name of the branch.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parent_branch": schema.StringAttribute{
				Description: "The name of the parent branch from which the branch was created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"production": schema.BoolAttribute{
				Description: "Whether or not the branch is a production branch.",
				Computed:    true, Optional: true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the branch.",
				Computed:    true,
			},
			"actor": schema.SingleNestedAttribute{
				Description: "The actor who created this branch.",
				Computed:    true,
				Attributes:  actorResourceSchemaAttribute,
			},
			"cluster_rate_name": schema.StringAttribute{
				Description: "The SKU representing the branch's cluster size.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "When the branch was created.",
				Computed:    true,
			},
			"html_url": schema.StringAttribute{
				Description: "Planetscale app URL for the branch.",
				Computed:    true,
			},
			"mysql_address": schema.StringAttribute{
				Description: "The MySQL address for the branch.",
				Computed:    true,
			},
			"mysql_edge_address": schema.StringAttribute{
				Description: "The address of the MySQL provider for the branch.",
				Computed:    true,
			},
			"region": schema.SingleNestedAttribute{
				Description: "The region in which this branch lives.",
				Computed:    true,
				Attributes:  regionResourceSchemaAttribute,
			},
			"ready": schema.BoolAttribute{
				Description: "Whether or not the branch is ready to serve queries.",
				Computed:    true,
			},
			"restore_checklist_completed_at": schema.StringAttribute{
				Description: "When a user last marked a backup restore checklist as completed.",
				Computed:    true,
			},
			"restored_from_branch": schema.SingleNestedAttribute{
				Description: "todo",
				Computed:    true,
				Attributes:  restoredFromBranchSchemaAttribute,
			},
			"schema_last_updated_at": schema.StringAttribute{
				Description: "When the schema for the branch was last updated.",
				Computed:    true,
			},
			"shard_count": schema.Float64Attribute{
				Description: "The number of shards in the branch.",
				Computed:    true,
			},
			"sharded": schema.BoolAttribute{
				Description: "Whether or not the branch is sharded.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "When the branch was last updated.",
				Computed:    true,
			},
		},
	}
}

// upgradeBranchStateV0toCurrent implements the logic for upgrading state from version 0 to the current version.
func upgradeBranchStateV0toCurrent(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var priorStateData branchResourceModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map prior state data (v0) to the current model (v1), omitting ClusterRateName
	upgradedStateData := branchResourceModelV1{
		Organization:                priorStateData.Organization,
		Database:                    priorStateData.Database,
		Name:                        priorStateData.Name,
		ParentBranch:                priorStateData.ParentBranch,
		Actor:                       priorStateData.Actor,
		CreatedAt:                   priorStateData.CreatedAt,
		HtmlUrl:                     priorStateData.HtmlUrl,
		Id:                          priorStateData.Id,
		MysqlAddress:                priorStateData.MysqlAddress,
		MysqlEdgeAddress:            priorStateData.MysqlEdgeAddress,
		Region:                      priorStateData.Region,
		Production:                  priorStateData.Production,
		Ready:                       priorStateData.Ready,
		RestoreChecklistCompletedAt: priorStateData.RestoreChecklistCompletedAt,
		RestoredFromBranch:          priorStateData.RestoredFromBranch,
		SchemaLastUpdatedAt:         priorStateData.SchemaLastUpdatedAt,
		ShardCount:                  priorStateData.ShardCount,
		Sharded:                     priorStateData.Sharded,
		UpdatedAt:                   priorStateData.UpdatedAt,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
}

// UpgradeState implements resource.ResourceWithUpgradeState.
func (r *branchResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// Upgrade v0 state to current version
		0: {
			PriorSchema:   branchSchemaV0(),
			StateUpgrader: upgradeBranchStateV0toCurrent,
		},
		// IMPORTANT!
		// Terraform does not iterate through each stateu upgrade function. It only runs
		// one upgrader func. Thus, when a new version is introduced, the new function
		// and each existing upgrade functions must be modified to upgrade to the new
		// current version.
	}
}
