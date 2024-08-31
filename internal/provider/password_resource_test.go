package provider

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPasswordResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-passwd-db")
	branchName := acctest.RandomWithPrefix("branch")
	passwdName := acctest.RandomWithPrefix("testacc-passwd")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "default_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_database.test", "cluster_size", "PS-10"),
					resource.TestCheckResourceAttr("planetscale_branch.test", "name", branchName),
					resource.TestCheckResourceAttr("planetscale_branch.test", "parent_branch", "main"),
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
			// TODO: Implement a test for password Update.
		},
	})
}

// TestAccPasswordResource_OutOfBandDelete tests the out-of-band deletion of a branch password.
// In this test we simulate the password has been deleted out of band, perhaps by
// a user on the console or using pscale CLI.
// https://github.com/planetscale/terraform-provider-planetscale/issues/53
func TestAccPasswordResource_OutOfBandDelete(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-passwd-db")
	branchName := acctest.RandomWithPrefix("branch")
	passwdName := acctest.RandomWithPrefix("testacc-passwd")

	passId := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName),
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

func testAccPasswordResourceConfig(dbName, branchName, passwdName string) string {
	const tmpl = `
resource "planetscale_database" "test" {
  name           = "{{.DBName}}"
  organization   = "{{.Org}}"
  cluster_size   = "PS-10"
  default_branch = "main"
}

resource "planetscale_branch" "test" {
  name          = "{{.BranchName}}"
  organization  = "{{.Org}}"
  database      = planetscale_database.test.name
  parent_branch = planetscale_database.test.default_branch
}

resource "planetscale_password" "test" {
  name         = "{{.PasswdName}}"
  organization = "{{.Org}}"
  database     = planetscale_database.test.name
  branch       = planetscale_branch.test.name
}
`
	t := template.Must(template.New("config").Parse(tmpl))
	var buf bytes.Buffer
	err := t.Execute(&buf, map[string]string{
		"Org":        testAccOrg,
		"DBName":     dbName,
		"BranchName": branchName,
		"PasswdName": passwdName,
	})
	if err != nil {
		return ""
	}
	return buf.String()
}

func getPasswordIDFromState(state *terraform.State, resourceName string) (string, error) {
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
