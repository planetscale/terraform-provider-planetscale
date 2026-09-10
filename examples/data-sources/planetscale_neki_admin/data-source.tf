data "planetscale_neki_admin" "my_nekiadmin" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  organization = "...my_organization..."
  parameters = {
    key = {
      # ...
    }
  }
}