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

resource "google_cloud_run_service" "default" {
  name     = "my-service"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Value is "us-central1"
output "location_from_id" {
  value = provider::google::location_from_id(google_cloud_run_service.default.id)
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

resource "google_cloud_run_service" "default" {
  # provider argument omitted - provisioning by google or google-beta doesn't impact this example
  name     = "my-service"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}


# Value is "us-central1"
output "location_from_id" {
  value = provider::google-beta::location_from_id(google_cloud_run_service.default.id)
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
