---
layout: "google"
page_title: "Google: google_compute_instance_group_manager"
sidebar_current: "docs-google-compute-instance-group-manager"
description: |-
  Manages an Instance Group within GCE.
---

# google\_compute\_instance\_group\_manager

The Google Compute Engine Instance Group Manager API creates and manages pools
of homogeneous Compute Engine virtual machine instances from a common instance
template. For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups/manager)
and [API](https://cloud.google.com/compute/docs/reference/latest/instanceGroupManagers)

~> **Note:** Use [google_compute_region_instance_group_manager](/docs/providers/google/r/compute_region_instance_group_manager.html) to create a regional (multi-zone) instance group manager.

## Example Usage

```hcl
resource "google_compute_instance_group_manager" "appserver" {
  name        = "appserver-igm"

  base_instance_name = "app"
  instance_template  = "${google_compute_instance_template.appserver.self_link}"
  update_strategy    = "NONE"
  zone               = "us-central1-a"

  target_pools = ["${google_compute_target_pool.appserver.self_link}"]
  target_size  = 2

  named_port {
    name = "customHTTP"
    port = 8888
  }
}
```

## Argument Reference

The following arguments are supported:

* `base_instance_name` - (Required) The base instance name to use for
    instances in this group. The value must be a valid
    [RFC1035](https://www.ietf.org/rfc/rfc1035.txt) name. Supported characters
    are lowercase letters, numbers, and hyphens (-). Instances are named by
    appending a hyphen and a random four-character string to the base instance
    name.

* `instance_template` - (Required) The full URL to an instance template from
    which all new instances will be created.

* `name` - (Required) The name of the instance group manager. Must be 1-63
    characters long and comply with
    [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Supported characters
    include lowercase letters, numbers, and hyphens.

* `zone` - (Required) The zone that instances in this group should be created
    in.

- - -

* `description` - (Optional) An optional textual description of the instance
    group manager.

* `named_port` - (Optional) The named port configuration. See the section below
    for details on configuration.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `update_strategy` - (Optional, Default `"RESTART"`) If the `instance_template`
    resource is modified, a value of `"NONE"` will prevent any of the managed
    instances from being restarted by Terraform. A value of `"RESTART"` will
    restart all of the instances at once. In the future, as the GCE API matures
    we will support `"ROLLING_UPDATE"` as well.

* `target_size` - (Optional) The target number of running instances for this managed
    instance group. This value should always be explicitly set unless this resource is attached to
     an autoscaler, in which case it should never be set. Defaults to `0`.

* `target_pools` - (Optional) The full URL of all target pools to which new
    instances in the group are added. Updating the target pools attribute does
    not affect existing instances.

---

* `auto_healing_policies` - (Optional, [Beta](/docs/providers/google/index.html#beta-features)) The autohealing policies for this managed instance
group. You can specify only one value. Structure is documented below.

The `named_port` block supports: (Include a `named_port` block for each named-port required).

* `name` - (Required) The name of the port.

* `port` - (Required) The port number.

The `auto_healing_policies` block supports:

* `health_check` - (Required) The health check that signals autohealing.

* `initial_delay_sec` - (Required) The number of seconds that the managed instance group waits before
 it applies autohealing policies to new instances or recently recreated instances. Between 0 and 3600.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `fingerprint` - The fingerprint of the instance group manager.

* `instance_group` - The full URL of the instance group created by the manager.

* `self_link` - The URL of the created resource.


## Import

Instance group managers can be imported using the `name`, e.g.

```
$ terraform import google_compute_instance_group_manager.appserver appserver-igm
```
