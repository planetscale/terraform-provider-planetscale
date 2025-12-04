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

func TestAccBranchSchemaDataSource(t *testing.T) {
	t.Parallel()
	t.Skipf("TODO: 404")

	name := fmt.Sprintf("testacc-%d", time.Now().Unix())
	resourceAddress := "data.planetscale_branch_schema.test"

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
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"html": knownvalue.NotNull(),
								"name": knownvalue.NotNull(),
								"raw":  knownvalue.NotNull(),
							}),
						}),
					),
				},
			},
		},
	})
}
