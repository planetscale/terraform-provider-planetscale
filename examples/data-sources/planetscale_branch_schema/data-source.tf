data "planetscale_branch_schema" "my_branchschema" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  keyspace     = "...my_keyspace..."
  namespace    = "...my_namespace..."
  organization = "...my_organization..."
}