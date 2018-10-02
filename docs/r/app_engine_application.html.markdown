---
layout: "google"
page_title: "Google: google_app_engine_application"
sidebar_current: "docs-google-app-engine-application"
description: |-
 Allows management of an App Engine application.
---

# google\_app_engine_application

Allows creation and management of an App Engine application.

~> App Engine applications cannot be deleted once they're created; you have to delete the
   entire project to delete the application. Terraform will force you to set the `ack_delete_noop`
   field to `true` to acknowledge this limitation before you can successfully delete an App Engine
   application. There's no harm in leaving the `ack_delete_noop` field set to true at all times.

## Example Usage

```hcl
resource "google_project" "my_project" {
  name = "My Project"
  project_id = "your-project-id"
  org_id     = "1234567"
}

resource "google_app_engine_application" "app" {
  project         = "${google_project.my_project.project_id}"
  location_id     = "us-central'
  ack_delete_noop = true
}
```

## Argument Reference

The following arguments are supported:

* `location_id` - (Required) The [location](https://cloud.google.com/appengine/docs/locations)
   to serve the app from.

* `ack_delete_noop` - (Optional) Set to true to allow Terraform to "delete" your application without error.
   Has no bearing except to indicate that you're aware that when Terraform says it deletes an application,
   the application has not actually been deleted. To delete an application, the entire project must be deleted.

* `auth_domain` - (Optional) The domain to authenticate users with when using App Engine's User API.

* `serving_status` - (Optional) The serving status of the app. Note that this can't be updated at the moment.

* `feature_settings` - (Optional) A block of optional settings to configure specific App Engine features:

  * `split_health_checks` - (Optional) Set to false to use the legacy health check instead of the readiness
    and liveness checks.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `name` - Unique name of the app, usually `apps/{PROJECT_ID}`

* `url_dispatch_rule` - A list of dispatch rule blocks. Each block has a `domain`, `path`, and `service` field.

* `code_bucket` - The GCS bucket code is being stored in for this app.

* `default_hostname` - The default hostname for this app.

* `default_bucket` - The GCS bucket content is being stored in for this app.

* `gcr_domain` - The GCR domain used for storing managed Docker images for this app.

## Import

Applications can be imported using the ID of the project the application belongs to, e.g.

```
$ terraform import google_app_engine_application.app your-project-id
```
