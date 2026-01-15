variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

resource "planetscale_vitess_branch" "test" {
  organization  = var.organization
  database      = var.database_name
  name          = "test"
}
