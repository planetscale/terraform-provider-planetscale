variable "database_name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10_AWS_ARM"
  kind         = "postgresql"
  name         = var.database_name
  organization = data.planetscale_organizations.test.data[0].name
}

resource "planetscale_branch" "test" {
  database      = planetscale_database.test.name
  name          = "test"
  organization  = planetscale_database.test.organization
  parent_branch = planetscale_database.test.default_branch
}

resource "planetscale_bouncer" "test" {
  branch       = planetscale_branch.test.name
  database     = planetscale_branch.test.database
  name         = "test"
  organization = planetscale_branch.test.organization
  target       = "primary"
}
