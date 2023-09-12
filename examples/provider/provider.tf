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

data "planetscale_organization_regions" "test" {
  organization = data.planetscale_organizations.test.organizations.0.name
}

# output "my_regions" {
#   value = data.planetscale_organization_regions.test
# }

data "planetscale_databases" "test" {
  organization = data.planetscale_organizations.test.organizations.0.name
}

# output "my_dbs" {
#   value = data.planetscale_databases.test
# }

resource "planetscale_database" "my_db" {
  organization = data.planetscale_organizations.test.organizations.0.name
  name = "again"
}


output "my_dbs_res" {
  value = resource.planetscale_database.my_db
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