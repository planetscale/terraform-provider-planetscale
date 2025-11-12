data "planetscale_organizations" "test" {}

data "planetscale_databases" "test" {
  organization = data.planetscale_organizations.test.data[0].name
}
