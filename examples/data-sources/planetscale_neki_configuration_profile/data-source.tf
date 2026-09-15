data "planetscale_neki_configuration_profile" "my_nekiconfigurationprofile" {
  branch   = "...my_branch..."
  database = "...my_database..."
  extensions = [
    "..."
  ]
  name         = "...my_name..."
  organization = "...my_organization..."
  parameters = {
    key = {
      # ...
    }
  }
}