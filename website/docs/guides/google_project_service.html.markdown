---
page_title: "User guide for google_project_service"
description: |-
  An advanced user guide for the google_project_service resource
---

# User Guide - google_project_service

## Enabling Multiple Services in a Config

Users may want to activate many services simultaneously within a single Terraform config. The `google_project_service.service` field only supports a single value, but Terraform itself provides list iteration with [for_each](https://developer.hashicorp.com/terraform/language/meta-arguments/for_each).
For example:

```
variable "services" {
  type = list(string)
}

resource "google_project_service" "services" {
  for_each = toset(var.services)
  project = "my-project"
  service = each.value
}
```

For a more robust example, Google recommends the [project_services module](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/project_services). The `project_services` module simplifies configuring multiple services on a project at once as shown in [this example](https://github.com/terraform-google-modules/terraform-google-project-factory/tree/master/modules/project_services#example-usage).

-> **Note** Even though `google_project_service` represents a single resource, the Google provider will batch multiple changes into a single request when possible. See the `batching` [reference documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference.html#batching) for details.

## Request Rate Errors

The service management API called by the google_project_service resource uses request rate quota on the project of the account used to call the API (i.e. against the Terraform credentials) by default. That project (or a fixed `billing_project`) may exceed your request rate quota in larger configurations. `google_project_service` batches multiple changes into single requests when possible, see the `batching` [reference documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference.html#batching) for details.

Minimizing the number of total resources in root modules will help maximize the provider’s ability to batch requests. Oversized root modules slow Terraform’s execution time and can cause same-type requests to miss the batch window set by [batching.send_after](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference.html#send_after). See Google’s [guidance on root modules](https://cloud.google.com/docs/terraform/best-practices-for-terraform#root-modules).

## Newly activated service errors

Users may run into [issues](https://github.com/hashicorp/terraform-provider-google/issues/8214) when activating a service with `google_project_service` and immediately using the service in another resource. The most common error when doing so will look like:

```
Error: Error creating <Resource>: googleapi: Error 403: <Service> API has not been used in project <Project> before or it is disabled.
```

Activating a service is eventually consistent in GCP. Terraform attempts to mitigate this in `google_project_service` by waiting for the activation API’s long-running operation to finish and verifying that the service appears in a list of activated services. Despite these checks, service activation is not guaranteed by the time the `google_project_service` resource is done provisioning.

At the time of writing, there is no way for the provider to completely verify service activation. The time before `google_project_service` returns successfully may vary depending on the service, GCP-internal caching, and other circumstances. In particular, just-created projects may experience longer service activation times. Further mitigations users can try are detailed below.

If you are experiencing significantly long activation time with a specific service, it would be best to file an issue in that service’s public issue tracker or speak to your customer representative. 

### Mitigation - Adding sleeps

A common way to deal with eventual consistency with Terraform is to implement a sleep in between resources.

#### Using the time provider

The [time provider](https://registry.terraform.io/providers/hashicorp/time/latest/docs) offers a [time_sleep](https://registry.terraform.io/providers/hashicorp/time/latest/docs/resources/sleep) resource. You can use the `create_duration` field to provide a sleep duration.

```
resource "google_project" "my_project" {
  name = "foo-bar-baz"
  project_id = "foo-bar-baz-test"
  org_id = var.organization_id
  billing_account = var.billing_account_id
}

resource "time_resource" "wait_30_seconds" {
  depends_on = [google_project.my_project]

  create_duration = "30s"
}

resource "google_project_service" "my_service" {
  project = google_project.my_project.id
  service = "firebase.googleapis.com"

  disable_dependent_services = true
  depends_on = [time_resource.wait_30_seconds]
}
```

#### Using the local-exec provisioner

Terraform provides the [local-exec provisioner](https://developer.hashicorp.com/terraform/language/resources/provisioners/local-exec) to invoke a local executable after resource creation. Depending on the machine running Terraform, you may invoke some command to sleep for a duration.

```
resource "google_project" "my_project" {
  name = "foo-bar-baz"
  project_id = "foo-bar-baz-test"
  org_id = var.organization_id
  billing_account = var.billing_account_id
}

resource "null_resource" "delay" {
  provisioner "local-exec" {
    command = "sleep 60"
  }
  triggers = {
    "project" = "${google_project.my_project.id}"
  }
}

resource "google_project_service" "my_service" {
  project = google_project.my_project.id
  service = "firebase.googleapis.com"

  disable_dependent_services = true
  depends_on = [null_resource.delay]
}
```