variable "database_name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  name         = var.database_name
  organization = data.planetscale_organizations.test.data[0].name
}

resource "planetscale_query_pattern_report" "test" {
  branch       = planetscale_database.test.default_branch
  database     = planetscale_database.test.name
  organization = planetscale_database.test.organization
}
