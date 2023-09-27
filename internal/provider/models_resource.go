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
	"avatar_url": schema.StringAttribute{
		Computed: true, Description: "The URL of the actor's avatar",
	},
	"display_name": schema.StringAttribute{
		Computed: true, Description: "The name of the actor",
	},
	"id": schema.StringAttribute{
		Computed: true, Description: "The ID of the actor",
	},
}

var actorResourceAttrTypes = map[string]attr.Type{
	"avatar_url":   types.StringType,
	"display_name": types.StringType,
	"id":           types.StringType,
}

var regionResourceSchemaAttribute = map[string]schema.Attribute{
	"display_name": schema.StringAttribute{
		Description: "Name of the region.",
		Computed:    true,
	},
	"enabled": schema.BoolAttribute{
		Description: "Whether or not the region is currently active.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID of the region.",
		Computed:    true,
	},
	"location": schema.StringAttribute{
		Description: "Location of the region.",
		Computed:    true,
	},
	"provider": schema.StringAttribute{
		Description: "Provider for the region (ex. AWS).",
		Computed:    true,
	},
	"public_ip_addresses": schema.ListAttribute{
		Description: "Public IP addresses for the region.",
		Computed:    true, ElementType: types.StringType,
	},
	"slug": schema.StringAttribute{
		Description: "The slug of the region.",
		Computed:    true,
	},
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
	"created_at": schema.StringAttribute{
		Description: "When the resource was created.",
		Computed:    true,
	},
	"deleted_at": schema.StringAttribute{
		Description: "When the resource was deleted, if deleted.",
		Computed:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID for the resource.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name for the resource.",
		Computed:    true,
	},
	"updated_at": schema.StringAttribute{
		Description: "When the resource was last updated.",
		Computed:    true,
	},
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
		Description: "The unit for the retention period of the backup policy.",
		Required:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
	"retention_value": schema.Float64Attribute{
		Description: "A number value for the retention period of the backup policy.",
		Required:    true,
		PlanModifiers: []planmodifier.Float64{
			float64planmodifier.RequiresReplace(),
		},
	},
	// read-only
	"id": schema.StringAttribute{
		Description: "The ID of the backup policy.",
		Computed:    true,
	},
	"created_at": schema.StringAttribute{
		Description: "When the backup policy was created.",
		Computed:    true,
	},
	"frequency_unit": schema.StringAttribute{
		Description: "The unit for the frequency of the backup policy.",
		Computed:    true,
	},
	"frequency_value": schema.Float64Attribute{
		Description: "A number value for the frequency of the backup policy.",
		Computed:    true,
	},
	"last_ran_at": schema.StringAttribute{
		Description: "When the backup was last run.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the backup policy.",
		Computed:    true,
	},
	"next_run_at": schema.StringAttribute{
		Description: "When the backup will next run.",
		Computed:    true,
	},
	"schedule_day": schema.StringAttribute{
		Description: "Day of the week that the backup is scheduled.",
		Computed:    true,
	},
	"schedule_week": schema.StringAttribute{
		Description: "Week of the month that the backup is scheduled.",
		Computed:    true,
	},
	"target": schema.StringAttribute{
		Description: "Whether the backup policy is for a production or development database, or for a database branch.",
		Computed:    true,
	},
	"updated_at": schema.StringAttribute{
		Description: "When the backup policy was last updated.",
		Computed:    true,
	},
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
