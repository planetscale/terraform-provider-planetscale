data "planetscale_branch_schema_lint" "example" {
  organization = "example.com"
  database     = "example_db"
  branch       = "main"
}

output "schema_lint" {
  value = data.planetscale_branch_schema_lint.example
}