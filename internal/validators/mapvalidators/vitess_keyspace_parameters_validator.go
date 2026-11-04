package mapvalidators

import (
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func VitessKeyspaceParametersValidator() validator.Map {
	return MapVitessParametersValidator{namespaces: []string{"vttablet", "mysqld"}}
}
