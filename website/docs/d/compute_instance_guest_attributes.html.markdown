---
subcategory: "Compute Engine"
description: |-
  Get GCE instance's guest attributes
---

# google_compute_instance_guest_attributes

Get information about a VM instance resource within GCE. For more information see
[the official documentation](https://cloud.google.com/compute/docs/instances)
and
[API](https://cloud.google.com/compute/docs/reference/latest/instances).

Get information about VM's guest attrubutes. For more information see [the official documentation](https://cloud.google.com/compute/docs/metadata/manage-guest-attributes)
and
[API](https://cloud.google.com/compute/docs/reference/rest/v1/instances/getGuestAttributes).

## Example Usage - get all attributes from a single namespace

```hcl
data "google_compute_instance_guest_attributes" "appserver_ga" {
  name       = "primary-application-server"
  zone       = "us-central1-a"
  query_path = "variables/"
}
```

## Example Usage - get a specific variable

```hcl
data "google_compute_instance_guest_attributes" "appserver_ga" {
  name         = "primary-application-server"
  zone         = "us-central1-a"
  variable_key = "variables/key1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name or self_link of the instance.

---

* `project` - (Optional) The ID of the project in which the resource belongs.
    If `self_link` is provided, this value is ignored.  If neither `self_link`
    nor `project` are provided, the provider project is used.

* `zone` - (Optional) The zone of the instance. If `self_link` is provided, this
    value is ignored.  If neither `self_link` nor `zone` are provided, the
    provider zone is used.

* `query_path` - (Optional) Path to query for the guest attributes. Consists of
  `namespace` name for the attributes followed with a `/`.

* `variable_key` - (Optional) Key of a variable to get the value of. Consists of
  `namespace` name and `key` name for the variable separated by a `/`.

## Attributes Reference

* `query_value` - Structure is [documented below](#nested_query_value).

* `variable_value` - Value of the queried guest_attribute.

---

<a name="nested_query_value"></a>The `query_value` block supports:

* `key` - Key of the guest_attribute.

* `namespace` - Namespace of the guest_attribute.

* `value` - Value of the guest_attribute.