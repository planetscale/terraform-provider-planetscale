# Workflow
* `speakeasy` is used to generate the Terraform provider from the PlanetScale OpenAPI schema and overlays in `schemas/`
* Changes to the provider are driven by Speakeasy configuration (`.speakeasy/`) and the OpenAPI spec, do not modify Go files directly outside of tests

# Bash commands
* make generate: regenerate the Terraform provider from the OpenAPI schema and overlays
* make download-openapi: download the latest OpenAPI schema from PlanetScale
* make update-speakeasy: update the `speakeasy` CLI

# Testing
* Acceptance tests use Terraform configs from `internal/provider/testdata/`
* Each test function has a matching directory: `TestAccFoo` → `testdata/TestAccFoo/`
* Tests use `config.TestNameDirectory()` to automatically load the matching testdata directory
* Prefer testify `require` over explicit checks and `t.Fatalf` in unit tests
