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

resource "planetscale_backup" "test" {
  branch          = "main"
  database        = planetscale_database.test.name
  name            = "test"
  organization    = planetscale_database.test.organization
  retention_unit  = "day"
  retention_value = 7
}

data "planetscale_backups" "test" {
  branch       = "main"
  database     = planetscale_backup.test.database
  organization = planetscale_backup.test.organization
}
