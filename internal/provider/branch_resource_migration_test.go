package provider

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestUpgradeBranchStateV0toCurrent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	v0Type, diags := schemaToObjectTFType(ctx, branchSchemaV0())
	require.False(t, diags.HasError(), "Failed to convert V0 schema to tftype")

	v1Type, diags := schemaToObjectTFType(ctx, branchSchemaV1())
	require.False(t, diags.HasError(), "Failed to convert V1 schema to tftype")

	req := resource.UpgradeStateRequest{
		State: &tfsdk.State{
			Raw: tftypes.NewValue(v0Type,
				map[string]tftypes.Value{
					"organization":                   tftypes.NewValue(tftypes.String, "test-org"),
					"database":                       tftypes.NewValue(tftypes.String, "test-db"),
					"name":                           tftypes.NewValue(tftypes.String, "test-name"),
					"parent_branch":                  tftypes.NewValue(tftypes.String, "test-parent-branch"),
					"actor":                          tftypes.NewValue(v0Type.AttributeTypes["actor"], nil),
					"cluster_rate_name":              tftypes.NewValue(tftypes.String, "test-cluster-rate-name"),
					"created_at":                     tftypes.NewValue(tftypes.String, "2023-01-01T00:00:00Z"),
					"html_url":                       tftypes.NewValue(tftypes.String, "https://planetscale.com/test-org/test-db/test-branch"),
					"id":                             tftypes.NewValue(tftypes.String, "test-id"),
					"mysql_address":                  tftypes.NewValue(tftypes.String, "test-mysql-address"),
					"mysql_edge_address":             tftypes.NewValue(tftypes.String, "test-mysql-edge-address"),
					"region":                         tftypes.NewValue(v0Type.AttributeTypes["region"], nil),
					"production":                     tftypes.NewValue(tftypes.Bool, true),
					"ready":                          tftypes.NewValue(tftypes.Bool, true),
					"restore_checklist_completed_at": tftypes.NewValue(tftypes.String, "2023-01-01T01:00:00Z"),
					"restored_from_branch":           tftypes.NewValue(v0Type.AttributeTypes["restored_from_branch"], nil),
					"schema_last_updated_at":         tftypes.NewValue(tftypes.String, "2023-01-01T02:00:00Z"),
					"shard_count":                    tftypes.NewValue(tftypes.Number, 1),
					"sharded":                        tftypes.NewValue(tftypes.Bool, false),
					"updated_at":                     tftypes.NewValue(tftypes.String, "2023-01-01T03:00:00Z"),
				}),
			Schema: branchSchemaV0(),
		},
	}

	resp := &resource.UpgradeStateResponse{
		State: tfsdk.State{
			Schema: branchSchemaV1(),
		},
	}

	// expectedResp is identical to req, but schema and types are now V1 and
	// cluster_rate_name is omitted.
	expectedResp := &resource.UpgradeStateResponse{
		State: tfsdk.State{
			Raw: tftypes.NewValue(v1Type,
				map[string]tftypes.Value{
					"organization":                   tftypes.NewValue(tftypes.String, "test-org"),
					"database":                       tftypes.NewValue(tftypes.String, "test-db"),
					"name":                           tftypes.NewValue(tftypes.String, "test-name"),
					"parent_branch":                  tftypes.NewValue(tftypes.String, "test-parent-branch"),
					"actor":                          tftypes.NewValue(v1Type.AttributeTypes["actor"], nil),
					"created_at":                     tftypes.NewValue(tftypes.String, "2023-01-01T00:00:00Z"),
					"html_url":                       tftypes.NewValue(tftypes.String, "https://planetscale.com/test-org/test-db/test-branch"),
					"id":                             tftypes.NewValue(tftypes.String, "test-id"),
					"mysql_address":                  tftypes.NewValue(tftypes.String, "test-mysql-address"),
					"mysql_edge_address":             tftypes.NewValue(tftypes.String, "test-mysql-edge-address"),
					"region":                         tftypes.NewValue(v1Type.AttributeTypes["region"], nil),
					"production":                     tftypes.NewValue(tftypes.Bool, true),
					"ready":                          tftypes.NewValue(tftypes.Bool, true),
					"restore_checklist_completed_at": tftypes.NewValue(tftypes.String, "2023-01-01T01:00:00Z"),
					"restored_from_branch":           tftypes.NewValue(v1Type.AttributeTypes["restored_from_branch"], nil),
					"schema_last_updated_at":         tftypes.NewValue(tftypes.String, "2023-01-01T02:00:00Z"),
					"shard_count":                    tftypes.NewValue(tftypes.Number, 1),
					"sharded":                        tftypes.NewValue(tftypes.Bool, false),
					"updated_at":                     tftypes.NewValue(tftypes.String, "2023-01-01T03:00:00Z"),
				}),
			Schema: branchSchemaV1(),
		},
	}

	upgradeBranchStateV0toCurrent(ctx, req, resp)

	require.False(t, resp.Diagnostics.HasError(), "Upgrade function reported errors: %v", resp.Diagnostics)

	if !cmp.Equal(expectedResp, resp) {
		t.Errorf("expected: %+v, got: %+v", expectedResp, resp)
	}
}
