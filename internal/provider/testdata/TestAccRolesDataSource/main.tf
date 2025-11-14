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

resource "planetscale_role" "test" {
  branch       = planetscale_database.test.default_branch
  database     = planetscale_database.test.name
  organization = planetscale_database.test.organization
}

data "planetscale_roles" "test" {
  branch       = planetscale_database.test.default_branch
  database     = planetscale_database.test.name
  organization = planetscale_database.test.organization

  depends_on = [planetscale_role.test]
}
