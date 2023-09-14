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
	_ datasource.DataSource              = &branchesDataSource{}
	_ datasource.DataSourceWithConfigure = &branchesDataSource{}
)

func newBranchesDataSource() datasource.DataSource {
	return &branchesDataSource{}
}

type branchesDataSource struct {
	client *planetscale.Client
}

type branchesDataSourceModel struct {
	Organization string                          `tfsdk:"organization"`
	Database     string                          `tfsdk:"database"`
	Branches     []branchesBranchDataSourceModel `tfsdk:"branches"`
}

type branchesApiActorDataSourceModel struct {
	AvatarUrl   string `tfsdk:"avatar_url"`
	DisplayName string `tfsdk:"display_name"`
	Id          string `tfsdk:"id"`
}

type branchesRegionDataSourceModel struct {
	DisplayName       string   `tfsdk:"display_name"`
	Enabled           bool     `tfsdk:"enabled"`
	Id                string   `tfsdk:"id"`
	Location          string   `tfsdk:"location"`
	Provider          string   `tfsdk:"provider"`
	PublicIpAddresses []string `tfsdk:"public_ip_addresses"`
	Slug              string   `tfsdk:"slug"`
}

type branchesRestoredFromBranchDataSourceModel struct {
	CreatedAt string `tfsdk:"created_at"`
	DeletedAt string `tfsdk:"deleted_at"`
	Id        string `tfsdk:"id"`
	Name      string `tfsdk:"name"`
	UpdatedAt string `tfsdk:"updated_at"`
}

type branchesBranchDataSourceModel struct {
	Id                          string                                     `tfsdk:"id"`
	Name                        string                                     `tfsdk:"name"`
	AccessHostUrl               *string                                    `tfsdk:"access_host_url"`
	ApiActor                    *branchesApiActorDataSourceModel           `tfsdk:"api_actor"`
	ClusterRateName             string                                     `tfsdk:"cluster_rate_name"`
	CreatedAt                   string                                     `tfsdk:"created_at"`
	HtmlUrl                     string                                     `tfsdk:"html_url"`
	InitialRestoreId            *string                                    `tfsdk:"initial_restore_id"`
	MysqlAddress                string                                     `tfsdk:"mysql_address"`
	MysqlEdgeAddress            string                                     `tfsdk:"mysql_edge_address"`
	ParentBranch                string                                     `tfsdk:"parent_branch"`
	PlanetscaleRegion           *branchesRegionDataSourceModel             `tfsdk:"region"`
	Production                  bool                                       `tfsdk:"production"`
	Ready                       bool                                       `tfsdk:"ready"`
	RestoreChecklistCompletedAt *string                                    `tfsdk:"restore_checklist_completed_at"`
	RestoredFromBranch          *branchesRestoredFromBranchDataSourceModel `tfsdk:"restored_from_branch"`
	SchemaLastUpdatedAt         string                                     `tfsdk:"schema_last_updated_at"`
	ShardCount                  *float64                                   `tfsdk:"shard_count"`
	Sharded                     bool                                       `tfsdk:"sharded"`
	UpdatedAt                   string                                     `tfsdk:"updated_at"`
}

func (d *branchesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branches"
}

func (d *branchesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"database":     schema.StringAttribute{Required: true},
		"branches": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id":                             schema.StringAttribute{Computed: true},
					"name":                           schema.StringAttribute{Computed: true},
					"access_host_url":                schema.StringAttribute{Computed: true},
					"cluster_rate_name":              schema.StringAttribute{Computed: true},
					"created_at":                     schema.StringAttribute{Computed: true},
					"html_url":                       schema.StringAttribute{Computed: true},
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
				},
			},
		},
	}}
}

func (d *branchesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *branchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data *branchesDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res200, err := d.client.ListBranches(ctx, data.Organization, data.Database, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database branches", err.Error())
		return
	}
	if res200 == nil {
		resp.Diagnostics.AddError("Unable to read database branches", "no data")
		return
	}
	state := branchesDataSourceModel{
		Organization: data.Organization,
		Database:     data.Database,
		Branches:     make([]branchesBranchDataSourceModel, 0, len(res200.Data)),
	}

	for _, item := range res200.Data {

		out := branchesBranchDataSourceModel{
			AccessHostUrl:               item.AccessHostUrl,
			ClusterRateName:             item.ClusterRateName,
			CreatedAt:                   item.CreatedAt,
			HtmlUrl:                     item.HtmlUrl,
			Id:                          item.Id,
			InitialRestoreId:            item.InitialRestoreId,
			MysqlAddress:                item.MysqlAddress,
			MysqlEdgeAddress:            item.MysqlEdgeAddress,
			Name:                        item.Name,
			ParentBranch:                item.ParentBranch,
			Production:                  item.Production,
			Ready:                       item.Ready,
			RestoreChecklistCompletedAt: item.RestoreChecklistCompletedAt,
			SchemaLastUpdatedAt:         item.SchemaLastUpdatedAt,
			ShardCount:                  item.ShardCount,
			Sharded:                     item.Sharded,
			UpdatedAt:                   item.UpdatedAt,
		}
		if item.ApiActor != nil {
			out.ApiActor = &branchesApiActorDataSourceModel{
				AvatarUrl:   item.ApiActor.AvatarUrl,
				DisplayName: item.ApiActor.DisplayName,
				Id:          item.ApiActor.Id,
			}
		}
		if item.PlanetscaleRegion != nil {
			out.PlanetscaleRegion = &branchesRegionDataSourceModel{
				DisplayName:       item.PlanetscaleRegion.DisplayName,
				Enabled:           item.PlanetscaleRegion.Enabled,
				Id:                item.PlanetscaleRegion.Id,
				Location:          item.PlanetscaleRegion.Location,
				Provider:          item.PlanetscaleRegion.Provider,
				PublicIpAddresses: item.PlanetscaleRegion.PublicIpAddresses,
				Slug:              item.PlanetscaleRegion.Slug,
			}
		}
		if item.RestoredFromBranch != nil {
			out.RestoredFromBranch = &branchesRestoredFromBranchDataSourceModel{
				CreatedAt: item.RestoredFromBranch.CreatedAt,
				DeletedAt: item.RestoredFromBranch.DeletedAt,
				Id:        item.RestoredFromBranch.Id,
				Name:      item.RestoredFromBranch.Name,
				UpdatedAt: item.RestoredFromBranch.UpdatedAt,
			}
		}

		state.Branches = append(state.Branches, out)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
