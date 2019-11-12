---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_instance"
sidebar_current: "docs-google-datasource-compute-instance-x"
description: |-
  Get a VM instance within GCE.
---

# google\_compute\_instance

Get information about a VM instance resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/instances)
and
[API](https://cloud.google.com/compute/docs/reference/latest/instances).


## Example Usage

```hcl
data "google_compute_instance" "appserver" {
  name = "primary-application-server"
  zone = "us-central1-a"
}
```

## Argument Reference

The following arguments are supported:

* `self_link` - (Optional) The self link of the instance. One of `name` or `self_link` must be provided.

* `name` - (Optional) The name of the instance. One of `name` or `self_link` must be provided.

---

* `project` - (Optional) The ID of the project in which the resource belongs.
    If `self_link` is provided, this value is ignored.  If neither `self_link`
    nor `project` are provided, the provider project is used.

* `zone` - (Optional) The zone of the instance. If `self_link` is provided, this
    value is ignored.  If neither `self_link` nor `zone` are provided, the
    provider zone is used.

## Attributes Reference

* `boot_disk` - The boot disk for the instance. Structure is documented below.

* `machine_type` - The machine type to create.

* `network_interface` - The networks attached to the instance. Structure is documented below.

* `attached_disk` - List of disks attached to the instance. Structure is documented below.

* `can_ip_forward` - Whether sending and receiving of packets with non-matching source or destination IPs is allowed.

* `description` - A brief description of the resource.

* `deletion_protection` - Whether deletion protection is enabled on this instance.

* `guest_accelerator` - List of the type and count of accelerator cards attached to the instance. Structure is documented below.

* `labels` - A set of key/value label pairs assigned to the instance.

* `metadata` - Metadata key/value pairs made available within the instance.

* `min_cpu_platform` - The minimum CPU platform specified for the VM instance.

* `scheduling` - The scheduling strategy being used by the instance.

* `scratch_disk` - The scratch disks attached to the instance. Structure is documented below.

* `service_account` - The service account to attach to the instance. Structure is documented below.

* `tags` - The list of tags attached to the instance.

* `instance_id` - The server-assigned unique identifier of this instance.

* `metadata_fingerprint` - The unique fingerprint of the metadata.

* `self_link` - The URI of the created resource.

* `tags_fingerprint` - The unique fingerprint of the tags.

* `label_fingerprint` - The unique fingerprint of the labels.

* `cpu_platform` - The CPU platform used by this instance.

* `shielded_instance_config` - The shielded vm config being used by the instance. Structure is documented below.

* `enable_display` -- Whether the instance has virtual displays enabled.

* `network_interface.0.network_ip` - The internal ip address of the instance, either manually or dynamically assigned.

* `network_interface.0.access_config.0.nat_ip` - If the instance has an access config, either the given external ip (in the `nat_ip` field) or the ephemeral (generated) ip (if you didn't provide one).

* `attached_disk.0.disk_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption) that protects this resource.

* `boot_disk.disk_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption) that protects this resource.

* `disk.0.disk_encryption_key_sha256` - The [RFC 4648 base64](https://tools.ietf.org/html/rfc4648#section-4)
    encoded SHA-256 hash of the [customer-supplied encryption key]
    (https://cloud.google.com/compute/docs/disks/customer-supplied-encryption) that protects this resource.

---

The `boot_disk` block supports:

* `auto_delete` - Whether the disk will be auto-deleted when the instance is deleted.

* `device_name` - Name with which attached disk will be accessible under `/dev/disk/by-id/`

* `initialize_params` - Parameters with which a disk was created alongside the instance.
    Structure is documented below.

* `source` - The name or self_link of an existing disk (such as those managed by
    `google_compute_disk`) that was attached to the instance.

The `initialize_params` block supports:

* `size` - The size of the image in gigabytes.

* `type` - The GCE disk type. One of `pd-standard` or `pd-ssd`.

* `image` - The image from which this disk was initialised.

The `scratch_disk` block supports:

* `interface` - The disk interface used for attaching this disk. One of `SCSI` or `NVME`.

The `attached_disk` block supports:

* `source` - The name or self_link of the disk attached to this instance.

* `device_name` - Name with which the attached disk is accessible
    under `/dev/disk/by-id/`

* `mode` - Read/write mode for the disk. One of `"READ_ONLY"` or `"READ_WRITE"`.

The `network_interface` block supports:

* `network` - The name or self_link of the network attached to this interface.

*  `subnetwork` - The name or self_link of the subnetwork attached to this interface.

*  `subnetwork_project` - The project in which the subnetwork belongs.

* `network_ip` - The private IP address assigned to the instance.

* `access_config` - Access configurations, i.e. IPs via which this
    instance can be accessed via the Internet. Structure documented below.

* `alias_ip_range` - An array of alias IP ranges for this network interface. Structure documented below.

The `access_config` block supports:

* `nat_ip` - The IP address that is be 1:1 mapped to the instance's
    network ip.

* `public_ptr_domain_name` - The DNS domain name for the public PTR record.

* `network_tier` - The [networking tier][network-tier] used for configuring this instance. One of `PREMIUM` or `STANDARD`.

The `alias_ip_range` block supports:

* `ip_cidr_range` - The IP CIDR range represented by this alias IP range.

* `subnetwork_range_name` - The subnetwork secondary range name specifying
    the secondary range from which to allocate the IP CIDR range for this alias IP
    range.

The `service_account` block supports:

* `email` - The service account e-mail address.

* `scopes` - A list of service scopes.

The `scheduling` block supports:

* `preemptible` - Whether the instance is preemptible.

* `on_host_maintenance` - Describes maintenance behavior for the
    instance. One of `MIGRATE` or `TERMINATE`, for more info, read
    [here](https://cloud.google.com/compute/docs/instances/setting-instance-scheduling-options)

* `automatic_restart` - Specifies if the instance should be
    restarted if it was terminated by Compute Engine (not a user).

The `guest_accelerator` block supports:

* `type` - The accelerator type resource exposed to this instance. E.g. `nvidia-tesla-k80`.

* `count` - The number of the guest accelerator cards exposed to this instance.

[network-tier]: https://cloud.google.com/network-tiers/docs/overview

The `shielded_instance_config` block supports:

* `enable_secure_boot` -- Whether secure boot is enabled for the instance.

* `enable_vtpm` -- Whether the instance uses vTPM.

* `enable_integrity_monitoring` -- Whether integrity monitoring is enabled for the instance.
