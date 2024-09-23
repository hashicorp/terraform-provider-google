subcategory: "Compute Engine"
page_title: "Google: google_compute_region_instance_group_manager"
description: |-
Get a Compute Region Instance Group within GCE.
---

# google\_compute\_region\_instance\_group\_manager

Get a Compute Region Instance Group Manager within GCE.
For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups/distributing-instances-with-regional-instance-groups)
and [API](https://cloud.google.com/compute/docs/reference/rest/v1/regionInstanceGroupManagers)

## Example Usage

```hcl
data "google_compute_region_instance_group_manager" "rigm" {
  name = "my-igm"
  region = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `self_link` - (Optional) The self link of the instance group. Either `name` or `self_link` must be provided.

* `name` - (Optional) The name of the instance group. Either `name` or `self_link` must be provided.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `Region` - (Optional) The region where the managed instance group resides. If not provided, the provider region is used.

---

## Attributes Reference

See [google_compute_region_instance_group_manager](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_region_instance_group_manager) resource for details of all the available attributes.