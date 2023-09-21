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

data "planetscale_organization" "test" {
  name = data.planetscale_organizations.test.organizations.0.name
}

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
#   name         = data.planetscale_database.test.name
# }

# output "database_regions" {
#   value = data.planetscale_database_regions.test
# }

# data "planetscale_database_read_only_regions" "test" {
#   organization = data.planetscale_database.test.organization
#   name         = data.planetscale_database.test.name
# }

# output "database_ro_regions" {
#   value = data.planetscale_database_regions.test
# }

# data "planetscale_branches" "test" {
#   organization = data.planetscale_database.test.organization
#   database     = data.planetscale_database.test.name
# }

# output "branches" {
#   value = data.planetscale_branches.test
# }

data "planetscale_branch" "test" {
  organization = data.planetscale_database.test.organization
  database     = data.planetscale_database.test.name
  name         = "world"
}

# output "branch" {
#   value = data.planetscale_branch.test
# }

# data "planetscale_branch_schema" "test" {
#   organization = data.planetscale_database.test.organization
#   database     = data.planetscale_database.test.name
#   branch       = data.planetscale_branch.test.name
# }

# output "branch_schema" {
#   value = data.planetscale_branch_schema.test
# }

# data "planetscale_branch_schema_lint" "test" {
#   organization = data.planetscale_database.test.organization
#   database     = data.planetscale_database.test.name
#   branch       = data.planetscale_branch.test.name
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

# data "planetscale_backups" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
#   database     = data.planetscale_database.test.name
#   branch       = data.planetscale_branch.test.name
# }

# output "backups" {
#   value = data.planetscale_backups.test
# }

# data "planetscale_backup" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
#   database     = data.planetscale_database.test.name
#   branch       = data.planetscale_branch.test.name
#   id           = data.planetscale_backups.test.backups.0.id
# }

# output "backup" {
#   value = data.planetscale_backup.test
# }

data "planetscale_passwords" "test" {
  organization = data.planetscale_organizations.test.organizations.0.name
  database     = data.planetscale_database.test.name
  branch       = data.planetscale_branch.test.name
}

output "passwords" {
  value = data.planetscale_passwords.test
}

# resource "planetscale_database" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
#   name         = "again"
# }

# resource "planetscale_branch" "test" {
#   organization  = data.planetscale_organizations.test.organizations.0.name
#   database      = data.planetscale_database.test.name
#   name          = "world"
#   parent_branch = "main"
# }

# resource "planetscale_backup" "test" {
#   organization = data.planetscale_organizations.test.organizations.0.name
#   database     = data.planetscale_database.test.name
#   branch       = resource.planetscale_branch.test.name
#   name         = "antoine_was_here"
#   backup_policy = {
#     retention_unit  = "day"
#     retention_value = 1
#   }
# }

resource "planetscale_password" "test" {
  organization = data.planetscale_organizations.test.organizations.0.name
  database     = data.planetscale_database.test.name
  branch       = data.planetscale_branch.test.name
  name         = "antoine-was-here"

  ttl_seconds = 120
}

output "password" {
  sensitive = true
  value = planetscale_password.test
}