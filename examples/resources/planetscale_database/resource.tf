resource "planetscale_database" "my_database" {
  cluster_size = "...my_cluster_size..."
  kind         = "mysql"
  name         = "...my_name..."
  organization = "...my_organization..."
  region       = "...my_region..."
}