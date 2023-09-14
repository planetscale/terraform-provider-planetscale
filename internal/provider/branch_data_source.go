package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &branchDataSource{}
	_ datasource.DataSourceWithConfigure = &branchDataSource{}
)

func newDatabaseBranchDataSource() datasource.DataSource {
	return &branchDataSource{}
}

type branchDataSource struct {
	client *planetscale.Client
}

type branchApiActorDataSourceModel struct {
	AvatarUrl   string `tfsdk:"avatar_url"`
	DisplayName string `tfsdk:"display_name"`
	Id          string `tfsdk:"id"`
}

type branchRegionDataSourceModel struct {
	DisplayName       string   `tfsdk:"display_name"`
	Enabled           bool     `tfsdk:"enabled"`
	Id                string   `tfsdk:"id"`
	Location          string   `tfsdk:"location"`
	Provider          string   `tfsdk:"provider"`
	PublicIpAddresses []string `tfsdk:"public_ip_addresses"`
	Slug              string   `tfsdk:"slug"`
}

type branchRestoredFromBranchDataSourceModel struct {
	CreatedAt string `tfsdk:"created_at"`
	DeletedAt string `tfsdk:"deleted_at"`
	Id        string `tfsdk:"id"`
	Name      string `tfsdk:"name"`
	UpdatedAt string `tfsdk:"updated_at"`
}

type branchDataSourceModel struct {
	Organization string `tfsdk:"organization"`
	Database     string `tfsdk:"database"`
	Name         string `tfsdk:"name"`

	AccessHostUrl               types.String                             `tfsdk:"access_host_url"`
	ClusterRateName             types.String                             `tfsdk:"cluster_rate_name"`
	CreatedAt                   types.String                             `tfsdk:"created_at"`
	HtmlUrl                     types.String                             `tfsdk:"html_url"`
	Id                          types.String                             `tfsdk:"id"`
	InitialRestoreId            types.String                             `tfsdk:"initial_restore_id"`
	MysqlAddress                types.String                             `tfsdk:"mysql_address"`
	MysqlEdgeAddress            types.String                             `tfsdk:"mysql_edge_address"`
	ParentBranch                types.String                             `tfsdk:"parent_branch"`
	Production                  types.Bool                               `tfsdk:"production"`
	Ready                       types.Bool                               `tfsdk:"ready"`
	RestoreChecklistCompletedAt types.String                             `tfsdk:"restore_checklist_completed_at"`
	SchemaLastUpdatedAt         types.String                             `tfsdk:"schema_last_updated_at"`
	ShardCount                  types.Float64                            `tfsdk:"shard_count"`
	Sharded                     types.Bool                               `tfsdk:"sharded"`
	UpdatedAt                   types.String                             `tfsdk:"updated_at"`
	ApiActor                    *branchApiActorDataSourceModel           `tfsdk:"api_actor"`
	PlanetscaleRegion           *branchRegionDataSourceModel             `tfsdk:"planetscale_region"`
	RestoredFromBranch          *branchRestoredFromBranchDataSourceModel `tfsdk:"restored_from_branch"`
}

func (d *branchDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

func (d *branchDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"database":     schema.StringAttribute{Required: true},
		"name":         schema.StringAttribute{Required: true},

		"access_host_url":                schema.StringAttribute{Computed: true},
		"cluster_rate_name":              schema.StringAttribute{Computed: true},
		"created_at":                     schema.StringAttribute{Computed: true},
		"html_url":                       schema.StringAttribute{Computed: true},
		"id":                             schema.StringAttribute{Computed: true},
		"initial_restore_id":             schema.StringAttribute{Computed: true},
		"mysql_address":                  schema.StringAttribute{Computed: true},
		"mysql_edge_address":             schema.StringAttribute{Computed: true},
		"parent_branch":                  schema.StringAttribute{Computed: true},
		"production":                     schema.BoolAttribute{Computed: true},
		"ready":                          schema.BoolAttribute{Computed: true},
		"restore_checklist_completed_at": schema.StringAttribute{Computed: true},
		"schema_last_updated_at":         schema.StringAttribute{Computed: true},
		"shard_count":                    schema.Float64Attribute{Computed: true},
		"sharded":                        schema.BoolAttribute{Computed: true},
		"updated_at":                     schema.StringAttribute{Computed: true},

		"api_actor": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"avatar_url":   schema.StringAttribute{Computed: true},
				"display_name": schema.StringAttribute{Computed: true},
				"id":           schema.StringAttribute{Computed: true},
			},
		},
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
	}}
}

func (d *branchDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*planetscale.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *planetscale.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *branchDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data *branchDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res200, err := d.client.GetBranch(ctx, data.Organization, data.Database, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database branch", err.Error())
		return
	}
	if res200 == nil {
		resp.Diagnostics.AddError("Unable to read database branch", "no data")
		return
	}
	state := branchDataSourceModel{
		Organization: data.Organization,
		Database:     data.Database,
		Name:         data.Name,

		AccessHostUrl:               types.StringPointerValue(res200.AccessHostUrl),
		ClusterRateName:             types.StringValue(res200.ClusterRateName),
		CreatedAt:                   types.StringValue(res200.CreatedAt),
		HtmlUrl:                     types.StringValue(res200.HtmlUrl),
		Id:                          types.StringValue(res200.Id),
		InitialRestoreId:            types.StringPointerValue(res200.InitialRestoreId),
		MysqlAddress:                types.StringValue(res200.MysqlAddress),
		MysqlEdgeAddress:            types.StringValue(res200.MysqlEdgeAddress),
		ParentBranch:                types.StringValue(res200.ParentBranch),
		Production:                  types.BoolValue(res200.Production),
		Ready:                       types.BoolValue(res200.Ready),
		RestoreChecklistCompletedAt: types.StringPointerValue(res200.RestoreChecklistCompletedAt),
		SchemaLastUpdatedAt:         types.StringValue(res200.SchemaLastUpdatedAt),
		ShardCount:                  types.Float64PointerValue(res200.ShardCount),
		Sharded:                     types.BoolValue(res200.Sharded),
		UpdatedAt:                   types.StringValue(res200.UpdatedAt),
	}
	if res200.ApiActor != nil {
		state.ApiActor = &branchApiActorDataSourceModel{
			AvatarUrl:   res200.ApiActor.AvatarUrl,
			DisplayName: res200.ApiActor.DisplayName,
			Id:          res200.ApiActor.Id,
		}
	}
	if res200.PlanetscaleRegion != nil {
		state.PlanetscaleRegion = &branchRegionDataSourceModel{
			DisplayName:       res200.PlanetscaleRegion.DisplayName,
			Enabled:           res200.PlanetscaleRegion.Enabled,
			Id:                res200.PlanetscaleRegion.Id,
			Location:          res200.PlanetscaleRegion.Location,
			Provider:          res200.PlanetscaleRegion.Provider,
			PublicIpAddresses: res200.PlanetscaleRegion.PublicIpAddresses,
			Slug:              res200.PlanetscaleRegion.Slug,
		}
	}
	if res200.RestoredFromBranch != nil {
		state.RestoredFromBranch = &branchRestoredFromBranchDataSourceModel{
			CreatedAt: res200.RestoredFromBranch.CreatedAt,
			DeletedAt: res200.RestoredFromBranch.DeletedAt,
			Id:        res200.RestoredFromBranch.Id,
			Name:      res200.RestoredFromBranch.Name,
			UpdatedAt: res200.RestoredFromBranch.UpdatedAt,
		}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
