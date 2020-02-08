---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_region_instance_group"
sidebar_current: "docs-google-datasource-compute-region-instance-group"
description: |-
  Get the instances inside a Compute Region Instance Group within GCE.
---

# google\_compute\_region\_instance\_group

Get a Compute Region Instance Group within GCE.
For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups/distributing-instances-with-regional-instance-groups) and [API](https://cloud.google.com/compute/docs/reference/latest/regionInstanceGroups).

```hcl
data "google_compute_region_instance_group" "group" {
  name = "instance-group-name"
}
```

The most common use of this datasource will be to fetch information about the instances inside regional managed instance groups, for instance:

```hcl
resource "google_compute_region_instance_group_manager" "foo" {
  name               = "some_name"
    ...
  base_instance_name = "foo"
    ...
  instance_template  = google_compute_instance_template.foo.self_link
  target_pools = [google_compute_target_pool.foo.self_link]
    ...
}

data "google_compute_region_instance_group" "data_source" {
  self_link = google_compute_region_instance_group_manager.foo.instance_group
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the instance group.  One of `name` or `self_link` must be provided.

* `self_link` - (Optional) The link to the instance group.  One of `name` or `self_link` must be provided.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs.
    If `self_link` is provided, this value is ignored.  If neither `self_link`
    nor `project` are provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs.  If `self_link`
    is provided, this value is ignored.  If neither `self_link` nor `region` are
    provided, the provider region is used.

## Attributes Reference

The following arguments are exported:

* `size` - The number of instances in the group.

* `instances` - List of instances in the group, as a list of resources, each containing:
    * `instance` - URL to the instance.
    * `named_ports` - List of named ports in the group, as a list of resources, each containing:
        * `port` - Integer port number
        * `name` - String port name
    * `status` - String description of current state of the instance.
