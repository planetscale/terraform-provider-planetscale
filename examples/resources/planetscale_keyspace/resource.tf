resource "planetscale_keyspace" "my_keyspace" {
  branch       = "...my_branch..."
  cluster_size = "...my_cluster_size..."
  database     = "...my_database..."
  name         = "...my_name..."
  organization = "...my_organization..."
}