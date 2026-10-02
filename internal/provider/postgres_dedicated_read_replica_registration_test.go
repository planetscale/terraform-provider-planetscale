package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

func TestPostgresDedicatedReadReplicaResourceRegistration(t *testing.T) {
	t.Parallel()

	typeNames := registeredResourceTypeNames((&PlanetscaleProvider{}).Resources(context.Background()))

	require.Equal(t, 1, typeNames["planetscale_postgres_dedicated_read_replica"])
	require.Equal(t, 1, typeNames["planetscale_postgres_read_only_replica"])
}

func TestPostgresDedicatedReadReplicaDataSourceRegistration(t *testing.T) {
	t.Parallel()

	typeNames := registeredDataSourceTypeNames((&PlanetscaleProvider{}).DataSources(context.Background()))

	require.Equal(t, 1, typeNames["planetscale_postgres_dedicated_read_replica"])
	require.Equal(t, 1, typeNames["planetscale_postgres_dedicated_read_replicas"])
	require.Equal(t, 1, typeNames["planetscale_postgres_read_only_replica"])
	require.Equal(t, 1, typeNames["planetscale_postgres_read_only_replicas"])
}

func registeredResourceTypeNames(factories []func() resource.Resource) map[string]int {
	typeNames := make(map[string]int, len(factories))

	for _, factory := range factories {
		response := &resource.MetadataResponse{}
		factory().Metadata(
			context.Background(),
			resource.MetadataRequest{ProviderTypeName: "planetscale"},
			response,
		)
		typeNames[response.TypeName]++
	}

	return typeNames
}

func registeredDataSourceTypeNames(factories []func() datasource.DataSource) map[string]int {
	typeNames := make(map[string]int, len(factories))

	for _, factory := range factories {
		response := &datasource.MetadataResponse{}
		factory().Metadata(
			context.Background(),
			datasource.MetadataRequest{ProviderTypeName: "planetscale"},
			response,
		)
		typeNames[response.TypeName]++
	}

	return typeNames
}
