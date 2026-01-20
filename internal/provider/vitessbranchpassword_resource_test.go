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

func TestAccVitessBranchPasswordResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-vitess"
	branchName := "main"
	passwordNameOriginal := "test-password"
	passwordNameRenamed := "test-password-renamed"
	resourceAddress := "planetscale_vitess_branch_password.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":  config.StringVariable(databaseName),
					"organization":   config.StringVariable(testAccOrg),
					"branch_name":    config.StringVariable(branchName),
					"password_name":  config.StringVariable(passwordNameOriginal),
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
						knownvalue.StringExact(passwordNameOriginal),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("username"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("plain_text"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("access_host_url"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("role"),
						knownvalue.StringExact("admin"),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":   config.StringVariable(testAccOrg),
					"database_name":  config.StringVariable(databaseName),
					"branch_name":    config.StringVariable(branchName),
					"password_name":  config.StringVariable(passwordNameRenamed),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("name"),
						knownvalue.StringExact(passwordNameRenamed),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":   config.StringVariable(testAccOrg),
					"database_name":  config.StringVariable(databaseName),
					"branch_name":    config.StringVariable(branchName),
					"password_name":  config.StringVariable(passwordNameRenamed),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"branch":       rs.Primary.Attributes["branch"],
						"database":     rs.Primary.Attributes["database"],
						"id":           rs.Primary.Attributes["id"],
						"organization": rs.Primary.Attributes["organization"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify: true,
				// Ignore plain_text as it is sensitive and not returned on read
				ImportStateVerifyIgnore: []string{"plain_text"},
			},
		},
	})
}
