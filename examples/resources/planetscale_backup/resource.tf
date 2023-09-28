resource "planetscale_backup" "example" {
  organization = "example"
  database     = "example_db"
  branch       = "main"
  name         = "antoine_was_here"
  backup_policy = {
    retention_unit  = "day"
    retention_value = 1
  }
}