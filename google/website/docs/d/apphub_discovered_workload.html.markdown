---
subcategory: "App Hub"
description: |-
  Get information about a discovered workload.
---

# google\_apphub\_discovered_workload

Get information about a discovered workload from its uri.


## Example Usage


```hcl
data "google_apphub_discovered_workload" "my-workload" {
  location = "us-central1"
  workload_uri = "my-workload-uri"
}
```

## Argument Reference

The following arguments are supported:

* `project` - The host project of the discovered workload.
* `workload_uri` - (Required)  The uri of the workload (instance group managed by the Instance Group Manager). Example: "//compute.googleapis.com/projects/1/regions/us-east1/instanceGroups/id1"
* `location` - (Required) The location of the discovered workload.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `name` - Resource name of a Workload. Format: "projects/{host-project-id}/locations/{location}/applications/{application-id}/workloads/{workload-id}".

* `workload_reference` - Reference to an underlying networking resource that can comprise a Workload. Structure is [documented below](#nested_workload_reference)

<a name="nested_workload_reference"></a>The `workload_reference` block supports:

* `uri` - The underlying resource URI.

* `workload_properties` - Properties of an underlying compute resource that can comprise a Workload. Structure is [documented below](#nested_workload_properties)

<a name="nested_workload_properties"></a>The `workload_properties` block supports:

* `gcp_project` - The service project identifier that the underlying cloud resource resides in.

* `location` - The location that the underlying resource resides in.

* `zone` - The location that the underlying resource resides in if it is zonal.
