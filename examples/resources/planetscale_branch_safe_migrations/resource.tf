# Enable safe_migrations on default branch:
resource "planetscale_database" "example" {
  organization   = "example"
  name           = "example"
  default_branch = "main"
}

resource "planetscale_branch_safe_migrations" "main" {
  organization = planetscale_branch.example.organization
  database     = planetscale_branch.example.database
  branch       = planetscale_branch.example.default_branch
  enabled      = true
}

# Enable safe_migrations on a branch:
resource "planetscale_branch" "staging" {
  organization  = planetscale_branch.example.organization
  database      = planetscale_branch.example.database
  parent_branch = planetscale_branch.example.default_branch
  name          = "staging"
}

resource "planetscale_branch_safe_migrations" "staging" {
  database     = planetscale_database.example.name
  organization = planetscale_database.example.organization
  branch       = planetscale_branch.staging.name
  enabled      = true
}
