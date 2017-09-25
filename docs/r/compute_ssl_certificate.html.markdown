---
layout: "google"
page_title: "Google: google_compute_ssl_certificate"
sidebar_current: "docs-google-compute-ssl-certificate"
description: |-
  Creates an SSL certificate resource necessary for HTTPS load balancing in GCE.
---

# google\_compute\_ssl\_certificate

Creates an SSL certificate resource necessary for HTTPS load balancing in GCE.
For more information see
[the official documentation](https://cloud.google.com/compute/docs/load-balancing/http/ssl-certificates) and
[API](https://cloud.google.com/compute/docs/reference/latest/sslCertificates).


## Example Usage

```hcl
resource "google_compute_ssl_certificate" "default" {
  name_prefix = "my-certificate-"
  description = "a description"
  private_key = "${file("path/to/private.key")}"
  certificate = "${file("path/to/certificate.crt")}"
}
```

## Using with Target HTTPS Proxies

SSL certificates cannot be updated after creation. In order to apply the
specified configuration, Terraform will destroy the existing resource
and create a replacement. To effectively use an SSL certificate resource
with a [Target HTTPS Proxy resource][1], it's recommended to specify
`create_before_destroy` in a [lifecycle][2] block. Either omit the
Instance Template `name` attribute, or specify a partial name with
`name_prefix`.  Example:

```hcl
resource "google_compute_ssl_certificate" "default" {
  name_prefix = "my-certificate-"
  description = "a description"
  private_key = "${file("path/to/private.key")}"
  certificate = "${file("path/to/certificate.crt")}"

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_target_https_proxy" "my_proxy" {
  name             = "public-proxy"
  url_map          = # ...
  ssl_certificates = ["${google_compute_ssl_certificate.default.self_link}"]
}
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Required) A local certificate file in PEM format. The chain
    may be at most 5 certs long, and must include at least one intermediate
    cert. Changing this forces a new resource to be created.

* `private_key` - (Required) Write only private key in PEM format.
    Changing this forces a new resource to be created.

- - -

* `name` - (Optional) A unique name for the SSL certificate. If you leave
  this blank, Terraform will auto-generate a unique name.

* `name_prefix` - (Optional) Creates a unique name beginning with the specified
  prefix. Conflicts with `name`.

* `description` - (Optional) An optional description of this resource.
    Changing this forces a new resource to be created.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `certificate_id` - A unique ID for the certificate, assigned by GCE.

* `self_link` - The URI of the created resource.

[1]: /docs/providers/google/r/compute_target_https_proxy.html
[2]: /docs/configuration/resources.html#lifecycle
