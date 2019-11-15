---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_organization"
sidebar_current: "docs-google-datasource-organization"
description: |-
  Get information about a Google Cloud Organization.
---

# google\_organization

Use this data source to get information about a Google Cloud Organization.

```hcl
data "google_organization" "org" {
  domain = "example.com"
}

resource "google_folder" "sales" {
  display_name = "Sales"
  parent       = data.google_organization.org.name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available Organizations.
The given filters must match exactly one Organizations whose data will be exported as attributes.
The following arguments are supported:

* `organization` (Optional) - The name of the Organization in the form `{organization_id}` or `organizations/{organization_id}`.
* `domain` (Optional) - The domain name of the Organization.

~> **NOTE:** One of `organization` or `domain` must be specified.

## Attributes Reference

The following additional attributes are exported:

* `id` - The Organization ID.
* `name` - The resource name of the Organization in the form `organizations/{organization_id}`.
* `directory_customer_id` - The Google for Work customer ID of the Organization.
* `create_time` - Timestamp when the Organization was created. A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example: "2014-10-02T15:01:23.045123456Z".
* `lifecycle_state` - The Organization's current lifecycle state.
