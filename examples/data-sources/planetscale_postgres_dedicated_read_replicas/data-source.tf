data "planetscale_postgres_dedicated_read_replicas" "my_postgresdedicatedreadreplicas" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  organization = "...my_organization..."
}