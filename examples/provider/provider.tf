terraform {
  required_providers {
    planetscale = {
      source  = "planetscale/planetscale"
      version = "1.0.0-rc2"
    }
  }
}

provider "planetscale" {
  server_url = "..." # Optional - can use PLANETSCALE_SERVER_URL environment variable
}