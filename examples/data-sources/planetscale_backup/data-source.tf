data "planetscale_backup" "example" {
  organization = "example.com"
  database     = "example_db"
  branch       = "main"
  id           = "k20nb1b7ut18"
}

output "backup" {
  value = data.planetscale_backup.example
}