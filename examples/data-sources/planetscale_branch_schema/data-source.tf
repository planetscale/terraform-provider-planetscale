data "planetscale_branch_schema" "example" {
  organization = "example.com"
  database     = "example_db"
  branch       = "main"
}

output "branch_schema" {
  value = data.planetscale_branch_schema.example
}