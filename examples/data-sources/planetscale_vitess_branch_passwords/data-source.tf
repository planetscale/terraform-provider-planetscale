data "planetscale_vitess_branch_passwords" "my_vitessbranchpasswords" {
  branch              = "...my_branch..."
  database            = "...my_database..."
  organization        = "...my_organization..."
  q                   = "...my_q..."
  read_only_region_id = "...my_read_only_region_id..."
  status              = "...my_status..."
}