---
subcategory: "Cloud DNS"
description: |-
  Provides access to a list of zones within Google Cloud DNS
---

# google\_dns\_managed\_zones

Provides access to a list of zones within Google Cloud DNS.
For more information see
[the official documentation](https://cloud.google.com/dns/zones/)
and
[API](https://cloud.google.com/dns/api/v1/managedZones).

```hcl
data "google_dns_managed_zones" "zones" {
  project = "my-project-id"
}
```

## Argument Reference

* `project` - (Optional) The ID of the project containing Google Cloud DNS zones. If this is not provided the default project will be used.

## Attributes Reference

The following attributes are exported:

* `managed_zones` - A list of managed zones.

To see the attributes available for each zone in the list, see the singular [google_dns_managed_zone](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/dns_managed_zone#attributes-reference) data source for details of the available attributes.