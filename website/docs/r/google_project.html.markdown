---
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

The service account used to run Terraform when creating a `google_project`
resource must have `roles/resourcemanager.projectCreator`. See the
[Access Control for Organizations Using IAM](https://cloud.google.com/resource-manager/docs/access-control-org)
doc for more information.

Note that prior to 0.8.5, `google_project` functioned like a data source,
meaning any project referenced by it had to be created and managed outside
Terraform. As of 0.8.5, `google_project` functions like any other Terraform
resource, with Terraform creating and managing the project. To replicate the old
behavior, either:

* Use the project ID directly in whatever is referencing the project, using the
  [google_project_iam_policy](/docs/providers/google/r/google_project_iam.html)
  to replace the old `policy_data` property.
* Use the [import](/docs/import/usage.html) functionality
  to import your pre-existing project into Terraform, where it can be referenced and
  used just like always, keeping in mind that Terraform will attempt to undo any changes
  made outside Terraform.

~> It's important to note that any project resources that were added to your Terraform config
prior to 0.8.5 will continue to function as they always have, and will not be managed by
Terraform. Only newly added projects are affected.

## Example Usage

```hcl
resource "google_project" "my_project" {
  name = "My Project"
  project_id = "your-project-id"
  org_id     = "1234567"
}
```

To create a project under a specific folder

```hcl
resource "google_project" "my_project-in-a-folder" {
  name = "My Project"
  project_id = "your-project-id"
  folder_id  = "${google_folder.department1.name}"
}

resource "google_folder" "department1" {
  display_name = "Department 1"
  parent     = "organizations/1234567"
}
```

To create a project with an App Engine app attached

```hcl
resource "google_project" "my-app-engine-app" {
  name = "App Engine Project"
  project_id = "app-engine-project"
  org_id = "1234567"

  app_engine {
    location_id = "us-central"
  }
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
    must have Billing Account Administrator privileges (`roles/billing.admin`) in
    the organization. See [Google Cloud Billing API Access Control](https://cloud.google.com/billing/v1/how-tos/access-control)
    for more details.

* `skip_delete` - (Optional) If true, the Terraform resource can be deleted
    without deleting the Project via the Google API.

* `policy_data` - (Deprecated) The IAM policy associated with the project.
    This argument is no longer supported, and will be removed in a future version
    of Terraform. It should be replaced with a `google_project_iam_policy` resource.

* `labels` - (Optional) A set of key/value label pairs to assign to the project.

* `auto_create_network` - (Optional) Create the 'default' network automatically.  Default true.
    Note: this might be more accurately described as "Delete Default Network", since the network
    is created automatically then deleted before project creation returns, but we choose this
    name to match the GCP Console UI. Setting this field to false will enable the Compute Engine
    API which is required to delete the network.

* `app_engine` - (Optional) A block of configuration to enable an App Engine app. Setting this
   field will enabled the App Engine Admin API, which is required to manage the app.

The `app_engine` block has the following configuration options:

* `location_id` - (Required) The [location](https://cloud.google.com/appengine/docs/locations)
   to serve the app from.
* `auth_domain` - (Optional) The domain to authenticate users with when using App Engine's User API.
* `serving_status` - (Optional) The serving status of the app. Note that this can't be updated at the moment.
* `feature_settings` - (Optional) A block of optional settings to configure specific App Engine features:
  * `split_health_checks` - (Optional) Set to false to use the legacy health check instead of the readiness
    and liveness checks.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `number` - The numeric identifier of the project.

* `policy_etag` - (Deprecated) The etag of the project's IAM policy, used to
    determine if the IAM policy has changed. Please use `google_project_iam_policy`'s
    `etag` property instead; future versions of Terraform will remove the `policy_etag`
    attribute

* `app_engine.0.name` - Unique name of the app, usually `apps/{PROJECT_ID}`
* `app_engine.0.url_dispatch_rule` - A list of dispatch rule blocks. Each block has a `domain`, `path`, and `service` field.
* `app_engine.0.code_bucket` - The GCS bucket code is being stored in for this app.
* `app_engine.0.default_hostname` - The default hostname for this app.
* `app_engine.0.default_bucket` - The GCS bucket content is being stored in for this app.
* `app_engine.0.gcr_domain` - The GCR domain used for storing managed Docker images for this app.

## Import

Projects can be imported using the `project_id`, e.g.

```
$ terraform import google_project.my_project your-project-id
```
