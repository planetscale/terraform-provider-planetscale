variable "cluster_size" {
  type = string
}

variable "name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = var.cluster_size
  database     = var.name
  name         = var.name
  organization = data.planetscale_organizations.test.data[0].name
}
