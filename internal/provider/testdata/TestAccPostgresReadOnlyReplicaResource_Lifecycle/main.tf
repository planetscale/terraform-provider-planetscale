variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "replica_name" {
  type = string
}

variable "replicas" {
  type    = number
  default = null
}

# The replica runs in the primary's region so the test does not need to
# hardcode a region slug.
data "planetscale_database_postgres" "test" {
  organization = var.organization
  id           = var.database_name
}

resource "planetscale_postgres_read_only_replica" "test" {
  organization = var.organization
  database     = var.database_name
  branch       = var.branch_name

  name     = var.replica_name
  region   = data.planetscale_database_postgres.test.region_data.slug
  replicas = var.replicas
}
