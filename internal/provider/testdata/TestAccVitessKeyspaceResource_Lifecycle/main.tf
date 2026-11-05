variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "keyspace_name" {
  type = string
}

variable "cluster_size" {
  type = string
}

variable "extra_replicas" {
  type = number
}

resource "planetscale_vitess_branch" "main" {
  organization = var.organization
  database     = var.database_name
  name         = var.branch_name
  # Default/main keyspace stays at PS_10; the Acc test manages an additional
  # keyspace (create + extra_replicas resize + import).
  cluster_size = "PS_10"
}

resource "planetscale_vitess_keyspace" "test" {
  organization   = var.organization
  database       = planetscale_vitess_branch.main.database
  branch         = planetscale_vitess_branch.main.name
  name           = var.keyspace_name
  cluster_size   = var.cluster_size
  extra_replicas = var.extra_replicas
}
