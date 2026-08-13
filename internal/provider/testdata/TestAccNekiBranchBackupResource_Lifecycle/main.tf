variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "backup_name" {
  type = string
}

resource "planetscale_neki_branch" "main" {
  organization = var.organization
  database     = var.database_name
  name         = "main"
  cluster_size = "PS_DEV_AWS_ARM"
}

resource "planetscale_neki_branch_backup" "test" {
  organization    = var.organization
  database        = planetscale_neki_branch.main.database
  branch          = planetscale_neki_branch.main.name
  name            = var.backup_name
  retention_value = 1
  retention_unit  = "day"
}
