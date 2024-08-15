---
subcategory: "Compute Engine"
description: |-
  Manages an Instance Group within GCE.
---

# google_compute_instance_group_manager

The Google Compute Engine Instance Group Manager API creates and manages pools
of homogeneous Compute Engine virtual machine instances from a common instance
template. For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups/manager)
and [API](https://cloud.google.com/compute/docs/reference/latest/instanceGroupManagers)

~> **Note:** Use [google_compute_region_instance_group_manager](/docs/providers/google/r/compute_region_instance_group_manager.html) to create a regional (multi-zone) instance group manager.

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

resource "google_compute_instance_group_manager" "appserver" {
  name = "appserver-igm"

  base_instance_name = "app"
  zone               = "us-central1-a"

  version {
    instance_template  = google_compute_instance_template.appserver.self_link_unique
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
    name = "customhttp"
    port = 8888
  }

  auto_healing_policies {
    health_check      = google_compute_health_check.autohealing.id
    initial_delay_sec = 300
  }
}
```

## Example Usage with multiple versions (`google-beta` provider)
```hcl
resource "google_compute_instance_group_manager" "appserver" {
  provider = google-beta
  name     = "appserver-igm"

  base_instance_name = "app"
  zone               = "us-central1-a"

  target_size = 5

  version {
    name              = "appserver"
    instance_template = google_compute_instance_template.appserver.self_link_unique
  }

  version {
    name              = "appserver-canary"
    instance_template = google_compute_instance_template.appserver-canary.self_link_unique
    target_size {
      fixed = 1
    }
  }
}
```

## Example Usage with standby policy (`google-beta` provider)
```hcl
resource "google_compute_instance_group_manager" "igm-sr" {
  provider = google-beta
  name = "tf-sr-igm"

  base_instance_name        = "tf-sr-igm-instance"
  zone                      = "us-central1-a"

  target_size               = 5

  version {
    instance_template = google_compute_instance_template.sr-igm.self_link
    name              = "primary"
  }

  standby_policy {
    initial_delay_sec           = 30
    mode                        = "MANUAL"
  }
  target_suspended_size         = 2
  target_stopped_size           = 1
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

* `zone` - (Required) The zone that instances in this group should be created
    in.

- - -

* `description` - (Optional) An optional textual description of the instance
    group manager.

* `named_port` - (Optional) The named port configuration. See the [section below](#nested_named_port)
    for details on configuration.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `target_size` - (Optional) The target number of running instances for this managed
    instance group. This value will fight with autoscaler settings when set, and generally shouldn't be set
    when using one. If a value is required, such as to specify a creation-time target size for the MIG,
    `lifecycle.ignore_changes` can be used to prevent Terraform from modifying the value. Defaults to `0`.

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
group. You can specify only one value. Structure is [documented below](#nested_auto_healing_policies). For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/creating-groups-of-managed-instances#monitoring_groups).

* `all_instances_config` - (Optional)
  Properties to set on all instances in the group. After setting
  allInstancesConfig on the group, you must update the group's instances to
  apply the configuration.

* `standby_policy` - (Optional [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) The standby policy for stopped and suspended instances. Structure is documented below. For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/suspended-and-stopped-vms-in-mig) and [API](https://cloud.google.com/compute/docs/reference/rest/beta/regionInstanceGroupManagers/patch)

* `target_suspended_size` - (Optional [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) The target number of suspended instances for this managed instance group.

* `target_stopped_size` - (Optional [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) The target number of stopped instances for this managed instance group.

* `stateful_disk` - (Optional) Disks created on the instances that will be preserved on instance delete, update, etc. Structure is [documented below](#nested_stateful_disk). For more information see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/configuring-stateful-disks-in-migs).

* `stateful_internal_ip` - (Optional) Internal network IPs assigned to the instances that will be preserved on instance delete, update, etc. This map is keyed with the network interface name. Structure is [documented below](#nested_stateful_internal_ip).

* `stateful_external_ip` - (Optional) External network IPs assigned to the instances that will be preserved on instance delete, update, etc. This map is keyed with the network interface name. Structure is [documented below](#nested_stateful_external_ip).

* `update_policy` - (Optional) The update policy for this managed instance group. Structure is [documented below](#nested_update_policy). For more information, see the [official documentation](https://cloud.google.com/compute/docs/instance-groups/updating-managed-instance-groups) and [API](https://cloud.google.com/compute/docs/reference/rest/v1/instanceGroupManagers/patch).

* `params` - (Optional [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) Input only additional params for instance group manager creation. Structure is [documented below](#nested_params). For more information, see [API](https://cloud.google.com/compute/docs/reference/rest/beta/instanceGroupManagers/insert).

- - -

The `standby_policy` block supports:

* `initial_delay_sec` - (Optional) - Specifies the number of seconds that the MIG should wait to suspend or stop a VM after that VM was created. The initial delay gives the initialization script the time to prepare your VM for a quick scale out. The value of initial delay must be between 0 and 3600 seconds. The default value is 0.
* `mode` - (Optional) - Defines how a MIG resumes or starts VMs from a standby pool when the group scales out. Valid options are: `MANUAL`, `SCALE_OUT_POOL`. If `MANUAL`(default), you have full control over which VMs are stopped and suspended in the MIG. If `SCALE_OUT_POOL`, the MIG uses the VMs from the standby pools to accelerate the scale out by resuming or starting them and then automatically replenishes the standby pool with new VMs to maintain the target sizes.
- - -

<a name="nested_update_policy"></a>The `update_policy` block supports:

```hcl
update_policy {
  type                           = "PROACTIVE"
  minimal_action                 = "REPLACE"
  most_disruptive_allowed_action = "REPLACE"
  max_surge_fixed                = 0
  max_unavailable_fixed          = 2
  min_ready_sec                  = 50
  replacement_method             = "RECREATE"
}
```

* `minimal_action` - (Required) - Minimal action to be taken on an instance. You can specify either `NONE` to forbid any actions, `REFRESH` to update without stopping instances, `RESTART` to restart existing instances or `REPLACE` to delete and create new instances from the target template. If you specify a `REFRESH`, the Updater will attempt to perform that action only. However, if the Updater determines that the minimal action you specify is not enough to perform the update, it might perform a more disruptive action.

* `most_disruptive_allowed_action` - (Optional) - Most disruptive action that is allowed to be taken on an instance. You can specify either NONE to forbid any actions, REFRESH to allow actions that do not need instance restart, RESTART to allow actions that can be applied without instance replacing or REPLACE to allow all possible actions. If the Updater determines that the minimal update action needed is more disruptive than most disruptive allowed action you specify it will not perform the update at all.

* `type` - (Required) - The type of update process. You can specify either `PROACTIVE` so that the instance group manager proactively executes actions in order to bring instances to their target versions or `OPPORTUNISTIC` so that no action is proactively executed but the update will be performed as part of other actions (for example, resizes or recreateInstances calls).

* `max_surge_fixed` - (Optional), Specifies a fixed number of VM instances. This must be a positive integer. Conflicts with `max_surge_percent`. Both cannot be 0.

* `max_surge_percent` - (Optional), Specifies a percentage of instances between 0 to 100%, inclusive. For example, specify 80 for 80%. Conflicts with `max_surge_fixed`.

* `max_unavailable_fixed` - (Optional), Specifies a fixed number of VM instances. This must be a positive integer.

* `max_unavailable_percent` - (Optional), Specifies a percentage of instances between 0 to 100%, inclusive. For example, specify 80 for 80%..

* `min_ready_sec` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)), Minimum number of seconds to wait for after a newly created instance becomes available. This value must be from range [0, 3600]

* `replacement_method` - (Optional), The instance replacement method for managed instance groups. Valid values are: "RECREATE", "SUBSTITUTE". If SUBSTITUTE (default), the group replaces VM instances with new instances that have randomly generated names. If RECREATE, instance names are preserved.  You must also set max_unavailable_fixed or max_unavailable_percent to be greater than 0.
- - -

<a name="nested_instance_lifecycle_policy"></a>The `instance_lifecycle_policy` block supports:

```hcl
instance_lifecycle_policy {
  force_update_on_repair    = "YES"
  default_action_on_failure = "DO_NOTHING"
}
```

* `force_update_on_repair` - (Optional), Specifies whether to apply the group's latest configuration when repairing a VM. Valid options are: `YES`, `NO`. If `YES` and you updated the group's instance template or per-instance configurations after the VM was created, then these changes are applied when VM is repaired. If `NO` (default), then updates are applied in accordance with the group's update policy type.
* `default_action_on_failure` - (Optional), Default behavior for all instance or health check failures. Valid options are: `REPAIR`, `DO_NOTHING`. If `DO_NOTHING` then instances will not be repaired. If `REPAIR` (default), then failed instances will be repaired.
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

* `metadata` - (Optional), The metadata key-value pairs that you want to patch onto the instance. For more information, see [Project and instance metadata](https://cloud.google.com/compute/docs/metadata#project_and_instance_metadata).

* `labels` - (Optional), The label key-value pairs that you want to patch onto the instance.

- - -

<a name="nested_named_port"></a>The `named_port` block supports: (Include a `named_port` block for each named-port required).

* `name` - (Required) The name of the port.

* `port` - (Required) The port number.
- - -

<a name="nested_auto_healing_policies"></a>The `auto_healing_policies` block supports:

* `health_check` - (Required) The health check resource that signals autohealing.

* `initial_delay_sec` - (Required) The number of seconds that the managed instance group waits before
 it applies autohealing policies to new instances or recently recreated instances. Between 0 and 3600.

<a name="nested_version"></a>The `version` block supports:

```hcl
version {
  name              = "appserver-canary"
  instance_template = google_compute_instance_template.appserver-canary.self_link_unique

  target_size {
    fixed = 1
  }
}
```

```hcl
version {
  name              = "appserver-canary"
  instance_template = google_compute_instance_template.appserver-canary.self_link_unique

  target_size {
    percent = 20
  }
}
```

* `name` - (Required) - Version name.

* `instance_template` - (Required) - The full URL to an instance template from which all new instances of this version will be created. It is recommended to reference instance templates through their unique id (`self_link_unique` attribute).

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

<a name="nested_stateful_internal_ip"></a>The `stateful_internal_ip` block supports:

* `interface_name` - (Required), The network interface name of the internal Ip. Possible value: `nic0`

* `delete_rule` - (Optional), A value that prescribes what should happen to the internal ip when the VM instance is deleted. The available options are `NEVER` and `ON_PERMANENT_INSTANCE_DELETION`. `NEVER` - detach the ip when the VM is deleted, but do not delete the ip. `ON_PERMANENT_INSTANCE_DELETION` will delete the internal ip when the VM is permanently deleted from the instance group.

<a name="nested_stateful_external_ip"></a>The `stateful_external_ip` block supports:

* `interface_name` - (Required), The network interface name of the external Ip. Possible value: `nic0`

* `delete_rule` - (Optional), A value that prescribes what should happen to the external ip when the VM instance is deleted. The available options are `NEVER` and `ON_PERMANENT_INSTANCE_DELETION`. `NEVER` - detach the ip when the VM is deleted, but do not delete the ip. `ON_PERMANENT_INSTANCE_DELETION` will delete the external ip when the VM is permanently deleted from the instance group.

<a name="nested_params"></a>The `params` block supports:

```hcl
params{
  resource_manager_tags = {
    "tagKeys/123": "tagValues/123"
  }
}
```

* `resource_manager_tags` - (Optional) Resource manager tags to bind to the managed instance group. The tags are key-value pairs. Keys must be in the format tagKeys/123 and values in the format tagValues/456. For more information, see [Manage tags for resources](https://cloud.google.com/compute/docs/tag-resources)

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{name}}`

* `creation_timestamp` - Creation timestamp in RFC3339 text format.

* `fingerprint` - The fingerprint of the instance group manager.

* `instance_group` - The full URL of the instance group created by the manager.

* `self_link` - The URL of the created resource.

* `status` - The status of this managed instance group.

The `status` block holds:

* `is_stable` - A bit indicating whether the managed instance group is in a stable state. A stable state means that: none of the instances in the managed instance group is currently undergoing any type of change (for example, creation, restart, or deletion); no future changes are scheduled for instances in the managed instance group; and the managed instance group itself is not being modified.

* `version_target` - A status of consistency of Instances' versions with their target version specified by version field on Instance Group Manager.

* `all_instances_config` - Status of all-instances configuration on the group.

* `stateful` - Stateful status of the given Instance Group Manager.

The `version_target` block holds:

* `version_target` - A bit indicating whether version target has been reached in this managed instance group, i.e. all instances are in their target version. Instances' target version are specified by version field on Instance Group Manager.

The `all_instances_config` block holds:

* `effective` -  A bit indicating whether this configuration has been applied to all managed instances in the group.

* `current_revision` - Current all-instances configuration revision. This value is in RFC3339 text format.

The `stateful` block holds:

* `has_stateful_config` - A bit indicating whether the managed instance group has stateful configuration, that is, if you have configured any items in a stateful policy or in per-instance configs. The group might report that it has no stateful config even when there is still some preserved state on a managed instance, for example, if you have deleted all PICs but not yet applied those deletions.

* `per_instance_configs` - Status of per-instance configs on the instances.

The `per_instance_configs` block holds:

* `all_effective` - A bit indicating if all of the group's per-instance configs (listed in the output of a listPerInstanceConfigs API call) have status `EFFECTIVE` or there are no per-instance-configs.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 15 minutes.
- `update` - Default is 15 minutes.
- `delete` - Default is 15 minutes.


## Import

Instance group managers can be imported using any of these accepted formats:

```
* `projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{name}}`
* `{{project}}/{{zone}}/{{name}}`
* `{{project}}/{{name}}`
* `{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import instance group managers using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{name}}"
  to = google_compute_instance_group_manager.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), instance group managers can be imported using one of the formats above. For example:

```
$ terraform import google_compute_instance_group_manager.default projects/{{project}}/zones/{{zone}}/instanceGroupManagers/{{name}}
$ terraform import google_compute_instance_group_manager.default {{project}}/{{zone}}/{{name}}
$ terraform import google_compute_instance_group_manager.default {{project}}/{{name}}
$ terraform import google_compute_instance_group_manager.default {{name}}
```
