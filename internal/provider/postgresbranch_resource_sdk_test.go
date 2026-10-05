package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/operations"
	"github.com/stretchr/testify/require"
)

func TestPostgresBranchExtensionsRequestAndRefresh(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	for _, tc := range []struct {
		name       string
		extensions []types.String
		present    bool
	}{
		{name: "unmanaged"},
		{name: "disable", extensions: []types.String{}, present: true},
		{name: "enable", extensions: []types.String{types.StringValue("hll")}, present: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			model := &PostgresBranchResourceModel{Extensions: tc.extensions}
			body, diags := model.ToOperationsApplyPostgresBranchTerraformChangesRequestBody(ctx, nil)
			require.False(t, diags.HasError(), diags.Errors())
			encoded, err := json.Marshal(body)
			require.NoError(t, err)
			var attributes map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(encoded, &attributes))
			value, present := attributes["extensions"]
			require.Equal(t, tc.present, present)
			if tc.present {
				require.Equal(t, len(tc.extensions), len(body.Extensions))
				if len(tc.extensions) == 0 {
					require.JSONEq(t, `[]`, string(value))
				}
			}
		})
	}

	model := &PostgresBranchResourceModel{Extensions: []types.String{types.StringValue("hll")}}
	diags := model.RefreshFromOperationsGetPostgresBranchManagedExtensionsResponseBody(ctx, &operations.GetPostgresBranchManagedExtensionsResponseBody{})
	require.False(t, diags.HasError(), diags.Errors())
	require.Nil(t, model.Extensions)

	diags = model.RefreshFromOperationsGetPostgresBranchManagedExtensionsResponseBody(ctx,
		&operations.GetPostgresBranchManagedExtensionsResponseBody{Extensions: []string{}})
	require.False(t, diags.HasError(), diags.Errors())
	require.Equal(t, []types.String{}, model.Extensions)
}
