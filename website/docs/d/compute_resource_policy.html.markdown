---
layout: "google"
subcategory: "Compute Engine"
page_title: "Google: google_compute_resource_policy"
sidebar_current: "docs-google-datasource-compute-resource-policy"
description: |-
  Provide access to a Resource Policy's attributes
---

# google\_compute\_resource\_policy

Provide access to a Resource Policy's attributes. For more information see [the official documentation](https://cloud.google.com/compute/docs/disks/scheduled-snapshots) or the [API](https://cloud.google.com/compute/docs/reference/rest/beta/resourcePolicies).

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

```hcl
provider "google-beta" {
  region = "us-central1"
  zone   = "us-central1-a"
}

data "google_compute_resource_policy" "daily" {
  provider = google-beta
  name     = "daily"
  region   = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` (Required) - The name of the Resource Policy.
* `project` (Optional) - Project from which to list the Resource Policy. Defaults to project declared in the provider.
* `region` (Required) - Region where the Resource Policy resides.

## Attributes Reference

The following attributes are exported:

* `description` - Description of this Resource Policy.
* `self_link` - The URI of the resource.
