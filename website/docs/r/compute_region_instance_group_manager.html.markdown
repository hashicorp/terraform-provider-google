---
layout: "google"
page_title: "Google: google_compute_region_instance_group_manager"
sidebar_current: "docs-google-compute-region_instance-group-manager"
description: |-
  Manages an Regional Instance Group within GCE.
---

# google\_compute\_region\_instance\_group\_manager

The Google Compute Engine Regional Instance Group Manager API creates and manages pools
of homogeneous Compute Engine virtual machine instances from a common instance
template. For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups/distributing-instances-with-regional-instance-groups)
and [API](https://cloud.google.com/compute/docs/reference/latest/regionInstanceGroupManagers)

~> **Note:** Use [google_compute_instance_group_manager](/docs/providers/google/r/compute_instance_group_manager.html) to create a single-zone instance group manager.

## Example Usage

```hcl
resource "google_compute_health_check" "autohealing" {
  name                = "autohealing-health-check"
  check_interval_sec  = 5
  timeout_sec         = 5
  healthy_threshold   = 2
  unhealthy_threshold = 10                         # 50 seconds

  http_health_check {
    request_path = "/healthz"
    port         = "8080"
  }
}

resource "google_compute_region_instance_group_manager" "appserver" {
  name = "appserver-igm"

  base_instance_name         = "app"
  instance_template          = "${google_compute_instance_template.appserver.self_link}"
  region                     = "us-central1"
  distribution_policy_zones  = ["us-central1-a", "us-central1-f"]

  target_pools = ["${google_compute_target_pool.appserver.self_link}"]
  target_size  = 2

  named_port {
    name = "custom"
    port = 8888
  }

  auto_healing_policies {
    health_check      = "${google_compute_health_check.autohealing.self_link}"
    initial_delay_sec = 300
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

* `region` - (Required) The region where the managed instance group resides.

- - -

* `description` - (Optional) An optional textual description of the instance
    group manager.

* `named_port` - (Optional) The named port configuration. See the section below
    for details on configuration.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `update_strategy` - (Optional, Default `"NONE"`) If the `instance_template`
    resource is modified, a value of `"NONE"` will prevent any of the managed
    instances from being restarted by Terraform. A value of `"ROLLING_UPDATE"` 
    is supported as [Beta feature]. A value of `"ROLLING_UPDATE"` requires 
    `rolling_update_policy` block to be set

* `target_size` - (Optional) The target number of running instances for this managed
    instance group. This value should always be explicitly set unless this resource is attached to
     an autoscaler, in which case it should never be set. Defaults to `0`.

* `target_pools` - (Optional) The full URL of all target pools to which new
    instances in the group are added. Updating the target pools attribute does
    not affect existing instances.

* `wait_for_instances` - (Optional) Whether to wait for all instances to be created/updated before
    returning. Note that if this is set to true and the operation does not succeed, Terraform will
    continue trying until it times out.

---

* `auto_healing_policies` - (Optional, [Beta](/docs/providers/google/index.html#beta-features)) The autohealing policies for this managed instance
group. You can specify only one value. Structure is documented below. For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/creating-groups-of-managed-instances#monitoring_groups).

* `rolling_update_policy` - (Optional, [Beta](/docs/providers/google/index.html#beta-features)) The update policy for this managed instance group. Structure is documented below. For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/updating-managed-instance-groups) and [API](https://cloud.google.com/compute/docs/reference/rest/beta/regionInstanceGroupManagers/patch)

* `distribution_policy_zones` - (Optional, [Beta](/docs/providers/google/index.html#beta-features)) The distribution policy for this managed instance
group. You can specify one or more values. For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/distributing-instances-with-regional-instance-groups#selectingzones).
- - -

The **rolling_update_policy** block supports:

```hcl
rolling_update_policy{
  type = "PROACTIVE"
  minimal_action = "REPLACE"
  max_surge_percent = 20
  max_unavailable_fixed = 2
  min_ready_sec = 50
}
```

* `minimal_action` - (Required) - Minimal action to be taken on an instance. Valid values are `"RESTART"`, `"REPLACE"`

* `type` - (Required) - The type of update. Valid values are `"OPPORTUNISTIC"`, `"PROACTIVE"`

* `max_surge_fixed` - (Optional), The maximum number of instances that can be created above the specified targetSize during the update process. Conflicts with `max_surge_percent`. It has to be either 0 or at least equal to the number of zones.  If fixed values are used, at least one of `max_unavailable_fixed` or `max_surge_fixed` must be greater than 0.

* `max_surge_percent` - (Optional), The maximum number of instances(calculated as percentage) that can be created above the specified targetSize during the update process. Conflicts with `max_surge_fixed`. Percent value is only allowed for regional managed instance groups with size at least 10.

* `max_unavailable_fixed` - (Optional), The maximum number of instances that can be unavailable during the update process. Conflicts with `max_unavailable_percent`. It has to be either 0 or at least equal to the number of zones. If fixed values are used, at least one of `max_unavailable_fixed` or `max_surge_fixed` must be greater than 0.

* `max_unavailable_percent` - (Optional), The maximum number of instances(calculated as percentage) that can be unavailable during the update process. Conflicts with `max_unavailable_fixed`. Percent value is only allowed for regional managed instance groups with size at least 10.

* `min_ready_sec` - (Optional), Minimum number of seconds to wait for after a newly created instance becomes available. This value must be from range [0, 3600]
- - -

The **named_port** block supports: (Include a `named_port` block for each named-port required).

* `name` - (Required) The name of the port.

* `port` - (Required) The port number.
- - -

The **auto_healing_policies** block supports:

* `health_check` - (Required) The health check resource that signals autohealing.

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
$ terraform import google_compute_region_instance_group_manager.appserver appserver-igm
```
