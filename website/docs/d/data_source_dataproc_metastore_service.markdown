---
subcategory: "Dataproc"
layout: "google"
page_title: "Google: google_dataproc_metastore_service"
sidebar_current: "docs-google-datasource-dataproc-metastore-service"
description: |-
  Get a Dataproc Metastore Service from Google Cloud
---

# google\_dataproc\_metastore\_service

Get a Dataproc Metastore service from Google Cloud by its id and location.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage

```tf
data "google_dataproc_metastore_service" "foo" {
  service_id = "foo-bar"
  location   = "global"  
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) The ID of the metastore service.
* `location` - (Required) The location where the metastore service resides.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_dataproc_metastore_service](https://www.terraform.io/docs/providers/google/r/dataproc_metastore_service.html) resource for details of all the available attributes.
