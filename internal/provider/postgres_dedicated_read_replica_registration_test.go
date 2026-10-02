package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

func TestPostgresDedicatedReadReplicaResourceMetadata(t *testing.T) {
	t.Parallel()

	response := resourceMetadata(NewPostgresDedicatedReadReplicaResource())

	require.Equal(t, "planetscale_postgres_dedicated_read_replica", response.TypeName)
}

func TestPostgresDedicatedReadReplicaDataSourceMetadata(t *testing.T) {
	t.Parallel()

	response := dataSourceMetadata(NewPostgresDedicatedReadReplicaDataSource())

	require.Equal(t, "planetscale_postgres_dedicated_read_replica", response.TypeName)
}

func TestPostgresDedicatedReadReplicasDataSourceMetadata(t *testing.T) {
	t.Parallel()

	response := dataSourceMetadata(NewPostgresDedicatedReadReplicasDataSource())

	require.Equal(t, "planetscale_postgres_dedicated_read_replicas", response.TypeName)
}

func resourceMetadata(terraformResource resource.Resource) *resource.MetadataResponse {
	response := &resource.MetadataResponse{}
	terraformResource.Metadata(
		context.Background(),
		resource.MetadataRequest{ProviderTypeName: "planetscale"},
		response,
	)

	return response
}

func dataSourceMetadata(terraformDataSource datasource.DataSource) *datasource.MetadataResponse {
	response := &datasource.MetadataResponse{}
	terraformDataSource.Metadata(
		context.Background(),
		datasource.MetadataRequest{ProviderTypeName: "planetscale"},
		response,
	)

	return response
}
