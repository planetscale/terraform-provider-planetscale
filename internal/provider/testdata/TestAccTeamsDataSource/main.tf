data "planetscale_organizations" "test" {}

data "planetscale_teams" "test" {
  organization_name = data.planetscale_organizations.test.data[0].name
}
