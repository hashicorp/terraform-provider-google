---
subcategory: "Cloud Platform"
page_title: "Google: google_service_networking_peered_dns_domain"
description: |-
 Allows management of a single peered DNS domain on a project.
---

# google\_project\_service\_peered\_dns\_domain

Allows management of a single peered DNS domain for an existing Google Cloud Platform project.

When using Google Cloud DNS to manage internal DNS, create peered DNS domains to make your DNS available to services like Google Cloud Build.

## Example Usage

```hcl
resource "google_service_networking_peered_dns_domain" "name" {
  project    = 10000000
  name       = "example-com"
  network    = "default"
  dns_suffix = "example.com."
  service    = "peering-service"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The producer project number. If not provided, the provider project is used.

* `name` - (Required) Internal name used for the peered DNS domain.

* `network` - (Required) The network in the consumer project.

* `dns_suffix` - (Required) The DNS domain suffix of the peered DNS domain. Make sure to suffix with a `.` (dot).

* `service` - (Optional) Private service connection between service and consumer network, defaults to `servicenetworking.googleapis.com`

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `services/{{service}}/projects/{{project}}/global/networks/{{network}}/peeredDnsDomains/{{name}}`

* `parent` - an identifier for the resource with format `services/{{service}}/projects/{{project}}/global/networks/{{network}}`

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `read`   - Default is 10 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Project peered DNS domains can be imported using the `service`, `project`, `network` and `name`, e.g.

```
$ terraform import google_service_networking_peered_dns_domain.my_domain services/{service}/projects/{project}/global/networks/{network}/peeredDnsDomains/{name}
```

Where:

- `service` is the service connection, defaults to `servicenetworking.googleapis.com`.
- `project` is the producer project name.
- `network` is the consumer network name.
- `name` is the name of your peered DNS domain.

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
