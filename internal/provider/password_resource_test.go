package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPasswordResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-passwd-db")
	passwdName := acctest.RandomWithPrefix("testacc-passwd")
	branchName := "main"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPasswordResourceConfig(dbName, passwdName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_password.test", "role", "admin"),
					resource.TestCheckResourceAttr("planetscale_password.test", "branch", branchName),
				),
			},
			// ImportState testing
			{
				ResourceName:      "planetscale_password.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Import requires: 'organization,database,branch,id' but 'id' of the password
				// is only known after creation. Use a func to retrieve the ID from the state:
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					id, err := getPasswordIDFromState(s, "planetscale_password.test")
					if err != nil {
						return "", err
					}
					// Import requires: 'organization,database,branch,id'
					return fmt.Sprintf("%s,%s,%s,%s", testAccOrg, dbName, branchName, id), nil
				},
				// The actual password is not returned by the API, so we can't verify it here:
				ImportStateVerifyIgnore: []string{"plaintext"},
			},
			// Update and Read testing
			// TODO: Implement a test for password Update. Best current idea is to update the
			//       password's branch to a new branch, but that requires the ability to
			//       create a database and branch in a single plan which is currently broken
			//       due to the async nature of planetscale_database.
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccPasswordResource_OutOfBandDelete tests the out-of-band deletion of a branch password.
// In this test we simulate the password has been deleted out of band, perhaps by
// a user on the console or using pscale CLI.
// https://github.com/planetscale/terraform-provider-planetscale/issues/53
func TestAccPasswordResource_OutOfBandDelete(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-passwd-db")
	passwdName := acctest.RandomWithPrefix("testacc-passwd")
	branchName := "main"

	passId := ""

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPasswordResourceConfig(dbName, passwdName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_password.test", "role", "admin"),
					resource.TestCheckResourceAttr("planetscale_password.test", "branch", branchName),
				),
			},
			// ImportState testing
			{
				ResourceName:      "planetscale_password.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					id, err := getPasswordIDFromState(s, "planetscale_password.test")
					if err != nil {
						return "", err
					}
					passId = id // save the ID for use in later steps' PreConfig func:
					// Import requires: 'organization,database,branch,id'
					return fmt.Sprintf("%s,%s,%s,%s", testAccOrg, dbName, branchName, id), nil
				},
				// The actual password is not returned by the API, so we can't verify it here:
				ImportStateVerifyIgnore: []string{"plaintext"},
			},
			// Test out-of-bands deletion of the database should produce a plan to recreate, not error.
			{
				ResourceName:       "planetscale_password.test",
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					ctx := context.Background()
					if _, err := testAccAPIClient.DeletePassword(ctx, testAccOrg, dbName, branchName, passId); err != nil {
						t.Fatalf("PreConfig: failed to delete password: %s", err)
					}
				},
			},
		},
	})
}

func testAccPasswordResourceConfig(dbName, passwdName string) string {
	return fmt.Sprintf(`
resource "planetscale_database" "test" {
  name           = "%s"
  organization   = "%s"
  cluster_size   = "PS-10"
  default_branch = "main"
}

# TODO: Uncomment when the issue with branch creation after db creation is solved, then we can
#       also expand the test coverage for password to include a change from one branch to another.
# resource "planetscale_branch" "two" {
#   name          = "TODO"
#   organization  = "TODO"
#   database      = planetscale_database.test.name
#   parent_branch = planetscale_database.test.default_branch
# }

resource "planetscale_password" "test" {
  name         = "%s"
  organization = "%s"
  database     = planetscale_database.test.name
  branch       = planetscale_database.test.default_branch
}
  `, dbName, testAccOrg, passwdName, testAccOrg)
}

func getPasswordIDFromState(state *terraform.State, resourceName string) (string, error) {
	// resourceName := "planetscale_password.test"
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}
	if rawState == nil {
		return "", fmt.Errorf("resource %s not found in state", resourceName)
	}
	return rawState["id"], nil
}
