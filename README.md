<div align="center">
    <img width="200px" src="https://avatars.githubusercontent.com/u/35612527?s=200&v=4">
    <h1>Terraform Provider</h1>
    <p>Terraform Provider for the *PlanetScale API*.</p>
    <a href="https://opensource.org/license/apache-2-0"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge" /></a>
</div>

<!-- Start Summary [summary] -->
## Summary

Terraform Provider for PlanetScale: Manage your PlanetScale resources with Terraform
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
  * [Installation](#installation)
  * [Authentication](#authentication)
  * [Available Resources and Data Sources](#available-resources-and-data-sources)
  * [Testing the provider locally](#testing-the-provider-locally)
* [Development](#development)
  * [Contributions](#contributions)

<!-- End Table of Contents [toc] -->

<!-- Start Installation [installation] -->
## Installation

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
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
```
<!-- End Installation [installation] -->

<!-- Start Authentication [security] -->
## Authentication

This provider supports authentication configuration via environment variables and provider configuration.

The configuration precedence is:

- Provider configuration
- Environment variables

Available configuration:

| Provider Attribute | Description |
|---|---|
| `service_token` | PlanetScale Service Token. Configurable via environment variable `PLANETSCALE_SERVICE_TOKEN`. |
| `service_token_id` | PlanetScale Service Token ID. Configurable via environment variable `PLANETSCALE_SERVICE_TOKEN_ID`. |
<!-- End Authentication [security] -->

<!-- Start Available Resources and Data Sources [operations] -->
## Available Resources and Data Sources

### Resources

* [planetscale_postgres_branch](docs/resources/postgres_branch.md)
* [planetscale_postgres_branch_role](docs/resources/postgres_branch_role.md)
* [planetscale_vitess_branch](docs/resources/vitess_branch.md)
* [planetscale_vitess_branch_password](docs/resources/vitess_branch_password.md)
### Data Sources

* [planetscale_database_postgres](docs/data-sources/database_postgres.md)
* [planetscale_databases](docs/data-sources/databases.md)
* [planetscale_database_vitess](docs/data-sources/database_vitess.md)
* [planetscale_organization](docs/data-sources/organization.md)
* [planetscale_organizations](docs/data-sources/organizations.md)
* [planetscale_postgres_branch](docs/data-sources/postgres_branch.md)
* [planetscale_postgres_branch_role](docs/data-sources/postgres_branch_role.md)
* [planetscale_postgres_branch_roles](docs/data-sources/postgres_branch_roles.md)
* [planetscale_user](docs/data-sources/user.md)
* [planetscale_vitess_branch](docs/data-sources/vitess_branch.md)
* [planetscale_vitess_branch_password](docs/data-sources/vitess_branch_password.md)
* [planetscale_vitess_branch_passwords](docs/data-sources/vitess_branch_passwords.md)
<!-- End Available Resources and Data Sources [operations] -->

<!-- Start Testing the provider locally [usage] -->
## Testing the provider locally

#### Local Provider

Should you want to validate a change locally, the `--debug` flag allows you to execute the provider against a terraform instance locally.

This also allows for debuggers (e.g. delve) to be attached to the provider.

```sh
go run main.go --debug
# Copy the TF_REATTACH_PROVIDERS env var
# In a new terminal
cd examples/your-example
TF_REATTACH_PROVIDERS=... terraform init
TF_REATTACH_PROVIDERS=... terraform apply
```

#### Compiled Provider

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

1. Execute `go build` to construct a binary called `terraform-provider-planetscale`
2. Ensure that the `.terraformrc` file is configured with a `dev_overrides` section such that your local copy of terraform can see the provider binary

Terraform searches for the `.terraformrc` file in your home directory and applies any configuration settings you set.

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/planetscale/planetscale" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
<!-- End Testing the provider locally [usage] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->

# Development

## Contributions

While we value open-source contributions to this terraform provider, this library is generated programmatically. Any manual changes added to internal files will be overwritten on the next generation.
We look forward to hearing your feedback. Feel free to open a PR or an issue with a proof of concept and we'll do our best to include it in a future release. 

### SDK Created by [Speakeasy](https://www.speakeasy.com/?utm_source=planetscale&utm_campaign=terraform)
