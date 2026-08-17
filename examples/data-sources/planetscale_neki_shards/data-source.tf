data "planetscale_neki_shards" "my_nekishards" {
  branch                        = "...my_branch..."
  database                      = "...my_database..."
  exclude_configuration_profile = "...my_exclude_configuration_profile..."
  organization                  = "...my_organization..."
  q                             = "...my_q..."
}