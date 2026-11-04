package mapvalidators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func vitessParametersValue(t *testing.T, namespaces map[string]attr.Value) types.Map {
	t.Helper()

	value, diags := types.MapValue(types.MapType{ElemType: types.StringType}, namespaces)
	require.False(t, diags.HasError(), diags)
	return value
}

func vitessNamespaceValue(t *testing.T, parameters map[string]attr.Value) types.Map {
	t.Helper()

	value, diags := types.MapValue(types.StringType, parameters)
	require.False(t, diags.HasError(), diags)
	return value
}

func TestVitessParametersValidator(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value      types.Map
		wantErrors []string
	}{
		"supported namespaces with values": {
			value: vitessParametersValue(t, map[string]attr.Value{
				"vttablet": vitessNamespaceValue(t, map[string]attr.Value{"queryserver-config-pool-size": types.StringValue("50")}),
				"mysqld":   vitessNamespaceValue(t, map[string]attr.Value{"wait_timeout": types.StringUnknown()}),
			}),
		},
		"unknown parameters": {
			value: types.MapUnknown(types.MapType{ElemType: types.StringType}),
		},
		"unknown namespace": {
			value: vitessParametersValue(t, map[string]attr.Value{
				"mysqld": types.MapUnknown(types.StringType),
			}),
		},
		"unsupported namespace": {
			value: vitessParametersValue(t, map[string]attr.Value{
				"vtgate": vitessNamespaceValue(t, map[string]attr.Value{"query-timeout": types.StringValue("30000")}),
			}),
			wantErrors: []string{"Unsupported parameter namespace"},
		},
		"empty namespace": {
			value: vitessParametersValue(t, map[string]attr.Value{
				"mysqld": vitessNamespaceValue(t, map[string]attr.Value{}),
			}),
			wantErrors: []string{"Empty parameter namespace"},
		},
		"null namespace": {
			value: vitessParametersValue(t, map[string]attr.Value{
				"mysqld": types.MapNull(types.StringType),
			}),
			wantErrors: []string{"Empty parameter namespace"},
		},
		"null parameter value": {
			value: vitessParametersValue(t, map[string]attr.Value{
				"mysqld": vitessNamespaceValue(t, map[string]attr.Value{"wait_timeout": types.StringNull()}),
			}),
			wantErrors: []string{"Null parameter value"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validator.MapResponse{}
			VitessKeyspaceParametersValidator().ValidateMap(context.Background(), validator.MapRequest{
				Path:        path.Root("parameters"),
				ConfigValue: tt.value,
			}, resp)

			var summaries []string
			for _, d := range resp.Diagnostics.Errors() {
				summaries = append(summaries, d.Summary())
			}
			require.Equal(t, tt.wantErrors, summaries)
		})
	}
}
