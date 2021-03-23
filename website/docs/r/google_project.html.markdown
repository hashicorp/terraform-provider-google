---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_project"
sidebar_current: "docs-google-project-x"
description: |-
 Allows management of a Google Cloud Platform project.
---

# google\_project

Allows creation and management of a Google Cloud Platform project.

Projects created with this resource must be associated with an Organization.
See the [Organization documentation](https://cloud.google.com/resource-manager/docs/quickstarts) for more details.

The user or service account that is running Terraform when creating a `google_project`
resource must have `roles/resourcemanager.projectCreator` on the specified organization. See the
[Access Control for Organizations Using IAM](https://cloud.google.com/resource-manager/docs/access-control-org)
doc for more information.

~> This resource reads the specified billing account on every terraform apply and plan operation so you must have permissions on the specified billing account.

~> It is recommended to use the `constraints/compute.skipDefaultNetworkCreation` [constraint](/docs/providers/google/r/google_organization_policy.html) to remove the default network instead of setting `auto_create_network` to false.

To get more information about projects, see:

* [API documentation](https://cloud.google.com/resource-manager/reference/rest/v1/projects)
* How-to Guides
    * [Creating and managing projects](https://cloud.google.com/resource-manager/docs/creating-managing-projects)

## Example Usage

```hcl
resource "google_project" "my_project" {
  name       = "My Project"
  project_id = "your-project-id"
  org_id     = "1234567"
}
```

To create a project under a specific folder

```hcl
resource "google_project" "my_project-in-a-folder" {
  name       = "My Project"
  project_id = "your-project-id"
  folder_id  = google_folder.department1.name
}

resource "google_folder" "department1" {
  display_name = "Department 1"
  parent       = "organizations/1234567"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The display name of the project.

* `project_id` - (Required) The project ID. Changing this forces a new project to be created.

* `org_id` - (Optional) The numeric ID of the organization this project belongs to.
    Changing this forces a new project to be created.  Only one of
    `org_id` or `folder_id` may be specified. If the `org_id` is
    specified then the project is created at the top level. Changing
    this forces the project to be migrated to the newly specified
    organization.

* `folder_id` - (Optional) The numeric ID of the folder this project should be
   created under. Only one of `org_id` or `folder_id` may be
   specified. If the `folder_id` is specified, then the project is
   created under the specified folder. Changing this forces the
   project to be migrated to the newly specified folder.

* `billing_account` - (Optional) The alphanumeric ID of the billing account this project
    belongs to. The user or service account performing this operation with Terraform
    must have at mininum Billing Account User privileges (`roles/billing.user`) on the billing account.
    See [Google Cloud Billing API Access Control](https://cloud.google.com/billing/docs/how-to/billing-access)
    for more details.

* `skip_delete` - (Optional) If true, the Terraform resource can be deleted
    without deleting the Project via the Google API.

* `labels` - (Optional) A set of key/value label pairs to assign to the project.

* `auto_create_network` - (Optional) Create the 'default' network automatically.  Default `true`.
    If set to `false`, the default network will be deleted.  Note that, for quota purposes, you
    will still need to have 1 network slot available to create the project successfully, even if
    you set `auto_create_network` to `false`, since the network will exist momentarily.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}`

* `number` - The numeric identifier of the project.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

Projects can be imported using the `project_id`, e.g.

```
$ terraform import google_project.my_project your-project-id
```
