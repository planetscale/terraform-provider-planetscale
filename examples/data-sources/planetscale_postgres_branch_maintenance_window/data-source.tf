data "planetscale_postgres_branch_maintenance_window" "my_postgresbranchmaintenancewindow" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  organization = "...my_organization..."
}