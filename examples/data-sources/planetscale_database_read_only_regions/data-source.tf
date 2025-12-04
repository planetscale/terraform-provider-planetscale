data "planetscale_database_read_only_regions" "my_databasereadonlyregions" {
  database     = "...my_database..."
  organization = "...my_organization..."
}