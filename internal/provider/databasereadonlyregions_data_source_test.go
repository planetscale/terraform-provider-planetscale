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

func TestAccDatabaseReadOnlyRegionsDataSource(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: Data is empty from API -- potentially needs updated Terraform configuration")

	name := fmt.Sprintf("testacc-%d", time.Now().Unix())
	resourceAddress := "data.planetscale_database_read_only_regions.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"name": config.StringVariable(name),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("data"),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"actor":        knownvalue.NotNull(),
								"created_at":   knownvalue.NotNull(),
								"display_name": knownvalue.NotNull(),
								"id":           knownvalue.NotNull(),
								"ready":        knownvalue.NotNull(),
								"ready_at":     knownvalue.NotNull(),
								"region":       knownvalue.NotNull(),
								"updated_at":   knownvalue.NotNull(),
							}),
						}),
					),
				},
			},
		},
	})
}
