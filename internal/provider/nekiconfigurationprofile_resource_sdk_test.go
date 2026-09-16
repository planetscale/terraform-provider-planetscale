package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/operations"
	"github.com/stretchr/testify/require"
)

func TestNekiConfigurationProfileUpdateOmitsUnmanagedExtensions(t *testing.T) {
	t.Parallel()
	model := &NekiConfigurationProfileResourceModel{}

	body, diags := model.ToOperationsUpdateNekiConfigurationProfileRequestBody(context.Background())

	require.False(t, diags.HasError(), diags.Errors())
	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	var attributes map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(encoded, &attributes))
	require.NotContains(t, attributes, "extensions")
}

func TestNekiConfigurationProfilePreservesEmptyExtensions(t *testing.T) {
	t.Parallel()
	model := &NekiConfigurationProfileResourceModel{Extensions: []types.String{}}

	body, diags := model.ToOperationsUpdateNekiConfigurationProfileRequestBody(context.Background())

	require.False(t, diags.HasError(), diags.Errors())
	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	var attributes map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(encoded, &attributes))
	require.JSONEq(t, "[]", string(attributes["extensions"]))

	diags = model.RefreshFromOperationsGetNekiConfigurationProfileParametersResponseBody(
		context.Background(), &operations.GetNekiConfigurationProfileParametersResponseBody{Extensions: []string{}},
	)

	require.False(t, diags.HasError(), diags.Errors())
	require.Equal(t, []types.String{}, model.Extensions)
}

func TestNekiConfigurationProfileReadOmitsUnmanagedExtensions(t *testing.T) {
	t.Parallel()
	model := &NekiConfigurationProfileResourceModel{}

	parameters, diags := model.ToOperationsGetNekiConfigurationProfileParametersRequest(context.Background())

	require.False(t, diags.HasError(), diags.Errors())
	require.Nil(t, parameters.Extensions)
}

func TestNekiConfigurationProfileRefreshPreservesNullExtensions(t *testing.T) {
	t.Parallel()
	model := &NekiConfigurationProfileResourceModel{Extensions: []types.String{types.StringValue("hll")}}

	diags := model.RefreshFromOperationsGetNekiConfigurationProfileParametersResponseBody(
		context.Background(), &operations.GetNekiConfigurationProfileParametersResponseBody{},
	)

	require.False(t, diags.HasError(), diags.Errors())
	require.Nil(t, model.Extensions)
}

func TestNekiConfigurationProfileRefreshUpdatesManagedExtensions(t *testing.T) {
	t.Parallel()
	model := &NekiConfigurationProfileResourceModel{Extensions: []types.String{types.StringValue("hll")}}

	diags := model.RefreshFromOperationsGetNekiConfigurationProfileParametersResponseBody(
		context.Background(), &operations.GetNekiConfigurationProfileParametersResponseBody{Extensions: []string{"auto_explain", "hll"}},
	)

	require.False(t, diags.HasError(), diags.Errors())
	require.Equal(t, []types.String{types.StringValue("auto_explain"), types.StringValue("hll")}, model.Extensions)
}
