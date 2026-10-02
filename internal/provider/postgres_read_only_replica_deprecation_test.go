package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

func TestDeprecatedPostgresReadOnlyReplicaResource(t *testing.T) {
	t.Parallel()

	terraformResource := NewDeprecatedPostgresReadOnlyReplicaResource()
	metadataResponse := resourceMetadata(terraformResource)
	schemaResponse := &resource.SchemaResponse{}
	terraformResource.Schema(context.Background(), resource.SchemaRequest{}, schemaResponse)

	require.Equal(t, "planetscale_postgres_read_only_replica", metadataResponse.TypeName)
	require.Equal(t, postgresReadOnlyReplicaResourceDeprecation, schemaResponse.Schema.DeprecationMessage)
}

func TestDeprecatedPostgresReadOnlyReplicaDataSource(t *testing.T) {
	t.Parallel()

	terraformDataSource := NewDeprecatedPostgresReadOnlyReplicaDataSource()
	metadataResponse := dataSourceMetadata(terraformDataSource)
	schemaResponse := &datasource.SchemaResponse{}
	terraformDataSource.Schema(context.Background(), datasource.SchemaRequest{}, schemaResponse)

	require.Equal(t, "planetscale_postgres_read_only_replica", metadataResponse.TypeName)
	require.Equal(t, postgresReadOnlyReplicaDataSourceDeprecation, schemaResponse.Schema.DeprecationMessage)
}

func TestDeprecatedPostgresReadOnlyReplicasDataSource(t *testing.T) {
	t.Parallel()

	terraformDataSource := NewDeprecatedPostgresReadOnlyReplicasDataSource()
	metadataResponse := dataSourceMetadata(terraformDataSource)
	schemaResponse := &datasource.SchemaResponse{}
	terraformDataSource.Schema(context.Background(), datasource.SchemaRequest{}, schemaResponse)

	require.Equal(t, "planetscale_postgres_read_only_replicas", metadataResponse.TypeName)
	require.Equal(t, postgresReadOnlyReplicasDataSourceDeprecation, schemaResponse.Schema.DeprecationMessage)
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
