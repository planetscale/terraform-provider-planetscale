data "planetscale_organization_cluster_size_skus" "my_organizationclustersizeskus" {
  engine = "mysql"
  name   = "...my_name..."
  rates  = true
  region = "...my_region..."
}