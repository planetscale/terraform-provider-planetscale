variable "organization" {
  type = string
}

data "planetscale_organization" "test" {
  organization = var.organization
}
