resource "planetscale_vitess_branch" "my_vitessbranch" {
  organization  = "my-organization"
  database      = "ru00w3vqvfr9"

  name          = "my-branch"
}