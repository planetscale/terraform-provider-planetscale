
data "planetscale_branch" "example" {
  organization = "example.com"
  database     = "example_db"
  name         = "main"
}

output "branch" {
  value = data.planetscale_branch.example
}