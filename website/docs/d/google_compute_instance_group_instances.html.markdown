---
layout: "google"
page_title: "Google: google_compute_instance_group_instances"
sidebar_current: "docs-google-datasource-compute-instance-group-instances"
description: |-
  Provides a list of Google Compute Instance Group Instances
---

# google\_compute\_instance\_group\_instances

Get a instances within GCE from instance group name and its zone.

```
data "google_compute_instance_group_instances" "all" {
	name = "instance-group-name"
	zone = "us-central1-a"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the instance group.

* `zone` - The zone of the instance group.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

The following attributes are exported:

* `instances` - A list of instance urls in the given instance group

* `names` - A list of instance names in the given instance group
