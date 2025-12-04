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

func TestAccPasswordsDataSource(t *testing.T) {
	t.Parallel()

	databaseName := fmt.Sprintf("testacc-%d", rand.Intn(1000000))
	resourceAddress := "data.planetscale_passwords.test"

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
								"access_host_regional_url":  knownvalue.NotNull(),
								"access_host_regional_urls": knownvalue.NotNull(),
								"access_host_url":           knownvalue.NotNull(),
								"actor":                     knownvalue.NotNull(),
								"cidrs":                     knownvalue.NotNull(),
								"created_at":                knownvalue.NotNull(),
								"database_branch":           knownvalue.NotNull(),
								"deleted_at":                knownvalue.NotNull(),
								"direct_vtgate":             knownvalue.NotNull(),
								"expired":                   knownvalue.NotNull(),
								"expires_at":                knownvalue.NotNull(),
								"id":                        knownvalue.NotNull(),
								"last_used_at":              knownvalue.NotNull(),
								"name":                      knownvalue.NotNull(),
								"plain_text":                knownvalue.NotNull(),
								"region":                    knownvalue.NotNull(),
								"renewable":                 knownvalue.NotNull(),
								"replica":                   knownvalue.NotNull(),
								"role":                      knownvalue.NotNull(),
								"ttl_seconds":               knownvalue.NotNull(),
								"username":                  knownvalue.NotNull(),
							}),
						}),
					),
				},
			},
		},
	})
}
