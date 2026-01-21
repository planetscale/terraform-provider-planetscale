# Contributing to the PlanetScale Terraform Provider

This provider is generated from the OpenAPI 3.0 spec which can be found at https://planetscale.com/docs/openapi.yaml.
Changes to that specification to support the Terraform provider should be done via [OpenAPI overlays](./schemas/).

## Workflow

For all contributors, we recommend the standard [GitHub flow](https://guides.github.com/introduction/flow/)
based on [forking and pull requests](https://guides.github.com/activities/forking/).

For significant changes, please [create an issue](https://github.com/planetscale/terraform-provider-planetscale/issues)
to let everyone know what you're planning to work on, and to track progress and design decisions.

## Testing

Acceptance tests create real resources in a PlanetScale organization.
These tests can be modified to run against your own organization rather than `planetscale-terraform-testing`,
but it is not required.
A PlanetScale maintainer can run tests for your changes after they've been reviewed before merging.
