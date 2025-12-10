data "planetscale_vitess_branch_password" "my_vitessbranchpassword" {
  branch       = "...my_branch..."
  database     = "...my_database..."
  id           = "...my_id..."
  organization = "...my_organization..."
}