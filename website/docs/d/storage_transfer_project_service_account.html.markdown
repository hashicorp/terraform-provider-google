---
subcategory: "Storage Transfer Service"
page_title: "Google: google_storage_transfer_project_service_account"
description: |-
  Retrieve default service account used by Storage Transfer Jobs running in this project
---

# google\_storage\_transfer\_project\_service\_account

Use this data source to retrieve Storage Transfer service account for this project

## Example Usage

```hcl
data "google_storage_transfer_project_service_account" "default" {
}

output "default_account" {
  value = data.google_storage_transfer_project_service_account.default.email
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project ID. If it is not provided, the provider project is used.


## Attributes Reference

The following attributes are exported:

* `email` - Email address of the default service account used by Storage Transfer Jobs running in this project.
* `subject_id` - Unique identifier for the service account.
* `member` - The Identity of the service account in the form `serviceAccount:{email}`. This value is often used to refer to the service account in order to grant IAM permissions.
