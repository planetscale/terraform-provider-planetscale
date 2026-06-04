terraform {
  required_providers {
    planetscale = {
      source  = "planetscale/planetscale"
      version = "1.1.0"
    }
  }
}

provider "planetscale" {
  server_url = "..." # Optional - can use PLANETSCALE_SERVER_URL environment variable
}