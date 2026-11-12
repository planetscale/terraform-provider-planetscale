# Allow connections to every schema, as any role, from these ranges.
resource "planetscale_postgres_ip_restriction" "office" {
  organization = "my-organization"
  database     = "my-database"

  cidrs       = ["203.0.113.0/24", "198.51.100.7/32"]
  description = "Office network and VPN gateway"
}

# Limit one role to the addresses of the hosts that use it.
resource "planetscale_postgres_ip_restriction" "reporting" {
  organization = "my-organization"
  database     = "my-database"

  role  = "reporting"
  cidrs = [for ip in ["192.0.2.10", "192.0.2.11"] : "${ip}/32"]
}
