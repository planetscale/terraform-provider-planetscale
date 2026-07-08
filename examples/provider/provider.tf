terraform {
  required_providers {
    planetscale = {
      source  = "planetscale/planetscale"
      version = "1.3.1"
    }
  }
}

provider "planetscale" {
  server_url = "..." # Optional - can use PLANETSCALE_SERVER_URL environment variable
}