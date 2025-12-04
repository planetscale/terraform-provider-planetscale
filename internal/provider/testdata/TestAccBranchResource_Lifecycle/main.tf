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

resource "planetscale_branch" "test" {
  branch        = "test"
  database      = planetscale_database.test.name
  organization  = planetscale_database.test.organization
  parent_branch = "main"
}
