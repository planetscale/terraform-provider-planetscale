data "planetscale_organizations" "example" {}

output "orgs" {
  value = data.planetscale_organizations.example
}