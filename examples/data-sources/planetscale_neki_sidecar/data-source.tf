data "planetscale_neki_sidecar" "my_nekisidecar" {
  branch                = "...my_branch..."
  configuration_profile = "...my_configuration_profile..."
  database              = "...my_database..."
  organization          = "...my_organization..."
  parameters = {
    key = {
      # ...
    }
  }
}