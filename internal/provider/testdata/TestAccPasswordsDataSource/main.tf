variable "database_name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  name         = var.database_name
  organization = data.planetscale_organizations.test.data[0].name
}

resource "planetscale_password" "test" {
  branch       = planetscale_database.test.default_branch
  database     = planetscale_database.test.name
  name         = "test"
  organization = planetscale_database.test.organization
  role         = "admin"
}

data "planetscale_passwords" "test" {
  branch       = planetscale_database.test.default_branch
  database     = planetscale_database.test.name
  organization = planetscale_database.test.organization

  depends_on = [planetscale_password.test]
}
