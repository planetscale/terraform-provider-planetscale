resource "planetscale_postgres_branch" "my_postgresbranch" {
  organization  = "my-organization"
  database      = "ru00w3vqvfr9"

  name          = "my-branch"
  cluster_size  = "PS_10_AWS_ARM"

  # Postgres parameter overrides, nested by namespace (pgconf, pgbouncer,
  # patroni). Omitted parameters are reset to their defaults.
  parameters = {
    pgconf = {
      max_connections = "200"
    }
  }
}
