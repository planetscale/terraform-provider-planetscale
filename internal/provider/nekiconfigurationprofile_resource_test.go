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

func TestAccNekiConfigurationProfileResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-neki-profile")
	profileNameOriginal := randomWithPrefix("test-profile")
	profileNameRenamed := randomWithPrefix("test-profile-renamed")
	resourceAddress := "planetscale_neki_configuration_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccNekiPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"profile_name":  config.StringVariable(profileNameOriginal),
					"cluster_size":  config.StringVariable("PS_DEV_AWS_ARM"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("name"), knownvalue.StringExact(profileNameOriginal)),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("cluster_size"), knownvalue.StringExact("PS_DEV_AWS_ARM")),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("state"), knownvalue.StringExact("ready")),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"profile_name":  config.StringVariable(profileNameRenamed),
					"cluster_size":  config.StringVariable("PS_DEV_AWS_ARM"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("name"), knownvalue.StringExact(profileNameRenamed)),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"profile_name":  config.StringVariable(profileNameRenamed),
					"cluster_size":  config.StringVariable("PS_DEV_AWS_ARM"),
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
