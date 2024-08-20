package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBranchResource(t *testing.T) {
	// TODO: This test currently fails because the provider returns immediately
	//       after DB creation but the DB is still pending and so the branch creation
	//       will fail. Unblock and finish this test once this issue is resolved.
	t.Skip()

	dbName := acctest.RandomWithPrefix("testacc-branch")
	branchName := "two"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBranchResourceConfig(dbName, branchName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.two", "name", branchName),
					resource.TestCheckResourceAttr("planetscale_branch.two", "parent_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_branch.two", "sharded", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "planetscale_branch.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			// TODO: Implement an update test.
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TODO: implement an out of bound deletion test like we have in the password and database tests.

func testAccBranchResourceConfig(dbName, branchName string) string {
	return fmt.Sprintf(`
resource "planetscale_database" "test" {
  organization   = "%s"
  name           = "%s"
  cluster_size   = "PS-10"
  default_branch = "main"
}

resource "planetscale_branch" "two" {
  organization  = "%s"
  database      = planetscale_database.test.name
  name          = "%s"
  parent_branch = planetscale_database.test.default_branch
}
  `, testAccOrg, dbName, testAccOrg, branchName)
}
