---
subcategory: "App Engine"
page_title: "Google: google_app_engine_default_service_account"
description: |-
  Retrieve the default App Engine service account used in this project
---

# google\_app_engine\_default\_service\_account

Use this data source to retrieve the default App Engine service account for the specified project.

## Example Usage

```hcl
data "google_app_engine_default_service_account" "default" {
}

output "default_account" {
  value = data.google_app_engine_default_service_account.default.email
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project ID. If it is not provided, the provider project is used.


## Attributes Reference

The following attributes are exported:

* `email` - Email address of the default service account used by App Engine in this project.

* `unique_id` - The unique id of the service account.

* `name` - The fully-qualified name of the service account.

* `display_name` - The display name for the service account.

* `member` - The Identity of the service account in the form `serviceAccount:{email}`. This value is often used to refer to the service account in order to grant IAM permissions.
