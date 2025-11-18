package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccOrganizationClusterSizeSkusDataSource(t *testing.T) {
	t.Parallel()

	resourceAddress := "data.planetscale_organization_cluster_size_skus.test"

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
								"cpu":          knownvalue.StringExact("1/16"),
								"display_name": knownvalue.StringExact("PS-DEV"),
								"name":         knownvalue.StringExact("PS_DEV"),
								"ram":          knownvalue.Int64Exact(536870912),
							}),
						}),
					),
				},
			},
		},
	})
}
