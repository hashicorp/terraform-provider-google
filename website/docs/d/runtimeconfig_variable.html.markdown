---
subcategory: "Runtime Configurator"
page_title: "Google: google_runtimeconfig_variable"
description: |-
  Get information about a Google Cloud RuntimeConfig variable.
---

# google\_runtimeconfig\_variable

To get more information about RuntimeConfigs, see:

* [API documentation](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/v1beta1/projects.configs)
* How-to Guides
    * [Runtime Configurator Fundamentals](https://cloud.google.com/deployment-manager/runtime-configurator/)

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta datasources.

## Example Usage

```hcl
data "google_runtimeconfig_variable" "run-service" {
  parent = "my-service"
  name   = "prod-variables/hostname"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Runtime Configurator configuration.
* `parent` - (Required) The name of the RuntimeConfig resource containing this variable.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_runtimeconfig_variable](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/runtimeconfig_variable#argument-reference) resource for details of the available attributes.
