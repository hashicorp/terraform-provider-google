---
subcategory: "Dataproc"
page_title: "Google: google_dataproc_metastore_service"
description: |-
  Get a Dataproc Metastore Service from Google Cloud
---

# google\_dataproc\_metastore\_service

Get a Dataproc Metastore service from Google Cloud by its id and location.

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

See [google_dataproc_metastore_service](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dataproc_metastore_service) resource for details of all the available attributes.
