---
subcategory: "Cloud Platform"
description: |-
 Generate service identity for a service.
---

# google_project_service_identity

~> **Warning:** These resources are in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

Generate service identity for a service.

~> **Note:** Once created, this resource cannot be updated or destroyed. These
actions are a no-op.

~> **Note:** This resource can be used to retrieve the emails of the [Google-managed service accounts](https://cloud.google.com/iam/docs/service-agents) 
of the APIs that Google has configured with a Service Identity. You can run `gcloud beta services identity create --service SERVICE_NAME.googleapis.com` to
verify if an API supports this.

To get more information about Service Identity, see:

* [API documentation](https://cloud.google.com/service-usage/docs/reference/rest/v1beta1/services/generateServiceIdentity)

## Example Usage - Service Identity Basic

```hcl
data "google_project" "project" {}

resource "google_project_service_identity" "hc_sa" {
  provider = google-beta

  project = data.google_project.project.project_id
  service = "healthcare.googleapis.com"
}

resource "google_project_iam_member" "hc_sa_bq_jobuser" {
  project = data.google_project.project.project_id
  role    = "roles/bigquery.jobUser"
  member  = google_project_service_identity.hc_sa.member
}
```

## Argument Reference

The following arguments are supported:

* `service` -
  (Required)
  The service to generate identity for.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `email` - The email address of the Google managed service account.
* `member` - The Identity of the Google managed service account in the form 'serviceAccount:{email}'. This value is often used to refer to the service account in order to grant IAM permissions.

## Import

This resource does not support import.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

* `create` - Default is 20 minutes.

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
