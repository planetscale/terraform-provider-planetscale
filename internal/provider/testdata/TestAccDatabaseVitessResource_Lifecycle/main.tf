variable "name" {
  type = string
}

variable "organization" {
  type = string
}

resource "planetscale_database_vitess" "test" {
  cluster_size = "PS_10"
  name         = var.name
  organization = var.organization
}


