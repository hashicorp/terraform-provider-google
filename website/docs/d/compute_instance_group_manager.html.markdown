---
subcategory: "Compute Engine"
page_title: "Google: google_compute_instance_group_manager"
description: |-
  Get a Compute Instance Group within GCE.
---

# google\_compute\_instance\_group\_manager

Get a Compute Instance Group Manager within GCE.
For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups#managed_instance_groups)
and [API](https://cloud.google.com/compute/docs/reference/latest/instanceGroupManagers)

## Example Usage

```hcl
data "google_compute_instance_group_manager" "igm1" {
  name = "my-igm"
  zone = "us-central1-a"
}

data "google_compute_instance_group_manager" "igm2" {
  self_link = "https://www.googleapis.com/compute/v1/projects/myproject/zones/us-central1-a/instanceGroupManagers/my-igm"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the instance group. Either `name` or `self_link` must be provided.

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `self_link` - (Optional) The self link of the instance group. Either `name` or `self_link` must be provided.

* `zone` - (Optional) The zone of the instance group. If referencing the instance group by name and `zone` is not provided, the provider zone is used.

---

## Attributes Reference

See [google_compute_instance_group_manager](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance_group_manager#argument-reference) resource for details of all the available attributes.
