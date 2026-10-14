data "planetscale_databases" "my_databases" {
  kind         = "neki"
  organization = "...my_organization..."
  q            = "...my_q..."
}