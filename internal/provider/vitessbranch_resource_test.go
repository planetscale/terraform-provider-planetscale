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
	branchNameOriginal := randomWithPrefix("test")
	branchNameRenamed := randomWithPrefix("test-renamed")
	resourceAddress := "planetscale_vitess_branch.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":      config.StringVariable(databaseName),
					"organization":       config.StringVariable(testAccOrg),
					"branch_name":        config.StringVariable(branchNameOriginal),
					"deletion_protected": config.BoolVariable(true),
					"safe_migrations":    config.BoolVariable(true),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("id"),
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
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("state"),
						knownvalue.StringExact("ready"),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("safe_migrations"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("deletion_protected"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":       config.StringVariable(testAccOrg),
					"database_name":      config.StringVariable(databaseName),
					"branch_name":        config.StringVariable(branchNameRenamed),
					"deletion_protected": config.BoolVariable(false),
					"safe_migrations":    config.BoolVariable(false),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(branchNameRenamed),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("safe_migrations"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("deletion_protected"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":       config.StringVariable(testAccOrg),
					"database_name":      config.StringVariable(databaseName),
					"branch_name":        config.StringVariable(branchNameRenamed),
					"deletion_protected": config.BoolVariable(false),
					"safe_migrations":    config.BoolVariable(false),
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

func TestAccVitessBranchResource_CreatesAndDeletesDatabase(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-vitess-lifecycle")
	branchName := "main"
	resourceAddress := "planetscale_vitess_branch.test"
	clusterSize := "PS_10"
	vtgateSize := "VTG_320"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":                  config.StringVariable(testAccOrg),
					"database_name":                 config.StringVariable(databaseName),
					"branch_name":                   config.StringVariable(branchName),
					"cluster_size":                  config.StringVariable(clusterSize),
					"safe_migrations":               config.BoolVariable(true),
					"vtgate_autoscaling":            config.BoolVariable(true),
					"vtgate_count":                  config.IntegerVariable(1),
					"vtgate_max_count":              config.IntegerVariable(2),
					"vtgate_size":                   config.StringVariable(vtgateSize),
					"vtgate_target_cpu_utilization": config.IntegerVariable(50),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(branchName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("ready"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("state"),
						knownvalue.StringExact("ready"),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("safe_migrations"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_autoscaling"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_count"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_max_count"),
						knownvalue.Int64Exact(2),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_size"),
						knownvalue.StringExact(vtgateSize),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_target_cpu_utilization"),
						knownvalue.Int64Exact(50),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":                  config.StringVariable(testAccOrg),
					"database_name":                 config.StringVariable(databaseName),
					"branch_name":                   config.StringVariable(branchName),
					"cluster_size":                  config.StringVariable(clusterSize),
					"safe_migrations":               config.BoolVariable(false),
					"vtgate_autoscaling":            config.BoolVariable(true),
					"vtgate_count":                  config.IntegerVariable(2),
					"vtgate_max_count":              config.IntegerVariable(3),
					"vtgate_size":                   config.StringVariable(vtgateSize),
					"vtgate_target_cpu_utilization": config.IntegerVariable(50),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("safe_migrations"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_count"),
						knownvalue.Int64Exact(2),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_max_count"),
						knownvalue.Int64Exact(3),
					),
				},
			},
		},
	})
}
