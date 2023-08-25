package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &organizationsDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationsDataSource{}
)

func newOrganizationsDataSource() datasource.DataSource {
	return &organizationsDataSource{}
}

type organizationsDataSource struct {
	client *planetscale.Client
}

type organizationsDataSourceModel struct {
	Organizations []organizationDataSourceModel `tfsdk:"organizations"`
}

func (d *organizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *organizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organizations": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"admin_only_production_access": schema.BoolAttribute{Computed: true},
					"billing_email":                schema.StringAttribute{Computed: true},
					"can_create_databases":         schema.BoolAttribute{Computed: true},
					"created_at":                   schema.StringAttribute{Computed: true},
					"database_count":               schema.Float64Attribute{Computed: true},
					"features": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"insights":       schema.StringAttribute{Computed: true},
							"single_tenancy": schema.StringAttribute{Computed: true},
							"sso":            schema.StringAttribute{Computed: true},
						},
					},
					"flags": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"example_flag": schema.StringAttribute{Computed: true},
						},
					},
					"free_databases_remaining": schema.Float64Attribute{Computed: true},
					"has_past_due_invoices":    schema.BoolAttribute{Computed: true},
					"id":                       schema.StringAttribute{Computed: true},
					"name":                     schema.StringAttribute{Computed: true},
					"plan":                     schema.StringAttribute{Computed: true},
					"single_tenancy":           schema.BoolAttribute{Computed: true},
					"sleeping_database_count":  schema.Float64Attribute{Computed: true},
					"sso":                      schema.BoolAttribute{Computed: true},
					"sso_directory":            schema.BoolAttribute{Computed: true},
					"sso_portal_url":           schema.StringAttribute{Computed: true},
					"updated_at":               schema.StringAttribute{Computed: true},
					"valid_billing_info":       schema.BoolAttribute{Computed: true},
				},
			},
		},
	}}
}

func (d *organizationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	res200, _, _, _, err := d.client.ListOrganizations(ctx, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read organizations", err.Error())
		return
	}
	orgs := make([]organizationDataSourceModel, 0, len(res200.Data))
	for _, org := range res200.Data {
		orgs = append(orgs, organizationDataSourceModel{
			AdminOnlyProductionAccess: org.AdminOnlyProductionAccess,
			BillingEmail:              org.BillingEmail,
			CanCreateDatabases:        org.CanCreateDatabases,
			CreatedAt:                 org.CreatedAt,
			DatabaseCount:             org.DatabaseCount,
			FreeDatabasesRemaining:    org.FreeDatabasesRemaining,
			HasPastDueInvoices:        org.HasPastDueInvoices,
			Id:                        org.Id,
			Name:                      org.Name,
			Plan:                      org.Plan,
			SingleTenancy:             org.SingleTenancy,
			SleepingDatabaseCount:     org.SleepingDatabaseCount,
			Sso:                       org.Sso,
			SsoDirectory:              org.SsoDirectory,
			SsoPortalUrl:              org.SsoPortalUrl,
			UpdatedAt:                 org.UpdatedAt,
			ValidBillingInfo:          org.ValidBillingInfo,
		})
	}
	state := organizationsDataSourceModel{Organizations: orgs}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
