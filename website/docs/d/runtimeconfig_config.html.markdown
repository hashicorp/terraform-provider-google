---
subcategory: "Runtime Configurator"
page_title: "Google: google_runtimeconfig_config"
description: |-
  Get information about a Google Cloud RuntimeConfig.
---

# google\_runtimeconfig\_config

To get more information about RuntimeConfigs, see:

* [API documentation](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/v1beta1/projects.configs)
* How-to Guides
    * [Runtime Configurator Fundamentals](https://cloud.google.com/deployment-manager/runtime-configurator/)

~> **Warning:** This datasource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta datasources.

## Example Usage

```hcl
data "google_runtimeconfig_config" "run-service" {
  name = "my-service"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Runtime Configurator configuration.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_runtimeconfig_config](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/runtimeconfig_config#argument-reference) resource for details of the available attributes.
