terraform {
  required_providers {
    planetscale = {
      source = "planetscale/planetscale"
    }
  }
}

// Only required if using a local `api-bb`
provider "planetscale" {
  server_url = "http://api.pscaledev.com:3000/v1/"
}

resource "planetscale_postgres_branch" "my_branch" {
  organization  = "big-bang"
  database      = "u6gngydvi8k0"
  parent_branch = "main"
  name          = "hello4"
  cluster_size  = "PS_5_AWS_X86"
}