<div align="center">
    <img width="200px" src="https://avatars.githubusercontent.com/u/35612527?s=200&v=4">
    <h1>Terraform Provider</h1>
    <p>Terraform Provider for the *PlanetScale API*.</p>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" /></a>
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
      version = "0.0.1"
    }
  }
}

provider "planetscale" {
  # Configuration options
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

* [planetscale_bouncer](docs/resources/bouncer.md)
* [planetscale_branch](docs/resources/branch.md)
* [planetscale_cidrs](docs/resources/cidrs.md)
* [planetscale_database](docs/resources/database.md)
* [planetscale_keyspace](docs/resources/keyspace.md)
* [planetscale_password](docs/resources/password.md)
* [planetscale_query_pattern_report](docs/resources/query_pattern_report.md)
* [planetscale_role](docs/resources/role.md)
* [planetscale_team](docs/resources/team.md)
* [planetscale_webhook](docs/resources/webhook.md)
* [planetscale_workflow](docs/resources/workflow.md)
### Data Sources

* [planetscale_bouncer](docs/data-sources/bouncer.md)
* [planetscale_branch](docs/data-sources/branch.md)
* [planetscale_cidrs](docs/data-sources/cidrs.md)
* [planetscale_database](docs/data-sources/database.md)
* [planetscale_databases](docs/data-sources/databases.md)
* [planetscale_keyspace](docs/data-sources/keyspace.md)
* [planetscale_organization](docs/data-sources/organization.md)
* [planetscale_organizations](docs/data-sources/organizations.md)
* [planetscale_password](docs/data-sources/password.md)
* [planetscale_query_pattern_report](docs/data-sources/query_pattern_report.md)
* [planetscale_role](docs/data-sources/role.md)
* [planetscale_team](docs/data-sources/team.md)
* [planetscale_teams](docs/data-sources/teams.md)
* [planetscale_webhook](docs/data-sources/webhook.md)
* [planetscale_workflow](docs/data-sources/workflow.md)
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
