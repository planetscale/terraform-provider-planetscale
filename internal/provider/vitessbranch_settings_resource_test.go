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

func TestAccVitessBranchSafeMigrationsResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-safe-migrations")
	resourceAddress := "planetscale_vitess_branch_safe_migrations.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":   config.StringVariable(databaseName),
					"organization":    config.StringVariable(testAccOrg),
					"safe_migrations": config.BoolVariable(true),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("safe_migrations"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":   config.StringVariable(databaseName),
					"organization":    config.StringVariable(testAccOrg),
					"safe_migrations": config.BoolVariable(false),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("safe_migrations"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":   config.StringVariable(databaseName),
					"organization":    config.StringVariable(testAccOrg),
					"safe_migrations": config.BoolVariable(false),
				},
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: branchSettingImportID(resourceAddress),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVitessBranchVTGateAutoscalingResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-vtgate-autoscaling")
	resourceAddress := "planetscale_vitess_branch_vtgate_autoscaling.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":      config.StringVariable(databaseName),
					"organization":       config.StringVariable(testAccOrg),
					"vtgate_autoscaling": config.BoolVariable(true),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_autoscaling"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_size"),
						knownvalue.StringExact("VTG_320"),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":      config.StringVariable(databaseName),
					"organization":       config.StringVariable(testAccOrg),
					"vtgate_autoscaling": config.BoolVariable(false),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("vtgate_autoscaling"),
						knownvalue.Bool(false),
					),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name":      config.StringVariable(databaseName),
					"organization":       config.StringVariable(testAccOrg),
					"vtgate_autoscaling": config.BoolVariable(false),
				},
				ResourceName:      resourceAddress,
				ImportState:       true,
				ImportStateIdFunc: branchSettingImportID(resourceAddress),
				ImportStateVerify: true,
			},
		},
	})
}

func branchSettingImportID(resourceAddress string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		attributes := state.RootModule().Resources[resourceAddress].Primary.Attributes
		value, err := json.Marshal(map[string]string{
			"organization": attributes["organization"],
			"database":     attributes["database"],
			"branch":       attributes["branch"],
		})
		return string(value), err
	}
}
