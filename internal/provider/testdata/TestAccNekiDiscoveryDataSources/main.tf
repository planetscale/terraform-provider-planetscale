variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

resource "planetscale_neki_branch" "main" {
  organization = var.organization
  database     = var.database_name
  name         = "main"
  cluster_size = "PS_DEV_AWS_ARM"
}

data "planetscale_neki_routers" "all" {
  organization = var.organization
  database     = planetscale_neki_branch.main.database
  branch       = planetscale_neki_branch.main.name
}

data "planetscale_neki_router_parameters" "default" {
  organization = var.organization
  database     = planetscale_neki_branch.main.database
  branch       = planetscale_neki_branch.main.name
  router       = "default"
}

data "planetscale_neki_configuration_profile_parameters" "default" {
  organization          = var.organization
  database              = planetscale_neki_branch.main.database
  branch                = planetscale_neki_branch.main.name
  configuration_profile = "default"
}

data "planetscale_neki_shards" "all" {
  organization = var.organization
  database     = planetscale_neki_branch.main.database
  branch       = planetscale_neki_branch.main.name
}

data "planetscale_neki_configuration_profile_shards" "default" {
  organization          = var.organization
  database              = planetscale_neki_branch.main.database
  branch                = planetscale_neki_branch.main.name
  configuration_profile = "default"
}
