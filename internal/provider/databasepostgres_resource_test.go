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

func TestAccDatabasePostgresResource_Lifecycle(t *testing.T) {
	t.Parallel()

	name := randomWithPrefix("testacc")
	resourceAddress := "planetscale_database_postgres.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"name":         config.StringVariable(name),
					"organization": config.StringVariable(testAccOrg),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("branches_url"),
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
						tfjsonpath.New("url"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"name":         config.StringVariable(name),
					"organization": config.StringVariable(testAccOrg),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"name":         rs.Primary.Attributes["name"],
						"organization": rs.Primary.Attributes["organization"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify: true,
				// Not returned by API, therefore cannot be imported correctly.
				ImportStateVerifyIgnore: []string{"cluster_size"},
			},
		},
	})
}

func TestAccDatabasePostgresResource_ClusterSize(t *testing.T) {
	t.Skip("TODO: This test is long, ~250-300s. Revisit this when we decide how to support resizes, or when we implement a deletion_protection attribute")
	t.Parallel()

	name := randomWithPrefix("testacc")
	clusterSizeOriginal := "PS_5_AWS_X86"
	clusterSizeUpdated := "PS_10_AWS_X86"
	resourceAddress := "planetscale_database_postgres.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"cluster_size": config.StringVariable(clusterSizeOriginal),
					"name":         config.StringVariable(name),
					"organization": config.StringVariable(testAccOrg),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"cluster_size": config.StringVariable(clusterSizeOriginal),
					"name":         config.StringVariable(name),
					"organization": config.StringVariable(testAccOrg),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"name":         rs.Primary.Attributes["name"],
						"organization": rs.Primary.Attributes["organization"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify: true,
				// Not returned by API, therefore cannot be imported correctly.
				ImportStateVerifyIgnore: []string{"cluster_size"},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"cluster_size": config.StringVariable(clusterSizeUpdated),
					"name":         config.StringVariable(name),
					"organization": config.StringVariable(testAccOrg),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("cluster_size"),
						knownvalue.StringExact(clusterSizeUpdated),
					),
				},
			},
		},
	})
}
