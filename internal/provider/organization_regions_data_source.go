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
	_ datasource.DataSource              = &organizationRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationRegionsDataSource{}
)

func newOrganizationRegionsDataSource() datasource.DataSource {
	return &organizationRegionsDataSource{}
}

type organizationRegionsDataSource struct {
	client *planetscale.Client
}

type organizationRegionDataSourceModel struct {
	DisplayName       string   `tfsdk:"display_name"`
	Enabled           bool     `tfsdk:"enabled"`
	Id                string   `tfsdk:"id"`
	Location          string   `tfsdk:"location"`
	Provider          string   `tfsdk:"provider"`
	PublicIpAddresses []string `tfsdk:"public_ip_addresses"`
	Slug              string   `tfsdk:"slug"`
}

type organizationRegionsDataSourceModel struct {
	Organization string                              `tfsdk:"organization"`
	Regions      []organizationRegionDataSourceModel `tfsdk:"regions"`
}

func (d *organizationRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_regions"
}

func (d *organizationRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"regions": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
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
	}}
}

func (d *organizationRegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *organizationRegionsDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgName := data.Organization

	res200, err := d.client.ListRegionsForOrganization(ctx, orgName, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read organization regions", err.Error())
		return
	}
	if res200 == nil {
		resp.Diagnostics.AddError("Unable to read organization regions", "no data")
		return
	}
	regions := make([]organizationRegionDataSourceModel, 0, len(res200.Data))
	for _, rg := range res200.Data {
		regions = append(regions, organizationRegionDataSourceModel{
			DisplayName:       rg.DisplayName,
			Enabled:           rg.Enabled,
			Id:                rg.Id,
			Location:          rg.Location,
			Provider:          rg.Provider,
			PublicIpAddresses: rg.PublicIpAddresses,
			Slug:              rg.Slug,
		})
	}
	state := organizationRegionsDataSourceModel{
		Organization: data.Organization,
		Regions:      regions,
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
