package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPostgresBouncersDataSource(t *testing.T) {
	t.Parallel()

	databaseName := "testacc-postgres"
	branchName := "main"
	bouncerName := randomWithPrefix("test-bouncer-ds")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(databaseName),
					"branch_name":   config.StringVariable(branchName),
					"bouncer_name":  config.StringVariable(bouncerName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.planetscale_postgres_bouncer.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(bouncerName),
					),
					statecheck.ExpectKnownValue(
						"data.planetscale_postgres_bouncer.test",
						tfjsonpath.New("bouncer_size"),
						knownvalue.StringExact("PGB_5"),
					),
					// Other acceptance tests may create bouncers on the same
					// branch in parallel, so only check the list is non-empty.
					statecheck.ExpectKnownValue(
						"data.planetscale_postgres_bouncers.test",
						tfjsonpath.New("data").AtSliceIndex(0),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
