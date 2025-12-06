resource "planetscale_role" "my_role" {
  branch   = "...my_branch..."
  database = "...my_database..."
  inherited_roles = [
    "pg_checkpoint"
  ]
  organization = "...my_organization..."
  successor    = "...my_successor..."
  ttl          = 5
}