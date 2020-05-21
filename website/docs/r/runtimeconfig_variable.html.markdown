---
subcategory: "Runtime Configurator"
layout: "google"
page_title: "Google: google_runtimeconfig_variable"
sidebar_current: "docs-google-runtimeconfig-variable"
description: |-
  Manages a RuntimeConfig variable in Google Cloud.
---

# google\_runtimeconfig\_variable

Manages a RuntimeConfig variable in Google Cloud. For more information, see the
[official documentation](https://cloud.google.com/deployment-manager/runtime-configurator/),
or the
[JSON API](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/).

## Example Usage

Example creating a RuntimeConfig variable.

```hcl
resource "google_runtimeconfig_config" "my-runtime-config" {
  name        = "my-service-runtime-config"
  description = "Runtime configuration values for my service"
}

resource "google_runtimeconfig_variable" "environment" {
  parent = google_runtimeconfig_config.my-runtime-config.name
  name   = "prod-variables/hostname"
  text   = "example.com"
}
```

You can also encode binary content using the `value` argument instead. The
value must be base64 encoded.

Example of using the `value` argument.

```hcl
resource "google_runtimeconfig_config" "my-runtime-config" {
  name        = "my-service-runtime-config"
  description = "Runtime configuration values for my service"
}

resource "google_runtimeconfig_variable" "my-secret" {
  parent = google_runtimeconfig_config.my-runtime-config.name
  name   = "secret"
  value  = base64encode(file("my-encrypted-secret.dat"))
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the variable to manage. Note that variable
names can be hierarchical using slashes (e.g. "prod-variables/hostname").

* `parent` - (Required) The name of the RuntimeConfig resource containing this
variable.

* `text` or `value` - (Required) The content to associate with the variable.
Exactly one of `text` or `variable` must be specified. If `text` is specified,
it must be a valid UTF-8 string and less than 4096 bytes in length. If `value`
is specified, it must be base64 encoded and less than 4096 bytes in length.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/configs/{{config}}/variables/{{name}}`

* `update_time` - (Computed) The timestamp in RFC3339 UTC "Zulu" format,
accurate to nanoseconds, representing when the variable was last updated.
Example: "2016-10-09T12:33:37.578138407Z".

## Import

Runtime Config Variables can be imported using the `name` or full variable name, e.g.

```
$ terraform import google_runtimeconfig_variable.myvariable myconfig/myvariable
```
```
$ terraform import google_runtimeconfig_variable.myvariable projects/my-gcp-project/configs/myconfig/variables/myvariable
```
When importing using only the name, the provider project must be set.
