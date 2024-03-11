---
page_title: region_from_zone Function - terraform-provider-google
description: |-
  Returns the region within a provided zone.
---

# Function: region_from_zone

Returns a region name derived from a provided zone.

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
  value = provider::google::region_from_zone("us-central1-b")
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
  value = provider::google-beta::region_from_zone("us-central1-b")
}
```

## Signature

```text
region_from_zone(zone string) string
```

## Arguments

1. `zone` (String) A string of a resource's zone
