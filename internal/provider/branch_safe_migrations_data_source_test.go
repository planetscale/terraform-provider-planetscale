package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBranchSafeMigrationsDataSource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-safe-mig-ds")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchSafeMigrationsDataSourceConfig(dbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.planetscale_branch_safe_migrations.test", "enabled", "false"),
				),
			},
		},
	})
}

func testAccBranchSafeMigrationsDataSourceConfig(dbName string) string {
	return fmt.Sprintf(`
resource "planetscale_database" "test" {
  organization   = "%s"
  name           = "%s"
  cluster_size   = "PS-10"
  default_branch = "main"
}

data "planetscale_branch_safe_migrations" "test" {
  organization = "%s"
  database     = planetscale_database.test.name
  branch       = planetscale_database.test.default_branch
}
`, testAccOrg, dbName, testAccOrg)
}
