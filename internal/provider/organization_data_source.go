package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &organizationDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationDataSource{}
)

func newOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

type organizationDataSource struct {
	client *planetscale.Client
}

type featuresDataSourceModel struct {
	Insights      *bool `json:"insights,omitempty" tfsdk:"insights"`
	SingleTenancy *bool `json:"single_tenancy,omitempty" tfsdk:"single_tenancy"`
	Sso           *bool `json:"sso,omitempty" tfsdk:"sso"`
}

type flagsDataSourceModel struct {
	ExampleFlag *string `tfsdk:"example_flag"`
}

type organizationDataSourceModel struct {
	AdminOnlyProductionAccess bool                     `tfsdk:"admin_only_production_access"`
	BillingEmail              *string                  `tfsdk:"billing_email"`
	CanCreateDatabases        bool                     `tfsdk:"can_create_databases"`
	CreatedAt                 string                   `tfsdk:"created_at"`
	DatabaseCount             float64                  `tfsdk:"database_count"`
	Features                  *featuresDataSourceModel `tfsdk:"features"`
	Flags                     *flagsDataSourceModel    `tfsdk:"flags"`
	FreeDatabasesRemaining    float64                  `tfsdk:"free_databases_remaining"`
	HasPastDueInvoices        bool                     `tfsdk:"has_past_due_invoices"`
	Id                        string                   `tfsdk:"id"`
	Name                      string                   `tfsdk:"name"`
	Plan                      string                   `tfsdk:"plan"`
	SingleTenancy             bool                     `tfsdk:"single_tenancy"`
	SleepingDatabaseCount     float64                  `tfsdk:"sleeping_database_count"`
	Sso                       bool                     `tfsdk:"sso"`
	SsoDirectory              bool                     `tfsdk:"sso_directory"`
	SsoPortalUrl              *string                  `tfsdk:"sso_portal_url"`
	UpdatedAt                 string                   `tfsdk:"updated_at"`
	ValidBillingInfo          bool                     `tfsdk:"valid_billing_info"`
}

func (d *organizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *organizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
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
	}}
}

func (d *organizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *organizationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res200, _, _, _, err := d.client.GetOrganization(ctx, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read organization", err.Error())
		return
	}
	if res200 == nil {
		resp.Diagnostics.AddError("Received a nil organization", "")
		return
	}
	state := organizationDataSourceModel{
		AdminOnlyProductionAccess: res200.AdminOnlyProductionAccess,
		BillingEmail:              res200.BillingEmail,
		CanCreateDatabases:        res200.CanCreateDatabases,
		CreatedAt:                 res200.CreatedAt,
		DatabaseCount:             res200.DatabaseCount,
		FreeDatabasesRemaining:    res200.FreeDatabasesRemaining,
		HasPastDueInvoices:        res200.HasPastDueInvoices,
		Id:                        res200.Id,
		Name:                      res200.Name,
		Plan:                      res200.Plan,
		SingleTenancy:             res200.SingleTenancy,
		SleepingDatabaseCount:     res200.SleepingDatabaseCount,
		Sso:                       res200.Sso,
		SsoDirectory:              res200.SsoDirectory,
		SsoPortalUrl:              res200.SsoPortalUrl,
		UpdatedAt:                 res200.UpdatedAt,
		ValidBillingInfo:          res200.ValidBillingInfo,
	}
	if res200.Flags != nil {
		state.Flags = &flagsDataSourceModel{
			ExampleFlag: res200.Flags.ExampleFlag,
		}
	}
	if res200.Features != nil {
		state.Features = &featuresDataSourceModel{
			Insights:      res200.Features.Insights,
			SingleTenancy: res200.Features.SingleTenancy,
			Sso:           res200.Features.Sso,
		}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
