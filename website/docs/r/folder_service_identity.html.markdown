---
subcategory: "Cloud Platform"
description: |-
 Generate folder service identity for a service.
---

# google_folder_service_identity

Generate folder service identity for a service.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

~> **Note:** Once created, this resource cannot be updated or destroyed. These
actions are a no-op.

~> **Note:** This resource can be used to retrieve the emails of the [Google-managed folder service accounts](https://cloud.google.com/iam/docs/service-agents) 
of the APIs that Google has configured with a Service Identity. You can run `gcloud beta services identity create --service SERVICE_NAME.googleapis.com --folder FOLDER` to
verify if an API supports this.

To get more information about Service Identity, see:

* [API documentation](https://cloud.google.com/service-usage/docs/reference/rest/v1beta1/services/generateServiceIdentity)

## Example Usage - Folder Service Identity Basic

```hcl
resource "google_folder" "my_folder" {
  parent = "organizations/1234567"
  display_name = "my-folder"
}

resource "google_folder_service_identity" "osconfig_sa" {
  provider = google-beta
  folder = google_folder.my_folder.folder_id
  service = "osconfig.googleapis.com"
}


resource "google_folder_iam_member" "admin" {
  folder = google_folder.my_folder.name
  role   = "roles/osconfig.serviceAgent"
  member = google_folder_service_identity.osconfig_sa.member
}
```

## Argument Reference

The following arguments are supported:

* `service` -
  (Required)
  The service to generate identity for.

- - -

* `folder` - (Required) The folder in which the resource belongs.

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
