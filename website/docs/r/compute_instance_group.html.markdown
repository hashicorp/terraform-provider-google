---
layout: "google"
page_title: "Google: google_compute_instance_group"
sidebar_current: "docs-google-compute-instance-group-x"
description: |-
  Manages an Instance Group within GCE.
---

# google\_compute\_instance\_group

The Google Compute Engine Instance Group API creates and manages pools
of homogeneous Compute Engine virtual machine instances from a common instance
template. For more information, see [the official documentation](https://cloud.google.com/compute/docs/instance-groups/#unmanaged_instance_groups)
and [API](https://cloud.google.com/compute/docs/reference/latest/instanceGroups)

## Example Usage

### Empty instance group

```hcl
resource "google_compute_instance_group" "test" {
  name        = "terraform-test"
  description = "Terraform test instance group"
  zone        = "us-central1-a"
  network     = "${google_compute_network.default.self_link}"
}
```

### With instances and named ports

```hcl
resource "google_compute_instance_group" "webservers" {
  name        = "terraform-webservers"
  description = "Terraform test instance group"

  instances = [
    "${google_compute_instance.test.self_link}",
    "${google_compute_instance.test2.self_link}",
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

* `project` - (Optional) The project in which the resource belongs. If it
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

* `self_link` - The URI of the created resource.

* `size` - The number of instances in the group.

## Import

Instance group can be imported using the `zone` and `name`, e.g.

```
$ terraform import google_compute_instance_group.webservers us-central1-a/terraform-webservers
```