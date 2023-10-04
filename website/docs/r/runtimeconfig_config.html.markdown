---
subcategory: "Runtime Configurator"
description: |-
  Manages a RuntimeConfig resource in Google Cloud.
---

# google\_runtimeconfig\_config

Manages a RuntimeConfig resource in Google Cloud.

To get more information about RuntimeConfigs, see:

* [API documentation](https://cloud.google.com/deployment-manager/runtime-configurator/reference/rest/v1beta1/projects.configs)
* How-to Guides
    * [Runtime Configurator Fundamentals](https://cloud.google.com/deployment-manager/runtime-configurator/)

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

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

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/configs/{{name}}`

## Import

Runtime Configs can be imported using the `name` or full config name, e.g.

* `projects/{{project_id}}/configs/{{name}}`
* `{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Runtime Configs using one of the formats above. For example:

```tf
import {
  id = "projects/{{project_id}}/configs/{{name}}"
  to = google_runtimeconfig_config.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Runtime Configs can be imported using one of the formats above. For example:

```
$ terraform import google_runtimeconfig_config.default projects/{{project_id}}/configs/{{name}}
$ terraform import google_runtimeconfig_config.default {{name}}
```

When importing using only the name, the provider project must be set.
