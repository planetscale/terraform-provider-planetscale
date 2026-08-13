variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "cluster_size" {
  type = string
}

variable "deletion_protected" {
  type = bool
}

variable "parameters" {
  type    = map(map(string))
  default = null
}

resource "planetscale_postgres_branch" "test" {
  organization       = var.organization
  database           = var.database_name
  name               = var.branch_name
  cluster_size       = var.cluster_size
  deletion_protected = var.deletion_protected
  parameters         = var.parameters
}
