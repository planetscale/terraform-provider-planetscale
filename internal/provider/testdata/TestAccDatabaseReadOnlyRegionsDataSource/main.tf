variable "name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  database     = var.name
  name         = var.name
  organization = data.planetscale_organizations.test.data[0].name
}

data "planetscale_database_read_only_regions" "test" {
  database     = planetscale_database.test.name
  organization = planetscale_database.test.organization
}
