package stringvalidators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = StringClusterSizeValidatorValidator{}

const (
	psNumericSizePattern  = `(?:5|[1-9][0-9]*0)`
	psSizePattern         = `(?:` + psNumericSizePattern + `|DEV)`
	qualifiedPSPattern    = `PS_` + psSizePattern + `_(?:AWS|GCP)_(?:ARM|X86|AMD)`
	vitessPSPattern       = `PS_` + psSizePattern
	instancePSPattern     = `PS_(?:(?:AWS|GCP)_)?[A-Z][0-9][A-Z0-9]*(?:_[A-Z0-9]+)+`
	qualifiedMetalPattern = `M[0-9]+_[0-9]+_(?:AWS|GCP)_(?:ARM|X86|AMD)_D_METAL_[0-9]+`
	vitessMetalPattern    = `M[0-9]*_[0-9]+(?:_(?:AWS|GCP)_(?:AMD|INTEL|X86))?_D_METAL_[0-9]+`
)

var clusterSizePattern = regexp.MustCompile(
	`^(?:` + qualifiedPSPattern + `|` + qualifiedMetalPattern + `)$`,
)

var vitessClusterSizePattern = regexp.MustCompile(
	`^(?:` + qualifiedPSPattern + `|` + qualifiedMetalPattern + `|` + vitessPSPattern + `|` + instancePSPattern + `|` + vitessMetalPattern + `)$`,
)

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
	validateClusterSize(req, resp, clusterSizePattern, "PS_10_AWS_ARM, PS_10_GCP_X86, or M1_10_AWS_AMD_D_METAL_10")
}

func validateClusterSize(req validator.StringRequest, resp *validator.StringResponse, pattern *regexp.Regexp, examples string) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if pattern.MatchString(value) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid cluster_size value",
		fmt.Sprintf(
			"%q is not a valid cluster_size. Use a valid cluster size such as %s. ",
			value,
			examples,
		),
	)
}

func ClusterSizeValidator() validator.String {
	return StringClusterSizeValidatorValidator{}
}
