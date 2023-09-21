data "planetscale_database_regions" "example" {
  organization = "example.com"
  name         = "example_db"
}

output "database_regions" {
  value = data.planetscale_database_regions.example
}