---
subcategory: "Cloud Platform"
page_title: "Google: google_project_service"
description: |-
 Allows management of a single API service for a Google Cloud Platform project.
---

# google\_project\_service

Allows management of a single API service for a Google Cloud Platform project. 

For a list of services available, visit the [API library page](https://console.cloud.google.com/apis/library)
or run `gcloud services list --available`.

This resource requires the [Service Usage API](https://console.cloud.google.com/apis/library/serviceusage.googleapis.com)
to use.

To get more information about `google_project_service`, see:

* [API documentation](https://cloud.google.com/service-usage/docs/reference/rest/v1/services)
* How-to Guides
    * [Enabling and Disabling Services](https://cloud.google.com/service-usage/docs/enable-disable)

## Example Usage

```hcl
resource "google_project_service" "project" {
  project = "your-project-id"
  service = "iam.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required) The service to enable.

* `project` - (Optional) The project ID. If not provided, the provider project
is used.

* `disable_dependent_services` - (Optional) If `true`, services that are enabled
and which depend on this service should also be disabled when this service is
destroyed. If `false` or unset, an error will be generated if any enabled
services depend on this service when destroying it.

* `disable_on_destroy` - (Optional) If true, disable the service when the
Terraform resource is destroyed. Defaults to true. May be useful in the event
that a project is long-lived but the infrastructure running in that project
changes frequently.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `{{project}}/{{service}}`

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `read`   - Default is 10 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Project services can be imported using the `project_id` and `service`, e.g.

```
$ terraform import google_project_service.my_project your-project-id/iam.googleapis.com
```

Note that unlike other resources that fail if they already exist,
`terraform apply` can be successfully used to verify already enabled services.
This means that when importing existing resources into Terraform, you can either
import the `google_project_service` resources or treat them as new
infrastructure and run `terraform apply` to add them to state.



## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
