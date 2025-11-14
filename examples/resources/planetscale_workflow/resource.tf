resource "planetscale_workflow" "my_workflow" {
  database             = "...my_database..."
  defer_secondary_keys = true
  global_keyspace      = "...my_global_keyspace..."
  name                 = "...my_name..."
  on_ddl               = "EXEC"
  organization         = "...my_organization..."
  source_keyspace      = "...my_source_keyspace..."
  tables = [
    "..."
  ]
  target_keyspace = "...my_target_keyspace..."
}