variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "policy_name" {
  type = string
}

resource "planetscale_neki_branch" "main" {
  organization = var.organization
  database     = var.database_name
  name         = "main"
  cluster_size = "PS_DEV_AWS_ARM"
}

resource "planetscale_neki_backup_policy" "test" {
  organization    = var.organization
  database        = planetscale_neki_branch.main.database
  name            = var.policy_name
  target          = "development"
  retention_value = 7
  retention_unit  = "day"
  frequency_value = 1
  frequency_unit  = "day"
  schedule_time   = "04:00"
}
