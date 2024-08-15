---
subcategory: "Cloud Platform"
description: |-
 Allows management of a single API service for a Google Cloud project.
---

# google_project_service

Allows management of a single API service for a Google Cloud project. 

For a list of services available, visit the [API library page](https://console.cloud.google.com/apis/library)
or run `gcloud services list --available`.

This resource requires the [Service Usage API](https://console.cloud.google.com/apis/library/serviceusage.googleapis.com)
to use.

To get more information about `google_project_service`, see:

* [API documentation](https://cloud.google.com/service-usage/docs/reference/rest/v1/services)
* How-to Guides
    * [Enabling and Disabling Services](https://cloud.google.com/service-usage/docs/enable-disable)
* Terraform guidance
    * [User Guide - google_project_service](/docs/providers/google/guides/google_project_service.html)

## Example Usage

```hcl
resource "google_project_service" "project" {
  project = "your-project-id"
  service = "iam.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_on_destroy = false
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) The service to enable.

* `project` - (Optional) The project ID. If not provided, the provider project
is used.

* `disable_on_destroy` - (Optional) If `true` or unset, disable the service when the
Terraform resource is destroyed. If `false`, the service will be left enabled when
the Terraform resource is destroyed. Defaults to `true`. Most configurations should
set this to `false`; it should generally only be `true` or unset in configurations
that manage the `google_project` resource itself.

* `disable_dependent_services` - (Optional) If `true`, services that are enabled
and which depend on this service should also be disabled when this service is
destroyed. If `false` or unset, an error will be returned if any enabled
services depend on this service when attempting to destroy it.

* `check_if_service_has_usage_on_destroy` - (Optional)
[Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)
If `true`, the usage of the service to be disabled will be checked and an error
will be returned if the service to be disabled has usage in last 30 days.
Defaults to `false`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `{{project}}/{{service}}`

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 20 minutes.
- `read`   - Default is 10 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Project services can be imported using the `project_id` and `service`, e.g.

* `{{project_id}}/{{service}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import project services using one of the formats above. For example:

```tf
import {
  id = "{{project_id}}/{{service}}"
  to = google_project_service.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), project services can be imported using one of the formats above. For example:

```
$ terraform import google_project_service.default {{project_id}}/{{service}}
```

Note that unlike other resources that fail if they already exist,
`terraform apply` can be successfully used to verify already enabled services.
This means that when importing existing resources into Terraform, you can either
import the `google_project_service` resources or treat them as new
infrastructure and run `terraform apply` to add them to state.

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
