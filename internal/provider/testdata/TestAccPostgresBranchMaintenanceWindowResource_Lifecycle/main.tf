variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "day" {
  type = number
}

variable "hour" {
  type = number
}

resource "planetscale_postgres_branch_maintenance_window" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name

  day  = var.day
  hour = var.hour
}

data "planetscale_postgres_branch_maintenance_window" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name

  # Read after the resource applies, so the data source sees the new window.
  depends_on = [planetscale_postgres_branch_maintenance_window.test]
}
