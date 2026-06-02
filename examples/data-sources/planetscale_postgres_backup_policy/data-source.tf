data "planetscale_postgres_backup_policy" "my_postgresbackuppolicy" {
  database     = "...my_database..."
  id           = "...my_id..."
  organization = "...my_organization..."
}