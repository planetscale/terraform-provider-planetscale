package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccNekiDiscoveryDataSources(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccNekiPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":  config.StringVariable(testAccOrg),
					"database_name": config.StringVariable(randomWithPrefix("testacc-neki-discovery")),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.planetscale_neki_routers.all", tfjsonpath.New("data"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_router_parameters.default", tfjsonpath.New("data"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_configuration_profile_parameters.default", tfjsonpath.New("data"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_shards.all", tfjsonpath.New("data"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_configuration_profile_shards.default", tfjsonpath.New("data"), knownvalue.NotNull()),
				},
			},
		},
	})
}
