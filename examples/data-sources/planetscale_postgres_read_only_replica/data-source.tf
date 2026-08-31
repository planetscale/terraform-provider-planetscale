data "planetscale_postgres_read_only_replica" "my_postgresreadonlyreplica" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  name         = "...my_name..."
  organization = "...my_organization..."
}