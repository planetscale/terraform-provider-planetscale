package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBranchDataSource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-branch-ds")
	branchName := acctest.RandomWithPrefix("branch")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchDataSourceConfig(dbName, branchName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.planetscale_branch.test", "name", branchName),
					resource.TestCheckResourceAttr("data.planetscale_branch.test", "database", dbName),
					resource.TestCheckResourceAttr("data.planetscale_branch.test", "organization", testAccOrg),
					resource.TestCheckResourceAttr("data.planetscale_branch.test", "parent_branch", "main"),
					resource.TestCheckResourceAttr("data.planetscale_branch.test", "production", "false"),
					resource.TestCheckResourceAttr("data.planetscale_branch.test", "sharded", "false"),
					resource.TestCheckResourceAttrSet("data.planetscale_branch.test", "id"),
					resource.TestCheckResourceAttrSet("data.planetscale_branch.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.planetscale_branch.test", "updated_at"),
					resource.TestCheckResourceAttrSet("data.planetscale_branch.test", "mysql_edge_address"),
				),
			},
		},
	})
}

func testAccBranchDataSourceConfig(dbName, branchName string) string {
	return fmt.Sprintf(`
resource "planetscale_database" "test" {
	organization   = "%s"
	name           = "%s"
	cluster_size   = "PS-10"
	default_branch = "main"
}

resource "planetscale_branch" "test" {
	organization  = planetscale_database.test.organization
	database      = planetscale_database.test.name
	name          = "%s"
	parent_branch = planetscale_database.test.default_branch
}

data "planetscale_branch" "test" {
	organization = planetscale_branch.test.organization
	database     = planetscale_branch.test.database
	name         = planetscale_branch.test.name
}
`, testAccOrg, dbName, branchName)
}
