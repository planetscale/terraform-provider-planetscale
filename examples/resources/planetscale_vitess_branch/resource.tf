resource "planetscale_vitess_branch" "my_vitessbranch" {
  backup_id     = "...my_backup_id..."
  database      = "...my_database..."
  name          = "...my_name..."
  organization  = "...my_organization..."
  parent_branch = "...my_parent_branch..."
  region        = "...my_region..."
  seed_data     = "last_successful_backup"
}