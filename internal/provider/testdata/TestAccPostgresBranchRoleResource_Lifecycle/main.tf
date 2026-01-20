variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "role_name" {
  type = string
}

resource "planetscale_postgres_branch_role" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name
  name         = var.role_name
}
