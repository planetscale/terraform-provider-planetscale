resource "planetscale_vitess_branch" "my_vitessbranch" {
  organization = "my-organization"
  database     = "ru00w3vqvfr9"

  name                          = "my-branch"
  safe_migrations               = true
  vtgate_autoscaling            = true
  vtgate_count                  = 1
  vtgate_max_count              = 2
  vtgate_size                   = "VTG_320"
  vtgate_target_cpu_utilization = 50
}
