---
subcategory: "Cloud Platform"
description: |-
  Get all organizations.
---


# google_organizations

Gets a list of all organizations.
See [the official documentation](https://cloud.google.com/resource-manager/docs/creating-managing-organization)
and [API](https://cloud.google.com/resource-manager/reference/rest/v1/organizations/search).

## Example Usage

```hcl
data "google_organizations" "example" {
  filter = "domain:example.com"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) An optional query string used to filter the Organizations to return in the response. Filter rules are case-insensitive. Further information can be found in the [REST API](https://cloud.google.com/resource-manager/reference/rest/v1/organizations/search#request-body).


## Attributes Reference

The following attributes are exported:

* `organizations` - A list of all retrieved organizations. Structure is [defined below](#nested_organizations).

<a name="nested_organizations"></a>The `organizations` block supports:

* `directory_customer_id` - The Google for Work customer ID of the Organization.

* `display_name` - A human-readable string that refers to the Organization in the Google Cloud console. The string will be set to the primary domain (for example, `"google.com"`) of the G Suite customer that owns the organization.

* `lifecycle_state` - The Organization's current lifecycle state.

* `name` - The resource name of the Organization in the form `organizations/{organization_id}`.

* `org_id` - The Organization ID.
