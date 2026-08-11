variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "safe_migrations" {
  type = bool
}

resource "planetscale_vitess_branch" "test" {
  organization    = var.organization
  database        = var.database_name
  name            = var.branch_name
  safe_migrations = var.safe_migrations
}
