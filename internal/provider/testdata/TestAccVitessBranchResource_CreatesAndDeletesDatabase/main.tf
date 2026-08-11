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

variable "safe_migrations" {
  type = bool
}

variable "vtgate_autoscaling" {
  type = bool
}

variable "vtgate_count" {
  type = number
}

variable "vtgate_max_count" {
  type = number
}

variable "vtgate_size" {
  type = string
}

variable "vtgate_target_cpu_utilization" {
  type = number
}

resource "planetscale_vitess_branch" "test" {
  organization                  = var.organization
  database                      = var.database_name
  name                          = var.branch_name
  cluster_size                  = var.cluster_size
  safe_migrations               = var.safe_migrations
  vtgate_autoscaling            = var.vtgate_autoscaling
  vtgate_count                  = var.vtgate_count
  vtgate_max_count              = var.vtgate_max_count
  vtgate_size                   = var.vtgate_size
  vtgate_target_cpu_utilization = var.vtgate_target_cpu_utilization
}
