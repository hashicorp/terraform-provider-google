---
subcategory: "App Hub"
description: |-
  Get information about a discovered service.
---

# google_apphub_discovered_service

Get information about a discovered service from its uri.


## Example Usage


```hcl
data "google_apphub_discovered_service" "my-service" {
  location = "my-location"
  service_uri = "my-service-uri"
}
```

## Argument Reference

The following arguments are supported:

* `project` - The host project of the discovered service.
* `service_uri` - (Required) The uri of the service.
* `location` - (Required) The location of the discovered service.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `name` - Resource name of a Service. Format: "projects/{host-project-id}/locations/{location}/applications/{application-id}/services/{service-id}".

* `service_reference` - Reference to an underlying networking resource that can comprise a Service. Structure is [documented below](#nested_service_reference)

<a name="nested_service_reference"></a>A `service_reference` object would contain the following fields:

* `uri` - The underlying resource URI.

* `path` - Additional path under the resource URI.

* `service_properties` - Properties of an underlying compute resource that can comprise a Service. Structure is [documented below](#nested_service_properties)

<a name="nested_service_properties"></a>A `service_properties` object would contain the following fields:

* `gcp_project` - The service project identifier that the underlying cloud resource resides in.

* `location` - The location that the underlying resource resides in.

* `zone` - The location that the underlying resource resides in if it is zonal.
