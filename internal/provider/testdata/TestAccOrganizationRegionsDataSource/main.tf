data "planetscale_organizations" "test" {}

data "planetscale_organization_regions" "test" {
  organization = data.planetscale_organizations.test.data[0].name
}
