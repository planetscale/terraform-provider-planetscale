variable "name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  name         = var.name
  organization = data.planetscale_organizations.test.data[0].name
}

data "planetscale_branch_schema" "test" {
  database     = planetscale_database.test.name
  name         = "main"
  organization = planetscale_database.test.organization
}
