data "planetscale_organizations" "test" {}

data "planetscale_organization_cluster_size_skus" "test" {
  name = data.planetscale_organizations.test.data[0].name
}
