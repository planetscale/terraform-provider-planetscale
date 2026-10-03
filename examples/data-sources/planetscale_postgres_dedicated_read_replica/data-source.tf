data "planetscale_postgres_dedicated_read_replica" "my_postgresdedicatedreadreplica" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  name         = "...my_name..."
  organization = "...my_organization..."
}