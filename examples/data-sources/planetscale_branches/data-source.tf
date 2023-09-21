data "planetscale_branches" "example" {
  organization = "example.com"
  database     = "example_db"
}

output "branches" {
  value = data.planetscale_branches.example
}