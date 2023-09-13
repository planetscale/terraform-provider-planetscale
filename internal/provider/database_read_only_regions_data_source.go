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
	_ datasource.DataSource              = &databaseReadOnlyRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &databaseReadOnlyRegionsDataSource{}
)

func newDatabaseReadOnlyRegionsDataSource() datasource.DataSource {
	return &databaseReadOnlyRegionsDataSource{}
}

type databaseReadOnlyRegionsDataSource struct {
	client *planetscale.Client
}

type databaseReadOnlyRegionsDataSourceModel struct {
	Organization string                                  `tfsdk:"organization"`
	Name         string                                  `tfsdk:"name"`
	Regions      []databaseReadOnlyRegionDataSourceModel `tfsdk:"regions"`
}

type databaseReadOnlyRegionActorDataSourceModel struct {
	AvatarUrl   string `tfsdk:"avatar_url"`
	DisplayName string `tfsdk:"display_name"`
	Id          string `tfsdk:"id"`
}

type databaseReadOnlyRegionRegionDataSourceModel struct {
	DisplayName       string   `tfsdk:"display_name"`
	Enabled           bool     `tfsdk:"enabled"`
	Id                string   `tfsdk:"id"`
	Location          string   `tfsdk:"location"`
	Provider          string   `tfsdk:"provider"`
	PublicIpAddresses []string `tfsdk:"public_ip_addresses"`
	Slug              string   `tfsdk:"slug"`
}

type databaseReadOnlyRegionDataSourceModel struct {
	Actor       databaseReadOnlyRegionActorDataSourceModel  `tfsdk:"actor"`
	CreatedAt   string                                      `tfsdk:"created_at"`
	DisplayName string                                      `tfsdk:"display_name"`
	Id          string                                      `tfsdk:"id"`
	Ready       bool                                        `tfsdk:"ready"`
	ReadyAt     string                                      `tfsdk:"ready_at"`
	Region      databaseReadOnlyRegionRegionDataSourceModel `tfsdk:"region"`
	UpdatedAt   string                                      `tfsdk:"updated_at"`
}

func (d *databaseReadOnlyRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_read_only_regions"
}

func (d *databaseReadOnlyRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"name":         schema.StringAttribute{Required: true},
		"regions": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"actor": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"avatar_url":   schema.StringAttribute{Computed: true},
							"display_name": schema.StringAttribute{Computed: true},
							"id":           schema.StringAttribute{Computed: true},
						},
					},
					"created_at":   schema.StringAttribute{Computed: true},
					"display_name": schema.StringAttribute{Computed: true},
					"id":           schema.StringAttribute{Computed: true},
					"ready":        schema.BoolAttribute{Computed: true},
					"ready_at":     schema.StringAttribute{Computed: true},
					"updated_at":   schema.StringAttribute{Computed: true},
					"region": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"slug":                schema.StringAttribute{Computed: true},
							"display_name":        schema.StringAttribute{Computed: true},
							"enabled":             schema.BoolAttribute{Computed: true},
							"id":                  schema.StringAttribute{Computed: true},
							"location":            schema.StringAttribute{Computed: true},
							"provider":            schema.StringAttribute{Computed: true},
							"public_ip_addresses": schema.ListAttribute{Computed: true, ElementType: types.StringType},
						},
					},
				},
			},
		},
	}}
}

func (d *databaseReadOnlyRegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *databaseReadOnlyRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *databaseReadOnlyRegionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res200, err := d.client.ListReadOnlyRegions(ctx, data.Organization, data.Name, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list database read only regions", err.Error())
		return
	}
	if res200 == nil {
		resp.Diagnostics.AddError("Received a nil database read only regions list", "")
		return
	}
	state := databaseReadOnlyRegionsDataSourceModel{
		Organization: data.Organization,
		Name:         data.Name,
		Regions:      make([]databaseReadOnlyRegionDataSourceModel, 0, len(res200.Data)),
	}
	for _, region := range res200.Data {
		state.Regions = append(state.Regions, databaseReadOnlyRegionDataSourceModel{
			Actor: databaseReadOnlyRegionActorDataSourceModel{
				AvatarUrl:   region.Actor.AvatarUrl,
				DisplayName: region.Actor.DisplayName,
				Id:          region.Actor.Id,
			},
			CreatedAt:   region.CreatedAt,
			DisplayName: region.DisplayName,
			Id:          region.Id,
			Ready:       region.Ready,
			ReadyAt:     region.ReadyAt,
			Region: databaseReadOnlyRegionRegionDataSourceModel{
				DisplayName:       region.Region.DisplayName,
				Enabled:           region.Region.Enabled,
				Id:                region.Region.Id,
				Location:          region.Region.Location,
				Provider:          region.Region.Provider,
				PublicIpAddresses: region.Region.PublicIpAddresses,
				Slug:              region.Region.Slug,
			},
			UpdatedAt: region.UpdatedAt,
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
