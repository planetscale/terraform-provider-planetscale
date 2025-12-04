package provider

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDatabaseDefaultBranchResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := fmt.Sprintf("testacc-%d", time.Now().Unix())
	resourceAddress := "planetscale_database_default_branch.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
			ConfigVariables: config.Variables{
				"database_name": config.StringVariable(databaseName),
				"organization":  config.StringVariable(testAccOrg),
			},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("branch"),
						knownvalue.StringExact("test"),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
			ConfigVariables: config.Variables{
				"database_name": config.StringVariable(databaseName),
				"organization":  config.StringVariable(testAccOrg),
			},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"database":     rs.Primary.Attributes["database"],
						"organization": rs.Primary.Attributes["organization"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "database",
			},
		},
	})
}
