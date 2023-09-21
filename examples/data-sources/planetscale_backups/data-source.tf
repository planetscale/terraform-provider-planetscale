data "planetscale_backups" "example" {
  organization = "example.com"
  database     = "example_db"
  branch       = "main"
}

output "backups" {
  value = data.planetscale_backups.example
}