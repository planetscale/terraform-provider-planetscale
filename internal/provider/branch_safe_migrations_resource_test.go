package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBranchSafeMigrationsResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-safe-mig-resource")
	branchName := "main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBranchSafeMigrationsResourceConfig(dbName, branchName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch_safe_migrations.test", "enabled", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "planetscale_branch_safe_migrations.test",
				ImportStateId:     fmt.Sprintf("%s,%s,%s", testAccOrg, dbName, branchName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccBranchSafeMigrationsResourceConfig(dbName, branchName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch_safe_migrations.test", "enabled", "false"),
				),
			},
		},
	})
}

// TestAccBranchSafeMigrationsResource_OutOfBandChange tests handling of out-of-band changes.
func TestAccBranchSafeMigrationsResource_OutOfBandChange(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-safe-mig-resource-oob")
	branchName := "main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchSafeMigrationsResourceConfig(dbName, branchName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch_safe_migrations.test", "enabled", "true"),
				),
			},
			{
				PreConfig: func() {
					ctx := context.Background()
					if _, err := testAccAPIClient.DisableSafeMigrations(ctx, testAccOrg, dbName, branchName); err != nil {
						t.Fatalf("PreConfig: failed to disable safe migrations: %s", err)
					}
				},
				Config: testAccBranchSafeMigrationsResourceConfig(dbName, branchName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch_safe_migrations.test", "enabled", "true"),
				),
			},
		},
	})
}

func testAccBranchSafeMigrationsResourceConfig(dbName string, branch string, enabled bool) string {
	return fmt.Sprintf(`
resource "planetscale_database" "test" {
  organization   = "%s"
  name           = "%s"
  cluster_size   = "PS-10"
	default_branch = "%s"
}

resource "planetscale_branch_safe_migrations" "test" {
  organization = planetscale_database.test.organization
  database     = planetscale_database.test.name
  branch       = planetscale_database.test.default_branch
  enabled      = %t
}
`, testAccOrg, dbName, branch, enabled)
}
