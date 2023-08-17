package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func regionSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"slug":         schema.StringAttribute{Computed: true},
		"display_name": schema.StringAttribute{Computed: true},
		"location":     schema.StringAttribute{Computed: true},
		"enabled":      schema.BoolAttribute{Computed: true},
	}
}

func regionListSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"data": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: regionSchemaAttributes(),
			},
		},
	}
}
