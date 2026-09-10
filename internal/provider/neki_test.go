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

// TestAccNeki_Lifecycle is the single Neki acceptance test. Creating a Neki
// branch provisions a whole cluster (default configuration profile, shard,
// router group, admin, and sidecar) and takes minutes, so the test creates one
// branch and covers everything else against it: role and backup policy
// create, update, and import, plus every Neki data source reading either the
// created objects or the provisioned defaults.
//
// The admin, sidecar, router, shard, and configuration profile resources are
// intentionally not exercised: their create and update paths wait on cluster
// provisioning or change requests. Their read paths are covered by the data
// sources here.
func TestAccNeki_Lifecycle(t *testing.T) {
	t.Parallel()

	databaseName := randomWithPrefix("testacc-neki")
	branchNameOriginal := "main"
	branchNameRenamed := randomWithPrefix("main-renamed")
	roleNameOriginal := randomWithPrefix("test-role")
	roleNameRenamed := randomWithPrefix("test-role-renamed")
	policyName := randomWithPrefix("test-backup-policy")
	policyNameUpdated := policyName + "-updated"

	branchAddress := "planetscale_neki_branch.test"
	roleAddress := "planetscale_neki_role.test"
	policyAddress := "planetscale_neki_backup_policy.test"

	originalVariables := config.Variables{
		"organization":    config.StringVariable(testAccOrg),
		"database_name":   config.StringVariable(databaseName),
		"branch_name":     config.StringVariable(branchNameOriginal),
		"role_name":       config.StringVariable(roleNameOriginal),
		"policy_name":     config.StringVariable(policyName),
		"target":          config.StringVariable("development"),
		"retention_value": config.IntegerVariable(7),
		"retention_unit":  config.StringVariable("day"),
		"frequency_value": config.IntegerVariable(1),
		"frequency_unit":  config.StringVariable("week"),
		"schedule_time":   config.StringVariable("04:00"),
		"schedule_day":    config.IntegerVariable(1),
		"schedule_week":   config.IntegerVariable(0),
	}
	updatedVariables := config.Variables{
		"organization":    config.StringVariable(testAccOrg),
		"database_name":   config.StringVariable(databaseName),
		"branch_name":     config.StringVariable(branchNameRenamed),
		"role_name":       config.StringVariable(roleNameRenamed),
		"policy_name":     config.StringVariable(policyNameUpdated),
		"target":          config.StringVariable("production"),
		"retention_value": config.IntegerVariable(14),
		"retention_unit":  config.StringVariable("week"),
		"frequency_value": config.IntegerVariable(2),
		"frequency_unit":  config.StringVariable("month"),
		"schedule_time":   config.StringVariable("05:30"),
		"schedule_day":    config.IntegerVariable(2),
		"schedule_week":   config.IntegerVariable(1),
	}

	// importIDFunc builds the JSON import ID from the named attributes of the
	// resource's current state.
	importIDFunc := func(address string, keys ...string) resource.ImportStateIdFunc {
		return func(s *terraform.State) (string, error) {
			rs := s.RootModule().Resources[address]
			id := make(map[string]string, len(keys))
			for _, key := range keys {
				id[key] = rs.Primary.Attributes[key]
			}
			jsonBytes, err := json.Marshal(id)
			return string(jsonBytes), err
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: originalVariables,
				ConfigStateChecks: []statecheck.StateCheck{
					// Created resources.
					statecheck.ExpectKnownValue(branchAddress, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(branchAddress, tfjsonpath.New("name"), knownvalue.StringExact(branchNameOriginal)),
					statecheck.ExpectKnownValue(branchAddress, tfjsonpath.New("state"), knownvalue.StringExact("ready")),
					statecheck.ExpectKnownValue(branchAddress, tfjsonpath.New("region"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(roleAddress, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(roleAddress, tfjsonpath.New("name"), knownvalue.StringExact(roleNameOriginal)),
					statecheck.ExpectKnownValue(roleAddress, tfjsonpath.New("username"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(roleAddress, tfjsonpath.New("password"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(roleAddress, tfjsonpath.New("access_host_url"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("name"), knownvalue.StringExact(policyName)),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("target"), knownvalue.StringExact("development")),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("retention_value"), knownvalue.Int64Exact(7)),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("required"), knownvalue.Bool(false)),

					// Data sources reading the created resources.
					statecheck.ExpectKnownValue("data.planetscale_neki_branch.test", tfjsonpath.New("name"), knownvalue.StringExact(branchNameOriginal)),
					statecheck.ExpectKnownValue("data.planetscale_neki_branch.test", tfjsonpath.New("state"), knownvalue.StringExact("ready")),
					statecheck.ExpectKnownValue("data.planetscale_neki_role.test", tfjsonpath.New("name"), knownvalue.StringExact(roleNameOriginal)),
					statecheck.ExpectKnownValue("data.planetscale_neki_role.test", tfjsonpath.New("username"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_roles.test", tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(roleNameOriginal)),
					statecheck.ExpectKnownValue("data.planetscale_neki_backup_policy.test", tfjsonpath.New("name"), knownvalue.StringExact(policyName)),
					statecheck.ExpectKnownValue("data.planetscale_neki_backup_policy.test", tfjsonpath.New("retention_value"), knownvalue.Int64Exact(7)),
					statecheck.ExpectKnownValue("data.planetscale_neki_backup_policies.test", tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),

					// Data sources reading the defaults that provisioning created.
					statecheck.ExpectKnownValue("data.planetscale_neki_admin.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_admin.test", tfjsonpath.New("admin_size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_sidecar.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_sidecar.test", tfjsonpath.New("configuration_profile"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("data.planetscale_neki_configuration_profile.test", tfjsonpath.New("name"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("data.planetscale_neki_configuration_profile.test", tfjsonpath.New("is_default"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.planetscale_neki_configuration_profile.test", tfjsonpath.New("cluster_size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_configuration_profiles.test", tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_shards.test", tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_shard.test", tfjsonpath.New("configuration_profile"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_router.test", tfjsonpath.New("name"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue("data.planetscale_neki_router.test", tfjsonpath.New("default"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.planetscale_neki_router.test", tfjsonpath.New("router_size"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.planetscale_neki_routers.test", tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
				},
			},
			{
				// Update in place: rename the branch and role, change every
				// backup policy setting. The role's branch attribute requires
				// replacement, so the rename also exercises role replacement.
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: updatedVariables,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(branchAddress, tfjsonpath.New("name"), knownvalue.StringExact(branchNameRenamed)),
					statecheck.ExpectKnownValue(branchAddress, tfjsonpath.New("state"), knownvalue.StringExact("ready")),
					statecheck.ExpectKnownValue(roleAddress, tfjsonpath.New("name"), knownvalue.StringExact(roleNameRenamed)),
					statecheck.ExpectKnownValue(roleAddress, tfjsonpath.New("branch"), knownvalue.StringExact(branchNameRenamed)),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("name"), knownvalue.StringExact(policyNameUpdated)),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("target"), knownvalue.StringExact("production")),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("retention_value"), knownvalue.Int64Exact(14)),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("retention_unit"), knownvalue.StringExact("week")),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("frequency_value"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("frequency_unit"), knownvalue.StringExact("month")),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("schedule_time"), knownvalue.StringExact("05:30")),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("schedule_day"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(policyAddress, tfjsonpath.New("schedule_week"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.planetscale_neki_branch.test", tfjsonpath.New("name"), knownvalue.StringExact(branchNameRenamed)),
					statecheck.ExpectKnownValue("data.planetscale_neki_backup_policy.test", tfjsonpath.New("name"), knownvalue.StringExact(policyNameUpdated)),
				},
			},
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   updatedVariables,
				ResourceName:      branchAddress,
				ImportState:       true,
				ImportStateIdFunc: importIDFunc(branchAddress, "organization", "database", "id"),
				ImportStateVerify: true,
			},
			{
				ConfigDirectory:         config.TestNameDirectory(),
				ConfigVariables:         updatedVariables,
				ResourceName:            roleAddress,
				ImportState:             true,
				ImportStateIdFunc:       importIDFunc(roleAddress, "organization", "database", "branch", "id"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   updatedVariables,
				ResourceName:      policyAddress,
				ImportState:       true,
				ImportStateIdFunc: importIDFunc(policyAddress, "organization", "database", "id"),
				ImportStateVerify: true,
			},
		},
	})
}
