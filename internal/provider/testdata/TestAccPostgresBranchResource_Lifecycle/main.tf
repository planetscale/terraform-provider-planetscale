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

variable "minimum_storage_bytes" {
  type = number
}

variable "maximum_storage_bytes" {
  type = number
}

variable "storage_autoscaling" {
  type = bool
}

variable "storage_iops" {
  type = number
}

variable "storage_throughput_mibs" {
  type = number
}

data "planetscale_database_postgres" "test" {
  organization = var.organization
  id           = var.database_name
}

resource "planetscale_postgres_branch" "test" {
  organization       = var.organization
  database           = var.database_name
  name               = var.branch_name
  cluster_size       = var.cluster_size
  deletion_protected = var.deletion_protected
  parameters         = var.parameters
  region             = data.planetscale_database_postgres.test.region_data.slug

  storage = {
    minimum_storage_bytes   = var.minimum_storage_bytes
    maximum_storage_bytes   = var.maximum_storage_bytes
    storage_autoscaling     = var.storage_autoscaling
    storage_iops            = var.storage_iops
    storage_throughput_mibs = var.storage_throughput_mibs
  }
}
