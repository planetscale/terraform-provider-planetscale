variable "organization" {
  type = string
}

data "planetscale_organization_regions" "test" {
  organization = var.organization
}
