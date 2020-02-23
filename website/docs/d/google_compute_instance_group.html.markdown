---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_instance_group"
sidebar_current: "docs-google-datasource-compute-instance-group"
description: |-
  Get a Compute Instance Group within GCE.
---

# google\_compute\_instance\_group

Get a Compute Instance Group within GCE.
For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups/#unmanaged_instance_groups)
and [API](https://cloud.google.com/compute/docs/reference/latest/instanceGroups)

```hcl
data "google_compute_instance_group" "all" {
	name = "instance-group-name"
	zone = "us-central1-a"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the instance group. Either `name` or `self_link` must be provided.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `self_link` - (Optional) The self link of the instance group. Either `name` or `self_link` must be provided.

* `zone` - (Optional) The zone of the instance group. If referencing the instance group by name
    and `zone` is not provided, the provider zone is used.

## Attributes Reference

The following arguments are exported:

* `description` - Textual description of the instance group.

* `instances` - List of instances in the group.

* `named_port` - List of named ports in the group.

* `network` - The URL of the network the instance group is in.

* `self_link` - The URI of the resource.

* `size` - The number of instances in the group.
