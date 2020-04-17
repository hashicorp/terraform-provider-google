---
subcategory: "Stackdriver Monitoring"
layout: "google"
page_title: "Google: google_monitoring_app_engine_service"
sidebar_current: "docs-google-datasource-monitoring-app-engine-service"
description: |-
  An Monitoring Service resource created automatically by GCP to monitor an
  App Engine service.
---

# google\_monitoring\_app\_engine\_service

A Monitoring Service is the root resource under which operational aspects of a
generic service are accessible. A service is some discrete, autonomous, and
network-accessible unit, designed to solve an individual concern

An App Engine monitoring service is automatically created by GCP to monitor
App Engine services.


To get more information about Service, see:

* [API documentation](https://cloud.google.com/monitoring/api/ref_v3/rest/v3/services)
* How-to Guides
    * [Service Monitoring](https://cloud.google.com/monitoring/service-monitoring)
    * [Monitoring API Documentation](https://cloud.google.com/monitoring/api/v3/)

## Example Usage - Monitoring App Engine Service


```hcl
# Monitors the default AppEngine service
data "google_monitoring_app_engine_service" "srv" {
  module_id = google_app_engine_standard_app_version.myapp.service
}

resource "google_app_engine_standard_app_version" "myapp" {
  version_id = "v1"
  service    = "myapp"
  runtime    = "nodejs10"

  entrypoint {
    shell = "node ./app.js"
  }

  deployment {
    zip {
      source_url = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${google_storage_bucket_object.object.name}"
    }
  }

  env_variables = {
    port = "8080"
  }

  delete_service_on_destroy = false
}

resource "google_storage_bucket" "bucket" {
  name = "appengine-static-content"
}

resource "google_storage_bucket_object" "object" {
  name   = "hello-world.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./test-fixtures/appengine/hello-world.zip"
}

```

## Argument Reference

The arguments of this data source act as filters for identifying a given App Engine-created service.

The given filters must match exactly one service whose data will be exported as attributes. The following arguments are supported:

One of the following fields must be specified:

* `module_id` - (Required) The ID of the App Engine module underlying this
  service. Corresponds to the moduleId resource label in the [gae_app](https://cloud.google.com/monitoring/api/resources#tag_gae_app) monitored resource, or the service/module name.

- - -

Other optional fields include:

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `name` -
  The full REST resource name for this channel. The syntax is:
  `projects/[PROJECT_ID]/services/[SERVICE_ID]`.

* `display_name` -
  Name used for UI elements listing this (Monitoring) Service.

* `telemetry` -
  Configuration for how to query telemetry on the Service. Structure is documented below.

The `telemetry` block includes:

* `resource_name` -
  (Optional)
  The full name of the resource that defines this service.
  Formatted as described in
  https://cloud.google.com/apis/design/resource_names.