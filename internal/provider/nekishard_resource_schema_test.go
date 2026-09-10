package provider

import (
	"context"
	"testing"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/require"
)

func TestNekiShardResourceSchemaIncludesComputedID(t *testing.T) {
	t.Parallel()

	resource := &NekiShardResource{}
	var resp frameworkresource.SchemaResponse
	resource.Schema(context.Background(), frameworkresource.SchemaRequest{}, &resp)

	require.False(t, resp.Diagnostics.HasError())
	id, ok := resp.Schema.Attributes["id"]
	require.True(t, ok, "Neki shard ID must be persisted for read, update, and delete")

	idAttribute, ok := id.(schema.StringAttribute)
	require.True(t, ok)
	require.True(t, idAttribute.Computed)
	require.False(t, idAttribute.Required)
	require.False(t, idAttribute.Optional)
}
