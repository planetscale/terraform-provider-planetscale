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

func TestAccPostgresReadOnlyReplicaResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-postgres"
	branchName := "main"
	replicaName := randomWithPrefix("tfrr")
	resourceAddress := "planetscale_postgres_read_only_replica.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			// Create with replicas and cluster_size omitted: the API defaults
			// to one instance at the primary's cluster size.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchName),
					"replica_name":  config.StringVariable(replicaName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(replicaName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("replicas"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("cluster_size"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("access_host_url"),
						knownvalue.NotNull(),
					),
				},
			},
			// Scale out in place.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchName),
					"replica_name":  config.StringVariable(replicaName),
					"replicas":      config.IntegerVariable(2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(replicaName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("replicas"),
						knownvalue.Int64Exact(2),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchName),
					"replica_name":  config.StringVariable(replicaName),
					"replicas":      config.IntegerVariable(2),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"organization": rs.Primary.Attributes["organization"],
						"database":     rs.Primary.Attributes["database"],
						"branch":       rs.Primary.Attributes["branch"],
						"name":         rs.Primary.Attributes["name"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}
