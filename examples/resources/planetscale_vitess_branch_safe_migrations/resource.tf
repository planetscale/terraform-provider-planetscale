resource "planetscale_vitess_branch" "main" {
  organization = "my-organization"
  database     = "my-database"
  name         = "main"
}

resource "planetscale_vitess_branch_safe_migrations" "main" {
  organization    = planetscale_vitess_branch.main.organization
  database        = planetscale_vitess_branch.main.database
  branch          = planetscale_vitess_branch.main.id
  safe_migrations = true
}
