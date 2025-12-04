variable "database_name" {
  type = string
}

variable "organization" {
  type = string
}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  database     = var.database_name
  name         = var.database_name
  organization = var.organization
}

resource "planetscale_password" "test" {
  branch       = "main"
  database     = planetscale_database.test.name
  name         = "test"
  organization = planetscale_database.test.organization
  role         = "admin"
}

data "planetscale_passwords" "test" {
  branch       = planetscale_password.test.branch
  database     = planetscale_password.test.database
  organization = planetscale_password.test.organization
}
