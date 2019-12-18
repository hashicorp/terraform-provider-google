---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_ssl_certificate"
sidebar_current: "docs-google-datasource-compute-ssl-certificate"
description: |-
  Get info about a Google Compute SSL Certificate.
---

# google\_compute\_ssl\_certificate

Get info about a Google Compute SSL Certificate from its name.

## Example Usage

```tf
data "google_compute_ssl_certificate" "my_cert" {
  name = "my-cert"
}

output "certificate" {
  value = data.google_compute_ssl_certificate.my_cert.certificate
}

output "certificate_id" {
  value = data.google_compute_ssl_certificate.my_cert.certificate_id
}

output "self_link" {
  value = data.google_compute_ssl_certificate.my_cert.self_link
}
```

## Argument Reference

The following arguments are supported:

* `name` (Required) - The name of the certificate.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_compute_ssl_certificate](https://www.terraform.io/docs/providers/google/r/compute_ssl_certificate.html) resource for details of the available attributes.
