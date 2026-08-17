variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "router_name" {
  type = string
}

variable "replicas_per_cell" {
  type = number
}

resource "planetscale_neki_branch" "main" {
  organization = var.organization
  database     = var.database_name
  name         = "main"
  cluster_size = "PS_DEV_AWS_ARM"
}

resource "planetscale_neki_router" "test" {
  organization      = var.organization
  database          = planetscale_neki_branch.main.database
  branch            = planetscale_neki_branch.main.name
  name              = var.router_name
  router_size       = "NKR_5"
  replicas_per_cell = var.replicas_per_cell
}
