# Workflow
* `speakeasy` is used to generate the Terraform provider from the PlanetScale OpenAPI schema and overlays in `schemas/`
* The generated Go code and docs should not be modified directly

# Bash commands
* make generate: regenerate the Terraform provider from the OpenAPI schema and overlays
* make download-openapi: download the latest OpenAPI schema from PlanetScale
* make update-speakeasy: update the `speakeasy` CLI