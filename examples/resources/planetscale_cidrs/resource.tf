resource "planetscale_cidrs" "my_cidrs" {
  cidrs = [
    "..."
  ]
  database_name = "...my_database_name..."
  organization  = "...my_organization..."
  role          = "...my_role..."
  schema        = "...my_schema..."
}