# requires a feature flag, contact support to enable it

data "planetscale_oauth_applications" "example" {
  organization = data.planetscale_organization.example.name
}

output "oauth_apps" {
  value = data.planetscale_oauth_applications.example
}