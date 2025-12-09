data "planetscale_vitess_branch_passwords" "my_vitessbranchpasswords" {
  branch              = "...my_branch..."
  database            = "...my_database..."
  organization        = "...my_organization..."
  read_only_region_id = "...my_read_only_region_id..."
}