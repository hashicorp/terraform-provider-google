---
layout: "google"
page_title: "Google: google_client_config"
sidebar_current: "docs-google-datasource-client-config"
description: |-
  Get information about the configuration of the Google Cloud provider.
---

# google\_client\_config

Use this data source to access the configuration of the Google Cloud provider.

## Example Usage

```tf
data "google_client_config" "current" {}

output "project" {
  value = "${data.google_client_config.current.project}"
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `project` - The ID of the project to apply any resources to.

* `region` - The region to operate under.
