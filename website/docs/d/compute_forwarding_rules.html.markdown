---
subcategory: "Compute Engine"
description: |-
  List forwarding rules in a region of a Google Cloud project.
---

# google_compute_forwarding_rules

List all networks in a specified Google Cloud project.

## Example Usage

```tf
data "google_compute_forwarding_rules" "my-forwarding-rules" {
  project = "my-cloud-project"
  region  = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The name of the project.

* `region`  - (Optional) The region you want to get the forwarding rules from.

These arguments must be set in either the provider or the resouce in order for the information to be queried.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format projects/{{project}}/region/{{region}}/forwardingRules

* `project` - The project name being queried.

* `region` - The region being queried.

* `rules` - This is a list of the forwarding rules in the project. Each forwarding rule will list the backend, description, ip address. name, network, self link, service label, service name, and subnet.

* `self_link` - The URI of the resource.
