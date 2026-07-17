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

variable "bouncer_size" {
  type = string
}

variable "pool_size" {
  type = string
}

resource "planetscale_postgres_bouncer" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name

  name              = var.bouncer_name
  target            = "primary"
  bouncer_size      = var.bouncer_size
  replicas_per_cell = 1

  parameters = {
    pgbouncer = {
      default_pool_size = var.pool_size
    }
  }
}
