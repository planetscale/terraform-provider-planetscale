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

func TestDeprecatedPostgresReadOnlyReplicaTypesDocumentTheReplacement(t *testing.T) {
	t.Parallel()

	resourceSchema := &resource.SchemaResponse{}
	NewDeprecatedPostgresReadOnlyReplicaResource().Schema(context.Background(), resource.SchemaRequest{}, resourceSchema)
	require.Contains(t, resourceSchema.Schema.MarkdownDescription, "**Deprecated.**")
	require.Contains(t, resourceSchema.Schema.MarkdownDescription, "import {\n  to = planetscale_postgres_dedicated_read_replica.example")
	require.Contains(t, resourceSchema.Schema.MarkdownDescription, "terraform import planetscale_postgres_dedicated_read_replica.example")

	dataSourceSchema := &datasource.SchemaResponse{}
	NewDeprecatedPostgresReadOnlyReplicaDataSource().Schema(context.Background(), datasource.SchemaRequest{}, dataSourceSchema)
	require.Contains(t, dataSourceSchema.Schema.MarkdownDescription, "**Deprecated.**")
	require.Contains(t, dataSourceSchema.Schema.MarkdownDescription, "planetscale_postgres_dedicated_read_replica")

	dataSourcesSchema := &datasource.SchemaResponse{}
	NewDeprecatedPostgresReadOnlyReplicasDataSource().Schema(context.Background(), datasource.SchemaRequest{}, dataSourcesSchema)
	require.Contains(t, dataSourcesSchema.Schema.MarkdownDescription, "**Deprecated.**")
	require.Contains(t, dataSourcesSchema.Schema.MarkdownDescription, "planetscale_postgres_dedicated_read_replicas")
}
