# PlanetScale Terraform Provider

Work in progress Terraform provider for PlanetScale.

Not ready for general consumption.

## Docs?

👉 [Docs!](https://registry.terraform.io/providers/planetscale/planetscale/latest/docs)

## Usage

To use this provider, include the following blocks in your Terraform code:

```hcl
terraform {
  required_providers {
    planetscale = {
      source = "registry.terraform.io/planetscale/planetscale"
    }
  }
}
```

You must then configure the provider to either use a service token, or an access token. Either can be configured in the provider config block, or via env vars (`PLANETSCALE_SERVICE_TOKEN_NAME` and `PLANETSCALE_SERVICE_TOKEN`, or `PLANETSCALE_ACCESS_TOKEN`):

```hcl
provider "planetscale" {
  # use a service token
  service_token_name = "..." # ID of the service token to use, e.g "8fbddg0zlq0r"
  service_token      = "..." # Secret for the service token.

  # or use an access token
  access_token       = "..." # Secret for the access token.
}
```

## Known limitations

- Support for deployments, deploy queues, deploy requests and reverts is not implemented at this time. If you have a use case for it, please let us know in the repository issues.
- Service tokens don't immediately have read/write access on the resources they create. For now, access must be granted via the UI or via the CLI (`pscale service-token add-access`)

## Contributing

Note that this provider builds on top of the OpenAPI part of the PlanetScale API. The client used by this project is currently different from the Client provided in `planetscale-go`, and is largely code-generated from the OpenAPI schema found at https://api.planetscale.com/v1/openapi-spec. This schema is first [vendored](openapi/openapi-spec.json), then [lightly transformed](internal/cmd/extractref) using [custom config rules](openapi/extract-ref-cfg.json) into [its final form](openapi-spec.json), and then used to [generate](internal/cmd/client_codegen) the [client code](internal/client/planetscale).

Contributions should follow this workflow, or integrate with it.

## License

MPL v2.0