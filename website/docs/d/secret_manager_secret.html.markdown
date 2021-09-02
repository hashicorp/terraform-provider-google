---
subcategory: "Secret Manager"
layout: "google"
page_title: "Google: google_secret_manager_secret"
sidebar_current: "docs-google-datasource-secret-manager-secret"
description: |-
  Get information about a Secret Manager Secret
---

# google\_secret\_manager\_secret

Use this data source to get information about a Secret Manager Secret

## Example Usage 


```hcl
data "google_secret_manager_secret" "qa" {
  secret_id = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `secret_id` - (required) The name of the secret.

* `project` - (optional) The ID of the project in which the resource belongs.

## Attributes Reference
See [google_secret_manager_secret](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret) resource for details of all the available attributes.
