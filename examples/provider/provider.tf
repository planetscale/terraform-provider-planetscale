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

# output "orgs" {
#   value = data.planetscale_organizations.test
# }

# data "planetscale_organization" "test" {
#   name = data.planetscale_organizations.test.organizations.0.name
# }

# output "org" {
#   value = data.planetscale_organization.test
# }

# data "planetscale_organization_regions" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
# }

# output "org_regions" {
#   value = data.planetscale_organization_regions.test
# }

# data "planetscale_databases" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
# }

# output "dbs" {
#   value = data.planetscale_databases.test
# }

data "planetscale_database" "test" {
  organization = data.planetscale_organizations.test.organizations.0.name
  name         = "again"
}

# output "db" {
#   value = data.planetscale_database.test
# }

# data "planetscale_database_regions" "test" {
#   organization = data.planetscale_database.test.organization
#   name = data.planetscale_database.test.name
# }

# output "database_regions" {
#   value = data.planetscale_database_regions.test
# }

# data "planetscale_database_read_only_regions" "test" {
#   organization = data.planetscale_database.test.organization
#   name = data.planetscale_database.test.name
# }

# output "database_ro_regions" {
#   value = data.planetscale_database_regions.test
# }

# data "planetscale_branches" "test" {
#   organization = data.planetscale_database.test.organization
#   database = data.planetscale_database.test.name
# }

# output "branches" {
#   value = data.planetscale_branches.test
# }

# data "planetscale_branch" "test" {
#   organization = data.planetscale_database.test.organization
#   database = data.planetscale_database.test.name
#   name = "main"
# }

# output "branch" {
#   value = data.planetscale_branch.test
# }

# data "planetscale_branch_schema" "test" {
#   organization = data.planetscale_database.test.organization
#   database = data.planetscale_database.test.name
#   branch = data.planetscale_branch.test.name
# }

# output "branch_schema" {
#   value = data.planetscale_branch_schema.test
# }

# data "planetscale_branch_schema_lint" "test" {
#   organization = data.planetscale_database.test.organization
#   database = data.planetscale_database.test.name
#   branch = data.planetscale_branch.test.name
# }

# output "schema_lint" {
#   value = data.planetscale_branch_schema_lint.test
# }

# requires a feature flag

# data "planetscale_oauth_applications" "test" {
#   organization = data.planetscale_organization.test.name
# }

# output "oauth_apps" {
#   value = data.planetscale_oauth_applications.test
# }

# doesn't work right now for some reason

# data "planetscale_user" "test" {}

# output "current_user" {
#   value = data.planetscale_user.test
# }


resource "planetscale_database" "test" {
  organization =  data.planetscale_organizations.test.organizations.0.name
  name = "again"
}

resource "planetscale_branch" "test" {
  organization  = data.planetscale_organizations.test.organizations.0.name
  database      = data.planetscale_database.test.name
  name          = "world"
  parent_branch = "main"
}
