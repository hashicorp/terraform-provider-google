---
subcategory: "Identity-Aware Proxy"
page_title: "Google: google_iap_client"
description: |-
  Contains the data that describes an Identity Aware Proxy owned client.
---
# google_iap_client

Get info about a Google Cloud IAP Client.

## Example Usage

```tf
data "google_project" "project" {
  project_id = "foobar"
}

data "google_iap_client" "project_client" {
  brand        =  "projects/${data.google_project.project.number}/brands/[BRAND_NUMBER]"
  client_id    = FOO.apps.googleusercontent.com
}

```

## Argument Reference

The following arguments are supported:

* `brand` - (Required) The name of the brand.

* `client_id` - (Required) The client_id of the brand.

## Attributes Reference

See [google_iap_client](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/iap_client) resource for details of the available attributes.
