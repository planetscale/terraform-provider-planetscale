package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccVitessBackupPolicyResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-vitess"
	policyName := randomWithPrefix("test-backup-policy")
	resourceAddress := "planetscale_vitess_backup_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"policy_name":   config.StringVariable(policyName),
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
						knownvalue.StringExact(policyName),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("target"),
						knownvalue.StringExact("development"),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("retention_value"),
						knownvalue.Int64Exact(7),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"policy_name":   config.StringVariable(policyName),
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
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVitessBackupPoliciesDataSource(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-vitess"
	policyName := randomWithPrefix("test-backup-policy")
	resourceAddress := "planetscale_vitess_backup_policy.test"
	dataSourceAddress := "data.planetscale_vitess_backup_policies.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"policy_name":   config.StringVariable(policyName),
				},
				Check: backupPolicyListedInDataSource(resourceAddress, dataSourceAddress),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataSourceAddress,
						tfjsonpath.New("data"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func backupPolicyListedInDataSource(resourceAddress, dataSourceAddress string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, ok := s.RootModule().Resources[resourceAddress]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceAddress)
		}
		dataSource, ok := s.RootModule().Resources[dataSourceAddress]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dataSourceAddress)
		}

		policyID := policy.Primary.Attributes["id"]
		countStr, ok := dataSource.Primary.Attributes["data.#"]
		if !ok {
			return fmt.Errorf("data source %s has no data attribute", dataSourceAddress)
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("parse data.# for %s: %w", dataSourceAddress, err)
		}

		for i := range count {
			if dataSource.Primary.Attributes[fmt.Sprintf("data.%d.id", i)] == policyID {
				return nil
			}
		}

		return fmt.Errorf("backup policy %s not found in %s list", policyID, dataSourceAddress)
	}
}
