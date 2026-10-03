package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// The read-only replica resource and data sources were renamed to dedicated
// read replica. These wrappers keep the old type names registered so existing
// configurations and state keep working, and they warn practitioners to move
// to the new names. Terraform's state mv command refuses to move state between
// resource types, so the resource deprecation message explains how to forget
// the old resource and import it under the new type instead.
const (
	postgresReadOnlyReplicaResourceDeprecation = "planetscale_postgres_read_only_replica is deprecated; use planetscale_postgres_dedicated_read_replica instead. " +
		"To switch an existing replica without recreating it, remove it from state with terraform state rm, " +
		"then import it as planetscale_postgres_dedicated_read_replica using the JSON import ID described in that resource's documentation."
	postgresReadOnlyReplicaDataSourceDeprecation  = "planetscale_postgres_read_only_replica is deprecated; use planetscale_postgres_dedicated_read_replica instead."
	postgresReadOnlyReplicasDataSourceDeprecation = "planetscale_postgres_read_only_replicas is deprecated; use planetscale_postgres_dedicated_read_replicas instead."

	postgresReadOnlyReplicaResourceDescription = "**Deprecated.** Use `planetscale_postgres_dedicated_read_replica` instead. " +
		"This resource keeps working for existing configurations.\n\n" +
		"To switch an existing replica to the new resource type without recreating it, " +
		"remove it from state and import it under the new type. In Terraform v1.7.0 and later, " +
		"a `removed` block and an `import` block can do this in a single apply:\n\n" +
		"```terraform\n" +
		"removed {\n" +
		"  from = planetscale_postgres_read_only_replica.example\n" +
		"\n" +
		"  lifecycle {\n" +
		"    destroy = false\n" +
		"  }\n" +
		"}\n" +
		"\n" +
		"import {\n" +
		"  to = planetscale_postgres_dedicated_read_replica.example\n" +
		"  id = jsonencode({\n" +
		"    organization = \"acme\"\n" +
		"    database     = \"app-db\"\n" +
		"    branch       = \"main\"\n" +
		"    name         = \"example\"\n" +
		"  })\n" +
		"}\n" +
		"```\n\n" +
		"The equivalent commands are:\n\n" +
		"```shell\n" +
		"terraform state rm planetscale_postgres_read_only_replica.example\n" +
		"terraform import planetscale_postgres_dedicated_read_replica.example '{\"organization\": \"acme\", \"database\": \"app-db\", \"branch\": \"main\", \"name\": \"example\"}'\n" +
		"```"
	postgresReadOnlyReplicaDataSourceDescription  = "**Deprecated.** Use `planetscale_postgres_dedicated_read_replica` instead."
	postgresReadOnlyReplicasDataSourceDescription = "**Deprecated.** Use `planetscale_postgres_dedicated_read_replicas` instead."
)

var _ resource.Resource = &DeprecatedPostgresReadOnlyReplicaResource{}
var _ resource.ResourceWithConfigure = &DeprecatedPostgresReadOnlyReplicaResource{}
var _ resource.ResourceWithImportState = &DeprecatedPostgresReadOnlyReplicaResource{}

type DeprecatedPostgresReadOnlyReplicaResource struct {
	*PostgresDedicatedReadReplicaResource
}

func NewDeprecatedPostgresReadOnlyReplicaResource() resource.Resource {
	return &DeprecatedPostgresReadOnlyReplicaResource{
		PostgresDedicatedReadReplicaResource: &PostgresDedicatedReadReplicaResource{},
	}
}

func (r *DeprecatedPostgresReadOnlyReplicaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_read_only_replica"
}

func (r *DeprecatedPostgresReadOnlyReplicaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.PostgresDedicatedReadReplicaResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = postgresReadOnlyReplicaResourceDeprecation
	resp.Schema.MarkdownDescription = postgresReadOnlyReplicaResourceDescription
}

var _ datasource.DataSource = &DeprecatedPostgresReadOnlyReplicaDataSource{}
var _ datasource.DataSourceWithConfigure = &DeprecatedPostgresReadOnlyReplicaDataSource{}

type DeprecatedPostgresReadOnlyReplicaDataSource struct {
	*PostgresDedicatedReadReplicaDataSource
}

func NewDeprecatedPostgresReadOnlyReplicaDataSource() datasource.DataSource {
	return &DeprecatedPostgresReadOnlyReplicaDataSource{
		PostgresDedicatedReadReplicaDataSource: &PostgresDedicatedReadReplicaDataSource{},
	}
}

func (d *DeprecatedPostgresReadOnlyReplicaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_read_only_replica"
}

func (d *DeprecatedPostgresReadOnlyReplicaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.PostgresDedicatedReadReplicaDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = postgresReadOnlyReplicaDataSourceDeprecation
	resp.Schema.MarkdownDescription = postgresReadOnlyReplicaDataSourceDescription
}

var _ datasource.DataSource = &DeprecatedPostgresReadOnlyReplicasDataSource{}
var _ datasource.DataSourceWithConfigure = &DeprecatedPostgresReadOnlyReplicasDataSource{}

type DeprecatedPostgresReadOnlyReplicasDataSource struct {
	*PostgresDedicatedReadReplicasDataSource
}

func NewDeprecatedPostgresReadOnlyReplicasDataSource() datasource.DataSource {
	return &DeprecatedPostgresReadOnlyReplicasDataSource{
		PostgresDedicatedReadReplicasDataSource: &PostgresDedicatedReadReplicasDataSource{},
	}
}

func (d *DeprecatedPostgresReadOnlyReplicasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_read_only_replicas"
}

func (d *DeprecatedPostgresReadOnlyReplicasDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	d.PostgresDedicatedReadReplicasDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = postgresReadOnlyReplicasDataSourceDeprecation
	resp.Schema.MarkdownDescription = postgresReadOnlyReplicasDataSourceDescription
}
