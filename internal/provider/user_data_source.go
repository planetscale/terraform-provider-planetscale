package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func newUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *planetscale.Client
}

type userDataSourceModel struct {
	Id                      *string `tfsdk:"id"`
	Name                    *string `tfsdk:"name"`
	AvatarUrl               *string `tfsdk:"avatar_url"`
	CreatedAt               *string `tfsdk:"created_at"`
	DefaultOrganizationId   *string `tfsdk:"default_organization_id"`
	DirectoryManaged        *bool   `tfsdk:"directory_managed"`
	DisplayName             *string `tfsdk:"display_name"`
	Email                   *string `tfsdk:"email"`
	EmailVerified           *bool   `tfsdk:"email_verified"`
	Managed                 *bool   `tfsdk:"managed"`
	Sso                     *bool   `tfsdk:"sso"`
	TwoFactorAuthConfigured *bool   `tfsdk:"two_factor_auth_configured"`
	UpdatedAt               *string `tfsdk:"updated_at"`
}

func (d *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"id":                         schema.StringAttribute{Computed: true},
		"name":                       schema.StringAttribute{Computed: true},
		"avatar_url":                 schema.StringAttribute{Computed: true},
		"created_at":                 schema.StringAttribute{Computed: true},
		"default_organization_id":    schema.StringAttribute{Computed: true},
		"directory_managed":          schema.BoolAttribute{Computed: true},
		"display_name":               schema.StringAttribute{Computed: true},
		"email":                      schema.StringAttribute{Computed: true},
		"email_verified":             schema.BoolAttribute{Computed: true},
		"managed":                    schema.BoolAttribute{Computed: true},
		"sso":                        schema.BoolAttribute{Computed: true},
		"two_factor_auth_configured": schema.BoolAttribute{Computed: true},
		"updated_at":                 schema.StringAttribute{Computed: true},
	}}
}

func (d *userDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *userDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res200, res403, res404, res500, err := d.client.GetCurrentUser(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read user", err.Error())
		return
	}
	switch {
	case res403 != nil:
		resp.Diagnostics.AddError("Unable to read user", "403 error, forbidden from getting current user")
		return
	case res404 != nil:
		resp.Diagnostics.AddError("Unable to read user", "current user not found")
		return
	case res500 != nil:
		resp.Diagnostics.AddError("Unable to read user", "500 error, try again later")
		return
	case res200 == nil:
		resp.Diagnostics.AddError("Unable to read user", "no data")
		return
	}
	state := userDataSourceModel{
		AvatarUrl:               res200.AvatarUrl,
		CreatedAt:               res200.CreatedAt,
		DefaultOrganizationId:   res200.DefaultOrganizationId,
		DirectoryManaged:        res200.DirectoryManaged,
		DisplayName:             res200.DisplayName,
		Email:                   res200.Email,
		EmailVerified:           res200.EmailVerified,
		Id:                      res200.Id,
		Managed:                 res200.Managed,
		Name:                    res200.Name,
		Sso:                     res200.Sso,
		TwoFactorAuthConfigured: res200.TwoFactorAuthConfigured,
		UpdatedAt:               res200.UpdatedAt,
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
