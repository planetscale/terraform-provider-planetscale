package boolplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.Bool = BoolDeletionProtectionPlanModifier{}

type BoolDeletionProtectionPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v BoolDeletionProtectionPlanModifier) Description(_ context.Context) string {
	return "Required parameter to prevent accidental deletions"
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v BoolDeletionProtectionPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the plan modification.
func (v BoolDeletionProtectionPlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// Ignore operations that are not a destroy
	if !req.Plan.Raw.IsNull() {
		return
	}

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

func DeletionProtection() planmodifier.Bool {
	return BoolDeletionProtectionPlanModifier{}
}
