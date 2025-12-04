variable "database_name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  database     = var.database_name
  name         = var.database_name
  organization = data.planetscale_organizations.test.data[0].name
}

resource "planetscale_password" "test" {
  branch       = "main"
  database     = planetscale_database.test.name
  name         = "test"
  organization = planetscale_database.test.organization
  role         = "admin"
}
