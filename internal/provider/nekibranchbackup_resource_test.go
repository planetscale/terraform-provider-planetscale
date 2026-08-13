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

func TestAccNekiBranchBackupResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-neki-backup")
	backupName := randomWithPrefix("test-backup")
	resourceAddress := "planetscale_neki_branch_backup.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccNekiPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"backup_name":   config.StringVariable(backupName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("name"), knownvalue.StringExact(backupName)),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("state"), knownvalue.StringExact("success")),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"backup_name":   config.StringVariable(backupName),
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
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retention_value", "retention_unit"},
			},
		},
	})
}
