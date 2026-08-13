variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "profile_name" {
  type = string
}

variable "cluster_size" {
  type = string
}

resource "planetscale_neki_branch" "main" {
  organization = var.organization
  database     = var.database_name
  name         = "main"
  cluster_size = var.cluster_size
}

resource "planetscale_neki_configuration_profile" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.main.database
  branch       = planetscale_neki_branch.main.name
  name         = var.profile_name
  cluster_size = var.cluster_size
}
