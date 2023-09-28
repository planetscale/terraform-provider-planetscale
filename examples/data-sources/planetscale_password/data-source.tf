data "planetscale_password" "example" {
  organization = "example"
  database     = "example_db"
  branch       = "main"
  name         = "antoine-was-here"
}

output "password" {
  value = data.planetscale_password.example
}