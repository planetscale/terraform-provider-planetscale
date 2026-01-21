variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

resource "planetscale_postgres_branch" "test" {
  organization = var.organization
  database     = var.database_name
  name         = var.branch_name
}
