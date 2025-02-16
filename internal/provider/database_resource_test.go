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

func TestAccDatabaseResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-db")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Initial creation with required fields
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization": testAccOrg,
					"name":         dbName,
					"cluster_size": "PS-10",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "cluster_size", "PS-10"),
					// Check defaults are set correctly
					resource.TestCheckResourceAttr("planetscale_database.test", "allow_data_branching", "false"),
				),
			},

			// ImportState testing
			{
				ResourceName:      "planetscale_database.test",
				ImportStateId:     fmt.Sprintf("%s,%s", testAccOrg, dbName),
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: API does not return cluster_size which causes a diff on import. When fixed, remove it.
				ImportStateVerifyIgnore: []string{"cluster_size", "updated_at"},
			},

			// Test updateable settings
			//
			// Enable data branching
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization":         testAccOrg,
					"name":                 dbName,
					"cluster_size":         "PS-10",
					"allow_data_branching": "true",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_database.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "allow_data_branching", "true"),
				),
			},
			// Enable insights raw queries
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization":         testAccOrg,
					"name":                 dbName,
					"cluster_size":         "PS-10",
					"allow_data_branching": "true",
					"insights_raw_queries": "true",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_database.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "insights_raw_queries", "true"),
				),
			},

			// Enable automatic migrations
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization":         testAccOrg,
					"name":                 dbName,
					"cluster_size":         "PS-10",
					"allow_data_branching": "true",
					"insights_raw_queries": "true",
					"automatic_migrations": "true",
					"migration_table_name": "schema_migrations",
					"migration_framework":  "rails",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_database.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "automatic_migrations", "true"),
					resource.TestCheckResourceAttr("planetscale_database.test", "migration_table_name", "schema_migrations"),
					resource.TestCheckResourceAttr("planetscale_database.test", "migration_framework", "rails"),
				),
			},
			// Change migration framework, rails -> other
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization":         testAccOrg,
					"name":                 dbName,
					"cluster_size":         "PS-10",
					"allow_data_branching": "true",
					"insights_raw_queries": "true",
					"automatic_migrations": "true",
					"migration_table_name": "schema_migrations",
					"migration_framework":  "other",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_database.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "automatic_migrations", "true"),
					resource.TestCheckResourceAttr("planetscale_database.test", "migration_table_name", "schema_migrations"),
					resource.TestCheckResourceAttr("planetscale_database.test", "migration_framework", "other"),
				),
			},

			// Disable data branching
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization":         testAccOrg,
					"name":                 dbName,
					"cluster_size":         "PS-10",
					"allow_data_branching": "false",
					"insights_raw_queries": "true",
					"automatic_migrations": "true",
					"migration_table_name": "schema_migrations",
					"migration_framework":  "other",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_database.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "allow_data_branching", "false"),
				),
			},
			// Disable insights raw queries
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization":         testAccOrg,
					"name":                 dbName,
					"cluster_size":         "PS-10",
					"allow_data_branching": "false",
					"insights_raw_queries": "false",
					"automatic_migrations": "true",
					"migration_table_name": "schema_migrations",
					"migration_framework":  "other",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_database.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "insights_raw_queries", "false"),
				),
			},
			// Disable automatic migrations
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization":         testAccOrg,
					"name":                 dbName,
					"cluster_size":         "PS-10",
					"allow_data_branching": "false",
					"insights_raw_queries": "false",
					"automatic_migrations": "false",
				}),
				ConfigPlanChecks: checkExpectUpdate("planetscale_database.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "automatic_migrations", "false"),
				),
			},

			// Change cluster_size should trigger a recreate.
			{
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization": testAccOrg,
					"name":         dbName,
					"cluster_size": "PS-20",
				}),
				// TODO: Update this test when the API supports in-place cluster_size changes: https://github.com/planetscale/terraform-provider-planetscale/issues/107
				ConfigPlanChecks: checkExpectRecreate("planetscale_database.test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_database.test", "cluster_size", "PS-20"),
				),
			},
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
				Config: testAccDatabaseResourceConfigTemplate(map[string]string{
					"organization": testAccOrg,
					"name":         dbName,
					"cluster_size": "PS-10",
				}),
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

func testAccDatabaseResourceConfigTemplate(settings map[string]string) string {
	const tmpl = `
resource "planetscale_database" "test" {
    organization  = "{{.organization}}"
    name          = "{{.name}}"
    cluster_size  = "{{.cluster_size}}"
    {{if .allow_data_branching}}allow_data_branching = {{.allow_data_branching}}{{end}}
    {{if .automatic_migrations}}automatic_migrations = {{.automatic_migrations}}{{end}}
    {{if .insights_raw_queries}}insights_raw_queries = {{.insights_raw_queries}}{{end}}
    {{if .migration_table_name}}migration_table_name = "{{.migration_table_name}}"{{end}}
    {{if .migration_framework}}migration_framework = "{{.migration_framework}}"{{end}}
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
