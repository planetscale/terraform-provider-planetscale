resource "planetscale_database" "example" {
  organization = "example"
  name         = "anotherdb"
  cluster_size = "PS_10"
}
