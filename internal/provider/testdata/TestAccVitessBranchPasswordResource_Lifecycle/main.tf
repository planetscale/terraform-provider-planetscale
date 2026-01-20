variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "password_name" {
  type = string
}

resource "planetscale_vitess_branch_password" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name
  name         = var.password_name
  role         = "admin"
}
