package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPostgresIPRestrictionResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-postgres"
	// A database allows one entry per schema and role, and an entry applies to
	// every branch. Scoping the entry to a role nothing connects as keeps it
	// unique and leaves the other acceptance tests' connections alone.
	role := randomWithPrefix("test-ip-restriction")
	resourceAddress := "planetscale_postgres_ip_restriction.test"

	variables := func(description string, cidrs ...string) config.Variables {
		cidrVariables := make([]config.Variable, 0, len(cidrs))
		for _, cidr := range cidrs {
			cidrVariables = append(cidrVariables, config.StringVariable(cidr))
		}
		return config.Variables{
			"organization":  config.StringVariable(testAccOrg),
			"database_name": config.StringVariable(databaseName),
			"role":          config.StringVariable(role),
			"cidrs":         config.ListVariable(cidrVariables...),
			"description":   config.StringVariable(description),
		}
	}
	stringList := func(values ...string) knownvalue.Check {
		checks := make([]knownvalue.Check, 0, len(values))
		for _, value := range values {
			checks = append(checks, knownvalue.StringExact(value))
		}
		return knownvalue.ListExact(checks)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables("Created by Terraform", "203.0.113.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("role"), knownvalue.StringExact(role)),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("schema"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("cidrs"), stringList("203.0.113.0/24")),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("description"), knownvalue.StringExact("Created by Terraform")),
					statecheck.ExpectKnownValue("data.planetscale_postgres_ip_restriction.test", tfjsonpath.New("role"), knownvalue.StringExact(role)),
					statecheck.ExpectKnownValue("data.planetscale_postgres_ip_restriction.test", tfjsonpath.New("cidrs"), stringList("203.0.113.0/24")),
					// The database may hold other entries, so only check the
					// list is non-empty.
					statecheck.ExpectKnownValue("data.planetscale_postgres_ip_restrictions.test", tfjsonpath.New("data").AtSliceIndex(0), knownvalue.NotNull()),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables("Updated by Terraform", "203.0.113.0/24", "198.51.100.7/32"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceAddress, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("cidrs"), stringList("203.0.113.0/24", "198.51.100.7/32")),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("description"), knownvalue.StringExact("Updated by Terraform")),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables("", "203.0.113.0/24", "198.51.100.7/32"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("description"), knownvalue.StringExact("")),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables("", "203.0.113.0/24", "198.51.100.7/32"),
				ResourceName:    resourceAddress,
				ImportState:     true,
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
