package mapplanmodifier

import "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

func WarnOnRemovedNekiParameters() planmodifier.Map {
	return MapWarnOnRemovedParametersPlanModifier{preservePreloads: true}
}
