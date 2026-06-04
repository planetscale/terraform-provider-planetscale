data "planetscale_postgres_branch_backup" "my_postgresbranchbackup" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  id           = "...my_id..."
  organization = "...my_organization..."
}