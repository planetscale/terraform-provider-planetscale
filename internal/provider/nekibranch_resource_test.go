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

func TestAccNekiBranchResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-neki")
	branchNameOriginal := "main"
	branchNameRenamed := randomWithPrefix("main-renamed")
	resourceAddress := "planetscale_neki_branch.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccNekiPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchNameOriginal),
					"cluster_size":  config.StringVariable("PS_DEV_AWS_ARM"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("name"), knownvalue.StringExact(branchNameOriginal)),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("state"), knownvalue.StringExact("ready")),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchNameRenamed),
					"cluster_size":  config.StringVariable("PS_DEV_AWS_ARM"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("name"), knownvalue.StringExact(branchNameRenamed)),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchNameRenamed),
					"cluster_size":  config.StringVariable("PS_DEV_AWS_ARM"),
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
