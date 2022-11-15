---
subcategory: "Compute Engine"
page_title: "Google: google_compute_instance"
description: |-
  Manages a VM instance resource within GCE.
---

# google\_compute\_instance

Manages a VM instance resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/instances)
and
[API](https://cloud.google.com/compute/docs/reference/latest/instances).

## Example Usage

```hcl
resource "google_service_account" "default" {
  account_id   = "service_account_id"
  display_name = "Service Account"
}

resource "google_compute_instance" "default" {
  name         = "test"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  tags = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
      labels = {
        my_label = "value"
      }
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral public IP
    }
  }

  metadata = {
    foo = "bar"
  }

  metadata_startup_script = "echo hi > /test.txt"

  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = google_service_account.default.email
    scopes = ["cloud-platform"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `boot_disk` - (Required) The boot disk for the instance.
    Structure is [documented below](#nested_boot_disk).

* `machine_type` - (Required) The machine type to create.

    **Note:** If you want to update this value (resize the VM) after initial creation, you must set [`allow_stopping_for_update`](#allow_stopping_for_update) to `true`.

    [Custom machine types][custom-vm-types] can be formatted as `custom-NUMBER_OF_CPUS-AMOUNT_OF_MEMORY_MB`, e.g. `custom-6-20480` for 6 vCPU and 20GB of RAM.

    There is a limit of 6.5 GB per CPU unless you add [extended memory][extended-custom-vm-type]. You must do this explicitly by adding the suffix `-ext`, e.g. `custom-2-15360-ext` for 2 vCPU and 15 GB of memory.

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

* `zone` - (Optional) The zone that the machine should be created in. If it is not provided, the provider zone is used.

* `network_interface` - (Required) Networks to attach to the instance. This can
    be specified multiple times. Structure is [documented below](#nested_network_interface).

- - -

* `allow_stopping_for_update` - (Optional) If true, allows Terraform to stop the instance to update its properties.
  If you try to update a property that requires stopping the instance without setting this field, the update will fail.

* `attached_disk` - (Optional) Additional disks to attach to the instance. Can be repeated multiple times for multiple disks. Structure is [documented below](#nested_attached_disk).

* `can_ip_forward` - (Optional) Whether to allow sending and receiving of
    packets with non-matching source or destination IPs.
    This defaults to false.

* `description` - (Optional) A brief description of this resource.

* `desired_status` - (Optional) Desired status of the instance. Either
`"RUNNING"` or `"TERMINATED"`.

* `deletion_protection` - (Optional) Enable deletion protection on this instance. Defaults to false.
    **Note:** you must disable deletion protection before removing the resource (e.g., via `terraform destroy`), or the instance cannot be deleted and the Terraform run will not complete successfully.

* `hostname` - (Optional) A custom hostname for the instance. Must be a fully qualified DNS name and RFC-1035-valid.
  Valid format is a series of labels 1-63 characters long matching the regular expression `[a-z]([-a-z0-9]*[a-z0-9])`, concatenated with periods.
  The entire hostname must not exceed 253 characters. Changing this forces a new resource to be created.

* `guest_accelerator` - (Optional) List of the type and count of accelerator cards attached to the instance. Structure [documented below](#nested_guest_accelerator).
    **Note:** GPU accelerators can only be used with [`on_host_maintenance`](#on_host_maintenance) option set to TERMINATE.
    **Note**: This field uses [attr-as-block mode](https://www.terraform.io/docs/configuration/attr-as-blocks.html) to avoid
    breaking users during the 0.12 upgrade. To explicitly send a list
    of zero objects you must use the following syntax:
    `example=[]`
    For more details about this behavior, see [this section](https://www.terraform.io/docs/configuration/attr-as-blocks.html#defining-a-fixed-object-collection-value).

* `labels` - (Optional) A map of key/value label pairs to assign to the instance.

* `metadata` - (Optional) Metadata key/value pairs to make available from
    within the instance. Ssh keys attached in the Cloud Console will be removed.
    Add them to your config in order to keep them attached to your instance. A
    list of default metadata values (e.g. ssh-keys) can be found [here](https://cloud.google.com/compute/docs/metadata/default-metadata-values)

-> Depending on the OS you choose for your instance, some metadata keys have
   special functionality.  Most linux-based images will run the content of
   `metadata.startup-script` in a shell on every boot.  At a minimum,
   Debian, CentOS, RHEL, SLES, Container-Optimized OS, and Ubuntu images
   support this key.  Windows instances require other keys depending on the format
   of the script and the time you would like it to run - see [this table](https://cloud.google.com/compute/docs/startupscript#providing_a_startup_script_for_windows_instances).
   For the convenience of the users of `metadata.startup-script`,
   we provide a special attribute, `metadata_startup_script`, which is documented below.

* `metadata_startup_script` - (Optional) An alternative to using the
startup-script metadata key, except this one forces the instance to be recreated
(thus re-running the script) if it is changed. This replaces the startup-script
metadata key on the created instance and thus the two mechanisms are not
allowed to be used simultaneously.  Users are free to use either mechanism - the
only distinction is that this separate attribute will cause a recreate on
modification.  On import, `metadata_startup_script` will not be set - if you
choose to specify it you will see a diff immediately after import causing a
destroy/recreate operation. If importing an instance and specifying this value
is desired, you will need to modify your state file manually using
`terraform state` commands.

* `min_cpu_platform` - (Optional) Specifies a minimum CPU platform for the VM instance. Applicable values are the friendly names of CPU platforms, such as
`Intel Haswell` or `Intel Skylake`. See the complete list [here](https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform).
    **Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `scheduling` - (Optional) The scheduling strategy to use. More details about
    this configuration option are [detailed below](#nested_scheduling).

* `scratch_disk` - (Optional) Scratch disks to attach to the instance. This can be
    specified multiple times for multiple scratch disks. Structure is [documented below](#nested_scratch_disk).

* `service_account` - (Optional) Service account to attach to the instance.
    Structure is [documented below](#nested_service_account).
    **Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

* `tags` - (Optional) A list of network tags to attach to the instance.

* `shielded_instance_config` - (Optional) Enable [Shielded VM](https://cloud.google.com/security/shielded-cloud/shielded-vm) on this instance. Shielded VM provides verifiable integrity to prevent against malware and rootkits. Defaults to disabled. Structure is [documented below](#nested_shielded_instance_config).
	**Note**: [`shielded_instance_config`](#shielded_instance_config) can only be used with boot images with shielded vm support. See the complete list [here](https://cloud.google.com/compute/docs/images#shielded-images).
  **Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

* `enable_display` - (Optional) Enable [Virtual Displays](https://cloud.google.com/compute/docs/instances/enable-instance-virtual-display#verify_display_driver) on this instance.
**Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

* `resource_policies` (Optional) -- A list of self_links of resource policies to attach to the instance. Modifying this list will cause the instance to recreate. Currently a max of 1 resource policy is supported.

* `reservation_affinity` - (Optional) Specifies the reservations that this instance can consume from.
    Structure is [documented below](#nested_reservation_affinity).

* `confidential_instance_config` (Optional) - Enable [Confidential Mode](https://cloud.google.com/compute/confidential-vm/docs/about-cvm) on this VM. Structure is [documented below](#nested_confidential_instance_config)

* `advanced_machine_features` (Optional) - Configure Nested Virtualisation and Simultaneous Hyper Threading  on this VM. Structure is [documented below](#nested_advanced_machine_features)

* `network_performance_config` (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)
    Configures network performance settings for the instance. Structure is
    [documented below](#nested_network_performance_config). **Note**: [`machine_type`](#machine_type) must be a [supported type](https://cloud.google.com/compute/docs/networking/configure-vm-with-high-bandwidth-configuration),
    the [`image`](#image) used must include the [`GVNIC`](https://cloud.google.com/compute/docs/networking/using-gvnic#create-instance-gvnic-image)
    in `guest-os-features`, and `network_interface.0.nic-type` must be `GVNIC`
    in order for this setting to take effect.

---

<a name="nested_boot_disk"></a>The `boot_disk` block supports:

* `auto_delete` - (Optional) Whether the disk will be auto-deleted when the instance
    is deleted. Defaults to true.

* `device_name` - (Optional) Name with which attached disk will be accessible.
    On the instance, this device will be `/dev/disk/by-id/google-{{device_name}}`.

* `mode` - (Optional) The mode in which to attach this disk, either `READ_WRITE`
  or `READ_ONLY`. If not specified, the default is to attach the disk in `READ_WRITE` mode.

* `disk_encryption_key_raw` - (Optional) A 256-bit [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption),
    encoded in [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    to encrypt this disk. Only one of `kms_key_self_link` and `disk_encryption_key_raw`
    may be set.

* `kms_key_self_link` - (Optional) The self_link of the encryption key that is
    stored in Google Cloud KMS to encrypt this disk. Only one of `kms_key_self_link`
    and `disk_encryption_key_raw` may be set.

* `initialize_params` - (Optional) Parameters for a new disk that will be created
    alongside the new instance. Either `initialize_params` or `source` must be set.
    Structure is [documented below](#nested_initialize_params).

* `source` - (Optional) The name or self_link of the existing disk (such as those managed by
    `google_compute_disk`) or disk image. To create an instance from a snapshot, first create a
    `google_compute_disk` from a snapshot and reference it here.

<a name="nested_initialize_params"></a>The `initialize_params` block supports:

* `size` - (Optional) The size of the image in gigabytes. If not specified, it
    will inherit the size of its base image.

* `type` - (Optional) The GCE disk type. Such as pd-standard, pd-balanced or pd-ssd.

* `image` - (Optional) The image from which to initialize this disk. This can be
    one of: the image's `self_link`, `projects/{project}/global/images/{image}`,
    `projects/{project}/global/images/family/{family}`, `global/images/{image}`,
    `global/images/family/{family}`, `family/{family}`, `{project}/{family}`,
    `{project}/{image}`, `{family}`, or `{image}`. If referred by family, the
    images names must include the family name. If they don't, use the
    [google_compute_image data source](/docs/providers/google/d/compute_image.html).
    For instance, the image `centos-6-v20180104` includes its family name `centos-6`.
    These images can be referred by family name here.

* `labels` - (Optional) A set of key/value label pairs assigned to the disk. This  
    field is only applicable for persistent disks.

<a name="nested_scratch_disk"></a>The `scratch_disk` block supports:

* `interface` - (Required) The disk interface to use for attaching this disk; either SCSI or NVME.

<a name="nested_attached_disk"></a>The `attached_disk` block supports:

* `source` - (Required) The name or self_link of the disk to attach to this instance.

* `device_name` - (Optional) Name with which the attached disk will be accessible
    under `/dev/disk/by-id/google-*`

* `mode` - (Optional) Either "READ_ONLY" or "READ_WRITE", defaults to "READ_WRITE"
    If you have a persistent disk with data that you want to share
    between multiple instances, detach it from any read-write instances and
    attach it to one or more instances in read-only mode.

* `disk_encryption_key_raw` - (Optional) A 256-bit [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption),
    encoded in [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    to encrypt this disk. Only one of `kms_key_self_link` and `disk_encryption_key_raw` may be set.

* `kms_key_self_link` - (Optional) The self_link of the encryption key that is
    stored in Google Cloud KMS to encrypt this disk. Only one of `kms_key_self_link`
    and `disk_encryption_key_raw` may be set.

<a name="nested_network_performance_config"></a>The `network_performance_config` block supports:

* `total_egress_bandwidth_tier` - (Optional) The egress bandwidth tier to enable.
    Possible values: TIER_1, DEFAULT

<a name="nested_network_interface"></a>The `network_interface` block supports:

* `network` - (Optional) The name or self_link of the network to attach this interface to.
    Either `network` or `subnetwork` must be provided. If network isn't provided it will
    be inferred from the subnetwork.

*  `subnetwork` - (Optional) The name or self_link of the subnetwork to attach this
    interface to. Either `network` or `subnetwork` must be provided. If network isn't provided
    it will be inferred from the subnetwork. The subnetwork must exist in the same region this
    instance will be created in. If the network resource is in
    [legacy](https://cloud.google.com/vpc/docs/legacy) mode, do not specify this field. If the
    network is in auto subnet mode, specifying the subnetwork is optional. If the network is
    in custom subnet mode, specifying the subnetwork is required.

*  `subnetwork_project` - (Optional) The project in which the subnetwork belongs.
   If the `subnetwork` is a self_link, this field is ignored in favor of the project
   defined in the subnetwork self_link. If the `subnetwork` is a name and this
   field is not provided, the provider project is used.

* `network_ip` - (Optional) The private IP address to assign to the instance. If
    empty, the address will be automatically assigned.

* `access_config` - (Optional) Access configurations, i.e. IPs via which this
    instance can be accessed via the Internet. Omit to ensure that the instance
    is not accessible from the Internet. If omitted, ssh provisioners will not
    work unless Terraform can send traffic to the instance's network (e.g. via
    tunnel or because it is running on another cloud instance on that network).
    This block can be repeated multiple times. Structure [documented below](#nested_access_config).

* `alias_ip_range` - (Optional) An
    array of alias IP ranges for this network interface. Can only be specified for network
    interfaces on subnet-mode networks. Structure [documented below](#nested_alias_ip_range).

* `nic_type` - (Optional) The type of vNIC to be used on this interface. Possible values: GVNIC, VIRTIO_NET.

* `stack_type` - (Optional) The stack type for this network interface to identify whether the IPv6 feature is enabled or not. Values are IPV4_IPV6 or IPV4_ONLY. If not specified, IPV4_ONLY will be used.

* `ipv6_access_config` - (Optional) An array of IPv6 access configurations for this interface.
Currently, only one IPv6 access config, DIRECT_IPV6, is supported. If there is no ipv6AccessConfig
specified, then this instance will have no external IPv6 Internet access. Structure [documented below](#nested_ipv6_access_config).

* `queue_count` - (Optional) The networking queue count that's specified by users for the network interface. Both Rx and Tx queues will be set to this number. It will be empty if not specified.


<a name="nested_access_config"></a>The `access_config` block supports:

* `nat_ip` - (Optional) The IP address that will be 1:1 mapped to the instance's
    network ip. If not given, one will be generated.

* `public_ptr_domain_name` - (Optional) The DNS domain name for the public PTR record.
    To set this field on an instance, you must be verified as the owner of the domain.
    See [the docs](https://cloud.google.com/compute/docs/instances/create-ptr-record) for how
    to become verified as a domain owner.

* `network_tier` - (Optional) The [networking tier][network-tier] used for configuring this instance.
    This field can take the following values: PREMIUM, FIXED_STANDARD or STANDARD. If this field is
    not specified, it is assumed to be PREMIUM.

<a name="nested_ipv6_access_config"></a>The `ipv6_access_config` block supports:

* `public_ptr_domain_name` - (Optional) The domain name to be used when creating DNSv6
    records for the external IPv6 ranges..

* `network_tier` - (Optional) The service-level to be provided for IPv6 traffic when the
    subnet has an external subnet. Only PREMIUM or STANDARD tier is valid for IPv6.

<a name="nested_alias_ip_range"></a>The `alias_ip_range` block supports:

* `ip_cidr_range` - The IP CIDR range represented by this alias IP range. This IP CIDR range
    must belong to the specified subnetwork and cannot contain IP addresses reserved by
    system or used by other network interfaces. This range may be a single IP address
    (e.g. 10.2.3.4), a netmask (e.g. /24) or a CIDR format string (e.g. 10.1.2.0/24).

* `subnetwork_range_name` - (Optional) The subnetwork secondary range name specifying
    the secondary range from which to allocate the IP CIDR range for this alias IP
    range. If left unspecified, the primary range of the subnetwork will be used.

<a name="nested_service_account"></a>The `service_account` block supports:

* `email` - (Optional) The service account e-mail address. If not given, the
    default Google Compute Engine service account is used.
    **Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

* `scopes` - (Required) A list of service scopes. Both OAuth2 URLs and gcloud
    short names are supported. To allow full access to all Cloud APIs, use the
    `cloud-platform` scope. See a complete list of scopes [here](https://cloud.google.com/sdk/gcloud/reference/alpha/compute/instances/set-scopes#--scopes).
    **Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

<a name="nested_scheduling"></a>The `scheduling` block supports:

* `preemptible` - (Optional) Specifies if the instance is preemptible.
    If this field is set to true, then `automatic_restart` must be
    set to false.  Defaults to false.

* `on_host_maintenance` - (Optional) Describes maintenance behavior for the
    instance. Can be MIGRATE or TERMINATE, for more info, read
    [here](https://cloud.google.com/compute/docs/instances/setting-instance-scheduling-options).

* `automatic_restart` - (Optional) Specifies if the instance should be
    restarted if it was terminated by Compute Engine (not a user).
    Defaults to true.

* `node_affinities` - (Optional) Specifies node affinities or anti-affinities
   to determine which sole-tenant nodes your instances and managed instance
   groups will use as host systems. Read more on sole-tenant node creation
   [here](https://cloud.google.com/compute/docs/nodes/create-nodes).
   Structure [documented below](#nested_node_affinities).

* `min_node_cpus` - (Optional) The minimum number of virtual CPUs this instance will consume when running on a sole-tenant node.

* `provisioning_model` - (Optional) Describe the type of preemptible VM. This field accepts the value `STANDARD` or `SPOT`. If the value is `STANDARD`, there will be no discount. If this   is set to `SPOT`, 
    `preemptible` should be `true` and `auto_restart` should be
    `false`. For more info about
    `SPOT`, read [here](https://cloud.google.com/compute/docs/instances/spot)
    
* `instance_termination_action` - (Optional) Describe the type of termination action for `SPOT` VM. Can be `STOP` or `DELETE`.  Read more on [here](https://cloud.google.com/compute/docs/instances/create-use-spot) 

<a name="nested_guest_accelerator"></a>The `guest_accelerator` block supports:

* `type` (Required) - The accelerator type resource to expose to this instance. E.g. `nvidia-tesla-k80`.

* `count` (Required) - The number of the guest accelerator cards exposed to this instance.

<a name="nested_node_affinities"></a>The `node_affinities` block supports:

* `key` (Required) - The key for the node affinity label.

* `operator` (Required) - The operator. Can be `IN` for node-affinities
    or `NOT_IN` for anti-affinities.

* `values` (Required) - The values for the node affinity label.

<a name="nested_shielded_instance_config"></a>The `shielded_instance_config` block supports:

* `enable_secure_boot` (Optional) -- Verify the digital signature of all boot components, and halt the boot process if signature verification fails. Defaults to false.
  **Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

* `enable_vtpm` (Optional) -- Use a virtualized trusted platform module, which is a specialized computer chip you can use to encrypt objects like keys and certificates. Defaults to true.
  **Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

* `enable_integrity_monitoring` (Optional) -- Compare the most recent boot measurements to the integrity policy baseline and return a pair of pass/fail results depending on whether they match or not. Defaults to true.
  **Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true or your instance must have a `desired_status` of `TERMINATED` in order to update this field.

<a name="nested_confidential_instance_config"></a>The `confidential_instance_config` block supports:

* `enable_confidential_compute` (Optional) Defines whether the instance should have confidential compute enabled. [`on_host_maintenance`](#on_host_maintenance) has to be set to TERMINATE or this will fail to create the VM.

<a name="nested_advanced_machine_features"></a>The `advanced_machine_features` block supports:

* `enable_nested_virtualization` (Optional) Defines whether the instance should have [nested virtualization](#on_host_maintenance)  enabled. Defaults to false.

* `threads_per_core` (Optional) he number of threads per physical core. To disable [simultaneous multithreading (SMT)](https://cloud.google.com/compute/docs/instances/disabling-smt) set this to 1.

* `visible_core_count` (Optional) The number of physical cores to expose to an instance. [visible cores info (VC)](https://cloud.google.com/compute/docs/instances/customize-visible-cores).

<a name="nested_reservation_affinity"></a>The `reservation_affinity` block supports:

* `type` - (Required) The type of reservation from which this instance can consume resources.

* `specific_reservation` - (Optional) Specifies the label selector for the reservation to use..
    Structure is [documented below](#nested_specific_reservation).

<a name="nested_specific_reservation"></a>The `specific_reservation` block supports:

* `key` - (Required) Corresponds to the label key of a reservation resource. To target a SPECIFIC_RESERVATION by name, specify compute.googleapis.com/reservation-name as the key and specify the name of your reservation as the only value.

* `values` - (Required) Corresponds to the label values of a reservation resource.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/zones/{{zone}}/instances/{{name}}`

* `instance_id` - The server-assigned unique identifier of this instance.

* `metadata_fingerprint` - The unique fingerprint of the metadata.

* `self_link` - The URI of the created resource.

* `tags_fingerprint` - The unique fingerprint of the tags.

* `label_fingerprint` - The unique fingerprint of the labels.

* `cpu_platform` - The CPU platform used by this instance.

* `ipv6_access_type` - One of EXTERNAL, INTERNAL to indicate whether the IP can be accessed from the Internet.
This field is always inherited from its subnetwork.

* `network_interface.0.network_ip` - The internal ip address of the instance, either manually or dynamically assigned.

* `network_interface.0.access_config.0.nat_ip` - If the instance has an access config, either the given external ip (in the `nat_ip` field) or the ephemeral (generated) ip (if you didn't provide one).

* `network_interface.0.ipv6_access_config.0.external_ipv6` - The first IPv6 address of the external IPv6 range
associated with this instance, prefix length is stored in externalIpv6PrefixLength in ipv6AccessConfig.
The field is output only, an IPv6 address from a subnetwork associated with the instance will be allocated dynamically.

* `network_interface.0.ipv6_access_config.0.external_ipv6_prefix_length` - The prefix length of the external IPv6 range.

* `attached_disk.0.disk_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption) that protects this resource.

* `boot_disk.disk_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption) that protects this resource.

* `disk.0.disk_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption) that protects this resource.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

~> **Note:** The fields `boot_disk.0.disk_encryption_raw` and `attached_disk.*.disk_encryption_key_raw` cannot be imported automatically. The API doesn't return this information. If you are setting one of these fields in your config, you will need to update your state manually after importing the resource.

-> **Note:** The `desired_status` field will not be set on import. If you have it set, Terraform will update the field on the next `terraform apply`, bringing your instance to the desired status.


Instances can be imported using any of these accepted formats:

```
$ terraform import google_compute_instance.default projects/{{project}}/zones/{{zone}}/instances/{{name}}
$ terraform import google_compute_instance.default {{project}}/{{zone}}/{{name}}
$ terraform import google_compute_instance.default {{name}}
```

[custom-vm-types]: https://cloud.google.com/dataproc/docs/concepts/compute/custom-machine-types
[network-tier]: https://cloud.google.com/network-tiers/docs/overview
[extended-custom-vm-type]: https://cloud.google.com/compute/docs/instances/creating-instance-with-custom-machine-type#extendedmemory
