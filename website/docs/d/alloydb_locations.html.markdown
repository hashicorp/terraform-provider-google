---
subcategory: "AlloyDB"
description: |-
  Fetches the details of available locations.
---

# google_alloydb_locations

Use this data source to get information about the available locations. For more details refer the [API docs](https://cloud.google.com/alloydb/docs/reference/rest/v1/projects.locations).

## Example Usage


```hcl
data "google_alloydb_locations" "qa" {
}
```

## Argument Reference

The following arguments are supported:

* `project` - (optional) The ID of the project.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `locations` - Contains a list of `location`, which contains the details about a particular location.

A `location` object would contain the following fields:-

* `name` - Resource name for the location, which may vary between implementations. For example: "projects/example-project/locations/us-east1".

* `location_id` - The canonical id for this location. For example: "us-east1"..

* `display_name` - The friendly name for this location, typically a nearby city name. For example, "Tokyo".

* `labels` - Cross-service attributes for the location. For example `{"cloud.googleapis.com/region": "us-east1"}`.

* `metadata` - Service-specific metadata. For example the available capacity at the given location.
