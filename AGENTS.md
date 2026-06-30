# Workflow
* `speakeasy` generates this provider from an OpenAPI schema and overlays; do not modify generated Go files directly outside of tests
* The upstream schema is `schemas/openapi.yaml`. Never hand-edit it — refresh it with `make download-openapi`
* All customization happens via overlay files in `schemas/` (`overlay-terraform-*.yaml`), composed in order by `.speakeasy/workflow.yaml`. `schemas/out.openapi.yaml` is the generated merge result; don't edit it
* After changing a schema or overlay, run `make generate`

# Bash commands
* make generate: regenerate the Terraform provider from the OpenAPI schema and overlays
* make download-openapi: download the latest OpenAPI schema from PlanetScale
* make update-speakeasy: update the `speakeasy` CLI

# Adding a new resource
* Add a new `overlay-terraform-<name>.yaml` (and a plural variant for the corresponding list data source, if any) and register it under `sources.overlays` in `.speakeasy/workflow.yaml`
* Tag the resource's schema with `x-speakeasy-entity: EntityName` (PascalCase) and tag each CRUD operation with `x-speakeasy-entity-operation: EntityName#create|read|update|delete`; list endpoints back a data source and use a plural entity name with `#read` (no `#list`)
* Look at an existing overlay (e.g. `schemas/overlay-terraform-vitess-branch.yaml`) for the conventions this repo follows

# Async operations
* If create/update doesn't return the resource in its final state, chain a second step onto the `read` operation: add an `entityOperation: EntityName#create#2` (or `#update#2`) entry with `options.polling.name` alongside the normal `#read` entry, then define that named poller via `x-speakeasy-polling` (delay/interval/limit + `successCriteria`/`failureCriteria` over `$statusCode`/`$response.body`)
* See `schemas/overlay-terraform-vitess-branch.yaml` (`WaitForReady`) and `schemas/overlay-terraform-postgres-branch-backup.yaml` (`WaitForComplete`) for working examples of both a state-based and a terminal-status-based poller

# Testing
* Acceptance tests use Terraform configs from `internal/provider/testdata/`
* Each test function has a matching directory: `TestAccFoo` → `testdata/TestAccFoo/`
* Tests use `config.TestNameDirectory()` to automatically load the matching testdata directory
* Add an acceptance test for every new resource/data source covering create, update-in-place, and import (see `internal/provider/vitessbranch_resource_test.go`); run with `make testacc` (creates real resources)
* Prefer testify `require` over explicit checks and `t.Fatalf` in unit tests; run with `make test`
