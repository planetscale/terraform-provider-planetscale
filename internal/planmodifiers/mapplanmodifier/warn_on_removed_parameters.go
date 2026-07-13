package mapplanmodifier

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Map = MapWarnOnRemovedParametersPlanModifier{}

type MapWarnOnRemovedParametersPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v MapWarnOnRemovedParametersPlanModifier) Description(_ context.Context) string {
	return "Warns when previously configured parameters are removed from the configuration, since removed parameters are reset to their defaults."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v MapWarnOnRemovedParametersPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyMap performs the plan modification.
func (v MapWarnOnRemovedParametersPlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	// No prior state to compare against (resource creation).
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	// Config null means the whole attribute was removed (or this is a destroy
	// plan). The attribute is Computed+Optional, so Terraform keeps the prior
	// state value in the plan and nothing is reset. Config unknown means we
	// can't know what will be sent.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var state map[string]map[string]string
	if diags := req.StateValue.ElementsAs(ctx, &state, false); diags.HasError() {
		return
	}

	// Only key presence matters, and nested config values may be unknown, so
	// walk the elements directly instead of using ElementsAs. A nil key set
	// marks a namespace whose contents are unknown.
	configKeys := map[string]map[string]struct{}{}
	for namespace, value := range req.ConfigValue.Elements() {
		inner, ok := value.(types.Map)
		if !ok || inner.IsUnknown() {
			configKeys[namespace] = nil
			continue
		}
		keys := map[string]struct{}{}
		if !inner.IsNull() {
			for key := range inner.Elements() {
				keys[key] = struct{}{}
			}
		}
		configKeys[namespace] = keys
	}

	var removed []string
	for namespace, params := range state {
		keys, namespaceInConfig := configKeys[namespace]
		if namespaceInConfig && keys == nil {
			continue
		}
		for key, value := range params {
			if _, ok := keys[key]; !ok {
				removed = append(removed, fmt.Sprintf("%s.%s (currently %q)", namespace, key, value))
			}
		}
	}
	if len(removed) == 0 {
		return
	}
	sort.Strings(removed)

	resp.Diagnostics.AddAttributeWarning(
		req.Path,
		"Removed parameters will be reset to their defaults",
		"The following parameters were removed from the configuration, so on apply each "+
			"will be reset to its default value:\n\n  - "+
			strings.Join(removed, "\n  - ")+
			"\n\nTo keep a parameter at its current value, add it back to the parameters attribute.",
	)
}

func WarnOnRemovedParameters() planmodifier.Map {
	return MapWarnOnRemovedParametersPlanModifier{}
}
