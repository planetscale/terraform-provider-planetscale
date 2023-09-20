package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
)

var (
	_ datasource.DataSource              = &backupsDataSource{}
	_ datasource.DataSourceWithConfigure = &backupsDataSource{}
)

func newBackupsDataSource() datasource.DataSource {
	return &backupsDataSource{}
}

type backupsDataSourceModel struct {
	Organization string                  `tfsdk:"organization"`
	Database     string                  `tfsdk:"database"`
	Branch       string                  `tfsdk:"branch"`
	Backups      []backupDataSourceModel `tfsdk:"backups"`
}

type backupsDataSource struct {
	client *planetscale.Client
}

func (d *backupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backups"
}

func (d *backupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{Required: true},
			"database":     schema.StringAttribute{Required: true},
			"branch":       schema.StringAttribute{Required: true},
			"backups": schema.ListNestedAttribute{Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: backupDataSourceSchemaAttribute(true),
				},
			},
		},
	}
}

func (d *backupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *backupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *backupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := d.client.ListBackups(ctx, data.Organization, data.Database, data.Branch, nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read database branch backups", err.Error())
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("Unable to read database branch backups", "no data")
		return
	}
	state := backupsDataSourceModel{
		Organization: data.Organization,
		Database:     data.Database,
		Branch:       data.Branch,
		Backups:      make([]backupDataSourceModel, 0, len(res.Data)),
	}
	for _, item := range res.Data {
		item := item
		state.Backups = append(state.Backups, *backupFromClient(&item, data.Organization, data.Database, data.Branch, resp.Diagnostics))
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
