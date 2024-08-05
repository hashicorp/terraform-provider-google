---
subcategory: "Compute Engine"
description: |-
  Get a VM instance template within GCE.
---

# google_compute_instance_template

-> **Note**: Global instance templates can be used in any region. To lower the impact of outages outside your region and gain data residency within your region, use [google_compute_region_instance_template](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_region_instance_template).

Get information about a VM instance template resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/instance-templates)
and
[API](https://cloud.google.com/compute/docs/reference/rest/v1/instanceTemplates).

## Example Usage

```hcl
# by name
data "google_compute_instance_template" "generic" {
  name    = "generic-tpl-20200107"
}

# using a filter
data "google_compute_instance_template" "generic-regex" {
  filter      = "name != generic-tpl-20200107"
  most_recent = true
}

# by unique ID
data "google_compute_instance_template" "generic" {
  self_link_unique    = "https://www.googleapis.com/compute/v1/projects/your-project-name/global/instanceTemplates/example-template-custom?uniqueId=1234"
}

```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the instance template. One of `name`, `filter` or `self_link_unique` must be provided.

- `filter` - (Optional) A filter to retrieve the instance templates.
    See [API filter parameter documentation](https://cloud.google.com/compute/docs/reference/rest/v1/instanceTemplates/list#body.QUERY_PARAMETERS.filter) for reference.
    If multiple instance templates match, either adjust the filter or specify `most_recent`.
	One of `name`, `filter` or `self_link_unique` must be provided.

- `self_link_unique` - (Optional) The self_link_unique URI of the instance template. One of `name`, `filter` or `self_link_unique` must be provided.

- `most_recent` - (Optional) If `filter` is provided, ensures the most recent template is returned when multiple instance templates match. One of `name`, `filter` or `self_link_unique` must be provided.

---

* `project` - (Optional) The ID of the project in which the resource belongs.
    If `project` is not provided, the provider project is used.

## Attributes Reference

* `disk` - Disks to attach to instances created from this template.
    This can be specified multiple times for multiple disks. Structure is
    [documented below](#nested_disk).

* `machine_type` - The machine type to create.

    To create a machine with a [custom type][custom-vm-types] (such as extended memory), format the value like `custom-VCPUS-MEM_IN_MB` like `custom-6-20480` for 6 vCPU and 20GB of RAM.

* `name` - The name of the instance template. If you leave
  this blank, Terraform will auto-generate a unique name.

* `name_prefix` - Creates a unique name beginning with the specified
  prefix. Conflicts with `name`.

* `can_ip_forward` - Whether to allow sending and receiving of
    packets with non-matching source or destination IPs. This defaults to false.

* `description` - A brief description of this resource.

* `instance_description` - A brief description to use for instances
    created from this template.

* `labels` - All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.

* `metadata` - Metadata key/value pairs to make available from
    within instances created from this template.

* `metadata_startup_script` - An alternative to using the
    startup-script metadata key, mostly to match the compute_instance resource.
    This replaces the startup-script metadata key on the created instance and
    thus the two mechanisms are not allowed to be used simultaneously.

* `network_interface` - Networks to attach to instances created from
    this template. This can be specified multiple times for multiple networks.
    Structure is [documented below](#nested_network_interface).

* `network_performance_config` - The network performance configuration setting
    for the instance, if set. Structure is [documented below](#nested_network_performance_config).

* `project` - The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - An instance template is a global resource that is not
    bound to a zone or a region. However, you can still specify some regional
    resources in an instance template, which restricts the template to the
    region where that resource resides. For example, a custom `subnetwork`
    resource is tied to a specific region. Defaults to the region of the
    Provider if no value is given.

* `scheduling` - The scheduling strategy to use. More details about
    this configuration option are detailed below.

* `service_account` - Service account to attach to the instance. Structure is [documented below](#nested_service_account).

* `tags` - Tags to attach to the instance.

* `guest_accelerator` - List of the type and count of accelerator cards attached to the instance. Structure [documented below](#nested_guest_accelerator).

* `min_cpu_platform` - Specifies a minimum CPU platform. Applicable values are the friendly names of CPU platforms, such as
`Intel Haswell` or `Intel Skylake`. See the complete list [here](https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform).

* `shielded_instance_config` - Enable [Shielded VM](https://cloud.google.com/security/shielded-cloud/shielded-vm) on this instance. Shielded VM provides verifiable integrity to prevent against malware and rootkits. Defaults to disabled. Structure is [documented below](#nested_shielded_instance_config).
	**Note**: [`shielded_instance_config`](#shielded_instance_config) can only be used with boot images with shielded vm support. See the complete list [here](https://cloud.google.com/compute/docs/images#shielded-images).

* `enable_display` - Enable [Virtual Displays](https://cloud.google.com/compute/docs/instances/enable-instance-virtual-display#verify_display_driver) on this instance.
**Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true in order to update this field.

* `confidential_instance_config` - Enable [Confidential Mode](https://cloud.google.com/compute/confidential-vm/docs/about-cvm) on this VM. Structure is [documented below](#nested_confidential_instance_config)

<a name="nested_disk"></a>The `disk` block supports:

* `auto_delete` - Whether or not the disk should be auto-deleted.
    This defaults to true.

* `boot` - Indicates that this is a boot disk.

* `device_name` - A unique device name that is reflected into the
    /dev/  tree of a Linux operating system running within the instance. If not
    specified, the server chooses a default device name to apply to this disk.

* `disk_name` - Name of the disk. When not provided, this defaults
    to the name of the instance.

* `provisioned_iops` - Indicates how many IOPS to provision for the disk. This
    sets the number of I/O operations per second that the disk can handle.
    Values must be between 10,000 and 120,000. For more details, see the
    [Extreme persistent disk documentation](https://cloud.google.com/compute/docs/disks/extreme-persistent-disk).

* `source_image` - The image from which to
    initialize this disk. This can be one of: the image's `self_link`,
    `projects/{project}/global/images/{image}`,
    `projects/{project}/global/images/family/{family}`, `global/images/{image}`,
    `global/images/family/{family}`, `family/{family}`, `{project}/{family}`,
    `{project}/{image}`, `{family}`, or `{image}`.
~> **Note:** Either `source` or `source_image` is **required** in a disk block unless the disk type is `local-ssd`. Check the API [docs](https://cloud.google.com/compute/docs/reference/rest/v1/instanceTemplates/insert) for details.

* `interface` - Specifies the disk interface to use for attaching this disk,
    which is either SCSI or NVME. The default is SCSI. Persistent disks must always use SCSI
    and the request will fail if you attempt to attach a persistent disk in any other format
    than SCSI. Local SSDs can use either NVME or SCSI.

* `mode` - The mode in which to attach this disk, either READ_WRITE
    or READ_ONLY. If you are attaching or creating a boot disk, this must
    read-write mode.

* `source` - The name (**not self_link**)
    of the disk (such as those managed by `google_compute_disk`) to attach.
~> **Note:** Either `source` or `source_image` is **required** in a disk block unless the disk type is `local-ssd`. Check the API [docs](https://cloud.google.com/compute/docs/reference/rest/v1/instanceTemplates/insert) for details.

* `disk_type` - The GCE disk type. Such as `"pd-ssd"`, `"local-ssd"`,
    `"pd-balanced"` or `"pd-standard"`.

* `disk_size_gb` - The size of the image in gigabytes. If not
    specified, it will inherit the size of its base image. For SCRATCH disks,
    the size must be exactly 375GB.

* `labels` - (Optional) A set of ket/value label pairs to assign to disk created from
    this template

* `type` - The type of GCE disk, can be either `"SCRATCH"` or
    `"PERSISTENT"`.

* `disk_encryption_key` - Encrypts or decrypts a disk using a customer-supplied encryption key.

    If you are creating a new disk, this field encrypts the new disk using an encryption key that you provide. If you are attaching an existing disk that is already encrypted, this field decrypts the disk using the customer-supplied encryption key.

    If you encrypt a disk using a customer-supplied key, you must provide the same key again when you attempt to use this resource at a later time. For example, you must provide the key when you create a snapshot or an image from the disk or when you attach the disk to a virtual machine instance.

    If you do not provide an encryption key, then the disk will be encrypted using an automatically generated key and you do not need to provide a key to use the disk later.

    Instance templates do not store customer-supplied encryption keys, so you cannot use your own keys to encrypt disks in a managed instance group.

* `resource_policies` (Optional) -- A list of short names of resource policies to attach to this disk for automatic snapshot creations. Currently a max of 1 resource policy is supported.

The `disk_encryption_key` block supports:

* `kms_key_self_link` - The self link of the encryption key that is stored in Google Cloud KMS

<a name="nested_network_interface"></a>The `network_interface` block supports:

* `network` - The name or self_link of the network to attach this interface to.
    Use `network` attribute for Legacy or Auto subnetted networks and
    `subnetwork` for custom subnetted networks.

* `subnetwork` - the name of the subnetwork to attach this interface
    to. The subnetwork must exist in the same `region` this instance will be
    created in. Either `network` or `subnetwork` must be provided.

* `network_interface` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) The URL of the network attachment that this interface should connect to in the following format: projects/{projectNumber}/regions/{region_name}/networkAttachments/{network_attachment_name}.  s

* `subnetwork_project` - The ID of the project in which the subnetwork belongs.
    If it is not provided, the provider project is used.

* `network_ip` - The private IP address to assign to the instance. If
    empty, the address will be automatically assigned.

* `access_config` - Access configurations, i.e. IPs via which this
    instance can be accessed via the Internet. Omit to ensure that the instance
    is not accessible from the Internet (this means that ssh provisioners will
    not work unless you are running Terraform can send traffic to the instance's
    network (e.g. via tunnel or because it is running on another cloud instance
    on that network). This block can be repeated multiple times. Structure [documented below](#nested_access_config).

* `alias_ip_range` - An
    array of alias IP ranges for this network interface. Can only be specified for network
    interfaces on subnet-mode networks. Structure [documented below](#nested_alias_ip_range).

<a name="nested_access_config"></a>The `access_config` block supports:

* `nat_ip` - The IP address that will be 1:1 mapped to the instance's
    network ip. If not given, one will be generated.

* `network_tier` - The [networking tier][network-tier] used for configuring
    this instance template. This field can take the following values: PREMIUM or
    STANDARD. If this field is not specified, it is assumed to be PREMIUM.

<a name="nested_alias_ip_range"></a>The `alias_ip_range` block supports:

* `ip_cidr_range` - The IP CIDR range represented by this alias IP range. This IP CIDR range
    must belong to the specified subnetwork and cannot contain IP addresses reserved by
    system or used by other network interfaces. At the time of writing only a
    netmask (e.g. /24) may be supplied, with a CIDR format resulting in an API
    error.

* `subnetwork_range_name` - The subnetwork secondary range name specifying
    the secondary range from which to allocate the IP CIDR range for this alias IP
    range. If left unspecified, the primary range of the subnetwork will be used.

<a name="nested_service_account"></a>The `service_account` block supports:

* `email` - The service account e-mail address. If not given, the
    default Google Compute Engine service account is used.

* `scopes` - A list of service scopes. Both OAuth2 URLs and gcloud
    short names are supported. To allow full access to all Cloud APIs, use the
    `cloud-platform` scope. See a complete list of scopes [here](https://cloud.google.com/sdk/gcloud/reference/alpha/compute/instances/set-scopes#--scopes).

    The [service accounts documentation](https://cloud.google.com/compute/docs/access/service-accounts#accesscopesiam)
    explains that access scopes are the legacy method of specifying permissions for your instance.
    If you are following best practices and using IAM roles to grant permissions to service accounts,
    then you can define this field as an empty list.

<a name="nested_scheduling"></a>The `scheduling` block supports:

* `automatic_restart` - Specifies whether the instance should be
    automatically restarted if it is terminated by Compute Engine (not
    terminated by a user). This defaults to true.

* `on_host_maintenance` - Defines the maintenance behavior for this
    instance.

* `preemptible` - Allows instance to be preempted. This defaults to
    false. Read more on this
    [here](https://cloud.google.com/compute/docs/instances/preemptible).

* `node_affinities` - Specifies node affinities or anti-affinities
   to determine which sole-tenant nodes your instances and managed instance
   groups will use as host systems. Read more on sole-tenant node creation
   [here](https://cloud.google.com/compute/docs/nodes/create-nodes).
   Structure [documented below](#nested_node_affinities).
   
* `provisioning_model` - Describe the type of preemptible VM. 

* `instance_termination_action` - Describe the type of termination action for `SPOT` VM. Can be `STOP` or `DELETE`.  Read more on [here](https://cloud.google.com/compute/docs/instances/create-use-spot) 

<a name="nested_guest_accelerator"></a>The `guest_accelerator` block supports:

* `type` - The accelerator type resource to expose to this instance. E.g. `nvidia-tesla-k80`.

* `count` - The number of the guest accelerator cards exposed to this instance.

<a name="nested_node_affinities"></a>The `node_affinities` block supports:

* `key` - The key for the node affinity label.

* `operator` - The operator. Can be `IN` for node-affinities
    or `NOT_IN` for anti-affinities.

* `value` - The values for the node affinity label.

<a name="nested_shielded_instance_config"></a>The `shielded_instance_config` block supports:

* `enable_secure_boot` -- Verify the digital signature of all boot components, and halt the boot process if signature verification fails. Defaults to false.

* `enable_vtpm` -- Use a virtualized trusted platform module, which is a specialized computer chip you can use to encrypt objects like keys and certificates. Defaults to true.

* `enable_integrity_monitoring` -- Compare the most recent boot measurements to the integrity policy baseline and return a pair of pass/fail results depending on whether they match or not. Defaults to true.

<a name="nested_confidential_instance_config"></a>The `confidential_instance_config` block supports:

* `enable_confidential_compute` Defines whether the instance should have confidential compute enabled. [`on_host_maintenance`](#on_host_maintenance) has to be set to TERMINATE or this will fail to create the VM.

<a name="nested_network_performance_config"></a>The `network_performance_config` block supports:

* `total_egress_bandwidth_tier` - The egress bandwidth tier for the instance.

---

* `id` - an identifier for the resource with format `projects/{{project}}/global/instanceTemplates/{{name}}`

* `metadata_fingerprint` - The unique fingerprint of the metadata.

* `self_link` - The URI of the created resource.

* `self_link_unique` - A special URI of the created resource that uniquely identifies this instance template with the following format: `projects/{{project}}/global/instanceTemplates/{{name}}?uniqueId={{uniqueId}}`
Referencing an instance template via this attribute prevents Time of Check to Time of Use attacks when the instance template resides in a shared/untrusted environment.

* `tags_fingerprint` - The unique fingerprint of the tags.

[1]: /docs/providers/google/r/compute_instance_group_manager.html
[2]: /docs/configuration/resources.html#lifecycle
