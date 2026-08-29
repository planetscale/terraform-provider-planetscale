data "planetscale_postgres_read_only_replicas" "my_postgresreadonlyreplicas" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  organization = "...my_organization..."
}