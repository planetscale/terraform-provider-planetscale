package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

var backupPolicyResourceAttribute = map[string]schema.Attribute{
	"retention_unit": schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
	"retention_value": schema.Float64Attribute{
		Required: true,
		PlanModifiers: []planmodifier.Float64{
			float64planmodifier.RequiresReplace(),
		},
	},
	// read-only
	"created_at":      schema.StringAttribute{Computed: true},
	"frequency_unit":  schema.StringAttribute{Computed: true},
	"frequency_value": schema.Float64Attribute{Computed: true},
	"id":              schema.StringAttribute{Computed: true},
	"last_ran_at":     schema.StringAttribute{Computed: true},
	"name":            schema.StringAttribute{Computed: true},
	"next_run_at":     schema.StringAttribute{Computed: true},
	"schedule_day":    schema.StringAttribute{Computed: true},
	"schedule_week":   schema.StringAttribute{Computed: true},
	"target":          schema.StringAttribute{Computed: true},
	"updated_at":      schema.StringAttribute{Computed: true},
}

var backupPolicyResourceAttrTypes = map[string]attr.Type{
	"created_at":      basetypes.StringType{},
	"frequency_unit":  basetypes.StringType{},
	"frequency_value": basetypes.Float64Type{},
	"id":              basetypes.StringType{},
	"last_ran_at":     basetypes.StringType{},
	"name":            basetypes.StringType{},
	"next_run_at":     basetypes.StringType{},
	"retention_unit":  basetypes.StringType{},
	"retention_value": basetypes.Float64Type{},
	"schedule_day":    basetypes.StringType{},
	"schedule_week":   basetypes.StringType{},
	"target":          basetypes.StringType{},
	"updated_at":      basetypes.StringType{},
}

var schemaSnapshotResourceAttribute = map[string]schema.Attribute{
	"created_at": schema.StringAttribute{Computed: true},
	"id":         schema.StringAttribute{Computed: true},
	"name":       schema.StringAttribute{Computed: true},
	"updated_at": schema.StringAttribute{Computed: true},
	"url":        schema.StringAttribute{Computed: true},
}

var schemaSnapshotResourceAttrTypes = map[string]attr.Type{
	"created_at": basetypes.StringType{},
	"id":         basetypes.StringType{},
	"name":       basetypes.StringType{},
	"updated_at": basetypes.StringType{},
	"url":        basetypes.StringType{},
}

var databaseBranchResourceAttribute = map[string]schema.Attribute{
	"access_host_url":    schema.StringAttribute{Computed: true},
	"id":                 schema.StringAttribute{Computed: true},
	"mysql_edge_address": schema.StringAttribute{Computed: true},
	"name":               schema.StringAttribute{Computed: true},
	"production":         schema.BoolAttribute{Computed: true},
}

var databaseBranchResourceAttrTypes = map[string]attr.Type{
	"access_host_url":    basetypes.StringType{},
	"id":                 basetypes.StringType{},
	"mysql_edge_address": basetypes.StringType{},
	"name":               basetypes.StringType{},
	"production":         basetypes.BoolType{},
}
