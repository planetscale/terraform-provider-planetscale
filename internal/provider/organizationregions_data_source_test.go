package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccOrganizationRegionsDataSource(t *testing.T) {
	t.Parallel()

	resourceAddress := "data.planetscale_organization_regions.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("data"),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"current_default":     knownvalue.NotNull(),
								"display_name":        knownvalue.StringExact("AWS us-east-1"),
								"enabled":             knownvalue.Bool(true),
								"id":                  knownvalue.StringExact("kc0e1ij8juzp"),
								"location":            knownvalue.StringExact("N. Virginia"),
								"provider":            knownvalue.StringExact("AWS"),
								"public_ip_addresses": knownvalue.NotNull(),
								"slug":                knownvalue.StringExact("us-east"),
							}),
						}),
					),
				},
			},
		},
	})
}
