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

func TestAccPostgresBouncerResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-postgres"
	branchName := "main"
	bouncerName := randomWithPrefix("test-bouncer")
	resourceAddress := "planetscale_postgres_bouncer.test"

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
					"bouncer_name":  config.StringVariable(bouncerName),
					"bouncer_size":  config.StringVariable("PGB_5"),
					"pool_size":     config.StringVariable("50"),
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
						knownvalue.StringExact(bouncerName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("target"),
						knownvalue.StringExact("primary"),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("replicas_per_cell"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("bouncer_size"),
						knownvalue.StringExact("PGB_5"),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("parameters").AtMapKey("pgbouncer").AtMapKey("default_pool_size"),
						knownvalue.StringExact("50"),
					),
				},
			},
			// Resize the bouncer and change a parameter in place. Changes roll
			// out asynchronously, but the API reflects the requested
			// configuration immediately.
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchName),
					"bouncer_name":  config.StringVariable(bouncerName),
					"bouncer_size":  config.StringVariable("PGB_10"),
					"pool_size":     config.StringVariable("100"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(bouncerName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("bouncer_size"),
						knownvalue.StringExact("PGB_10"),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("parameters").AtMapKey("pgbouncer").AtMapKey("default_pool_size"),
						knownvalue.StringExact("100"),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchName),
					"bouncer_name":  config.StringVariable(bouncerName),
					"bouncer_size":  config.StringVariable("PGB_10"),
					"pool_size":     config.StringVariable("100"),
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
