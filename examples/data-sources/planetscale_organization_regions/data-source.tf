data "planetscale_organization_regions" "example" {
  organization = "example"
}

output "org_regions" {
  value = data.planetscale_organization_regions.example
}
