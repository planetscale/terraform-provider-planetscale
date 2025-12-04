#
# Example MySQL-based Terraform configuration for generated PlanetScale provider
#
# This configuration:
#   - Creates a PlanetScale MySQL database
#   - Creates two branches on that database: release and staging
#   - Sets the default branch to release
#   - Creates a password for the staging branch
#

data "planetscale_organizations" "example" {}

resource "planetscale_database_vitess" "example" {
  cluster_size = "PS_10"
  name         = "e2e-mysql-example"
  organization = data.planetscale_organizations.example.data[0].name
}

resource "planetscale_database_default_branch" "example" {
  branch       = planetscale_branch.release.name
  database     = planetscale_database_vitess.example.name
  organization = planetscale_database_vitess.example.organization
}

resource "planetscale_branch" "release" {
  database      = planetscale_database_vitess.example.name
  name          = "release"
  organization  = planetscale_database_vitess.example.organization
  parent_branch = "main"
}

resource "planetscale_branch" "staging" {
  database      = planetscale_database_vitess.example.name
  name          = "staging"
  organization  = planetscale_database_vitess.example.organization
  parent_branch = "main"
}

resource "planetscale_password" "staging-ci" {
  branch       = planetscale_branch.staging.name
  database     = planetscale_branch.staging.database
  name         = "staging-ci"
  organization = planetscale_branch.staging.organization
  role         = "admin"
}
