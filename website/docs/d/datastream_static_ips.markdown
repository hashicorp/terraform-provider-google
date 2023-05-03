---
subcategory: "Datastream"
description: |-
  Returns the list of IP addresses that Datastream connects from.
---

# google\_datastream\_static\_ips

Returns the list of IP addresses that Datastream connects from. For more information see
the [official documentation](https://cloud.google.com/datastream/docs/ip-allowlists-and-regions).

## Example Usage

```hcl
data "google_datastream_static_ips" "datastream_ips" {
  location       = "us-west1"
  project        = "my-project"
}


output "ip_list" {
  value = data.google_datastream_static_ips.datastream_ips.static_ips
}
```

## Argument Reference

The following arguments are supported:

* `location` - (required) The location to list Datastream IPs for. For example: `us-east1`.

* `project` (Optional) - Project from which to list static IP addresses. Defaults to project declared in the provider.

## Attributes Reference

The following attributes are exported:

* `static_ips` - A list of static IP addresses that Datastream will connect from.
