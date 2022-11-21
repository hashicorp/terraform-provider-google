---
subcategory: "Compute Engine"
page_title: "Google: google_compute_instance_from_machine_image"
description: |-
  Manages a VM instance resource within GCE.
---

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

# google\_compute\_instance\_from\_machine\_image

Manages a VM instance resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/instances)
and
[API](https://cloud.google.com/compute/docs/reference/latest/instances).

This resource is specifically to create a compute instance from a given
`source_machine_image`. To create an instance without a machine image, use the
`google_compute_instance` resource.


## Example Usage

```hcl
resource "google_compute_instance_from_machine_image" "tpl" {
  provider = google-beta
  name     = "instance-from-machine-image"
  zone     = "us-central1-a"

  source_machine_image = "projects/PROJECT-ID/global/machineImages/NAME"

  // Override fields from machine image
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

* `source_machine_image` - (Required) Name or self link of a machine
  image to create the instance based on.

- - -

* `zone` - (Optional) The zone that the machine should be created in. If not
  set, the provider zone is used.

In addition to these, most* arguments from `google_compute_instance` are supported
as a way to override the properties in the machine image. All exported attributes
from `google_compute_instance` are likewise exported here.

~> **Warning:** *Due to API limitations, disk overrides are currently disabled. This includes the "boot_disk", "attached_disk", and "scratch_disk" fields.

## Attributes Reference

All exported attributes from `google_compute_instance` are exported here.
See https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance#attributes-reference
for details.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 6 minutes.
- `update` - Default is 6 minutes.
- `delete` - Default is 6 minutes.
