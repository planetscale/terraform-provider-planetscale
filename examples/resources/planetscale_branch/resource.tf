resource "planetscale_branch" "example" {
  organization  = "example"
  database      = "example_db"
  name          = "antoinewritescode"
  parent_branch = "main"
}