---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_instance_template"
sidebar_current: "docs-google-compute-instance-template"
description: |-
  Manages a VM instance template resource within GCE.
---

# google\_compute\_instance\_template

Manages a VM instance template resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/instance-templates)
and
[API](https://cloud.google.com/compute/docs/reference/latest/instanceTemplates).


## Example Usage

```hcl
resource "google_compute_instance_template" "default" {
  name        = "appserver-template"
  description = "This template is used to create app server instances."

  tags = ["foo", "bar"]

  labels = {
    environment = "dev"
  }

  instance_description = "description assigned to instances"
  machine_type         = "n1-standard-1"
  can_ip_forward       = false

  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
  }

  // Create a new boot disk from an image
  disk {
    source_image = "debian-cloud/debian-9"
    auto_delete  = true
    boot         = true
  }

  // Use an existing disk resource
  disk {
    // Instance Templates reference disks by name, not self link
    source      = google_compute_disk.foobar.name
    auto_delete = false
    boot        = false
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
}

data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_disk" "foobar" {
  name  = "existing-disk"
  image = data.google_compute_image.my_image.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}
```

## Using with Instance Group Manager

Instance Templates cannot be updated after creation with the Google
Cloud Platform API. In order to update an Instance Template, Terraform will
destroy the existing resource and create a replacement. In order to effectively
use an Instance Template resource with an [Instance Group Manager resource][1],
it's recommended to specify `create_before_destroy` in a [lifecycle][2] block.
Either omit the Instance Template `name` attribute, or specify a partial name
with `name_prefix`.  Example:

```hcl
resource "google_compute_instance_template" "instance_template" {
  name_prefix  = "instance-template-"
  machine_type = "n1-standard-1"
  region       = "us-central1"

  // boot disk
  disk {
    # ...
  }

  // networking
  network_interface {
    # ...
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_instance_group_manager" "instance_group_manager" {
  name               = "instance-group-manager"
  instance_template  = google_compute_instance_template.instance_template.id
  base_instance_name = "instance-group-manager"
  zone               = "us-central1-f"
  target_size        = "1"
}
```

With this setup Terraform generates a unique name for your Instance
Template and can then update the Instance Group manager without conflict before
destroying the previous Instance Template.

## Deploying the Latest Image

A common way to use instance templates and managed instance groups is to deploy the
latest image in a family, usually the latest build of your application. There are two
ways to do this in Terraform, and they have their pros and cons. The difference ends
up being in how "latest" is interpreted. You can either deploy the latest image available
when Terraform runs, or you can have each instance check what the latest image is when
it's being created, either as part of a scaling event or being rebuilt by the instance
group manager.

If you're not sure, we recommend deploying the latest image available when Terraform runs,
because this means all the instances in your group will be based on the same image, always,
and means that no upgrades or changes to your instances happen outside of a `terraform apply`.
You can achieve this by using the [`google_compute_image`](../d/compute_image.html)
data source, which will retrieve the latest image on every `terraform apply`, and will update
the template to use that specific image:

```tf
data "google_compute_image" "my_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance_template" "instance_template" {
  name_prefix  = "instance-template-"
  machine_type = "n1-standard-1"
  region       = "us-central1"

  // boot disk
  disk {
    source_image = google_compute_image.my_image.self_link
  }
}
```

To have instances update to the latest on every scaling event or instance re-creation,
use the family as the image for the disk, and it will use GCP's default behavior, setting
the image for the template to the family:

```tf
resource "google_compute_instance_template" "instance_template" {
  name_prefix  = "instance-template-"
  machine_type = "n1-standard-1"
  region       = "us-central1"

  // boot disk
  disk {
    source_image = "debian-cloud/debian-9"
  }
}
```

## Argument Reference

Note that changing any field for this resource forces a new resource to be created.

The following arguments are supported:

* `disk` - (Required) Disks to attach to instances created from this template.
    This can be specified multiple times for multiple disks. Structure is
    documented below.

* `machine_type` - (Required) The machine type to create.

    To create a machine with a [custom type][custom-vm-types] (such as extended memory), format the value like `custom-VCPUS-MEM_IN_MB` like `custom-6-20480` for 6 vCPU and 20GB of RAM.

- - -
* `name` - (Optional) The name of the instance template. If you leave
  this blank, Terraform will auto-generate a unique name.

* `name_prefix` - (Optional) Creates a unique name beginning with the specified
  prefix. Conflicts with `name`.

* `can_ip_forward` - (Optional) Whether to allow sending and receiving of
    packets with non-matching source or destination IPs. This defaults to false.

* `description` - (Optional) A brief description of this resource.

* `instance_description` - (Optional) A brief description to use for instances
    created from this template.

* `labels` - (Optional) A set of key/value label pairs to assign to instances
    created from this template,

* `metadata` - (Optional) Metadata key/value pairs to make available from
    within instances created from this template.

* `metadata_startup_script` - (Optional) An alternative to using the
    startup-script metadata key, mostly to match the compute_instance resource.
    This replaces the startup-script metadata key on the created instance and
    thus the two mechanisms are not allowed to be used simultaneously.

* `network_interface` - (Required) Networks to attach to instances created from
    this template. This can be specified multiple times for multiple networks.
    Structure is documented below.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) An instance template is a global resource that is not
    bound to a zone or a region. However, you can still specify some regional
    resources in an instance template, which restricts the template to the
    region where that resource resides. For example, a custom `subnetwork`
    resource is tied to a specific region. Defaults to the region of the
    Provider if no value is given.

* `scheduling` - (Optional) The scheduling strategy to use. More details about
    this configuration option are detailed below.

* `service_account` - (Optional) Service account to attach to the instance. Structure is documented below.

* `tags` - (Optional) Tags to attach to the instance.

* `guest_accelerator` - (Optional) List of the type and count of accelerator cards attached to the instance. Structure documented below.

* `min_cpu_platform` - (Optional) Specifies a minimum CPU platform. Applicable values are the friendly names of CPU platforms, such as
`Intel Haswell` or `Intel Skylake`. See the complete list [here](https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform).

* `shielded_instance_config` - (Optional) Enable [Shielded VM](https://cloud.google.com/security/shielded-cloud/shielded-vm) on this instance. Shielded VM provides verifiable integrity to prevent against malware and rootkits. Defaults to disabled. Structure is documented below.
	**Note**: [`shielded_instance_config`](#shielded_instance_config) can only be used with boot images with shielded vm support. See the complete list [here](https://cloud.google.com/compute/docs/images#shielded-images).

* `enable_display` - (Optional) Enable [Virtual Displays](https://cloud.google.com/compute/docs/instances/enable-instance-virtual-display#verify_display_driver) on this instance.
**Note**: [`allow_stopping_for_update`](#allow_stopping_for_update) must be set to true in order to update this field.

The `disk` block supports:

* `auto_delete` - (Optional) Whether or not the disk should be auto-deleted.
    This defaults to true.

* `boot` - (Optional) Indicates that this is a boot disk.

* `device_name` - (Optional) A unique device name that is reflected into the
    /dev/  tree of a Linux operating system running within the instance. If not
    specified, the server chooses a default device name to apply to this disk.

* `disk_name` - (Optional) Name of the disk. When not provided, this defaults
    to the name of the instance.

* `source_image` - (Optional) The image from which to
    initialize this disk. This can be one of: the image's `self_link`,
    `projects/{project}/global/images/{image}`,
    `projects/{project}/global/images/family/{family}`, `global/images/{image}`,
    `global/images/family/{family}`, `family/{family}`, `{project}/{family}`,
    `{project}/{image}`, `{family}`, or `{image}`.
~> **Note:** Either `source` or `source_image` is **required** when creating a new instance except for when creating a local SSD. Check the API [docs](https://cloud.google.com/compute/docs/reference/rest/v1/instanceTemplates/insert) for details.

* `interface` - (Optional) Specifies the disk interface to use for attaching
    this disk.

* `mode` - (Optional) The mode in which to attach this disk, either READ_WRITE
    or READ_ONLY. If you are attaching or creating a boot disk, this must
    read-write mode.

* `source` - (Optional) The name (**not self_link**)
    of the disk (such as those managed by `google_compute_disk`) to attach.
~> **Note:** Either `source` or `source_image` is **required** when creating a new instance except for when creating a local SSD. Check the API [docs](https://cloud.google.com/compute/docs/reference/rest/v1/instanceTemplates/insert) for details.

* `disk_type` - (Optional) The GCE disk type. Can be either `"pd-ssd"`,
    `"local-ssd"`, or `"pd-standard"`.

* `disk_size_gb` - (Optional) The size of the image in gigabytes. If not
    specified, it will inherit the size of its base image. For SCRATCH disks,
    the size must be exactly 375GB.

* `type` - (Optional) The type of GCE disk, can be either `"SCRATCH"` or
    `"PERSISTENT"`.

* `disk_encryption_key` - (Optional) Encrypts or decrypts a disk using a customer-supplied encryption key.

    If you are creating a new disk, this field encrypts the new disk using an encryption key that you provide. If you are attaching an existing disk that is already encrypted, this field decrypts the disk using the customer-supplied encryption key.

    If you encrypt a disk using a customer-supplied key, you must provide the same key again when you attempt to use this resource at a later time. For example, you must provide the key when you create a snapshot or an image from the disk or when you attach the disk to a virtual machine instance.

    If you do not provide an encryption key, then the disk will be encrypted using an automatically generated key and you do not need to provide a key to use the disk later.

    Instance templates do not store customer-supplied encryption keys, so you cannot use your own keys to encrypt disks in a managed instance group.

The `disk_encryption_key` block supports:

* `kms_key_self_link` - (Required) The self link of the encryption key that is stored in Google Cloud KMS

The `network_interface` block supports:

* `network` - (Optional) The name or self_link of the network to attach this interface to.
    Use `network` attribute for Legacy or Auto subnetted networks and
    `subnetwork` for custom subnetted networks.

* `subnetwork` - (Optional) the name of the subnetwork to attach this interface
    to. The subnetwork must exist in the same `region` this instance will be
    created in. Either `network` or `subnetwork` must be provided.

* `subnetwork_project` - (Optional) The ID of the project in which the subnetwork belongs.
    If it is not provided, the provider project is used.

* `network_ip` - (Optional) The private IP address to assign to the instance. If
    empty, the address will be automatically assigned.

* `access_config` - (Optional) Access configurations, i.e. IPs via which this
    instance can be accessed via the Internet. Omit to ensure that the instance
    is not accessible from the Internet (this means that ssh provisioners will
    not work unless you are running Terraform can send traffic to the instance's
    network (e.g. via tunnel or because it is running on another cloud instance
    on that network). This block can be repeated multiple times. Structure documented below.

* `alias_ip_range` - (Optional) An
    array of alias IP ranges for this network interface. Can only be specified for network
    interfaces on subnet-mode networks. Structure documented below.

The `access_config` block supports:

* `nat_ip` - (Optional) The IP address that will be 1:1 mapped to the instance's
    network ip. If not given, one will be generated.

* `network_tier` - (Optional) The [networking tier][network-tier] used for configuring
    this instance template. This field can take the following values: PREMIUM or
    STANDARD. If this field is not specified, it is assumed to be PREMIUM.

The `alias_ip_range` block supports:

* `ip_cidr_range` - The IP CIDR range represented by this alias IP range. This IP CIDR range
    must belong to the specified subnetwork and cannot contain IP addresses reserved by
    system or used by other network interfaces. At the time of writing only a
    netmask (e.g. /24) may be supplied, with a CIDR format resulting in an API
    error.

* `subnetwork_range_name` - (Optional) The subnetwork secondary range name specifying
    the secondary range from which to allocate the IP CIDR range for this alias IP
    range. If left unspecified, the primary range of the subnetwork will be used.

The `service_account` block supports:

* `email` - (Optional) The service account e-mail address. If not given, the
    default Google Compute Engine service account is used.

* `scopes` - (Required) A list of service scopes. Both OAuth2 URLs and gcloud
    short names are supported. To allow full access to all Cloud APIs, use the
    `cloud-platform` scope. See a complete list of scopes [here](https://cloud.google.com/sdk/gcloud/reference/alpha/compute/instances/set-scopes#--scopes).

    The [service accounts documentation](https://cloud.google.com/compute/docs/access/service-accounts#accesscopesiam)
    explains that access scopes are the legacy method of specifying permissions for your instance.
    If you are following best practices and using IAM roles to grant permissions to service accounts,
    then you can define this field as an empty list.

The `scheduling` block supports:

* `automatic_restart` - (Optional) Specifies whether the instance should be
    automatically restarted if it is terminated by Compute Engine (not
    terminated by a user). This defaults to true.

* `on_host_maintenance` - (Optional) Defines the maintenance behavior for this
    instance.

* `preemptible` - (Optional) Allows instance to be preempted. This defaults to
    false. Read more on this
    [here](https://cloud.google.com/compute/docs/instances/preemptible).

* `node_affinities` - (Optional) Specifies node affinities or anti-affinities
   to determine which sole-tenant nodes your instances and managed instance
   groups will use as host systems. Read more on sole-tenant node creation
   [here](https://cloud.google.com/compute/docs/nodes/create-nodes).
   Structure documented below.

The `guest_accelerator` block supports:

* `type` (Required) - The accelerator type resource to expose to this instance. E.g. `nvidia-tesla-k80`.

* `count` (Required) - The number of the guest accelerator cards exposed to this instance.

The `node_affinities` block supports:

* `key` (Required) - The key for the node affinity label.

* `operator` (Required) - The operator. Can be `IN` for node-affinities
    or `NOT_IN` for anti-affinities.

* `value` (Required) - The values for the node affinity label.

The `shielded_instance_config` block supports:

* `enable_secure_boot` (Optional) -- Verify the digital signature of all boot components, and halt the boot process if signature verification fails. Defaults to false.

* `enable_vtpm` (Optional) -- Use a virtualized trusted platform module, which is a specialized computer chip you can use to encrypt objects like keys and certificates. Defaults to true.

* `enable_integrity_monitoring` (Optional) -- Compare the most recent boot measurements to the integrity policy baseline and return a pair of pass/fail results depending on whether they match or not. Defaults to true.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/global/instanceTemplates/{{name}}`

* `metadata_fingerprint` - The unique fingerprint of the metadata.

* `self_link` - The URI of the created resource.

* `tags_fingerprint` - The unique fingerprint of the tags.

[1]: /docs/providers/google/r/compute_instance_group_manager.html
[2]: /docs/configuration/resources.html#lifecycle

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import

Instance templates can be imported using any of these accepted formats:

```
$ terraform import google_compute_instance_template.default projects/{{project}}/global/instanceTemplates/{{name}}
$ terraform import google_compute_instance_template.default {{project}}/{{name}}
$ terraform import google_compute_instance_template.default {{name}}
```

[custom-vm-types]: https://cloud.google.com/dataproc/docs/concepts/compute/custom-machine-types
[network-tier]: https://cloud.google.com/network-tiers/docs/overview
