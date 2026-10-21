---
page_title: "Move Neki sidecar parameters to configuration profiles"
subcategory: ""
description: |-
  Replace planetscale_neki_sidecar with configuration profile parameters.
---

# Move Neki sidecar parameters to configuration profiles

The `planetscale_neki_sidecar` resource and data source have been removed. Sidecar parameters now belong to `planetscale_neki_configuration_profile`.

Complete these steps before upgrading the provider:

1. Move any sidecar parameter values you want to keep into the configuration profile's `parameters.sidecar` map. Keep its other parameter namespaces as well. If Terraform does not already manage the profile, [import the existing profile](../resources/neki_configuration_profile.md#import) first.

   For example, add this namespace to the profile's existing `parameters` map:

   ```hcl
   sidecar = {
     "pool-max-lifetime" = "1h"
   }
   ```

2. Remove the old `planetscale_neki_sidecar` resource and data source blocks from your configuration. Replace references to them with the configuration profile resource or data source.

3. Use `terraform state list` to find the old sidecar addresses, then remove each from state. For example:

   ```shell
   terraform state rm 'planetscale_neki_sidecar.example'
   ```

   If a sidecar data source is also in state, remove its address too:

   ```shell
   terraform state rm 'data.planetscale_neki_sidecar.example'
   ```

   Use your actual addresses, including any module path or instance key. Removing these state entries does not change the running connection pool or its parameters. Simply deleting the configuration blocks can still trigger a refresh against the retired Sidecars API and fail with `410 Gone`.

4. Upgrade the provider and run `terraform plan`. Check that the configuration profile retains the parameter values you intended before applying.
