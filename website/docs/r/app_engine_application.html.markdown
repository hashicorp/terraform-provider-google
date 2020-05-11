---
subcategory: "App Engine"
layout: "google"
page_title: "Google: google_app_engine_application"
sidebar_current: "docs-google-app-engine-application"
description: |-
 Allows management of an App Engine application.
---

# google_app_engine_application

Allows creation and management of an App Engine application.

~> App Engine applications cannot be deleted once they're created; you have to delete the
   entire project to delete the application. Terraform will report the application has been
   successfully deleted; this is a limitation of Terraform, and will go away in the future.
   Terraform is not able to delete App Engine applications.

~> **Warning:** All arguments including `iap.oauth2_client_secret` will be stored in the raw
state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "google_project" "my_project" {
  name       = "My Project"
  project_id = "your-project-id"
  org_id     = "1234567"
}

resource "google_app_engine_application" "app" {
  project     = google_project.my_project.project_id
  location_id = "us-central"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project ID to create the application under.
   ~>**NOTE**: GCP only accepts project ID, not project number. If you are using number,
   you may get a "Permission denied" error.

* `location_id` - (Required) The [location](https://cloud.google.com/appengine/docs/locations)
   to serve the app from.

* `auth_domain` - (Optional) The domain to authenticate users with when using App Engine's User API.

* `serving_status` - (Optional) The serving status of the app.

* `feature_settings` - (Optional) A block of optional settings to configure specific App Engine features:

  * `split_health_checks` - (Required) Set to false to use the legacy health check instead of the readiness
    and liveness checks.

* `iap` - (Optional) Settings for enabling Cloud Identity Aware Proxy

  * `oauth2_client_id` - (Required) OAuth2 client ID to use for the authentication flow.

  * `oauth2_client_secret` - (Required) OAuth2 client secret to use for the authentication flow.
    The SHA-256 hash of the value is returned in the oauth2ClientSecretSha256 field.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `{{project}}`

* `name` - Unique name of the app, usually `apps/{PROJECT_ID}`

* `app_id` - Identifier of the app, usually `{PROJECT_ID}`

* `url_dispatch_rule` - A list of dispatch rule blocks. Each block has a `domain`, `path`, and `service` field.

* `code_bucket` - The GCS bucket code is being stored in for this app.

* `default_hostname` - The default hostname for this app.

* `default_bucket` - The GCS bucket content is being stored in for this app.

* `gcr_domain` - The GCR domain used for storing managed Docker images for this app.

* `iap` - Settings for enabling Cloud Identity Aware Proxy

  * `oauth2_client_secret_sha256` - Hex-encoded SHA-256 hash of the client secret.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 4 minutes.
- `update` - Default is 4 minutes.

## Import

Applications can be imported using the ID of the project the application belongs to, e.g.

```
$ terraform import google_app_engine_application.app your-project-id
```
