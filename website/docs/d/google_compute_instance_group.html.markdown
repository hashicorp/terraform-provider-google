---
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

```
data "google_compute_instance_group" "all" {
	name = "instance-group-name"
	zone = "us-central1-a"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the instance group.

* `zone` - (Required) The zone of the instance group.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

The following arguments are exported:

* `description` - Textual description of the instance group.

* `instances` - List of instances in the group.

* `named_port` - List of named ports in the group. 

* `network` - The URL of the network the instance group is in.

* `self_link` - The URI of the resource.

* `size` - The number of instances in the group.
