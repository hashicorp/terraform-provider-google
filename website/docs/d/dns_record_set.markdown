---
subcategory: "Cloud DNS"
description: |-
  Get a DNS record set within Google Cloud DNS
---

# google_dns_record_set

Get a DNS record set within Google Cloud DNS
For more information see
[the official documentation](https://cloud.google.com/dns/docs/records)
and
[API](https://cloud.google.com/dns/docs/reference/v1/resourceRecordSets)

## Example Usage

```tf
data "google_dns_managed_zone" "sample" {
  name = "sample-zone"
}

data "google_dns_record_set" "rs" {
  managed_zone = data.google_dns_managed_zone.sample.name
  name = "my-record.${data.google_dns_managed_zone.sample.dns_name}"
  type = "A"
}
```

## Argument Reference

The following arguments are supported:

* `managed_zone` - (Required) The Name of the zone.

* `name` - (Required) The DNS name for the resource.

* `type` - (Required) The RRSet type. [See this table for supported types](https://cloud.google.com/dns/docs/records#record_type).

* `project` - (Optional) The ID of the project for the Google Cloud.

## Attributes Reference

The following attributes are exported:

* `rrdatas` - The string data for the records in this record set.

* `ttl` - The time-to-live of this record set (seconds).
