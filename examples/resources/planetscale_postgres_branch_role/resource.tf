resource "planetscale_postgres_branch_role" "my_postgresbranchrole" {
  branch   = "...my_branch..."
  database = "...my_database..."
  inherited_roles = [
    "pg_maintain"
  ]
  organization = "...my_organization..."
  successor    = "...my_successor..."
  ttl          = 1
}