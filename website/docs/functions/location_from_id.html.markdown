---
page_title: location_from_id Function - terraform-provider-google
description: |-
  Returns the location within a provided resource id, self link, or OP style resource name.
---

# Function: location_from_id

Returns the location within a provided resource's id, resource URI, self link, or full resource name.

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

# Value is "us-central1"
output "function_output" {
  value = provider::google::location_from_id("https://run.googleapis.com/v2/projects/my-project/locations/us-central1/services/my-service")
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

# Value is "us-central1"
output "function_output" {
  value = provider::google-beta::location_from_id("https://run.googleapis.com/v2/projects/my-project/locations/us-central1/services/my-service")
}
```

## Signature

```text
location_from_id(id string) string
```

## Arguments

1. `id` (String) A string of a resource's id, resource URI, self link, or full resource name. For example, these are all valid values:

* `"projects/my-project/locations/us-central1/services/my-service"`
* `"https://run.googleapis.com/v2/projects/my-project/locations/us-central1/services/my-service"`
* `"//run.googleapis.com/v2/projects/my-project/locations/us-central1/services/my-service"`
