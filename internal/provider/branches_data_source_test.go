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

func TestAccBranchesDataSource(t *testing.T) {
	t.Parallel()

	databaseName := fmt.Sprintf("terraform-testing-%d", rand.Intn(1000000))
	resourceAddress := "data.planetscale_branches.test"

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
						knownvalue.ListSizeExact(2),
					),
					statecheck.ExpectKnownValue(
						resourceAddress,
						tfjsonpath.New("data").AtSliceIndex(0),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"actor":                          knownvalue.NotNull(),
							"cluster_iops":                   knownvalue.NotNull(),
							"cluster_name":                   knownvalue.NotNull(),
							"created_at":                     knownvalue.NotNull(),
							"deleted_at":                     knownvalue.NotNull(),
							"direct_vtgate":                  knownvalue.NotNull(),
							"has_read_only_replicas":         knownvalue.NotNull(),
							"has_replicas":                   knownvalue.NotNull(),
							"html_url":                       knownvalue.NotNull(),
							"id":                             knownvalue.NotNull(),
							"kind":                           knownvalue.NotNull(),
							"metal":                          knownvalue.NotNull(),
							"mysql_address":                  knownvalue.NotNull(),
							"mysql_edge_address":             knownvalue.NotNull(),
							"name":                           knownvalue.NotNull(),
							"parent_branch":                  knownvalue.NotNull(),
							"private_edge_connectivity":      knownvalue.NotNull(),
							"production":                     knownvalue.NotNull(),
							"ready":                          knownvalue.NotNull(),
							"region":                         knownvalue.NotNull(),
							"restore_checklist_completed_at": knownvalue.NotNull(),
							"restored_from_branch":           knownvalue.NotNull(),
							"safe_migrations":                knownvalue.NotNull(),
							"schema_last_updated_at":         knownvalue.NotNull(),
							"schema_ready":                   knownvalue.NotNull(),
							"shard_count":                    knownvalue.NotNull(),
							"sharded":                        knownvalue.NotNull(),
							"stale_schema":                   knownvalue.NotNull(),
							"state":                          knownvalue.NotNull(),
							"updated_at":                     knownvalue.NotNull(),
							"url":                            knownvalue.NotNull(),
							"vtgate_count":                   knownvalue.NotNull(),
							"vtgate_size":                    knownvalue.NotNull(),
						}),
					),
				},
			},
		},
	})
}
