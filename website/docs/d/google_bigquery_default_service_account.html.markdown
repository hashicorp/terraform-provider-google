---
subcategory: "BigQuery"
layout: "google"
page_title: "Google: google_bigquery_default_service_account"
sidebar_current: "docs-google-datasource-bigquery-default-service-account"
description: |-
  Retrieve default service account used by bigquery encryption in this project
---

# google\_bigquery\_default\_service\_account

Use this data source to retrieve default service account for this project

## Example Usage

```hcl
data "google_bigquery_default_service_account" "default" { }

output "default_account" {
  value = "${data.google_bigquery_default_service_account.default.email}"
} 
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project ID. If it is not provided, the provider project is used.


## Attributes Reference

The following attributes are exported:

* `email` - Email address of the default service account used by bigquery encryption in this project
