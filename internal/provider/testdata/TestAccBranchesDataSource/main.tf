variable "database_name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  name         = var.database_name
  organization = data.planetscale_organizations.test.data[0].name
}

resource "planetscale_branch" "test" {
  database      = planetscale_database.test.name
  name          = "test"
  organization  = planetscale_database.test.organization
  parent_branch = "main"
}

data "planetscale_branches" "test" {
  database     = planetscale_branch.test.database
  organization = planetscale_branch.test.organization
}
