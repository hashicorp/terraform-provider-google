---
subcategory: "Compute Engine"
page_title: "Google: google_compute_global_forwarding_rule"
description: |-
  Get a global forwarding rule within GCE.
---

# google\_compute\_global_\forwarding\_rule

Get a global forwarding rule within GCE from its name.

## Example Usage

```tf
data "google_compute_global_forwarding_rule" "my-forwarding-rule" {
  name = "forwarding-rule-global"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the global forwarding rule.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference
See [google_compute_global_forwarding_rule](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_global_forwarding_rule) resource for details of the available attributes.
