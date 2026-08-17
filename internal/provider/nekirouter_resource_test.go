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

func TestAccNekiRouterResource_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-neki-router")
	routerName := randomWithPrefix("test-router")
	resourceAddress := "planetscale_neki_router.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccNekiPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":      config.StringVariable(testAccOrg),
					"database_name":     config.StringVariable(databaseName),
					"router_name":       config.StringVariable(routerName),
					"replicas_per_cell": config.IntegerVariable(1),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("name"), knownvalue.StringExact(routerName)),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("replicas_per_cell"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("state"), knownvalue.StringExact("ready")),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":      config.StringVariable(testAccOrg),
					"database_name":     config.StringVariable(databaseName),
					"router_name":       config.StringVariable(routerName),
					"replicas_per_cell": config.IntegerVariable(2),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("replicas_per_cell"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(resourceAddress, tfjsonpath.New("state"), knownvalue.StringExact("ready")),
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"organization":      config.StringVariable(testAccOrg),
					"database_name":     config.StringVariable(databaseName),
					"router_name":       config.StringVariable(routerName),
					"replicas_per_cell": config.IntegerVariable(2),
				},
				ResourceName: resourceAddress,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs := s.RootModule().Resources[resourceAddress]
					jsonBytes, err := json.Marshal(map[string]string{
						"branch":       rs.Primary.Attributes["branch"],
						"database":     rs.Primary.Attributes["database"],
						"name":         rs.Primary.Attributes["name"],
						"organization": rs.Primary.Attributes["organization"],
					})
					return string(jsonBytes), err
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"router_size"},
			},
		},
	})
}
