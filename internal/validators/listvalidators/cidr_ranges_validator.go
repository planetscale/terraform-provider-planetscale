package listvalidators

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.List = ListCidrRangesValidatorValidator{}

// ListCidrRangesValidatorValidator requires every element to be an IPv4 range
// written the way the API stores it. The API rewrites anything else (a bare
// address gains /32, host bits are cleared), which Terraform would report as
// an inconsistent result after apply.
type ListCidrRangesValidatorValidator struct{}

// Description describes the validation in plain text formatting.
func (v ListCidrRangesValidatorValidator) Description(_ context.Context) string {
	return "each value must be an IPv4 range in CIDR notation, such as 203.0.113.0/24 or 198.51.100.7/32"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v ListCidrRangesValidatorValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v ListCidrRangesValidatorValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	for i, element := range req.ConfigValue.Elements() {
		value, ok := element.(types.String)
		if !ok || value.IsNull() || value.IsUnknown() {
			continue
		}

		if detail := cidrRangeProblem(value.ValueString()); detail != "" {
			resp.Diagnostics.AddAttributeError(req.Path.AtListIndex(i), "Invalid IP range", detail)
		}
	}
}

func cidrRangeProblem(value string) string {
	if addr, err := netip.ParseAddr(value); err == nil && addr.Is4() {
		return fmt.Sprintf("%q is a single address. Write it as a range: %q.", value, value+"/32")
	}

	prefix, err := netip.ParsePrefix(value)
	if err != nil || !prefix.Addr().Is4() {
		return fmt.Sprintf("%q is not an IPv4 range in CIDR notation. Use a range such as \"203.0.113.0/24\", or \"198.51.100.7/32\" for a single address.", value)
	}

	if canonical := prefix.Masked().String(); canonical != value {
		return fmt.Sprintf("%q does not start at the first address of its range. Use %q instead.", value, canonical)
	}

	return ""
}

func CidrRangesValidator() validator.List {
	return ListCidrRangesValidatorValidator{}
}
