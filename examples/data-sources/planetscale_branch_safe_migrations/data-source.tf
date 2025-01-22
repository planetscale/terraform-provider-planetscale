data "planetscale_branch_safe_migrations" "example" {
  organization = "example.com"
  database     = "example_db"
  branch       = "main"
}

output "safe_migrations" {
  value = data.planetscale_branch_safe_migrations.example
}
