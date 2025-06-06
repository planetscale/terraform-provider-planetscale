package provider

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBranchResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-branch-db")
	branchName := acctest.RandomWithPrefix("branch")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Initial creation with required fields
			{
				Config: testAccBranchResourceConfigTemplate(map[string]string{
					"organization":  testAccOrg,
					"database":      dbName,
					"name":          branchName,
					"parent_branch": "main",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.test", "name", branchName),
					resource.TestCheckResourceAttr("planetscale_branch.test", "parent_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_branch.test", "production", "false"),
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
			// Update in-place to production branch
			{
				Config: testAccBranchResourceConfigTemplate(map[string]string{
					"organization":  testAccOrg,
					"database":      dbName,
					"name":          branchName,
					"parent_branch": "main",
					"production":    "true",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_branch.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.test", "production", "true"),
				),
			},
		},
	})
}

// TestAccBranchResource_WithSeedData tests the creation of a branch with the seed data
func TestAccBranchResource_SeedData(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-branch-db")
	branchName := acctest.RandomWithPrefix("branch")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Initial creation with required fields
			{
				Config: testAccBranchResourceConfigTemplate(map[string]string{
					"organization":  testAccOrg,
					"database":      dbName,
					"name":          branchName,
					"parent_branch": "main",
					"seed_data":     "last_successful_backup",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.test", "name", branchName),
					resource.TestCheckResourceAttr("planetscale_branch.test", "parent_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_branch.test", "production", "true"),
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
			// Update in-place to development branch
			{
				Config: testAccBranchResourceConfigTemplate(map[string]string{
					"organization":  testAccOrg,
					"database":      dbName,
					"name":          branchName,
					"parent_branch": "main",
					"production":    "false",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_branch.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.test", "production", "false"),
				),
			},
		},
	})
}

// TestAccBranchResource_ProductionCreate tests the creation of a branch with
// the production flag set to true.
func TestAccBranchResource_ProductionCreate(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-branch-db")
	branchName := acctest.RandomWithPrefix("branch")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBranchResourceConfigTemplate(map[string]string{
					"organization":  testAccOrg,
					"database":      dbName,
					"name":          branchName,
					"parent_branch": "main",
					"production":    "true",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.test", "production", "true"),
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
				Config: testAccBranchResourceConfigTemplate(map[string]string{
					"organization":  testAccOrg,
					"database":      dbName,
					"name":          branchName,
					"parent_branch": "main",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_branch.test", "name", branchName),
					resource.TestCheckResourceAttr("planetscale_branch.test", "parent_branch", "main"),
					resource.TestCheckResourceAttr("planetscale_branch.test", "sharded", "false"),
				),
			},
			// Test out-of-bands deletion of the branch should produce a plan to recreate, not error
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

func testAccBranchResourceConfigTemplate(settings map[string]string) string {
	const tmpl = `
resource "planetscale_database" "test" {
    organization = "{{.organization}}"
    name         = "{{.database}}"
    cluster_size = "PS-10"
}

resource "planetscale_branch" "test" {
    organization  = "{{.organization}}"
    database      = planetscale_database.test.name
    name          = "{{.name}}"
    parent_branch = "{{.parent_branch}}"
    {{if .production}}production = {{.production}}{{end}}
    {{if .seed_data}}seed_data = "{{.seed_data}}"{{end}}
}
`
	t := template.Must(template.New("config").Parse(tmpl))
	var buf bytes.Buffer
	err := t.Execute(&buf, settings)
	if err != nil {
		return ""
	}
	return buf.String()
}
