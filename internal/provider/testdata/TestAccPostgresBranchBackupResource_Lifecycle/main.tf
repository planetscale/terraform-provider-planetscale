variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "backup_name" {
  type = string
}

resource "planetscale_postgres_branch_backup" "test" {
  organization    = var.organization
  database        = var.database_name
  branch          = var.branch_name
  name            = var.backup_name
  retention_value = 1
  retention_unit  = "day"
}
