---
layout: "google"
page_title: "Provider: Google Cloud"
sidebar_current: "docs-google-index"
description: |-
  The Google Cloud provider is used to interact with Google Cloud services. The provider needs to be configured with the proper credentials before it can be used.
---

# Google Cloud Provider

The Google Cloud provider is used to interact with
[Google Cloud services](https://cloud.google.com/). The provider needs
to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
// Configure the Google Cloud provider
provider "google" {
  credentials = "${file("account.json")}"
  project     = "my-gce-project-id"
  region      = "us-central1"
}

// Create a new instance
resource "google_compute_instance" "default" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `credentials` - (Optional) Contents of a file that contains your service
  account private key in JSON format. You can download your existing
  [Google Cloud service account file]
  from the Google Cloud Console, or you can create a new one from the same page.

  Credentials can also be specified using any of the following environment
  variables (listed in order of precedence):

    * `GOOGLE_CREDENTIALS`
    * `GOOGLE_CLOUD_KEYFILE_JSON`
    * `GCLOUD_KEYFILE_JSON`

  The [`GOOGLE_APPLICATION_CREDENTIALS`][adc]
  environment variable can also contain the path of a file to obtain credentials
  from.

  If no credentials are specified, the provider will fall back to using the
  [Google Application Default Credentials][adc].
  If you are running Terraform from a GCE instance, see [Creating and Enabling
  Service Accounts for Instances][gce-service-account] for details.

  On your computer, if you have made your identity available as the
  Application Default Credentials by running [`gcloud auth application-default
  login`][gcloud adc], the provider will use your identity.

  ~> **Warning:** The gcloud method is not guaranteed to work for all APIs, and
  [service accounts] or [GCE metadata] should be used if possible.

* `project` - (Optional) The ID of the project to apply any resources to.  This
  can also be specified using any of the following environment variables (listed
  in order of precedence):

    * `GOOGLE_PROJECT`
    * `GOOGLE_CLOUD_PROJECT`
    * `GCLOUD_PROJECT`
    * `CLOUDSDK_CORE_PROJECT`

* `region` - (Optional) The region to operate under, if not specified by a given resource.
  This can also be specified using any of the following environment variables (listed in order of
  precedence):

    * `GOOGLE_REGION`
    * `GCLOUD_REGION`
    * `CLOUDSDK_COMPUTE_REGION`

* `zone` - (Optional) The zone to operate under, if not specified by a given resource.
  This can also be specified using any of the following environment variables (listed in order of
  precedence):

    * `GOOGLE_ZONE`
    * `GCLOUD_ZONE`
    * `CLOUDSDK_COMPUTE_ZONE`


## Beta Features

Some Google Provider resources contain Beta features; Beta GCP Features have no
deprecation policy, and no SLA, but are otherwise considered to be feature-complete
with only minor outstanding issues after their Alpha period. Beta is when a GCP feature
is publicly announced, and is when they generally become publicly available. For
more information see [the official documentation](https://cloud.google.com/terms/launch-stages).

Terraform resources that support beta features will always use the Beta APIs to provision
the resource. Importing a resource that supports beta features will always import those
features, even if the resource was created in a matter that was not explicitly beta.

[Google Cloud service account file]: https://console.cloud.google.com/apis/credentials/serviceaccountkey
[adc]: https://cloud.google.com/docs/authentication/production
[gce-service-account]: https://cloud.google.com/compute/docs/authentication
[gcloud adc]: https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login
[service accounts]: https://cloud.google.com/docs/authentication/getting-started
[GCE metadata]: https://cloud.google.com/docs/authentication/production#obtaining_credentials_on_compute_engine_kubernetes_engine_app_engine_flexible_environment_and_cloud_functions
