---
page_title: zone_from_id Function - terraform-provider-google
description: |-
  Returns the project within a provided resource id, self link, or OP style resource name.
---

# Function: zone_from_id

Returns the zone within a provided resource's id, resource URI, self link, or full resource name.

For more information about using provider-defined functions with Terraform [see the official documentation](https://developer.hashicorp.com/terraform/plugin/framework/functions/concepts).

## Example Usage

### Use with the `google` provider

```terraform
terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

resource "google_compute_disk" "default" {
  name  = "my-disk"
  zone  = "us-central1-c"
}

# Value is "us-central1-c"
output "zone_from_id" {
  value = provider::google::zone_from_id(google_compute_disk.default.id)
}

# Value is "us-central1-c"
output "zone_from_self_link" {
  value = provider::google::zone_from_id(google_compute_disk.default.self_link)
}
```

### Use with the `google-beta` provider

```terraform
terraform {
  required_providers {
    google-beta = {
      source = "hashicorp/google-beta"
    }
  }
}

resource "google_compute_disk" "default" {
  # provider argument omitted - provisioning by google or google-beta doesn't impact this example
  name  = "my-disk"
  zone  = "us-central1-c"
}

# Value is "us-central1-c"
output "zone_from_id" {
  value = provider::google-beta::zone_from_id(google_compute_disk.default.id)
}

# Value is "us-central1-c"
output "zone_from_self_link" {
  value = provider::google-beta::zone_from_id(google_compute_disk.default.self_link)
}
```

## Signature

```text
zone_from_id(id string) string
```

## Arguments

1. `id` (String) A string of a resource's id, resource URI, self link, or full resource name. For example, these are all valid values:

* `"projects/my-project/zones/us-central1-c/instances/my-instance"`
* `"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance"`
* `"//gkehub.googleapis.com/projects/my-project/locations/us-central1/memberships/my-membership"`