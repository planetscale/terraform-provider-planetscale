variable "cluster_size" {
  type = string
}

variable "name" {
  type = string
}

variable "organization" {
  type = string
}

resource "planetscale_database_vitess" "test" {
  cluster_size = var.cluster_size
  name         = var.name
  organization = var.organization
}


