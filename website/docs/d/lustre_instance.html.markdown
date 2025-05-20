---
subcategory: "Lustre"
description: |-
  Fetches the details of a Lustre instance.
---

# google_lustre_instance

Use this data source to get information about a Lustre instance. For more information see the [API docs](https://cloud.google.com/filestore/docs/lustre/reference/rest/v1/projects.locations.instances).

## Example Usage

```hcl
data "google_lustre_instance" "instance" {
  name   = "my-instance"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The instance id of the Lustre instance.

* `zone` - (Optional) The ID of the zone in which the resource belongs. If it is not provided, the provider zone is used.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

## Attributes Reference

See [google_lustre_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/lustre_instance) resource for details of all the available attributes.
