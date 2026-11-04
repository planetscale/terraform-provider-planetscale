package mapvalidators

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Map = MapVitessParametersValidator{}

// MapVitessParametersValidator checks a Vitess `parameters` map at plan time:
// namespaces must be components of the resource so typos fail before the
// provider reads the branch, and namespaces must not be empty or contain null
// values because state can't keep either, so they would drift on every plan.
type MapVitessParametersValidator struct {
	namespaces []string
}

// Description describes the validation in plain text formatting.
func (v MapVitessParametersValidator) Description(_ context.Context) string {
	return "namespaces must be " + v.allowed() + " and contain at least one parameter value"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v MapVitessParametersValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateMap performs the validation.
func (v MapVitessParametersValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()
	namespaces := make([]string, 0, len(elements))
	for namespace := range elements {
		namespaces = append(namespaces, namespace)
	}
	slices.Sort(namespaces)

	for _, namespace := range namespaces {
		namespacePath := req.Path.AtMapKey(namespace)
		if !slices.Contains(v.namespaces, namespace) {
			resp.Diagnostics.AddAttributeError(
				namespacePath,
				"Unsupported parameter namespace",
				fmt.Sprintf("%q is not a parameter namespace for this resource. Use %s.", namespace, v.allowed()),
			)
			continue
		}

		parameters, ok := elements[namespace].(types.Map)
		if !ok || parameters.IsUnknown() {
			continue
		}
		if parameters.IsNull() || len(parameters.Elements()) == 0 {
			resp.Diagnostics.AddAttributeError(
				namespacePath,
				"Empty parameter namespace",
				fmt.Sprintf("The %q namespace has no parameters. Set at least one parameter or remove the namespace.", namespace),
			)
			continue
		}

		for name, value := range parameters.Elements() {
			if value.IsNull() {
				resp.Diagnostics.AddAttributeError(
					namespacePath.AtMapKey(name),
					"Null parameter value",
					fmt.Sprintf("Parameter %q in the %q namespace is null. Remove it to reset it to its default.", name, namespace),
				)
			}
		}
	}
}

func (v MapVitessParametersValidator) allowed() string {
	quoted := make([]string, 0, len(v.namespaces))
	for _, namespace := range v.namespaces {
		quoted = append(quoted, fmt.Sprintf("%q", namespace))
	}
	return strings.Join(quoted, " or ")
}

func VitessVtgateParametersValidator() validator.Map {
	return MapVitessParametersValidator{namespaces: []string{"vtgate"}}
}
