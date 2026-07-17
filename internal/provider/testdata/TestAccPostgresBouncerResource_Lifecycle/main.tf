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
  type    = string
  default = null
}

resource "planetscale_postgres_bouncer" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name

  name              = var.bouncer_name
  target            = "primary"
  bouncer_size      = var.bouncer_size
  replicas_per_cell = 1

  # A bouncer only accepts one configuration change at a time, and changes
  # take effect over hours, so the test makes a single update: parameters are
  # added in the same step that resizes the bouncer.
  parameters = var.pool_size == null ? null : {
    pgbouncer = {
      default_pool_size = var.pool_size
    }
  }
}
