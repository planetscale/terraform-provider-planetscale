package stringvalidators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = StringVitessClusterSizeValidatorValidator{}

type StringVitessClusterSizeValidatorValidator struct{}

// Description describes the validation in plain text formatting.
func (v StringVitessClusterSizeValidatorValidator) Description(_ context.Context) string {
	return "value must be a valid cluster_size SKU"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v StringVitessClusterSizeValidatorValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v StringVitessClusterSizeValidatorValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateClusterSize(req, resp, vitessClusterSizePattern, "PS_10, PS_1400, PS_DEV, or M7_160_D_METAL_460")
}

func VitessClusterSizeValidator() validator.String {
	return StringVitessClusterSizeValidatorValidator{}
}
