data "planetscale_organizations" "test" {}

data "planetscale_organization_cluster_size_skus" "test" {
  organization = data.planetscale_organizations.test.data[0].name
}
