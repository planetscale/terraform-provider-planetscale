# doesn't work right now for some reason

data "planetscale_user" "example" {}

output "current_user" {
  value = data.planetscale_user.example
}