package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPostgresBranchReplicasOmittedWhenUnsetAndSentWhenZero(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	config := &PostgresBranchResourceModel{
		BackupID:     types.StringNull(),
		RestorePoint: types.StringNull(),
	}
	opts := &PostgresBranchResourceModelOptions{Config: config}

	unset := &PostgresBranchResourceModel{
		Name:        types.StringValue("main"),
		ClusterSize: types.StringValue("PS_10_AWS_ARM"),
		Replicas:    types.Int64Null(),
	}

	createUnset, diags := unset.ToOperationsCreatePostgresBranchRequestBody(ctx, opts)
	require.False(t, diags.HasError())
	createUnsetJSON, err := json.Marshal(createUnset)
	require.NoError(t, err)
	require.NotContains(t, string(createUnsetJSON), `"replicas"`)

	applyUnset, diags := unset.ToOperationsApplyPostgresBranchChangesRequestBody(ctx, opts)
	require.False(t, diags.HasError())
	applyUnsetJSON, err := json.Marshal(applyUnset)
	require.NoError(t, err)
	require.NotContains(t, string(applyUnsetJSON), `"replicas"`)

	zero := &PostgresBranchResourceModel{
		Name:        types.StringValue("main"),
		ClusterSize: types.StringValue("PS_10_AWS_ARM"),
		Replicas:    types.Int64Value(0),
	}

	createZero, diags := zero.ToOperationsCreatePostgresBranchRequestBody(ctx, opts)
	require.False(t, diags.HasError())
	createZeroJSON, err := json.Marshal(createZero)
	require.NoError(t, err)
	require.Contains(t, string(createZeroJSON), `"replicas":0`)

	applyZero, diags := zero.ToOperationsApplyPostgresBranchChangesRequestBody(ctx, opts)
	require.False(t, diags.HasError())
	applyZeroJSON, err := json.Marshal(applyZero)
	require.NoError(t, err)
	require.Contains(t, string(applyZeroJSON), `"replicas":0`)
}
