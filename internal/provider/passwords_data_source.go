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
	_ datasource.DataSource              = &passwordsDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordsDataSource{}
)

func newPasswordsDataSource() datasource.DataSource {
	return &passwordsDataSource{}
}

type passwordsDataSource struct {
	client *planetscale.Client
}

type passwordsDataSourceModel struct {
	Organization     types.String              `tfsdk:"organization"`
	Database         types.String              `tfsdk:"database"`
	Branch           types.String              `tfsdk:"branch"`
	ReadOnlyRegionId types.String              `tfsdk:"read_only_region_id"`
	Passwords        []passwordDataSourceModel `tfsdk:"passwords"`
}

func (d *passwordsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_passwords"
}

func (d *passwordsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization":        schema.StringAttribute{Required: true},
			"database":            schema.StringAttribute{Required: true},
			"branch":              schema.StringAttribute{Required: true},
			"read_only_region_id": schema.StringAttribute{Optional: true},
			"passwords": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: passwordDataSourceSchemaAttribute(true),
				},
			},
		},
	}
}

func (d *passwordsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *passwordsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *passwordsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.ListPasswords(
		ctx,
		data.Organization.ValueString(),
		data.Database.ValueString(),
		data.Branch.ValueString(),
		data.ReadOnlyRegionId.ValueStringPointer(),
		nil,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database passwords", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read database passwords", "no data")
		return
	}
	state := passwordsDataSourceModel{
		Organization:     data.Organization,
		Database:         data.Database,
		Branch:           data.Branch,
		ReadOnlyRegionId: data.ReadOnlyRegionId,
		Passwords:        make([]passwordDataSourceModel, 0, len(res.Data)),
	}
	for _, item := range res.Data {
		item := item
		state.Passwords = append(state.Passwords, *passwordFromClient(
			&item,
			data.Organization.ValueString(),
			data.Database.ValueString(),
			data.Branch.ValueString(),
			data.ReadOnlyRegionId.ValueStringPointer(),
			resp.Diagnostics,
		))
		if resp.Diagnostics.HasError() {
			return
		}
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
