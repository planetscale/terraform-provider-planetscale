package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PostgresReadOnlyReplicaCompatibilityResource{}
var _ resource.ResourceWithConfigure = &PostgresReadOnlyReplicaCompatibilityResource{}
var _ resource.ResourceWithImportState = &PostgresReadOnlyReplicaCompatibilityResource{}

type PostgresReadOnlyReplicaCompatibilityResource struct {
	*PostgresDedicatedReadReplicaResource
}

func NewPostgresReadOnlyReplicaCompatibilityResource() resource.Resource {
	return &PostgresReadOnlyReplicaCompatibilityResource{
		PostgresDedicatedReadReplicaResource: &PostgresDedicatedReadReplicaResource{},
	}
}

func (r *PostgresReadOnlyReplicaCompatibilityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_read_only_replica"
}

var _ datasource.DataSource = &PostgresReadOnlyReplicaCompatibilityDataSource{}
var _ datasource.DataSourceWithConfigure = &PostgresReadOnlyReplicaCompatibilityDataSource{}

type PostgresReadOnlyReplicaCompatibilityDataSource struct {
	*PostgresDedicatedReadReplicaDataSource
}

func NewPostgresReadOnlyReplicaCompatibilityDataSource() datasource.DataSource {
	return &PostgresReadOnlyReplicaCompatibilityDataSource{
		PostgresDedicatedReadReplicaDataSource: &PostgresDedicatedReadReplicaDataSource{},
	}
}

func (d *PostgresReadOnlyReplicaCompatibilityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_read_only_replica"
}

var _ datasource.DataSource = &PostgresReadOnlyReplicasCompatibilityDataSource{}
var _ datasource.DataSourceWithConfigure = &PostgresReadOnlyReplicasCompatibilityDataSource{}

type PostgresReadOnlyReplicasCompatibilityDataSource struct {
	*PostgresDedicatedReadReplicasDataSource
}

func NewPostgresReadOnlyReplicasCompatibilityDataSource() datasource.DataSource {
	return &PostgresReadOnlyReplicasCompatibilityDataSource{
		PostgresDedicatedReadReplicasDataSource: &PostgresDedicatedReadReplicasDataSource{},
	}
}

func (d *PostgresReadOnlyReplicasCompatibilityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_read_only_replicas"
}
