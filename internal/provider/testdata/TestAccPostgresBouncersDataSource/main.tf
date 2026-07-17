variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "bouncer_name" {
  type = string
}

resource "planetscale_postgres_bouncer" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name

  name         = var.bouncer_name
  target       = "primary"
  bouncer_size = "PGB_5"
}

data "planetscale_postgres_bouncer" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name
  name         = planetscale_postgres_bouncer.test.name
}

data "planetscale_postgres_bouncers" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name

  depends_on = [planetscale_postgres_bouncer.test]
}
