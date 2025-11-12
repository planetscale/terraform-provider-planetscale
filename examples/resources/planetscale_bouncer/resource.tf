resource "planetscale_bouncer" "my_bouncer" {
  bouncer_size      = "...my_bouncer_size..."
  branch            = "...my_branch..."
  database          = "...my_database..."
  name              = "...my_name..."
  organization      = "...my_organization..."
  replicas_per_cell = 1.89
  target            = "...my_target..."
}