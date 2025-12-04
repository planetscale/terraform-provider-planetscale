data "planetscale_organization_cluster_size_skus" "my_organizationclustersizeskus" {
  engine       = "mysql"
  organization = "...my_organization..."
  rates        = true
  region       = "...my_region..."
}