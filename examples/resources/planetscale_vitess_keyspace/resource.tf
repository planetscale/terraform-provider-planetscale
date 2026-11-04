# Creating a planetscale_vitess_branch already provisions the branch's default
# keyspace. Use this resource for additional keyspaces, or import the default
# keyspace if you need to manage its size, replicas, or parameters in
# Terraform.
#
# cluster_size and extra_replicas update in place (resize), and parameters
# update in place through config changes. Apply can take several minutes while
# provisioning, resizing, or a parameter rollout finishes.
resource "planetscale_vitess_branch" "example" {
  organization = "example"
  database     = "example"
  name         = "main"
  cluster_size = "PS_10"
}

resource "planetscale_vitess_keyspace" "example" {
  organization   = planetscale_vitess_branch.example.organization
  database       = planetscale_vitess_branch.example.database
  branch         = planetscale_vitess_branch.example.name
  name           = "metrics"
  cluster_size   = "PS_10"
  extra_replicas = 0

  # VTTablet and MySQL parameter overrides. Once set, overrides missing from
  # this map are reset to their defaults; leave it unset to leave parameters
  # unchanged.
  parameters = {
    vttablet = {
      "vreplication-parallel-insert-workers" = "4"
    }
    mysqld = {
      wait_timeout = "600"
    }
  }
}
