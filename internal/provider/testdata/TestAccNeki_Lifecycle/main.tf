variable "organization" {
  type = string
}

variable "database_name" {
  type = string
}

variable "branch_name" {
  type = string
}

variable "role_name" {
  type = string
}

variable "policy_name" {
  type = string
}

variable "target" {
  type = string
}

variable "retention_value" {
  type = number
}

variable "retention_unit" {
  type = string
}

variable "frequency_value" {
  type = number
}

variable "frequency_unit" {
  type = string
}

variable "schedule_time" {
  type = string
}

variable "schedule_day" {
  type = number
}

variable "schedule_week" {
  type = number
}

# Provisioning the branch (and its database) is the slow part. Everything
# below runs against this one branch.
resource "planetscale_neki_branch" "test" {
  organization = var.organization
  database     = var.database_name
  name         = var.branch_name
}

resource "planetscale_neki_role" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
  name         = var.role_name
}

resource "planetscale_neki_backup_policy" "test" {
  organization    = var.organization
  database        = planetscale_neki_branch.test.database
  name            = var.policy_name
  target          = var.target
  retention_value = var.retention_value
  retention_unit  = var.retention_unit
  frequency_value = var.frequency_value
  frequency_unit  = var.frequency_unit
  schedule_time   = var.schedule_time
  schedule_day    = var.schedule_day
  schedule_week   = var.schedule_week
}

data "planetscale_neki_branch" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  id           = planetscale_neki_branch.test.id
}

# Branch provisioning creates a default admin, sidecar, configuration
# profile, shard, and router group; the data sources read those defaults.
data "planetscale_neki_admin" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
}

data "planetscale_neki_sidecar" "test" {
  organization          = var.organization
  database              = planetscale_neki_branch.test.database
  branch                = planetscale_neki_branch.test.name
  configuration_profile = "default"
}

data "planetscale_neki_configuration_profile" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
  name         = "default"
}

data "planetscale_neki_configuration_profiles" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
}

data "planetscale_neki_shards" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
}

data "planetscale_neki_shard" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
  id           = data.planetscale_neki_shards.test.data[0].id
}

data "planetscale_neki_router" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
  name         = "default"
}

data "planetscale_neki_routers" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
}

data "planetscale_neki_role" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
  id           = planetscale_neki_role.test.id
}

data "planetscale_neki_roles" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database
  branch       = planetscale_neki_branch.test.name
  q            = planetscale_neki_role.test.name
}

data "planetscale_neki_backup_policy" "test" {
  organization = var.organization
  database     = planetscale_neki_backup_policy.test.database
  id           = planetscale_neki_backup_policy.test.id
}

# No root outputs: the import steps below evaluate the config against a state
# holding only the imported resource, and terraform-plugin-testing cannot shim
# an output that evaluates to null there.
data "planetscale_neki_backup_policies" "test" {
  organization = var.organization
  database     = planetscale_neki_branch.test.database

  depends_on = [planetscale_neki_backup_policy.test]
}
