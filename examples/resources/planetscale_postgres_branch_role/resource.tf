resource "planetscale_postgres_branch_role" "my_postgresbranchrole" {
  organization = "my-organization"
  database = "ru00w3vqvfr9"
  branch   = "2474dzfubrf3"

  name         = "application-role"

  inherited_roles = [
    "pg_read_all_data",
    "pg_write_all_data",
  ]
}