resource "planetscale_password" "my_password" {
  branch = "...my_branch..."
  cidrs = [
    "..."
  ]
  database      = "...my_database..."
  direct_vtgate = false
  name          = "...my_name..."
  organization  = "...my_organization..."
  replica       = true
  role          = "admin"
  ttl           = 8.92
}