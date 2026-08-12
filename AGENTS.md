# Workflow
* `speakeasy` generates this provider from an OpenAPI schema and overlays; do not modify generated Go files directly outside of tests
* The upstream schema is `schemas/openapi.yaml`. Never hand-edit it — refresh it with `make download-openapi`
* All customization happens via overlay files in `schemas/` (`overlay-terraform-*.yaml`), composed in order by `.speakeasy/workflow.yaml`. `schemas/out.openapi.yaml` is the generated merge result; don't edit it
* OpenAPI update PRs only refresh the schema. Check them out, adjust overlays, run `make generate`, and review the complete generated diff before merging
* Use overlays to keep unrelated or UI-only API fields out of Terraform

# Bash commands
* make generate: regenerate the Terraform provider from the OpenAPI schema and overlays
* make download-openapi: download the latest OpenAPI schema from PlanetScale
* make update-speakeasy: update the `speakeasy` CLI

# Terraform design
* Put settings on the resource that owns them; add a separate resource only for an object with its own long-lived lifecycle
* Use public IDs for stable resource identity
* API operations used in generated create or update chains must safely handle unchanged desired state

# Adding a new resource
* Add a new `overlay-terraform-<name>.yaml` (and a plural variant for the corresponding list data source, if any) and register it under `sources.overlays` in `.speakeasy/workflow.yaml`
* Tag the resource's schema with `x-speakeasy-entity: EntityName` (PascalCase) and tag each CRUD operation with `x-speakeasy-entity-operation: EntityName#create|read|update|delete`; list endpoints back a data source and use a plural entity name with `#read` (no `#list`)
* Look at an existing overlay (e.g. `schemas/overlay-terraform-vitess-branch.yaml`) for the conventions this repo follows

# Async operations
* Poll after the write that starts the async work
* If the resource read reports completion, chain it after create or update and configure `x-speakeasy-polling`
* If the write returns a separate operation, poll that operation by ID and define explicit success and failure criteria
* See `schemas/overlay-terraform-vitess-branch.yaml` and `schemas/overlay-terraform-postgres-branch-backup.yaml` for working examples

# Testing
* Acceptance tests use Terraform configs from `internal/provider/testdata/`
* Each test function has a matching directory: `TestAccFoo` → `testdata/TestAccFoo/`
* Tests use `config.TestNameDirectory()` to automatically load the matching testdata directory
* Add an acceptance test for every new resource/data source covering create, update-in-place, and import (see `internal/provider/vitessbranch_resource_test.go`); run with `make testacc` (creates real resources)
* Prefer testify `require` over explicit checks and `t.Fatalf` in unit tests; run with `make test`
