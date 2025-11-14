resource "planetscale_webhook" "my_webhook" {
  database = "...my_database..."
  enabled  = false
  events = [
    "..."
  ]
  organization = "...my_organization..."
  url          = "...my_url..."
}