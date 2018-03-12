---
layout: "google"
page_title: "Google: google_compute_address"
sidebar_current: "docs-google-compute-address"
description: |-
  Creates a static IP address resource for Google Compute Engine.
---

# google\_compute\_address

Creates a static IP address resource for Google Compute Engine. For more information see
the official documentation for
[external](https://cloud.google.com/compute/docs/instances-and-network) and
[internal](https://cloud.google.com/compute/docs/ip-addresses/reserve-static-internal-ip-address)
static IP reservations, as well as the
[API](https://cloud.google.com/compute/docs/reference/beta/addresses/insert).


## Example Usage

```hcl
resource "google_compute_address" "default" {
  name = "test-address"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The Region in which the created address should reside.
    If it is not provided, the provider region is used.

* `address_type` - (Optional) The Address Type that should be configured.
    Specify INTERNAL to reserve an internal static IP address EXTERNAL to
    specify an external static IP address. Defaults to EXTERNAL if omitted.

* `subnetwork` - (Optional) The self link URI of the subnetwork in which to
    create the address. A subnetwork may only be specified for INTERNAL
    address types.

* `address` - (Optional) The IP address to reserve. An address may only be
    specified for INTERNAL address types. The IP address must be inside the
    specified subnetwork, if any.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.
* `address` - The IP of the created resource.

## Import

Addresses can be imported using the `project`, `region` and `name`, e.g.

```
$ terraform import google_compute_address.default gcp-project/us-central1/test-address
```

If `project` is omitted, the default project set for the provider is used:

```
$ terraform import google_compute_address.default us-central1/test-address
```

If `project` and `region` are omitted, the default project and region set for the provider are used.

```
$ terraform import google_compute_address.default test-address
```

Alternatively, addresses can be imported using a full or partial `self_link`.

```
$ terraform import google_compute_address.default https://www.googleapis.com/compute/v1/projects/gcp-project/regions/us-central1/addresses/test-address

$ terraform import google_compute_address.default projects/gcp-project/regions/us-central1/addresses/test-address
```
