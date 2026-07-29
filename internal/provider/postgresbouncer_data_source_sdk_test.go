package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/operations"
	"github.com/stretchr/testify/require"
)

func TestPostgresBouncerDataSourceRefreshFromTerraformStateParameters(t *testing.T) {
	t.Parallel()

	model := &PostgresBouncerDataSourceModel{
		Parameters: map[string]map[string]types.String{
			"pgbouncer": {
				"default_pool_size": types.StringValue("50"),
			},
		},
	}

	diags := model.RefreshFromOperationsGetPostgresBouncerTerraformStateResponseBody(
		context.Background(),
		&operations.GetPostgresBouncerTerraformStateResponseBody{
			ID:              "bouncer-id",
			Name:            "my-bouncer",
			Target:          operations.GetPostgresBouncerTerraformStateTargetPrimary,
			BouncerSize:     "PGB_5",
			ReplicasPerCell: 1,
			Actor:           operations.GetPostgresBouncerTerraformStateActor{ID: "actor-id"},
			Parameters: map[string]map[string]string{
				"pgbouncer": {
					"default_pool_size": "100",
				},
			},
		},
	)

	require.False(t, diags.HasError(), diags.Errors())
	require.Equal(t, "100", model.Parameters["pgbouncer"]["default_pool_size"].ValueString())
	require.Equal(t, "actor-id", model.Actor.ID.ValueString())

	diags = model.RefreshFromOperationsGetPostgresBouncerTerraformStateResponseBody(
		context.Background(),
		&operations.GetPostgresBouncerTerraformStateResponseBody{
			ID:              "bouncer-id",
			Name:            "my-bouncer",
			Target:          operations.GetPostgresBouncerTerraformStateTargetPrimary,
			BouncerSize:     "PGB_5",
			ReplicasPerCell: 1,
			Actor:           operations.GetPostgresBouncerTerraformStateActor{ID: "actor-id"},
			Parameters:      map[string]map[string]string{},
		},
	)

	require.False(t, diags.HasError(), diags.Errors())
	require.NotNil(t, model.Parameters)
	require.Empty(t, model.Parameters)
}
