data "planetscale_passwords" "example" {
  organization = "example"
  database     = "example_db"
  branch       = "main"
}

output "passwords" {
  value = data.planetscale_passwords.example
}