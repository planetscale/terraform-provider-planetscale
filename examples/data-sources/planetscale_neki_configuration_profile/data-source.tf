data "planetscale_neki_configuration_profile" "my_nekiconfigurationprofile" {
  organization = "acme"
  database = "app"
  branch = "main"
  name = "default"
}

output "enabled_extensions" {
  value = data.planetscale_neki_configuration_profile.my_nekiconfigurationprofile.extensions
}
