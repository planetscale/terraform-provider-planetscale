resource "planetscale_postgres_branch_maintenance_window" "my_postgresbranchmaintenancewindow" {
  organization = "my-organization"
  database     = "my-database"
  branch       = "main"

  # Sundays at 03:00 UTC
  day  = 0
  hour = 3
}
