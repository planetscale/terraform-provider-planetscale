data "planetscale_database" "example" {
  organization = "example"
  name         = "again"
}

output "db" {
  value = data.planetscale_database.example
}