variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "role" {
  type = string
}

variable "cidrs" {
  type = list(string)
}

variable "description" {
  type = string
}

resource "planetscale_postgres_ip_restriction" "test" {
  organization = var.organization
  database     = var.database_name

  role        = var.role
  cidrs       = var.cidrs
  description = var.description
}

data "planetscale_postgres_ip_restriction" "test" {
  organization = var.organization
  database     = var.database_name
  id           = planetscale_postgres_ip_restriction.test.id
}

data "planetscale_postgres_ip_restrictions" "test" {
  organization = var.organization
  database     = var.database_name

  depends_on = [planetscale_postgres_ip_restriction.test]
}
