terraform {
  required_providers {
    planetscale = {
      source  = "planetscale/planetscale"
      version = "1.0.0"
    }
  }
}

provider "planetscale" {
  server_url = "..." # Optional
}