package stringvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = StringClusterSizeValidatorValidator{}

var clusterSizePattern = regexp.MustCompile(`^(?:PS_(?:[0-9]+|DEV)_(?:AWS|GCP)_(?:ARM|X86|AMD)|M[0-9]+_[0-9]+_(?:AWS|GCP)_(?:ARM|X86|AMD)_D_METAL_[0-9]+)$`)

type StringClusterSizeValidatorValidator struct{}

// Description describes the validation in plain text formatting.
func (v StringClusterSizeValidatorValidator) Description(_ context.Context) string {
	return "value must be a valid cluster_size SKU"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v StringClusterSizeValidatorValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v StringClusterSizeValidatorValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if clusterSizePattern.MatchString(value) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid cluster_size value",
		fmt.Sprintf(
			"%q is not a valid cluster_size. Use a fully qualified cluster size such as "+
				"PS_10_AWS_ARM, PS_10_GCP_X86, or M1_10_AWS_AMD_D_METAL_10. ",
			value,
		),
	)
}

func ClusterSizeValidator() validator.String {
	return StringClusterSizeValidatorValidator{}
}
