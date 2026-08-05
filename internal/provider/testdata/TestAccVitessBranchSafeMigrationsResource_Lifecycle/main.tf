variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "safe_migrations" {
  type = bool
}

resource "planetscale_vitess_branch" "test" {
  organization = var.organization
  database     = var.database_name
  name         = "main"
}

resource "planetscale_vitess_branch_safe_migrations" "test" {
  organization    = var.organization
  database        = var.database_name
  branch          = planetscale_vitess_branch.test.id
  safe_migrations = var.safe_migrations
}
