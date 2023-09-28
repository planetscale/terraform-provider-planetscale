resource "planetscale_password" "example" {
  organization = "example"
  database     = "example_db"
  branch       = "main"
  name         = "a-password-for-antoine"
}

output "password" {
  sensitive = true
  value     = planetscale_password.example
}