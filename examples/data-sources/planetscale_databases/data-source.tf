data "planetscale_databases" "example" {
  organization = "example"
}

output "dbs" {
  value = data.planetscale_databases.example
}