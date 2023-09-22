---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "planetscale_oauth_applications Data Source - terraform-provider-planetscale"
subcategory: ""
description: |-
  A list of PlanetScale OAuth applications. (requires feature flag)
---

# planetscale_oauth_applications (Data Source)

A list of PlanetScale OAuth applications. (requires feature flag)

## Example Usage

```terraform
# requires a feature flag, contact support to enable it

data "planetscale_oauth_applications" "example" {
  organization = data.planetscale_organization.example.name
}

output "oauth_apps" {
  value = data.planetscale_oauth_applications.example
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `organization` (String)

### Read-Only

- `applications` (Attributes List) (see [below for nested schema](#nestedatt--applications))

<a id="nestedatt--applications"></a>
### Nested Schema for `applications`

Read-Only:

- `avatar` (String) The image source for the OAuth application's avatar.
- `client_id` (String) The OAuth application's unique client id.
- `created_at` (String) When the OAuth application was created.
- `domain` (String) The domain of the OAuth application. Used for verification of a valid redirect uri.
- `id` (String) The ID of the OAuth application.
- `name` (String) The name of the OAuth application.
- `redirect_uri` (String) The redirect URI of the OAuth application.
- `scopes` (List of String) The scopes that the OAuth application requires on a user's accout.
- `tokens` (Number) The number of tokens issued by the OAuth application.
- `updated_at` (String) When the OAuth application was last updated.