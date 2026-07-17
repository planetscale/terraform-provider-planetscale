resource "planetscale_postgres_bouncer" "my_bouncer" {
  organization = "my-organization"
  database     = "my-database"
  branch       = "main"

  name              = "my-bouncer"
  target            = "primary"
  bouncer_size      = "PGB_5"
  replicas_per_cell = 1

  parameters = {
    pgbouncer = {
      default_pool_size = "100"
    }
  }
}
