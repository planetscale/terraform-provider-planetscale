variable "organization" {
  type = string
}

data "planetscale_databases" "test" {
  organization = var.organization
}
