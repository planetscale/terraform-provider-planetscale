data "planetscale_database_read_only_regions" "example" {
  organization = "example.com"
  name         = "example_db"
}

output "database_ro_regions" {
  value = data.planetscale_database_regions.example
}