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

func TestPostgresBranchResource_ClusterSizeValidation(t *testing.T) {
	t.Parallel()

	r := NewPostgresBranchResource()
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
