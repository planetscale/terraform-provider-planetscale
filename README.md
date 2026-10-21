# PlanetScale Terraform Provider

<!-- Start Summary [summary] -->
## Summary

Terraform Provider for PlanetScale: Manage your PlanetScale resources with Terraform
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
* [PlanetScale Terraform Provider](#planetscale-terraform-provider)
  * [Installation](#installation)
  * [Authentication](#authentication)
  * [Conceptual model](#conceptual-model)
  * [Available Resources and Data Sources](#available-resources-and-data-sources)
  * [Testing the provider locally](#testing-the-provider-locally)
* [Contributing](#contributing)

<!-- End Table of Contents [toc] -->

<!-- Start Installation [installation] -->
## Installation

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    planetscale = {
      source  = "planetscale/planetscale"
      version = "1.15.0"
    }
  }
}

provider "planetscale" {
  server_url = "..." # Optional - can use PLANETSCALE_SERVER_URL environment variable
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

<!--
  Hand-maintained section. It mirrors https://planetscale.com/docs/terraform#conceptual-model
  and lives outside the Speakeasy "Start ... / End ..." markers, so `make generate` leaves it
  alone. Update it here and in the public docs together.
-->
## Conceptual model

The PlanetScale Terraform provider is organized around a clear distinction between **Vitess**, **Neki**, and **Postgres** resources and is focused on long-lived infrastructure objects.

* There are separate resources for each **database kind**—Vitess, Neki, or Postgres.
* The provider maps cleanly onto PlanetScale's public API while exposing Terraform-friendly fields and lifecycle behavior.

PlanetScale databases do not have a dedicated Terraform resource. Their lifecycle is managed implicitly through branch resources:

* Applying a [`planetscale_vitess_branch`](docs/resources/vitess_branch.md), [`planetscale_neki_branch`](docs/resources/neki_branch.md), or [`planetscale_postgres_branch`](docs/resources/postgres_branch.md) creates the parent database if it does not already exist.
* Destroying the last branch in a database also destroys the database. Set `deletion_protected = true` on the branch resource to guard against this.

For example, this creates the `my-db` Postgres database (if it does not already exist) along with its `main` branch:

```hcl
resource "planetscale_postgres_branch" "main" {
  organization = "my-org"
  database     = "my-db"
  name         = "main"
}
```

Creating a Neki branch also provisions a default [configuration profile](docs/resources/neki_configuration_profile.md), [shard](docs/resources/neki_shard.md), [router group](docs/resources/neki_router.md), and [admin](docs/resources/neki_admin.md). Import those objects if Terraform should manage them. Do not recreate the defaults with `planetscale_neki_configuration_profile` or `planetscale_neki_shard`. Creating a Vitess branch likewise provisions the branch's default [keyspace](docs/resources/vitess_keyspace.md); import it to manage its size or replicas.

Sidecar connection-pool parameters belong to the configuration profile. If you used `planetscale_neki_sidecar`, follow the [sidecar migration guide](docs/guides/migrate-neki-sidecars.md) before upgrading.

### Credential models

**Vitess**, **Neki**, and **Postgres** use different terminology for database credentials:

| Database Type | Resource | How permissions work |
| - | - | - |
| Vitess | [`planetscale_vitess_branch_password`](docs/resources/vitess_branch_password.md) | Set the `role` attribute to `reader`, `writer`, `readwriter`, or `admin` |
| Neki | [`planetscale_neki_role`](docs/resources/neki_role.md) | Use `inherited_roles` to inherit from built-in Postgres roles like `pg_read_all_data` or `pg_write_all_data`. |
| Postgres | [`planetscale_postgres_branch_role`](docs/resources/postgres_branch_role.md) | Use `inherited_roles` to inherit from built-in Postgres roles like `pg_read_all_data` or `pg_write_all_data`. |
| Postgres | [`planetscale_postgres_redacted_branch_role`](docs/resources/postgres_redacted_branch_role.md) | Same as `planetscale_postgres_branch_role`, but the password is not stored in Terraform state. Use this when you manage passwords in a secret manager. |

<!-- Start Available Resources and Data Sources [operations] -->
## Available Resources and Data Sources

### Managed Resources

* [planetscale_neki_admin](docs/resources/neki_admin.md)
* [planetscale_neki_backup_policy](docs/resources/neki_backup_policy.md)
* [planetscale_neki_branch](docs/resources/neki_branch.md)
* [planetscale_neki_configuration_profile](docs/resources/neki_configuration_profile.md)
* [planetscale_neki_role](docs/resources/neki_role.md)
* [planetscale_neki_router](docs/resources/neki_router.md)
* [planetscale_neki_shard](docs/resources/neki_shard.md)
* [planetscale_postgres_backup_policy](docs/resources/postgres_backup_policy.md)
* [planetscale_postgres_bouncer](docs/resources/postgres_bouncer.md)
* [planetscale_postgres_branch](docs/resources/postgres_branch.md)
* [planetscale_postgres_branch_backup](docs/resources/postgres_branch_backup.md)
* [planetscale_postgres_branch_role](docs/resources/postgres_branch_role.md)
* [planetscale_postgres_dedicated_read_replica](docs/resources/postgres_dedicated_read_replica.md)
* [planetscale_postgres_redacted_branch_role](docs/resources/postgres_redacted_branch_role.md)
* [planetscale_vitess_backup_policy](docs/resources/vitess_backup_policy.md)
* [planetscale_vitess_branch](docs/resources/vitess_branch.md)
* [planetscale_vitess_branch_backup](docs/resources/vitess_branch_backup.md)
* [planetscale_vitess_branch_password](docs/resources/vitess_branch_password.md)
* [planetscale_vitess_keyspace](docs/resources/vitess_keyspace.md)

### Data Sources

* [planetscale_database_postgres](docs/data-sources/database_postgres.md)
* [planetscale_database_vitess](docs/data-sources/database_vitess.md)
* [planetscale_databases](docs/data-sources/databases.md)
* [planetscale_neki_admin](docs/data-sources/neki_admin.md)
* [planetscale_neki_backup_policies](docs/data-sources/neki_backup_policies.md)
* [planetscale_neki_backup_policy](docs/data-sources/neki_backup_policy.md)
* [planetscale_neki_branch](docs/data-sources/neki_branch.md)
* [planetscale_neki_configuration_profile](docs/data-sources/neki_configuration_profile.md)
* [planetscale_neki_configuration_profiles](docs/data-sources/neki_configuration_profiles.md)
* [planetscale_neki_role](docs/data-sources/neki_role.md)
* [planetscale_neki_roles](docs/data-sources/neki_roles.md)
* [planetscale_neki_router](docs/data-sources/neki_router.md)
* [planetscale_neki_routers](docs/data-sources/neki_routers.md)
* [planetscale_neki_shard](docs/data-sources/neki_shard.md)
* [planetscale_neki_shards](docs/data-sources/neki_shards.md)
* [planetscale_organization](docs/data-sources/organization.md)
* [planetscale_organizations](docs/data-sources/organizations.md)
* [planetscale_postgres_backup_policies](docs/data-sources/postgres_backup_policies.md)
* [planetscale_postgres_backup_policy](docs/data-sources/postgres_backup_policy.md)
* [planetscale_postgres_bouncer](docs/data-sources/postgres_bouncer.md)
* [planetscale_postgres_bouncers](docs/data-sources/postgres_bouncers.md)
* [planetscale_postgres_branch](docs/data-sources/postgres_branch.md)
* [planetscale_postgres_branch_backup](docs/data-sources/postgres_branch_backup.md)
* [planetscale_postgres_branch_backups](docs/data-sources/postgres_branch_backups.md)
* [planetscale_postgres_branch_role](docs/data-sources/postgres_branch_role.md)
* [planetscale_postgres_branch_roles](docs/data-sources/postgres_branch_roles.md)
* [planetscale_postgres_dedicated_read_replica](docs/data-sources/postgres_dedicated_read_replica.md)
* [planetscale_postgres_dedicated_read_replicas](docs/data-sources/postgres_dedicated_read_replicas.md)
* [planetscale_postgres_redacted_branch_role](docs/data-sources/postgres_redacted_branch_role.md)
* [planetscale_vitess_backup_policies](docs/data-sources/vitess_backup_policies.md)
* [planetscale_vitess_backup_policy](docs/data-sources/vitess_backup_policy.md)
* [planetscale_vitess_branch](docs/data-sources/vitess_branch.md)
* [planetscale_vitess_branch_backup](docs/data-sources/vitess_branch_backup.md)
* [planetscale_vitess_branch_backups](docs/data-sources/vitess_branch_backups.md)
* [planetscale_vitess_branch_password](docs/data-sources/vitess_branch_password.md)
* [planetscale_vitess_branch_passwords](docs/data-sources/vitess_branch_passwords.md)
* [planetscale_vitess_keyspace](docs/data-sources/vitess_keyspace.md)
* [planetscale_vitess_keyspaces](docs/data-sources/vitess_keyspaces.md)
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

# Contributing

See [Contributing to the PlanetScale Terraform Provider](./CONTRIBUTING.md)
