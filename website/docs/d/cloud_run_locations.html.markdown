---
subcategory: "Cloud Run"
description: |-
  Get Cloud Run locations available for a project.
---

# google_cloud_run_locations

Get Cloud Run locations available for a project. 

To get more information about Cloud Run, see:

* [API documentation](https://cloud.google.com/run/docs/reference/rest/v1/projects.locations)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/run/docs/)
    
## Example Usage

```hcl
data "google_cloud_run_locations" "available" {
}
```

## Example Usage: Multi-regional Cloud Run deployment

```hcl
data "google_cloud_run_locations" "available" {
}

resource "google_cloud_run_service" "service_one" {
  name     = "service-one"
  location = data.google_cloud_run_locations.available.locations[0]

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

resource "google_cloud_run_service" "service_two" {
  name     = "service-two"
  location = data.google_cloud_run_locations.available.locations[1]

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello""
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project to list versions for. If it
    is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `locations` - The list of Cloud Run locations available for the given project.
