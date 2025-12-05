variable "name" {
  type = string
}

variable "organization" {
  type = string
}

resource "planetscale_database_postgres" "test" {
  cluster_size = "PS_5_AWS_X86"
  name         = var.name
  organization = var.organization
}

