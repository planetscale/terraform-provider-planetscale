package provider

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPostgresBranchMaintenanceWindowResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-postgres"
	branchName := "main"
	resourceAddress := "planetscale_postgres_branch_maintenance_window.test"
	dataSourceAddress := "data.planetscale_postgres_branch_maintenance_window.test"

	// The API refuses a change while the window is open, so keep both windows
	// away from the hour the test runs.
	now := time.Now().UTC()
	windowAt := func(offset time.Duration) (int, int) {
		start := now.Add(offset)
		return int(start.Weekday()), start.Hour()
	}
	firstDay, firstHour := windowAt(12 * time.Hour)
	secondDay, secondHour := windowAt(-6 * time.Hour)

	variables := func(day, hour int) config.Variables {
		return config.Variables{
			"organization":  config.StringVariable(testAccOrg),
			"database_name": config.StringVariable(databaseName),
			"branch_name":   config.StringVariable(branchName),
			"day":           config.IntegerVariable(day),
			"hour":          config.IntegerVariable(hour),
		}
	}

	windowChecks := func(day, hour int) []statecheck.StateCheck {
		checks := []statecheck.StateCheck{
			statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("id"), knownvalue.NotNull()),
			statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("duration"), knownvalue.Int64Exact(2)),
		}
		for _, address := range []string{resourceAddress, dataSourceAddress} {
			checks = append(checks,
				statecheck.ExpectKnownValue(address, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(day))),
				statecheck.ExpectKnownValue(address, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(hour))),
			)
		}
		return checks
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   variables(firstDay, firstHour),
				ConfigStateChecks: windowChecks(firstDay, firstHour),
			},
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   variables(secondDay, secondHour),
				ConfigStateChecks: windowChecks(secondDay, secondHour),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables(secondDay, secondHour),
				ResourceName:    resourceAddress,
				ImportState:     true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"branch":       rs.Primary.Attributes["branch"],
						"database":     rs.Primary.Attributes["database"],
						"organization": rs.Primary.Attributes["organization"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify: true,
			},
		},
	})
}
