variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "policy_name" {
  type = string
}

resource "planetscale_postgres_backup_policy" "test" {
  organization    = var.organization
  database        = var.database_name
  name            = var.policy_name
  target          = "development"
  retention_value = 7
  retention_unit  = "day"
  frequency_value = 1
  frequency_unit  = "day"
  schedule_time   = "04:00"
}

data "planetscale_postgres_backup_policies" "test" {
  organization = var.organization
  database     = var.database_name
}
