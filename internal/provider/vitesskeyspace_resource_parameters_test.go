package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestVitessKeyspaceResourceReadsParameterOverrides(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	api, client := newFakeVitessAPI(t)
	r := &VitessKeyspaceResource{client: client}
	var schemaResponse resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)

	state := tfsdk.State{Schema: schemaResponse.Schema, Raw: terraformValue(t, ctx, schemaResponse.Schema, &VitessKeyspaceResourceModel{
		Organization: types.StringValue("org"),
		Database:     types.StringValue("db"),
		Branch:       types.StringValue("main"),
		Name:         types.StringValue("commerce"),
	})}
	response := resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, &response)
	require.False(t, response.Diagnostics.HasError(), response.Diagnostics.Errors())

	require.Equal(t, map[string]map[string]string{
		"mysqld": {"wait_timeout": "600"},
	}, stateParameters(t, ctx, response.State))
	require.Zero(t, api.draftCount())
}
