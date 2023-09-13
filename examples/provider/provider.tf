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

resource "planetscale_database" "my_db" {
  organization = data.planetscale_organizations.test.organizations.0.name
  name = "again"
}

# data "planetscale_database_regions" "my_regions" {
#   organization = resource.planetscale_database.my_db.organization
#   name = resource.planetscale_database.my_db.name
# }

data "planetscale_database_read_only_regions" "my_ro_regions" {
  organization = resource.planetscale_database.my_db.organization
  name = resource.planetscale_database.my_db.name
}

output "my_dbs_ro_regions" {
  value = data.planetscale_database_read_only_regions.my_ro_regions
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