package listvalidators

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func NekiExtensionsValidator() validator.List {
	return listvalidator.ConflictsWith(
		path.MatchRoot("parameters").AtMapKey("pgconf").AtMapKey("shared_preload_libraries"),
		path.MatchRoot("parameters").AtMapKey("pgconf").AtMapKey("session_preload_libraries"),
	)
}
