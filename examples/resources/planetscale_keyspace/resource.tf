resource "planetscale_keyspace" "my_keyspace" {
  branch         = "...my_branch..."
  cluster_size   = "...my_cluster_size..."
  database       = "...my_database..."
  extra_replicas = 2.01
  name           = "...my_name..."
  organization   = "...my_organization..."
  shards         = 9.66
}