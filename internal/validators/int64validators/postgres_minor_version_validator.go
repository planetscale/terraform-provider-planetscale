package int64validators

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// PostgresMinorVersionValidator requires postgres_major_version whenever a
// PostgreSQL minor version is configured.
func PostgresMinorVersionValidator() validator.Int64 {
	return int64validator.AlsoRequires(path.MatchRoot("postgres_major_version"))
}
