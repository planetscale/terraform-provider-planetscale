terraform {
  required_providers {
    planetscale = {
      source  = "planetscale/planetscale"
      version = "1.14.1"
    }
  }
}

provider "planetscale" {
  server_url = "..." # Optional - can use PLANETSCALE_SERVER_URL environment variable
}