package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PostgresDedicatedReadReplicaResource{}
var _ resource.ResourceWithConfigure = &PostgresDedicatedReadReplicaResource{}
var _ resource.ResourceWithImportState = &PostgresDedicatedReadReplicaResource{}

type PostgresDedicatedReadReplicaResource struct {
	*PostgresReadOnlyReplicaResource
}

func NewPostgresDedicatedReadReplicaResource() resource.Resource {
	return &PostgresDedicatedReadReplicaResource{
		PostgresReadOnlyReplicaResource: &PostgresReadOnlyReplicaResource{},
	}
}

func (r *PostgresDedicatedReadReplicaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_dedicated_read_replica"
}

var _ datasource.DataSource = &PostgresDedicatedReadReplicaDataSource{}
var _ datasource.DataSourceWithConfigure = &PostgresDedicatedReadReplicaDataSource{}

type PostgresDedicatedReadReplicaDataSource struct {
	*PostgresReadOnlyReplicaDataSource
}

func NewPostgresDedicatedReadReplicaDataSource() datasource.DataSource {
	return &PostgresDedicatedReadReplicaDataSource{
		PostgresReadOnlyReplicaDataSource: &PostgresReadOnlyReplicaDataSource{},
	}
}

func (d *PostgresDedicatedReadReplicaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_dedicated_read_replica"
}

var _ datasource.DataSource = &PostgresDedicatedReadReplicasDataSource{}
var _ datasource.DataSourceWithConfigure = &PostgresDedicatedReadReplicasDataSource{}

type PostgresDedicatedReadReplicasDataSource struct {
	*PostgresReadOnlyReplicasDataSource
}

func NewPostgresDedicatedReadReplicasDataSource() datasource.DataSource {
	return &PostgresDedicatedReadReplicasDataSource{
		PostgresReadOnlyReplicasDataSource: &PostgresReadOnlyReplicasDataSource{},
	}
}

func (d *PostgresDedicatedReadReplicasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_dedicated_read_replicas"
}
