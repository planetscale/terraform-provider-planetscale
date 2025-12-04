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

resource "planetscale_database_default_branch" "test" {
  branch       = planetscale_branch.test.name
  database     = planetscale_branch.test.database
  organization = planetscale_branch.test.organization
}
