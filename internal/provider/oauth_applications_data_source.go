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
	_ datasource.DataSource              = &oauthApplicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthApplicationsDataSource{}
)

func newOAuthApplicationsDataSource() datasource.DataSource {
	return &oauthApplicationsDataSource{}
}

type oauthApplicationsDataSource struct {
	client *planetscale.Client
}

type oauthApplicationDataSourceModel struct {
	Avatar      *string  `tfsdk:"avatar"`
	ClientId    string   `tfsdk:"client_id"`
	CreatedAt   string   `tfsdk:"created_at"`
	Domain      string   `tfsdk:"domain"`
	Id          string   `tfsdk:"id"`
	Name        string   `tfsdk:"name"`
	RedirectUri string   `tfsdk:"redirect_uri"`
	Scopes      []string `tfsdk:"scopes"`
	Tokens      float64  `tfsdk:"tokens"`
	UpdatedAt   string   `tfsdk:"updated_at"`
}

type oauthApplicationsDataSourceModel struct {
	Organization string                            `tfsdk:"organization"`
	Applications []oauthApplicationDataSourceModel `tfsdk:"applications"`
}

func (d *oauthApplicationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_applications"
}

func (d *oauthApplicationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"organization": schema.StringAttribute{Required: true},
		"applications": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"avatar":       schema.StringAttribute{Computed: true},
					"client_id":    schema.StringAttribute{Computed: true},
					"created_at":   schema.StringAttribute{Computed: true},
					"domain":       schema.StringAttribute{Computed: true},
					"id":           schema.StringAttribute{Computed: true},
					"name":         schema.StringAttribute{Computed: true},
					"redirect_uri": schema.StringAttribute{Computed: true},
					"scopes":       schema.ListAttribute{Computed: true, ElementType: types.StringType},
					"tokens":       schema.Float64Attribute{Computed: true},
					"updated_at":   schema.StringAttribute{Computed: true},
				},
			},
		},
	}}
}

func (d *oauthApplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *oauthApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *oauthApplicationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res200, res403, res404, res500, err := d.client.ListOauthApplications(ctx, data.Organization, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list oauth applications", err.Error())
		return
	}
	switch {
	case res403 != nil:
		resp.Diagnostics.AddError("Unable to list oauth applications", fmt.Sprintf("403 error, forbidden from listing oauth applications for organization %q", data.Organization))
		return
	case res404 != nil:
		resp.Diagnostics.AddError("Unable to list oauth applications", fmt.Sprintf("organization %q not found", data.Organization))
		return
	case res500 != nil:
		resp.Diagnostics.AddError("Unable to list oauth applications", "500 error, try again later")
		return
	case res200 == nil:
		resp.Diagnostics.AddError("Unable to list oauth applications", "no data")
		return
	}

	var list []oauthApplicationDataSourceModel
	for _, el := range res200.Data {
		list = append(list, oauthApplicationDataSourceModel{
			Avatar:      el.Avatar,
			ClientId:    el.ClientId,
			CreatedAt:   el.CreatedAt,
			Domain:      el.Domain,
			Id:          el.Id,
			Name:        el.Name,
			RedirectUri: el.RedirectUri,
			Scopes:      el.Scopes,
			Tokens:      el.Tokens,
			UpdatedAt:   el.UpdatedAt,
		})
	}

	state := oauthApplicationsDataSourceModel{
		Organization: data.Organization,
		Applications: list,
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
