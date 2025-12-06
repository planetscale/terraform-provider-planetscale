resource "planetscale_database_postgres" "my_databasepostgres" {
  cluster_size  = "...my_cluster_size..."
  major_version = "...my_major_version..."
  name          = "...my_name..."
  organization  = "...my_organization..."
  region        = "...my_region..."
  replicas      = 8
}