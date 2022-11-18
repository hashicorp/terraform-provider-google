---
subcategory: "Compute Engine"
page_title: "Google: google_compute_region_instance_group_manager"
description: |-
  Manages an Regional Instance Group within GCE.
---

# google\_compute\_region\_instance\_group\_manager

The Google Compute Engine Regional Instance Group Manager API creates and manages pools
of homogeneous Compute Engine virtual machine instances from a common instance
template.

To get more information about regionInstanceGroupManagers, see:

* [API documentation](https://cloud.google.com/compute/docs/reference/latest/regionInstanceGroupManagers)
* How-to Guides
    * [Regional Instance Groups Guide](https://cloud.google.com/compute/docs/instance-groups/distributing-instances-with-regional-instance-groups)

~> **Note:** Use [google_compute_instance_group_manager](/docs/providers/google/r/compute_instance_group_manager.html) to create a zonal instance group manager.

## Example Usage with top level instance template (`google` provider)

```hcl
resource "google_compute_health_check" "autohealing" {
  name                = "autohealing-health-check"
  check_interval_sec  = 5
  timeout_sec         = 5
  healthy_threshold   = 2
  unhealthy_threshold = 10 # 50 seconds

  http_health_check {
    request_path = "/healthz"
    port         = "8080"
  }
}

resource "google_compute_region_instance_group_manager" "appserver" {
  name = "appserver-igm"

  base_instance_name         = "app"
  region                     = "us-central1"
  distribution_policy_zones  = ["us-central1-a", "us-central1-f"]

  version {
    instance_template = google_compute_instance_template.appserver.id
  }

  all_instances_config {
    metadata = {
      metadata_key = "metadata_value"
    }
    labels = {
      label_key = "label_value"
    }
  }

  target_pools = [google_compute_target_pool.appserver.id]
  target_size  = 2

  named_port {
    name = "custom"
    port = 8888
  }

  auto_healing_policies {
    health_check      = google_compute_health_check.autohealing.id
    initial_delay_sec = 300
  }
}
```

## Example Usage with multiple versions
```hcl
resource "google_compute_region_instance_group_manager" "appserver" {
  name = "appserver-igm"

  base_instance_name = "app"
  region             = "us-central1"

  target_size = 5

  version {
    instance_template = google_compute_instance_template.appserver.id
  }

  version {
    instance_template = google_compute_instance_template.appserver-canary.id
    target_size {
      fixed = 1
    }
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

* `version` - (Required) Application versions managed by this instance group. Each
    version deals with a specific instance template, allowing canary release scenarios.
    Structure is [documented below](#nested_version).

* `name` - (Required) The name of the instance group manager. Must be 1-63
    characters long and comply with
    [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Supported characters
    include lowercase letters, numbers, and hyphens.

* `region` - (Optional) The region where the managed instance group resides. If not provided, the provider region is used.

- - -

* `description` - (Optional) An optional textual description of the instance
    group manager.

* `named_port` - (Optional) The named port configuration. See the [section below](#nested_named_port)
    for details on configuration.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `target_size` - (Optional) The target number of running instances for this managed
    instance group. This value should always be explicitly set unless this resource is attached to
     an autoscaler, in which case it should never be set. Defaults to `0`.

* `list_managed_instances_results` - (Optional) Pagination behavior of the `listManagedInstances` API
    method for this managed instance group. Valid values are: `PAGELESS`, `PAGINATED`.
    If `PAGELESS` (default), Pagination is disabled for the group's `listManagedInstances` API method.
    `maxResults` and `pageToken` query parameters are ignored and all instances are returned in a single
    response. If `PAGINATED`, pagination is enabled, `maxResults` and `pageToken` query parameters are
    respected.

* `target_pools` - (Optional) The full URL of all target pools to which new
    instances in the group are added. Updating the target pools attribute does
    not affect existing instances.

* `wait_for_instances` - (Optional) Whether to wait for all instances to be created/updated before
    returning. Note that if this is set to true and the operation does not succeed, Terraform will
    continue trying until it times out.

* `wait_for_instances_status` - (Optional) When used with `wait_for_instances` it specifies the status to wait for.
    When `STABLE` is specified this resource will wait until the instances are stable before returning. When `UPDATED` is
    set, it will wait for the version target to be reached and any per instance configs to be effective as well as all
    instances to be stable before returning. The possible values are `STABLE` and `UPDATED`

---

* `auto_healing_policies` - (Optional) The autohealing policies for this managed instance
group. You can specify only one value. Structure is documented below. For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/creating-groups-of-managed-instances#monitoring_groups).

* `all_instances_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
  Properties to set on all instances in the group. After setting
  allInstancesConfig on the group, you must update the group's instances to
  apply the configuration.

* `update_policy` - (Optional) The update policy for this managed instance group. Structure is [documented below](#nested_update_policy). For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/updating-managed-instance-groups) and [API](https://cloud.google.com/compute/docs/reference/rest/beta/regionInstanceGroupManagers/patch)

* `distribution_policy_zones` - (Optional) The distribution policy for this managed instance
group. You can specify one or more values. For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/distributing-instances-with-regional-instance-groups#selectingzones).

* `distribution_policy_target_shape` - (Optional) The shape to which the group converges either proactively or on resize events (depending on the value set in update_policy.0.instance_redistribution_type). For more information see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/regional-mig-distribution-shape).

* `stateful_disk` - (Optional) Disks created on the instances that will be preserved on instance delete, update, etc. Structure is [documented below](#nested_stateful_disk). For more information see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/configuring-stateful-disks-in-migs). Proactive cross zone instance redistribution must be disabled before you can update stateful disks on existing instance group managers. This can be controlled via the `update_policy`.

- - -

<a name="nested_update_policy"></a>The `update_policy` block supports:

```hcl
update_policy {
  type                           = "PROACTIVE"
  instance_redistribution_type   = "PROACTIVE"
  minimal_action                 = "REPLACE"
  most_disruptive_allowed_action = "REPLACE"
  max_surge_percent              = 20
  max_unavailable_fixed          = 2
  min_ready_sec                  = 50
  replacement_method             = "RECREATE"
}
```

* `minimal_action` - (Required) - Minimal action to be taken on an instance. You can specify either `REFRESH` to update without stopping instances, `RESTART` to restart existing instances or `REPLACE` to delete and create new instances from the target template. If you specify a `REFRESH`, the Updater will attempt to perform that action only. However, if the Updater determines that the minimal action you specify is not enough to perform the update, it might perform a more disruptive action.

* `most_disruptive_allowed_action` - (Optional) - Most disruptive action that is allowed to be taken on an instance. You can specify either NONE to forbid any actions, REFRESH to allow actions that do not need instance restart, RESTART to allow actions that can be applied without instance replacing or REPLACE to allow all possible actions. If the Updater determines that the minimal update action needed is more disruptive than most disruptive allowed action you specify it will not perform the update at all.

* `type` - (Required) - The type of update process. You can specify either `PROACTIVE` so that the instance group manager proactively executes actions in order to bring instances to their target versions or `OPPORTUNISTIC` so that no action is proactively executed but the update will be performed as part of other actions (for example, resizes or recreateInstances calls).

* `instance_redistribution_type` - (Optional) - The instance redistribution policy for regional managed instance groups. Valid values are: `"PROACTIVE"`, `"NONE"`. If `PROACTIVE` (default), the group attempts to maintain an even distribution of VM instances across zones in the region. If `NONE`, proactive redistribution is disabled.

* `max_surge_fixed` - (Optional), The maximum number of instances that can be created above the specified targetSize during the update process. Conflicts with `max_surge_percent`. It has to be either 0 or at least equal to the number of zones.  If fixed values are used, at least one of `max_unavailable_fixed` or `max_surge_fixed` must be greater than 0.

* `max_surge_percent` - (Optional), The maximum number of instances(calculated as percentage) that can be created above the specified targetSize during the update process. Conflicts with `max_surge_fixed`. Percent value is only allowed for regional managed instance groups with size at least 10.

* `max_unavailable_fixed` - (Optional), The maximum number of instances that can be unavailable during the update process. Conflicts with `max_unavailable_percent`. It has to be either 0 or at least equal to the number of zones. If fixed values are used, at least one of `max_unavailable_fixed` or `max_surge_fixed` must be greater than 0.

* `max_unavailable_percent` - (Optional), The maximum number of instances(calculated as percentage) that can be unavailable during the update process. Conflicts with `max_unavailable_fixed`. Percent value is only allowed for regional managed instance groups with size at least 10.

* `min_ready_sec` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)), Minimum number of seconds to wait for after a newly created instance becomes available. This value must be from range [0, 3600]

* `replacement_method` - (Optional), The instance replacement method for managed instance groups. Valid values are: "RECREATE", "SUBSTITUTE". If SUBSTITUTE (default), the group replaces VM instances with new instances that have randomly generated names. If RECREATE, instance names are preserved.  You must also set max_unavailable_fixed or max_unavailable_percent to be greater than 0.
- - -

<a name="nested_all_instances_config"></a>The `all_instances_config` block supports:

```hcl
all_instances_config {
  metadata = {
    metadata_key = "metadata_value"
  }
  labels = {
    label_key = "label_Value"
  }
}
```

* `metadata` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)), The metadata key-value pairs that you want to patch onto the instance. For more information, see [Project and instance metadata](https://cloud.google.com/compute/docs/metadata#project_and_instance_metadata).

* `labels` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)), The label key-value pairs that you want to patch onto the instance.

- - -

<a name="nested_named_port"></a>The `named_port` block supports: (Include a `named_port` block for each named-port required).

* `name` - (Required) The name of the port.

* `port` - (Required) The port number.
- - -

The `auto_healing_policies` block supports:

* `health_check` - (Required) The health check resource that signals autohealing.

* `initial_delay_sec` - (Required) The number of seconds that the managed instance group waits before
 it applies autohealing policies to new instances or recently recreated instances. Between 0 and 3600.

<a name="nested_version"></a>The `version` block supports:

```hcl
version {
  name              = "appserver-canary"
  instance_template = google_compute_instance_template.appserver-canary.id

  target_size {
    fixed = 1
  }
}
```

```hcl
version {
  name              = "appserver-canary"
  instance_template = google_compute_instance_template.appserver-canary.id

  target_size {
    percent = 20
  }
}
```

* `name` - (Required) - Version name.

* `instance_template` - (Required) - The full URL to an instance template from which all new instances of this version will be created.

* `target_size` - (Optional) - The number of instances calculated as a fixed number or a percentage depending on the settings. Structure is [documented below](#nested_target_size).

-> Exactly one `version` you specify must not have a `target_size` specified. During a rolling update, the instance group manager will fulfill the `target_size`
constraints of every other `version`, and any remaining instances will be provisioned with the version where `target_size` is unset.

<a name="nested_target_size"></a>The `target_size` block supports:

* `fixed` - (Optional), The number of instances which are managed for this version. Conflicts with `percent`.

* `percent` - (Optional), The number of instances (calculated as percentage) which are managed for this version. Conflicts with `fixed`.
Note that when using `percent`, rounding will be in favor of explicitly set `target_size` values; a managed instance group with 2 instances and 2 `version`s,
one of which has a `target_size.percent` of `60` will create 2 instances of that `version`.

<a name="nested_stateful_disk"></a>The `stateful_disk` block supports: (Include a `stateful_disk` block for each stateful disk required).

* `device_name` - (Required), The device name of the disk to be attached.

* `delete_rule` - (Optional), A value that prescribes what should happen to the stateful disk when the VM instance is deleted. The available options are `NEVER` and `ON_PERMANENT_INSTANCE_DELETION`. `NEVER` - detach the disk when the VM is deleted, but do not delete the disk. `ON_PERMANENT_INSTANCE_DELETION` will delete the stateful disk when the VM is permanently deleted from the instance group. The default is `NEVER`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `{{disk.name}}`

* `fingerprint` - The fingerprint of the instance group manager.

* `instance_group` - The full URL of the instance group created by the manager.

* `self_link` - The URL of the created resource.

The `status` block holds:

* `is_stable` - A bit indicating whether the managed instance group is in a stable state. A stable state means that: none of the instances in the managed instance group is currently undergoing any type of change (for example, creation, restart, or deletion); no future changes are scheduled for instances in the managed instance group; and the managed instance group itself is not being modified.

* `version_target` - A status of consistency of Instances' versions with their target version specified by version field on Instance Group Manager.

The `version_target` block holds:

* `version_target` - A bit indicating whether version target has been reached in this managed instance group, i.e. all instances are in their target version. Instances' target version are specified by version field on Instance Group Manager.

* `stateful` - Stateful status of the given Instance Group Manager.

The `stateful` block holds:

* `has_stateful_config` - A bit indicating whether the managed instance group has stateful configuration, that is, if you have configured any items in a stateful policy or in per-instance configs. The group might report that it has no stateful config even when there is still some preserved state on a managed instance, for example, if you have deleted all PICs but not yet applied those deletions.

* `per_instance_configs` - Status of per-instance configs on the instance.

The `per_instance_configs` block holds:

* `all_effective` - A bit indicating if all of the group's per-instance configs (listed in the output of a listPerInstanceConfigs API call) have status `EFFECTIVE` or there are no per-instance-configs.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 15 minutes.
- `update` - Default is 15 minutes.
- `delete` - Default is 15 minutes.


## Import

Instance group managers can be imported using the `name`, e.g.

```
$ terraform import google_compute_region_instance_group_manager.appserver appserver-igm
```
