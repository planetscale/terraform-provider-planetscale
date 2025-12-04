resource "planetscale_database_vitess" "my_databasevitess" {
  cluster_size = "...my_cluster_size..."
  database     = "...my_database..."
  name         = "...my_name..."
  organization = "...my_organization..."
  region       = "...my_region..."
  replicas     = 6.94
}