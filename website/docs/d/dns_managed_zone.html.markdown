---
subcategory: "Cloud DNS"
layout: "google"
page_title: "Google: google_dns_managed_zone"
sidebar_current: "docs-google-datasource-dns-managed-zone"
description: |-
  Provides access to the attributes of a zone within Google Cloud DNS
---

# google\_dns\_managed\_zone

Provides access to a zone's attributes within Google Cloud DNS.
For more information see
[the official documentation](https://cloud.google.com/dns/zones/)
and
[API](https://cloud.google.com/dns/api/v1/managedZones).

```hcl
data "google_dns_managed_zone" "env_dns_zone" {
  name = "qa-zone"
}

resource "google_dns_record_set" "dns" {
  name = "my-address.${data.google_dns_managed_zone.env_dns_zone.dns_name}"
  type = "TXT"
  ttl  = 300

  managed_zone = data.google_dns_managed_zone.env_dns_zone.name

  rrdatas = ["test"]
}
```

## Argument Reference

* `name` - (Required) A unique name for the resource.

* `project` - (Optional) The ID of the project for the Google Cloud DNS zone.

## Attributes Reference

The following attributes are exported:

* `dns_name` - The fully qualified DNS name of this zone, e.g. `terraform.io.`.

* `description` - A textual description field.

* `name_servers` - The list of nameservers that will be authoritative for this
    domain. Use NS records to redirect from your DNS provider to these names,
    thus making Google Cloud DNS authoritative for this zone.

* `visibility` - The zone's visibility: public zones are exposed to the Internet,
    while private zones are visible only to Virtual Private Cloud resources.
