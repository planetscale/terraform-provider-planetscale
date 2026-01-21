resource "planetscale_postgres_branch" "my_postgresbranch" {
  organization  = "my-organization"
  database      = "ru00w3vqvfr9"

  name          = "my-branch"
  cluster_size  = "PS-10"
}
