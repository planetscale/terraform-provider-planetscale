package listvalidators

import "github.com/hashicorp/terraform-plugin-framework/schema/validator"

func PostgresExtensionsValidator() validator.List {
	return extensionsValidator()
}
