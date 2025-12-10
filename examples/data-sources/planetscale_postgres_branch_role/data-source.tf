data "planetscale_postgres_branch_role" "my_postgresbranchrole" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  id           = "...my_id..."
  organization = "...my_organization..."
}