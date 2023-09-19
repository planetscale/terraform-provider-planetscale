package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var actorResourceSchemaAttribute = map[string]schema.Attribute{
	"avatar_url":   schema.StringAttribute{Computed: true},
	"display_name": schema.StringAttribute{Computed: true},
	"id":           schema.StringAttribute{Computed: true},
}

var actorResourceAttrTypes = map[string]attr.Type{
	"avatar_url":   types.StringType,
	"display_name": types.StringType,
	"id":           types.StringType,
}

var regionResourceSchemaAttribute = map[string]schema.Attribute{
	"display_name":        schema.StringAttribute{Computed: true},
	"enabled":             schema.BoolAttribute{Computed: true},
	"id":                  schema.StringAttribute{Computed: true},
	"location":            schema.StringAttribute{Computed: true},
	"provider":            schema.StringAttribute{Computed: true},
	"public_ip_addresses": schema.ListAttribute{Computed: true, ElementType: types.StringType},
	"slug":                schema.StringAttribute{Computed: true},
}

var regionResourceAttrTypes = map[string]attr.Type{
	"display_name":        types.StringType,
	"enabled":             types.BoolType,
	"id":                  types.StringType,
	"location":            types.StringType,
	"provider":            types.StringType,
	"public_ip_addresses": types.ListType{ElemType: types.StringType},
	"slug":                types.StringType,
}

var restoredFromBranchSchemaAttribute = map[string]schema.Attribute{
	"created_at": schema.StringAttribute{Computed: true},
	"deleted_at": schema.StringAttribute{Computed: true},
	"id":         schema.StringAttribute{Computed: true},
	"name":       schema.StringAttribute{Computed: true},
	"updated_at": schema.StringAttribute{Computed: true},
}

var restoredFromBranchResourceAttrTypes = map[string]attr.Type{
	"created_at": types.StringType,
	"deleted_at": types.StringType,
	"id":         types.StringType,
	"name":       types.StringType,
	"updated_at": types.StringType,
}

type restoredFromBranchResource struct {
	CreatedAt types.String `tfsdk:"created_at"`
	DeletedAt types.String `tfsdk:"deleted_at"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

var importDataSourceResourceAttrTypes = map[string]attr.Type{
	"database": basetypes.StringType{},
	"hostname": basetypes.StringType{},
	"port":     basetypes.StringType{},
}

var importResourceAttrTypes = map[string]attr.Type{
	"data_source":         basetypes.ObjectType{AttrTypes: importDataSourceResourceAttrTypes},
	"finished_at":         basetypes.StringType{},
	"import_check_errors": basetypes.StringType{},
	"started_at":          basetypes.StringType{},
	"state":               basetypes.StringType{},
}
