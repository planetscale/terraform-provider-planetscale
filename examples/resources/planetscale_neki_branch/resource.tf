resource "planetscale_neki_branch" "my_neki_branch" {
  organization = "my-organization"
  database     = "my-neki-database"

  name               = "main"
  deletion_protected = true
}
