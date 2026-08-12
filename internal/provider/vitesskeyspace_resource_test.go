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

func TestAccVitessKeyspaceResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-vitess-ks")
	branchName := "main"
	keyspaceName := randomWithPrefix("tfks")
	resourceAddress := "planetscale_vitess_keyspace.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":   config.StringVariable(testAccOrg),
					"database_name":  config.StringVariable(databaseName),
					"branch_name":    config.StringVariable(branchName),
					"keyspace_name":  config.StringVariable(keyspaceName),
					"cluster_size":   config.StringVariable("PS_10"),
					"extra_replicas": config.IntegerVariable(0),
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
						knownvalue.StringExact(keyspaceName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("cluster_size"),
						knownvalue.StringExact("PS_10"),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("ready"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("resizing"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":   config.StringVariable(testAccOrg),
					"database_name":  config.StringVariable(databaseName),
					"branch_name":    config.StringVariable(branchName),
					"keyspace_name":  config.StringVariable(keyspaceName),
					"cluster_size":   config.StringVariable("PS_10"),
					"extra_replicas": config.IntegerVariable(1),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("extra_replicas"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("ready"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("resizing"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("resize_pending"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":   config.StringVariable(testAccOrg),
					"database_name":  config.StringVariable(databaseName),
					"branch_name":    config.StringVariable(branchName),
					"keyspace_name":  config.StringVariable(keyspaceName),
					"cluster_size":   config.StringVariable("PS_10"),
					"extra_replicas": config.IntegerVariable(1),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"branch":       rs.Primary.Attributes["branch"],
						"database":     rs.Primary.Attributes["database"],
						"name":         rs.Primary.Attributes["name"],
						"organization": rs.Primary.Attributes["organization"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify: true,
			},
		},
	})
}
