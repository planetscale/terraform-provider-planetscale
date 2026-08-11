resource "planetscale_vitess_branch" "main" {
  organization = "my-organization"
  database     = "my-database"
  name         = "main"
}

resource "planetscale_vitess_branch_vtgate" "main" {
  organization                  = planetscale_vitess_branch.main.organization
  database                      = planetscale_vitess_branch.main.database
  branch                        = planetscale_vitess_branch.main.id
  vtgate_autoscaling            = true
  vtgate_size                   = "VTG_320"
  vtgate_count                  = 2
  vtgate_max_count              = 8
  vtgate_target_cpu_utilization = 50
}
