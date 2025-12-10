data "planetscale_branches" "my_branches" {
  database        = "...my_database..."
  order           = "asc"
  organization    = "...my_organization..."
  production      = false
  q               = "...my_q..."
  safe_migrations = true
}