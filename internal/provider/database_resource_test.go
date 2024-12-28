package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-db")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatabaseResourceConfig(dbName, "PS-10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "production_branches_count", "1"),
					resource.TestCheckResourceAttr("planetscale_database.test", "default_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_database.test", "cluster_size", "PS-10"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "planetscale_database.test",
				ImportStateId:     fmt.Sprintf("%s,%s", testAccOrg, dbName),
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: API does not return cluster_size which causes a diff on import. When fixed, remove this:
				ImportStateVerifyIgnore: []string{"cluster_size", "updated_at"},
			},
			// Update and Read testing
			{
				Config: testAccDatabaseResourceConfig(dbName, "PS-20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "production_branches_count", "1"),
					resource.TestCheckResourceAttr("planetscale_database.test", "default_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_database.test", "cluster_size", "PS-20"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccDatabaseResource_outOfBandDelete tests the out-of-band deletion of the database.
// In this test we simulate the remote database has been deleted out of band, perhaps by
// a user on the console or using pscale CLI.
// https://github.com/planetscale/terraform-provider-planetscale/issues/53
func TestAccDatabaseResource_OutOfBandDelete(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-db-oob")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatabaseResourceConfig(dbName, "PS-10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "production_branches_count", "1"),
					resource.TestCheckResourceAttr("planetscale_database.test", "default_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_database.test", "cluster_size", "PS-10"),
				),
			},
			// Test out-of-bands deletion of the database should produce a plan to recreate, not error.
			{
				ResourceName:       "planetscale_database.test",
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					ctx := context.Background()
					_, err := testAccAPIClient.DeleteDatabase(ctx, testAccOrg, dbName)
					if err != nil {
						t.Fatalf("PreConfig: failed to delete database: %s", err)
					}
				},
			},
		},
	})
}

func testAccDatabaseResourceConfig(dbName string, clusterSize string) string {
	return fmt.Sprintf(`
resource "planetscale_database" "test" {
  organization   = "%s"
  name           = "%s"
  cluster_size   = "%s"
}
`, testAccOrg, dbName, clusterSize)
}
