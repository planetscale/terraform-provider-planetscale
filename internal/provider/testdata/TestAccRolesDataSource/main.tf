variable "database_name" {
  type = string
}

variable "organization" {
  type = string
}

resource "planetscale_database" "test" {
  cluster_size = "PS_10_AWS_ARM"
  database     = var.database_name
  kind         = "postgresql"
  name         = var.database_name
  organization = var.organization
}

resource "planetscale_role" "test" {
  branch       = "main"
  database     = planetscale_database.test.name
  organization = planetscale_database.test.organization
}

data "planetscale_roles" "test" {
  branch       = planetscale_role.test.branch
  database     = planetscale_role.test.database
  organization = planetscale_role.test.organization
}
