variable "organization" {
  type = string
}

data "planetscale_organization_cluster_size_skus" "test" {
  organization = var.organization
}
