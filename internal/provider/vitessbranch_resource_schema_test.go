package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestVitessBranchResource_ClusterSizeValidation(t *testing.T) {
	t.Parallel()

	r := NewVitessBranchResource()
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)

	attr, ok := schemaResp.Schema.Attributes["cluster_size"]
	require.True(t, ok)

	clusterSizeAttr, ok := attr.(schema.StringAttribute)
	require.True(t, ok)

	require.NotEmpty(t, clusterSizeAttr.Validators)

	testCases := []struct {
		value string
		valid bool
	}{
		{value: "PS_5_AWS_X86", valid: true},
		{value: "PS_5_GCP_ARM", valid: true},
		{value: "PS_DEV_AWS_X86", valid: true},
		{value: "M1_10_AWS_AMD_D_METAL_10", valid: true},
		{value: "M4_160_D_METAL_460", valid: true},
		{value: "M_160_GCP_AMD_D_METAL_742", valid: true},
		{value: "M_160_GCP_INTEL_D_METAL_742", valid: true},
		{value: "PS_AWS_R6I_2XLARGE", valid: true},
		{value: "PS_1400", valid: true},
		{value: "PS_DEV", valid: true},
		{value: "PS_99", valid: false},
		{value: "PS-5", valid: false},
		{value: "PS-5-AWS-X86", valid: false},
		{value: "PS_5_GCP", valid: false},
		{value: "PS_5_AMAZON_X86", valid: false},
		{value: "PS_DEV_AWS_x86", valid: false},
		{value: "M_10_AWS_ARM_D_METAL_10", valid: false},
		{value: "M1_10_AWS_ARM", valid: false},
	}

	for _, tc := range testCases {
		t.Run(tc.value, func(t *testing.T) {
			req := validator.StringRequest{
				Path:        path.Root("cluster_size"),
				ConfigValue: types.StringValue(tc.value),
			}
			var resp validator.StringResponse
			for _, v := range clusterSizeAttr.Validators {
				v.ValidateString(context.Background(), req, &resp)
			}

			hasErr := resp.Diagnostics.HasError()
			if tc.valid {
				require.False(t, hasErr)
				return
			}

			require.True(t, hasErr)
			errs := resp.Diagnostics.Errors()
			require.NotEmpty(t, errs)
		})
	}
}

func TestVitessBranchResource_BranchSettingsAreConfigurable(t *testing.T) {
	t.Parallel()

	r := NewVitessBranchResource()
	var schemaResp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)

	boolAttributes := []string{"safe_migrations", "vtgate_autoscaling"}
	for _, name := range boolAttributes {
		attr, ok := schemaResp.Schema.Attributes[name].(schema.BoolAttribute)
		require.True(t, ok, name)
		require.True(t, attr.Optional, name)
		require.True(t, attr.Computed, name)
	}

	intAttributes := []string{"vtgate_count", "vtgate_max_count", "vtgate_target_cpu_utilization"}
	for _, name := range intAttributes {
		attr, ok := schemaResp.Schema.Attributes[name].(schema.Int64Attribute)
		require.True(t, ok, name)
		require.True(t, attr.Optional, name)
		require.True(t, attr.Computed, name)
	}

	vtgateSize, ok := schemaResp.Schema.Attributes["vtgate_size"].(schema.StringAttribute)
	require.True(t, ok)
	require.True(t, vtgateSize.Optional)
	require.True(t, vtgateSize.Computed)
}
