variable "team_name" {
  type = string
}

data "planetscale_organizations" "test" {}

resource "planetscale_team" "test" {
  name              = var.team_name
  organization_name = data.planetscale_organizations.test.data[0].name
  description       = "Terraform testing team"
}
