package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	postgresReadOnlyReplicaResourceDeprecation    = "planetscale_postgres_read_only_replica is deprecated; use planetscale_postgres_dedicated_read_replica instead and move existing state with terraform state mv."
	postgresReadOnlyReplicaDataSourceDeprecation  = "planetscale_postgres_read_only_replica is deprecated; use planetscale_postgres_dedicated_read_replica instead."
	postgresReadOnlyReplicasDataSourceDeprecation = "planetscale_postgres_read_only_replicas is deprecated; use planetscale_postgres_dedicated_read_replicas instead."
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
}
