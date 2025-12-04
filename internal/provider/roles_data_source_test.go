package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccRolesDataSource(t *testing.T) {
	t.Parallel()

	databaseName := fmt.Sprintf("testacc-%d", rand.Intn(1000000))
	resourceAddress := "data.planetscale_roles.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: config.Variables{
					"database_name": config.StringVariable(databaseName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"access_host_url":       knownvalue.NotNull(),
								"created_at":            knownvalue.NotNull(),
								"database_name":         knownvalue.NotNull(),
								"default":               knownvalue.NotNull(),
								"expired":               knownvalue.NotNull(),
								"id":                    knownvalue.NotNull(),
								"name":                  knownvalue.NotNull(),
								"ttl":                   knownvalue.NotNull(),
								"updated_at":            knownvalue.NotNull(),
								"username":              knownvalue.NotNull(),
							}),
						}),
					),
				},
			},
		},
	})
}
