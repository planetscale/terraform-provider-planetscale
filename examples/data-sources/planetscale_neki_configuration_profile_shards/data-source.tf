data "planetscale_neki_configuration_profile_shards" "my_nekiconfigurationprofileshards" {
  branch                = "...my_branch..."
  configuration_profile = "...my_configuration_profile..."
  database              = "...my_database..."
  organization          = "...my_organization..."
  q                     = "...my_q..."
}