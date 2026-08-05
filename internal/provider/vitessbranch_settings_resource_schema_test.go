package provider

import (
	"testing"

	"github.com/planetscale/terraform-provider-planetscale/internal/sdk"
	"github.com/stretchr/testify/require"
)

func TestConfigurationFromResizeUsesPublicVTGateNames(t *testing.T) {
	t.Parallel()

	maxCount := int64(8)
	targetCPU := int64(50)
	current := configurationFromResize("branch-id", sdk.VitessBranchResizeRequest{
		State:                      "completed",
		VTGateSize:                 "vg.c1.large",
		VTGateName:                 "VTG_320",
		VTGateCount:                2,
		VTGateMaxCount:             &maxCount,
		VTGateAutoscaling:          true,
		VTGateTargetCPUUtilization: &targetCPU,
	})

	require.Equal(t, "branch-id", current.branchID)
	require.Equal(t, "VTG_320", current.size)
	require.Equal(t, int64(2), current.count)
	require.Equal(t, maxCount, *current.maxCount)
	require.True(t, current.autoscaling)
	require.Equal(t, targetCPU, *current.targetCPUUtilization)
}

func TestConfigurationFromCanceledResizeUsesPreviousValues(t *testing.T) {
	t.Parallel()

	previousMaxCount := int64(4)
	current := configurationFromResize("branch-id", sdk.VitessBranchResizeRequest{
		State:                     "canceled",
		VTGateName:                "VTG_320",
		PreviousVTGateName:        "VTG_80",
		VTGateAutoscaling:         true,
		PreviousVTGateAutoscaling: false,
		PreviousVTGateCount:       1,
		PreviousVTGateMaxCount:    &previousMaxCount,
	})

	require.Equal(t, "VTG_80", current.size)
	require.Equal(t, int64(1), current.count)
	require.Equal(t, previousMaxCount, *current.maxCount)
	require.False(t, current.autoscaling)
}
