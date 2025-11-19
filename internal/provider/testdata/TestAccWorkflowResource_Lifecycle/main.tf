variable "database_name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_database" "test" {
  cluster_size = "PS_10"
  name         = var.database_name
  organization = data.planetscale_organizations.test.data[0].name
}

resource "planetscale_keyspace" "source" {
  branch       = "main"
  cluster_size = "PS_10"
  database     = planetscale_database.test.name
  name         = "source"
  organization = planetscale_database.test.organization
}

resource "planetscale_keyspace" "target" {
  branch       = "main"
  cluster_size = "PS_10"
  database     = planetscale_database.test.name
  name         = "target"
  organization = planetscale_database.test.organization
}

resource "planetscale_workflow" "test" {
  database        = planetscale_database.test.name
  name            = "test-workflow"
  organization    = planetscale_database.test.organization
  source_keyspace = planetscale_keyspace.source.name
  target_keyspace = planetscale_keyspace.target.name
  tables          = ["table1"]
}
