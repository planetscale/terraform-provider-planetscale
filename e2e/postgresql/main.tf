#
# Example PostgreSQL-based Terraform configuration for generated PlanetScale provider
#
# This configuration:
#   - Creates a PlanetScale PostgreSQL database
#   - Creates two branches on that database: release and staging
#   - Sets the default branch to release
#   - Creates a role for the staging branch
#

data "planetscale_organizations" "example" {}

resource "planetscale_database" "example" {
  cluster_size = "PS_10_AWS_ARM"
  kind         = "postgresql"
  name         = "e2e-postgres-example"
  organization = data.planetscale_organizations.example.data[0].name
}

resource "planetscale_database_default_branch" "example" {
  branch       = planetscale_branch.release.name
  database     = planetscale_database.example.name
  organization = planetscale_database.example.organization
}

resource "planetscale_branch" "release" {
  database      = planetscale_database.example.name
  name          = "release"
  organization  = planetscale_database.example.organization
  parent_branch = "main"
}

resource "planetscale_branch" "staging" {
  database      = planetscale_database.example.name
  name          = "staging"
  organization  = planetscale_database.example.organization
  parent_branch = "main"
}

resource "planetscale_role" "staging-ci" {
  branch       = planetscale_branch.staging.name
  database     = planetscale_branch.staging.database
  organization = planetscale_branch.staging.organization
}
