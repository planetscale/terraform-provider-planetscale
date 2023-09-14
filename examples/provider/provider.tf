terraform {
  required_providers {
    planetscale = {
      source = "registry.terraform.io/planetscale/planetscale"
    }
  }
}

provider "planetscale" {
  service_token_name = "luq1jk0pjccp"
}

data "planetscale_organizations" "test" {}

# data "planetscale_organization_regions" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
# }

# output "my_regions" {
#   value = data.planetscale_organization_regions.test
# }

# data "planetscale_databases" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
# }

# output "my_dbs" {
#   value = data.planetscale_databases.test
# }

data "planetscale_database" "my_db" {
  organization = data.planetscale_organizations.test.organizations.0.name
  name = "again"
}

# resource "planetscale_database" "my_db" {
#   organization = data.planetscale_organizations.test.organizations.0.name
#   name = "again"
# }

# data "planetscale_database_regions" "my_regions" {
#   organization = resource.planetscale_database.my_db.organization
#   name = resource.planetscale_database.my_db.name
# }

# data "planetscale_database_read_only_regions" "my_ro_regions" {
#   organization = resource.planetscale_database.my_db.organization
#   name = resource.planetscale_database.my_db.name
# }

# data "planetscale_branches" "my_branches" {
#   organization = data.planetscale_database.my_db.organization
#   database = data.planetscale_database.my_db.name
# }

# output "my_branches" {
#   value = data.planetscale_branches.my_branches
# }

# data "planetscale_branch" "my_branch" {
#   organization = data.planetscale_database.my_db.organization
#   database = data.planetscale_database.my_db.name
#   name = "world"
# }

# output "my_branch" {
#   value = data.planetscale_branch.my_branch
# }

# data "planetscale_branch_schema" "my_schema" {
#   organization = resource.planetscale_database.my_db.organization
#   database = resource.planetscale_database.my_db.name
#   branch = data.planetscale_branch.my_branch.name
# }

# output "my_schema" {
#   value = data.planetscale_branch_schema.my_schema
# }

# data "planetscale_branch_schema_lint" "my_schema_lint" {
#   organization = resource.planetscale_database.my_db.organization
#   database = resource.planetscale_database.my_db.name
#   branch = data.planetscale_branch.my_branch.name
# }

# output "my_schema_lint" {
#   value = data.planetscale_branch_schema_lint.my_schema_lint
# }

resource "planetscale_branch" "test" {
  organization = data.planetscale_organizations.test.organizations.0.name
  database = data.planetscale_database.my_db.name
  name = "hello"
  parent_branch = "main"
}

# data "planetscale_oauth_applications" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
# }

# output "my_oauth_apps" {
#   value = data.planetscale_oauth_applications.test
# }

# data "planetscale_user" "test" {}

# output "current_user" {
#   value = data.planetscale_user.test
# }