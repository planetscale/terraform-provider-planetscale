package mapplanmodifier

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

var parametersElemType = types.MapType{ElemType: types.StringType}

func parametersValue(t *testing.T, params map[string]map[string]string) types.Map {
	t.Helper()

	outer := map[string]attr.Value{}
	for namespace, kv := range params {
		inner := map[string]attr.Value{}
		for k, v := range kv {
			inner[k] = types.StringValue(v)
		}
		outer[namespace] = types.MapValueMust(types.StringType, inner)
	}
	return types.MapValueMust(parametersElemType, outer)
}

func TestWarnOnRemovedParameters(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		state        types.Map
		config       types.Map
		wantInDetail []string
	}{
		"create with null state": {
			state:  types.MapNull(parametersElemType),
			config: parametersValue(t, map[string]map[string]string{"pgconf": {"max_connections": "200"}}),
		},
		"no changes": {
			state:  parametersValue(t, map[string]map[string]string{"pgconf": {"max_connections": "200"}}),
			config: parametersValue(t, map[string]map[string]string{"pgconf": {"max_connections": "200"}}),
		},
		"value changed": {
			state:  parametersValue(t, map[string]map[string]string{"pgconf": {"max_connections": "200"}}),
			config: parametersValue(t, map[string]map[string]string{"pgconf": {"max_connections": "300"}}),
		},
		"one key removed": {
			state: parametersValue(t, map[string]map[string]string{
				"pgconf": {"max_connections": "200", "statement_timeout": "5s"},
			}),
			config: parametersValue(t, map[string]map[string]string{
				"pgconf": {"max_connections": "200"},
			}),
			wantInDetail: []string{`pgconf.statement_timeout (currently "5s")`},
		},
		"whole namespace removed": {
			state: parametersValue(t, map[string]map[string]string{
				"pgconf":    {"max_connections": "200"},
				"pgbouncer": {"pool_mode": "transaction"},
			}),
			config: parametersValue(t, map[string]map[string]string{
				"pgconf": {"max_connections": "200"},
			}),
			wantInDetail: []string{`pgbouncer.pool_mode (currently "transaction")`},
		},
		"removals across namespaces are sorted": {
			state: parametersValue(t, map[string]map[string]string{
				"pgconf":  {"max_connections": "200"},
				"patroni": {"loop_wait": "10"},
			}),
			config: parametersValue(t, map[string]map[string]string{}),
			wantInDetail: []string{
				"- patroni.loop_wait (currently \"10\")\n  - pgconf.max_connections (currently \"200\")",
			},
		},
		"key added only": {
			state: parametersValue(t, map[string]map[string]string{
				"pgconf": {"max_connections": "200"},
			}),
			config: parametersValue(t, map[string]map[string]string{
				"pgconf": {"max_connections": "200", "statement_timeout": "5s"},
			}),
		},
		"whole attribute removed": {
			state:  parametersValue(t, map[string]map[string]string{"pgconf": {"max_connections": "200"}}),
			config: types.MapNull(parametersElemType),
		},
		"unknown config": {
			state:  parametersValue(t, map[string]map[string]string{"pgconf": {"max_connections": "200"}}),
			config: types.MapUnknown(parametersElemType),
		},
		"unknown namespace in config": {
			state: parametersValue(t, map[string]map[string]string{"pgconf": {"max_connections": "200"}}),
			config: types.MapValueMust(parametersElemType, map[string]attr.Value{
				"pgconf": types.MapUnknown(types.StringType),
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.MapRequest{
				Path:        path.Root("parameters"),
				StateValue:  tc.state,
				ConfigValue: tc.config,
				PlanValue:   tc.config,
			}
			resp := &planmodifier.MapResponse{PlanValue: req.PlanValue}
			WarnOnRemovedParameters().PlanModifyMap(context.Background(), req, resp)

			require.Equal(t, req.PlanValue, resp.PlanValue)
			require.False(t, resp.Diagnostics.HasError())

			if len(tc.wantInDetail) == 0 {
				require.Empty(t, resp.Diagnostics)
				return
			}

			require.Len(t, resp.Diagnostics, 1)
			d := resp.Diagnostics[0]
			require.Equal(t, diag.SeverityWarning, d.Severity())
			require.Contains(t, d.Summary(), "reset to their defaults")
			for _, want := range tc.wantInDetail {
				require.Contains(t, d.Detail(), want)
			}
		})
	}
}
