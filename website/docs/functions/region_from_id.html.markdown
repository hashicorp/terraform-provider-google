---
page_title: region_from_id Function - terraform-provider-google
description: |-
  Returns the region within a provided resource id, self link, or OP style resource name.
---

# Function: region_from_id

Returns the region within a provided resource's id, resource URI, self link, or full resource name.

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

resource "google_compute_node_template" "default" {
  name = "my-node-template"
  region = "us-central1"
}

# Value is "us-central1"
output "region_from_id" {
  value = provider::google::region_from_id(google_compute_node_template.default.id)
}

# Value is "us-central1"
output "region_from_self_link" {
  value = provider::google::region_from_id(google_compute_node_template.default.self_link)
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

resource "google_compute_node_template" "default" {
  # provider argument omitted - provisioning by google or google-beta doesn't impact this example
  name = "my-node-template"
  region = "us-central1"
}

# Value is "us-central1"
output "region_from_id" {
  value = provider::google-beta::region_from_id(google_compute_node_template.default.id)
}

# Value is "us-central1"
output "region_from_self_link" {
  value = provider::google-beta::region_from_id(google_compute_node_template.default.self_link)
}
```

## Signature

```text
region_from_id(id string) string
```

## Arguments

1. `id` (String) A string of a resource's id, resource URI, self link, or full resource name. For example, these are all valid values:

* `"projects/my-project/regions/us-central1/subnetworks/my-subnetwork"`
* `"https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/subnetworks/my-subnetwork"`
* `"//compute.googleapis.com/projects/my-project/regions/us-central1/subnetworks/my-subnetwork"`