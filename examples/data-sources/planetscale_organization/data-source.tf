data "planetscale_organization" "example" {
  name = "example"
}

output "org" {
  value = data.planetscale_organization.example
}