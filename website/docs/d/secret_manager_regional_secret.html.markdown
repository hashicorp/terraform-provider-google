---
subcategory: "Secret Manager"
description: |-
  Get information about a Secret Manager Regional Secret
---

# google_secret_manager_regional_secret

Use this data source to get information about a Secret Manager Regional Secret

## Example Usage 


```hcl
data "google_secret_manager_regional_secret" "secret_datasource" {
  secret_id = "secretname"
  location  = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `secret_id` - (required) The name of the regional secret.

* `location` - (required) The location of the regional secret. eg us-central1

* `project` - (optional) The ID of the project in which the resource belongs.

## Attributes Reference
See [google_secret_manager_regional_secret](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_regional_secret) resource for details of all the available attributes.
