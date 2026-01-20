package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccVitessBranchResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-vitess"
	branchNameOriginal := "test"
	branchNameRenamed := "test-renamed"
	resourceAddress := "planetscale_vitess_branch.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name": config.StringVariable(databaseName),
					"organization":  config.StringVariable(testAccOrg),
					"branch_name":   config.StringVariable(branchNameOriginal),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("cluster_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("html_url"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("mysql_address"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("mysql_edge_address"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(branchNameOriginal),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("ready"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("url"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchNameRenamed),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(branchNameRenamed),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchNameRenamed),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"database":     rs.Primary.Attributes["database"],
						"id":           rs.Primary.Attributes["id"],
						"organization": rs.Primary.Attributes["organization"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify: true,
			},
		},
	})
}
