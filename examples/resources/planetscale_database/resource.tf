resource "planetscale_database" "my_database" {
  cluster_size  = "...my_cluster_size..."
  database      = "...my_database..."
  kind          = "mysql"
  major_version = "...my_major_version..."
  name          = "...my_name..."
  organization  = "...my_organization..."
  region        = "...my_region..."
  replicas      = 3.39
}