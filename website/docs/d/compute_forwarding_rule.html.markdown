---
subcategory: "Compute Engine"
page_title: "Google: google_compute_forwarding_rule"
description: |-
  Get a regional forwarding rule within GCE.
---

# google\_compute\_forwarding\_rule

Get a forwarding rule within GCE from its name.

## Example Usage

```tf
data "google_compute_forwarding_rule" "my-forwarding-rule" {
  name = "forwarding-rule-us-east1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the forwarding rule.


- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `region` - (Optional) The region in which the resource belongs. If it
    is not provided, the project region is used.

## Attributes Reference
See [google_compute_forwarding_rule](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_forwarding_rule) resource for details of the available attributes.
