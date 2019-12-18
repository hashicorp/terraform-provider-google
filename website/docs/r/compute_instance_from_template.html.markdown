---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_instance_from_template"
sidebar_current: "docs-google-compute-instance-from-template"
description: |-
  Manages a VM instance resource within GCE.
---

# google\_compute\_instance\_from\_template

Manages a VM instance resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/instances)
and
[API](https://cloud.google.com/compute/docs/reference/latest/instances).

This resource is specifically to create a compute instance from a given
`source_instance_template`. To create an instance without a template, use the
`google_compute_instance` resource.


## Example Usage

```hcl
resource "google_compute_instance_template" "tpl" {
  name         = "template"
  machine_type = "n1-standard-1"

  disk {
    source_image = "debian-cloud/debian-9"
    auto_delete  = true
    disk_size_gb = 100
    boot         = true
  }

  network_interface {
    network = "default"
  }

  metadata = {
    foo = "bar"
  }

  can_ip_forward = true
}

resource "google_compute_instance_from_template" "tpl" {
  name = "instance-from-template"
  zone = "us-central1-a"

  source_instance_template = google_compute_instance_template.tpl.self_link

  // Override fields from instance template
  can_ip_forward = false
  labels = {
    my_key = "my_value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

* `source_instance_template` - (Required) Name or self link of an instance
  template to create the instance based on.

- - -

* `zone` - (Optional) The zone that the machine should be created in. If not
  set, the provider zone is used.

In addition to these, all arguments from `google_compute_instance` are supported
as a way to override the properties in the template. All exported attributes
from `google_compute_instance` are likewise exported here.

To support removal of Optional/Computed fields in Terraform 0.12 the following fields
are marked [Attributes as Blocks](/docs/configuration/attr-as-blocks.html):
* `attached_disk`
* `guest_accelerator`
* `service_account`
* `scratch_disk`
* `network_interface.alias_ip_range`
* `network_interface.access_config`

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 6 minutes.
- `update` - Default is 6 minutes.
- `delete` - Default is 6 minutes.

