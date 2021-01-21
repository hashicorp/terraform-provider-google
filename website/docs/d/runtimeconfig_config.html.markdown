---
subcategory: "Runtime Configurator"
layout: "google"
page_title: "Google: google_runtimeconfig_config"
sidebar_current: "docs-google-datasource-runtimeconfig-config"
description: |-
  Get information about a Google Cloud RuntimeConfig.
---

# google\_runtimeconfig\_config

To get more information about RuntimeConfigs, see:

* [API documentation](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/v1beta1/projects.configs)
* How-to Guides
    * [Runtime Configurator Fundamentals](https://cloud.google.com/deployment-manager/runtime-configurator/)

## Example Usage

```hcl
data "google_runtimeconfig_config" "run-service" {
  name = "my-service"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Run Service.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_runtimeconfig_config](https://www.terraform.io/docs/providers/google/r/runtimeconfig_config.html#argument-reference) resource for details of the available attributes.
