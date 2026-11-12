data "planetscale_postgres_ip_restriction" "my_postgresiprestriction" {
  database     = "...my_database..."
  id           = "...my_id..."
  organization = "...my_organization..."
}