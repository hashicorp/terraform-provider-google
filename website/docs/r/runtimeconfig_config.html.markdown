---
subcategory: "Cloud Runtime Configuration"
layout: "google"
page_title: "Google: google_runtimeconfig_config"
sidebar_current: "docs-google-runtimeconfig-config"
description: |-
  Manages a RuntimeConfig resource in Google Cloud.
---

# google\_runtimeconfig\_config

Manages a RuntimeConfig resource in Google Cloud. For more information, see the
[official documentation](https://cloud.google.com/deployment-manager/runtime-configurator/),
or the
[JSON API](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/).

## Example Usage

Example creating a RuntimeConfig resource.

```hcl
resource "google_runtimeconfig_config" "my-runtime-config" {
  name        = "my-service-runtime-config"
  description = "Runtime configuration values for my service"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the runtime config.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
is not provided, the provider project is used.

* `description` - (Optional) The description to associate with the runtime
config.

## Import

Runtime Configs can be imported using the `name` or full config name, e.g.

```
$ terraform import google_runtimeconfig_config.myconfig myconfig
```
```
$ terraform import google_runtimeconfig_config.myconfig projects/my-gcp-project/configs/myconfig
```
When importing using only the name, the provider project must be set.
