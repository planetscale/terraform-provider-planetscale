package boolvalidators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Bool = BoolDeletionProtectionValidator{}

type BoolDeletionProtectionValidator struct{}

// Description describes the validation in plain text formatting.
func (v BoolDeletionProtectionValidator) Description(_ context.Context) string {
	return "Required parameter to prevent accidental deletions"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v BoolDeletionProtectionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v BoolDeletionProtectionValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	// Deletes are only allowed if the value is explicitly set to false
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() && !req.ConfigValue.ValueBool() {
		return
	}
	
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Deletion protection must be set to false to delete the branch",
		req.Path.String()+": "+v.Description(ctx),
	)
}

func DeletionProtection() validator.Bool {
	return BoolDeletionProtectionValidator{}
}
