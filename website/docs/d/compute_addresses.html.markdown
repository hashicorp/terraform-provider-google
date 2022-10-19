---
subcategory: "Compute Engine"
page_title: "Google: google_compute_addresses"
description: |-
  List google compute addresses.
---

# google\_compute\_addresses

List IP addresses in a project. For more information see
the official API [list](https://cloud.google.com/compute/docs/reference/latest/addresses/list) and 
[aggregated lsit](https://cloud.google.com/compute/docs/reference/rest/v1/addresses/aggregatedList) documentation.

## Example Usage

```hcl
data "google_compute_addresses" "my_addresses" {
    filter = "name:test-*"
}

resource "google_dns_record_set" "frontend" {
  name = "frontend.${google_dns_managed_zone.prod.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = google_dns_managed_zone.prod.name

  rrdatas = data.google_compute_addresses.my_addresses[*].address
}

resource "google_dns_managed_zone" "prod" {
  name     = "prod-zone"
  dns_name = "prod.mydomain.com."
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The google project in which addresses are listed.
    Defaults to provider's configuration if missing.

* `region` - (Optional) Region that should be considered to search addresses.
    All regions are considered if missing.

* `filter` - (Optional) A filter expression that
    filters resources listed in the response. The expression must specify
    the field name, an operator, and the value that you want to use for
    filtering. The value must be a string, a number, or a boolean. The
    operator must be either "=", "!=", ">", "<", "<=", ">=" or ":". For
    example, if you are filtering Compute Engine instances, you can
    exclude instances named "example-instance" by specifying "name !=
    example-instance". The ":" operator can be used with string fields to
    match substrings. For non-string fields it is equivalent to the "="
    operator. The ":*" comparison can be used to test whether a key has
    been defined. For example, to find all objects with "owner" label
    use: """ labels.owner:* """ You can also filter nested fields. For
    example, you could specify "scheduling.automaticRestart = false" to
    include instances only if they are not scheduled for automatic
    restarts. You can use filtering on nested fields to filter based on
    resource labels. To filter on multiple expressions, provide each
    separate expression within parentheses. For example: """
    (scheduling.automaticRestart = true) (cpuPlatform = "Intel Skylake")
    """ By default, each expression is an "AND" expression. However, you
    can include "AND" and "OR" expressions explicitly. For example: """
    (cpuPlatform = "Intel Skylake") OR (cpuPlatform = "Intel Broadwell")
    AND (scheduling.automaticRestart = true)

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `addresses` - A list of addresses matching the filter. Structure is [defined below](#nested_addresses).

<a name="nested_addresses"></a>The `addresses` block supports:

* `name` - The IP address name.
* `address` - The IP address (for example `1.2.3.4`).
* `address_type` - The IP address type, can be `EXTERNAL` or `INTERNAL`.
* `description` - The IP address description.
* `status` - Indicates if the address is used. Possible values are: RESERVED or IN_USE.
* `labels` - (Beta only) A map containing IP labels.
* `region` - The region in which the address resides.
* `self_link` - The URI of the created resource.
