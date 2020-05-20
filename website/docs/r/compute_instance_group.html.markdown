---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_instance_group"
sidebar_current: "docs-google-compute-instance-group-x"
description: |-
  Manages an Instance Group within GCE.
---

# google\_compute\_instance\_group

Creates a group of dissimilar Compute Engine virtual machine instances.
For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups/#unmanaged_instance_groups)
and [API](https://cloud.google.com/compute/docs/reference/latest/instanceGroups)

-> Recreating an instance group that's in use by another resource will give a
`resourceInUseByAnotherResource` error. You can avoid this error with a
Terraform `lifecycle` block as outlined in the example below.

## Example Usage - Empty instance group

```hcl
resource "google_compute_instance_group" "test" {
  name        = "terraform-test"
  description = "Terraform test instance group"
  zone        = "us-central1-a"
  network     = google_compute_network.default.id
}
```

### Example Usage - With instances and named ports

```hcl
resource "google_compute_instance_group" "webservers" {
  name        = "terraform-webservers"
  description = "Terraform test instance group"

  instances = [
    google_compute_instance.test.id,
    google_compute_instance.test2.id,
  ]

  named_port {
    name = "http"
    port = "8080"
  }

  named_port {
    name = "https"
    port = "8443"
  }

  zone = "us-central1-a"
}
```

### Example Usage - Recreating an instance group in use
Recreating an instance group that's in use by another resource will give a
`resourceInUseByAnotherResource` error. Use `lifecycle.create_before_destroy`
as shown in this example to avoid this type of error.

```hcl
resource "google_compute_instance_group" "staging_group" {
  name      = "staging-instance-group"
  zone      = "us-central1-c"
  instances = [google_compute_instance.staging_vm.id]
  named_port {
    name = "http"
    port = "8080"
  }

  named_port {
    name = "https"
    port = "8443"
  }

  lifecycle {
    create_before_destroy = true
  }
}

data "google_compute_image" "debian_image" {
  family  = "debian-9"
  project = "debian-cloud"
}

resource "google_compute_instance" "staging_vm" {
  name         = "staging-vm"
  machine_type = "n1-standard-1"
  zone         = "us-central1-c"
  boot_disk {
    initialize_params {
      image = data.google_compute_image.debian_image.self_link
    }
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_backend_service" "staging_service" {
  name      = "staging-service"
  port_name = "https"
  protocol  = "HTTPS"

  backend {
    group = google_compute_instance_group.staging_group.id
  }

  health_checks = [
    google_compute_https_health_check.staging_health.id,
  ]
}

resource "google_compute_https_health_check" "staging_health" {
  name         = "staging-health"
  request_path = "/health_check"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the instance group. Must be 1-63
    characters long and comply with
    [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Supported characters
    include lowercase letters, numbers, and hyphens.

* `zone` - (Required) The zone that this instance group should be created in.

- - -

* `description` - (Optional) An optional textual description of the instance
    group.

* `instances` - (Optional) List of instances in the group. They should be given
    as self_link URLs. When adding instances they must all be in the same
    network and zone as the instance group.

* `named_port` - (Optional) The named port configuration. See the section below
    for details on configuration.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `network` - (Optional) The URL of the network the instance group is in. If
    this is different from the network where the instances are in, the creation
    fails. Defaults to the network where the instances are in (if neither
    `network` nor `instances` is specified, this field will be blank).

The `named_port` block supports:

* `name` - (Required) The name which the port will be mapped to.

* `port` - (Required) The port number to map the name to.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}/zones/{{zone}}/instanceGroups/{{name}}`

* `self_link` - The URI of the created resource.

* `size` - The number of instances in the group.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is `6 minutes`
- `update` - Default is `6 minutes`
- `delete` - Default is `6 minutes`

## Import

Instance group can be imported using the `zone` and `name` with an optional `project`, e.g.

```
$ terraform import google_compute_instance_group.webservers us-central1-a/terraform-webservers
$ terraform import google_compute_instance_group.webservers big-project/us-central1-a/terraform-webservers
$ terraform import google_compute_instance_group.webservers projects/big-project/zones/us-central1-a/instanceGroups/terraform-webservers
```
