variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "vtgate_autoscaling" {
  type = bool
}

resource "planetscale_vitess_branch" "test" {
  organization = var.organization
  database     = var.database_name
  name         = "main"
  cluster_size = "PS_10"
}

resource "planetscale_vitess_branch_vtgate" "test" {
  organization                  = var.organization
  database                      = var.database_name
  branch                        = planetscale_vitess_branch.test.id
  vtgate_autoscaling            = var.vtgate_autoscaling
  vtgate_size                   = "VTG_320"
  vtgate_count                  = 1
  vtgate_max_count              = 2
  vtgate_target_cpu_utilization = 50
}
