---
subcategory: "Certificate Authority Service"
layout: "google"
page_title: "Google: google_privateca_certificate_authority"
sidebar_current: "docs-google-datasource-privateca-certificate-authority"
description: |-
  Contains the data that describes a Certificate Authority
---
# google_privateca_certificate_authority

Get info about a Google Cloud IAP Client.

## Example Usage

```tf
data "google_privateca_certificate_authority" "default" {
  location = "us-west1"
  pool = "pool-name"
  certificate_authority_id = "ca-id"
}

output "csr" {
  value = data.google_privateca_certificate_authority.default.pem_csr
}

```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location the certificate authority exists in.

* `pool` - (Required) The name of the pool the certificate authority belongs to.

* `certificate_authority_id` - (Required) ID of the certificate authority.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_privateca_certificate_authority](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/privateca_certificate_authority) resource for details of the available attributes.

* `pem_csr` - The PEM-encoded signed certificate signing request (CSR). This is only set on subordinate certificate authorities.