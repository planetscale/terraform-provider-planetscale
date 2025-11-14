variable "database_name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10_AWS_ARM"
  kind         = "postgresql"
  name         = var.database_name
  organization = data.planetscale_organizations.test.data[0].name
}

resource "planetscale_cidrs" "test" {
  cidrs         = ["192.168.1.0/24"]
  database_name = planetscale_database.test.name
  organization  = planetscale_database.test.organization
}
