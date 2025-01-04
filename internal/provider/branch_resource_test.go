package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBranchResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-branch")
	branchName := acctest.RandomWithPrefix("branch")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBranchResourceConfig(dbName, branchName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.test", "name", branchName),
					resource.TestCheckResourceAttr("planetscale_branch.test", "parent_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_branch.test", "sharded", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "planetscale_branch.test",
				ImportStateId:           fmt.Sprintf("%s,%s,%s", testAccOrg, dbName, branchName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated_at", "schema_last_updated_at"},
			},
			// Update and Read testing
			// TODO: Implement an update test.
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccBranchResource_OutOfBandDelete tests the out-of-band deletion of a branch.
// In this test we simulate the branch has been deleted out of band, perhaps by
// a user on the console or using pscale CLI.
// https://github.com/planetscale/terraform-provider-planetscale/issues/53
func TestAccBranchResource_OutOfBandDelete(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-branch-db")
	branchName := acctest.RandomWithPrefix("branch")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBranchResourceConfig(dbName, branchName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.test", "name", branchName),
					resource.TestCheckResourceAttr("planetscale_branch.test", "parent_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_branch.test", "sharded", "false"),
				),
			},
			// Test out-of-bands deletion of the database should produce a plan to recreate, not error.
			{
				ResourceName:       "planetscale_branch.test",
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					ctx := context.Background()
					if _, err := testAccAPIClient.DeleteBranch(ctx, testAccOrg, dbName, branchName); err != nil {
						t.Fatalf("PreConfig: failed to delete branch: %s", err)
					}
				},
			},
		},
	})
}

func testAccBranchResourceConfig(dbName, branchName string) string {
	return fmt.Sprintf(`
resource "planetscale_database" "test" {
  organization   = "%s"
  name           = "%s"
  cluster_size   = "PS-10"
  default_branch = "main"
}

resource "planetscale_branch" "test" {
  organization  = "%s"
  database      = planetscale_database.test.name
  name          = "%s"
  parent_branch = planetscale_database.test.default_branch
}
  `, testAccOrg, dbName, testAccOrg, branchName)
}
