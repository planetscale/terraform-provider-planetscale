resource "planetscale_postgres_branch" "my_postgresbranch" {
  backup_id     = "...my_backup_id..."
  cluster_size  = "...my_cluster_size..."
  database      = "...my_database..."
  major_version = "...my_major_version..."
  name          = "...my_name..."
  organization  = "...my_organization..."
  parameters = {
    key = jsonencode("value")
  }
  parent_branch = "...my_parent_branch..."
  region        = "...my_region..."
  replicas      = 0
  restore_point = "...my_restore_point..."
}