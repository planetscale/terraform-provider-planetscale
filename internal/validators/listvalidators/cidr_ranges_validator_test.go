package listvalidators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCidrRangesValidator(t *testing.T) {
	t.Parallel()

	list := func(values ...attr.Value) types.List {
		return types.ListValueMust(types.StringType, values)
	}

	tests := map[string]struct {
		value      types.List
		wantErrors []string
	}{
		"ranges and single addresses": {
			value: list(types.StringValue("203.0.113.0/24"), types.StringValue("198.51.100.7/32"), types.StringValue("0.0.0.0/0")),
		},
		"null list": {
			value: types.ListNull(types.StringType),
		},
		"unknown list": {
			value: types.ListUnknown(types.StringType),
		},
		"unknown element": {
			value: list(types.StringUnknown(), types.StringValue("203.0.113.0/24")),
		},
		"bare address": {
			value:      list(types.StringValue("198.51.100.7")),
			wantErrors: []string{`"198.51.100.7" is a single address. Write it as a range: "198.51.100.7/32".`},
		},
		"host bits set": {
			value:      list(types.StringValue("10.0.0.5/24")),
			wantErrors: []string{`"10.0.0.5/24" does not start at the first address of its range. Use "10.0.0.0/24" instead.`},
		},
		"IPv6": {
			value:      list(types.StringValue("2001:db8::/32"), types.StringValue("2001:db8::1")),
			wantErrors: []string{`"2001:db8::/32" is not an IPv4 range`, `"2001:db8::1" is not an IPv4 range`},
		},
		"not an address": {
			value:      list(types.StringValue("203.0.113.0/24"), types.StringValue("office")),
			wantErrors: []string{`"office" is not an IPv4 range`},
		},
		"leading zeros": {
			value:      list(types.StringValue("010.0.0.0/8")),
			wantErrors: []string{`"010.0.0.0/8" is not an IPv4 range`},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validator.ListResponse{}
			CidrRangesValidator().ValidateList(context.Background(), validator.ListRequest{
				Path:        path.Root("cidrs"),
				ConfigValue: tt.value,
			}, resp)

			require.Len(t, resp.Diagnostics.Errors(), len(tt.wantErrors))
			for i, want := range tt.wantErrors {
				require.Contains(t, resp.Diagnostics.Errors()[i].Detail(), want)
			}
		})
	}
}
