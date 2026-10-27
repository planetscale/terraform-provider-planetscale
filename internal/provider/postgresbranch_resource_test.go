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

func TestAccPostgresBranchResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-postgres"
	branchNameOriginal := randomWithPrefix("test")
	branchNameRenamed := randomWithPrefix("test-renamed")
	resourceAddress := "planetscale_postgres_branch.test"
	clusterSize := "PS_DEV_AWS_ARM"
	minimumStorageBytesOriginal := int64(20 * 1024 * 1024 * 1024)
	minimumStorageBytesUpdated := int64(30 * 1024 * 1024 * 1024)
	maximumStorageBytes := int64(1024 * 1024 * 1024 * 1024)
	storageIopsOriginal := int64(3000)
	storageIopsUpdated := int64(3500)
	storageThroughputMibsOriginal := int64(125)
	storageThroughputMibsUpdated := int64(150)
	parameters := config.MapVariable(map[string]config.Variable{
		"pgconf": config.MapVariable(map[string]config.Variable{
			"max_connections": config.StringVariable("50"),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":           config.StringVariable(databaseName),
					"organization":            config.StringVariable(testAccOrg),
					"branch_name":             config.StringVariable(branchNameOriginal),
					"cluster_size":            config.StringVariable(clusterSize),
					"deletion_protected":      config.BoolVariable(true),
					"parameters":              parameters,
					"minimum_storage_bytes":   config.IntegerVariable(minimumStorageBytesOriginal),
					"maximum_storage_bytes":   config.IntegerVariable(maximumStorageBytes),
					"storage_autoscaling":     config.BoolVariable(true),
					"storage_iops":            config.IntegerVariable(storageIopsOriginal),
					"storage_throughput_mibs": config.IntegerVariable(storageThroughputMibsOriginal),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("parameters"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"pgconf": knownvalue.MapExact(map[string]knownvalue.Check{
								"max_connections": knownvalue.StringExact("50"),
							}),
						}),
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
						tfjsonpath.New("deletion_protected"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("minimum_storage_bytes"),
						knownvalue.Int64Exact(minimumStorageBytesOriginal),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("maximum_storage_bytes"),
						knownvalue.Int64Exact(maximumStorageBytes),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("storage_autoscaling"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("storage_iops"),
						knownvalue.Int64Exact(storageIopsOriginal),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("storage_throughput_mibs"),
						knownvalue.Int64Exact(storageThroughputMibsOriginal),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":            config.StringVariable(testAccOrg),
					"database_name":           config.StringVariable(databaseName),
					"branch_name":             config.StringVariable(branchNameRenamed),
					"cluster_size":            config.StringVariable(clusterSize),
					"deletion_protected":      config.BoolVariable(false),
					"parameters":              parameters,
					"minimum_storage_bytes":   config.IntegerVariable(minimumStorageBytesUpdated),
					"maximum_storage_bytes":   config.IntegerVariable(maximumStorageBytes),
					"storage_autoscaling":     config.BoolVariable(true),
					"storage_iops":            config.IntegerVariable(storageIopsUpdated),
					"storage_throughput_mibs": config.IntegerVariable(storageThroughputMibsUpdated),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(branchNameRenamed),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("parameters"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"pgconf": knownvalue.MapExact(map[string]knownvalue.Check{
								"max_connections": knownvalue.StringExact("50"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("deletion_protected"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("minimum_storage_bytes"),
						knownvalue.Int64Exact(minimumStorageBytesUpdated),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("maximum_storage_bytes"),
						knownvalue.Int64Exact(maximumStorageBytes),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("storage_autoscaling"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("storage_iops"),
						knownvalue.Int64Exact(storageIopsUpdated),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("storage").AtMapKey("storage_throughput_mibs"),
						knownvalue.Int64Exact(storageThroughputMibsUpdated),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":            config.StringVariable(testAccOrg),
					"database_name":           config.StringVariable(databaseName),
					"branch_name":             config.StringVariable(branchNameRenamed),
					"cluster_size":            config.StringVariable(clusterSize),
					"deletion_protected":      config.BoolVariable(false),
					"parameters":              parameters,
					"minimum_storage_bytes":   config.IntegerVariable(minimumStorageBytesUpdated),
					"maximum_storage_bytes":   config.IntegerVariable(maximumStorageBytes),
					"storage_autoscaling":     config.BoolVariable(true),
					"storage_iops":            config.IntegerVariable(storageIopsUpdated),
					"storage_throughput_mibs": config.IntegerVariable(storageThroughputMibsUpdated),
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
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_size"},
			},
		},
	})
}

func TestAccPostgresBranchResource_CreatesAndDeletesDatabase(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-pg-lifecycle")
	branchName := "main"
	resourceAddress := "planetscale_postgres_branch.test"
	clusterSize := "PS_DEV_AWS_ARM"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchName),
					"cluster_size":  config.StringVariable(clusterSize),
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
				},
			},
		},
	})
}
