resource "planetscale_vitess_branch_password" "my_vitessbranchpassword" {
  branch = "...my_branch..."
  cidrs = [
    "..."
  ]
  database      = "...my_database..."
  direct_vtgate = true
  name          = "...my_name..."
  organization  = "...my_organization..."
  replica       = false
  role          = "admin"
  ttl           = 8.45
}