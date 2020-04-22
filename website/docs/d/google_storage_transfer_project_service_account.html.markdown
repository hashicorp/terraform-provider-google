---
subcategory: "Storage Transfer Service"
layout: "google"
page_title: "Google: google_storage_transfer_project_service_account"
sidebar_current: "docs-google-datasource-storage-transfer-project-service-account"
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

* `email` - Email address of the default service account used by Storage Transfer Jobs running in this project
