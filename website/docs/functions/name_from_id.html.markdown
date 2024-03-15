---
page_title: name_from_id Function - terraform-provider-google
description: |-
  Returns the project within a provided resource id, self link, or OP style resource name.
---

# Function: name_from_id

Returns the short-form name within a provided resource's id, resource URI, self link, or full resource name.

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

resource "google_pubsub_topic" "default" {
  name = "my-topic"
}

# Value is "my-topic"
output "function_output" {
  value = provider::google::name_from_id(google_pubsub_topic.default.id)
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

resource "google_pubsub_topic" "default" {
  # provider argument omitted - provisioning by google or google-beta doesn't impact this example
  name = "my-topic"
}

# Value is "my-topic"
output "function_output" {
  value = provider::google-beta::name_from_id(google_pubsub_topic.default.id)
}
```

## Signature

```text
name_from_id(id string) string
```

## Arguments

1. `id` (String) A string of a resource's id, resource URI, self link, or full resource name. For example, these are all valid values:

* `"projects/my-project/zones/us-central1-c/instances/my-instance"`
* `"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance"`
* `"//gkehub.googleapis.com/projects/my-project/locations/us-central1/memberships/my-membership"`
