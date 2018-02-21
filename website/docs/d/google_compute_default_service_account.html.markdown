---
layout: "google"
page_title: "Google: google_compute_default_service_account"
sidebar_current: "docs-google-datasource-compute-default-service-account"
description: |-
  Retrieve default service account used by VMs running in this project
---

# google\_compute\_default\_service\_account

Use this data source to retrieve default service account for this project

## Example Usage

```hcl
data "google_compute_default_service_account" "default" { }

output "default_account" {
  value = "${google_compute_default_service_account.default.id}"
} 
```

## Argument Reference

There are no arguments available for this data source.


## Attributes Reference

The following attributes are exported:

* `id` - Email address of the default service account used by VMs running in this project
