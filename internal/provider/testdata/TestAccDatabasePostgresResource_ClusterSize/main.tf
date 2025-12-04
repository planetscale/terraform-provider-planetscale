variable "cluster_size" {
  type = string
}

variable "name" {
  type = string
}

variable "organization" {
  type = string
}

resource "planetscale_database_postgres" "test" {
  cluster_size = var.cluster_size
  database     = var.name
  name         = var.name
  organization = var.organization
}

