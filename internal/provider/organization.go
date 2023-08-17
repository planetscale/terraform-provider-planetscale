package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func orgSchemaAttributes(nameIsArg bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name":                     schema.StringAttribute{Computed: !nameIsArg, Required: nameIsArg},
		"created_at":               schema.StringAttribute{Computed: true},
		"updated_at":               schema.StringAttribute{Computed: true},
		"free_databases_remaining": schema.Int64Attribute{Computed: true},
	}
}

func orgListSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"organizations": schema.ListNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: orgSchemaAttributes(false),
			},
		},
	}
}
