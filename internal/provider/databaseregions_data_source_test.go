package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDatabaseRegionsDataSource(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("testacc-%d", time.Now().Unix())
	resourceAddress := "data.planetscale_database_regions.test"

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
						tfjsonpath.New("data"),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"display_name":        knownvalue.NotNull(),
								"enabled":             knownvalue.NotNull(),
								"id":                  knownvalue.NotNull(),
								"location":            knownvalue.NotNull(),
								"provider":            knownvalue.NotNull(),
								"public_ip_addresses": knownvalue.NotNull(),
								"slug":                knownvalue.NotNull(),
							}),
						}),
					),
				},
			},
		},
	})
}
