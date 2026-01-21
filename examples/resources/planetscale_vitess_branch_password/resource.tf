resource "planetscale_vitess_branch_password" "my_vitessbranchpassword" {
  organization = "my-organization"
  database = "ru00w3vqvfr9"
  branch   = "2474dzfubrf3"

  role          = "admin"
}