resource "planetscale_branch" "my_branch" {
  backup_id     = "...my_backup_id..."
  branch        = "...my_branch..."
  cluster_size  = "...my_cluster_size..."
  database      = "...my_database..."
  major_version = "...my_major_version..."
  organization  = "...my_organization..."
  parent_branch = "...my_parent_branch..."
  region        = "...my_region..."
  restore_point = "...my_restore_point..."
  seed_data     = "last_successful_backup"
}