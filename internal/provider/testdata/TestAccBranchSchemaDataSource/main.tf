variable "name" {
  type = string
}

variable "organization" {
  type = string
}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  database     = var.name
  name         = var.name
  organization = var.organization
}

data "planetscale_branch_schema" "test" {
  database     = planetscale_database.test.name
  name         = "main"
  organization = planetscale_database.test.organization
}
