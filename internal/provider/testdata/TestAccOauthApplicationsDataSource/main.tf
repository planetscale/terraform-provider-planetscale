data "planetscale_organizations" "test" {}

data "planetscale_oauth_applications" "test" {
  organization = data.planetscale_organizations.test.data[0].name
}
