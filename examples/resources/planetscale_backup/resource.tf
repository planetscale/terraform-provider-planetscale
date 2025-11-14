resource "planetscale_backup" "my_backup" {
  branch          = "...my_branch..."
  database        = "...my_database..."
  emergency       = true
  name            = "...my_name..."
  organization    = "...my_organization..."
  retention_unit  = "day"
  retention_value = 4.09
}