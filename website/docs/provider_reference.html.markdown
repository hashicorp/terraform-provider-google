---
layout: "google"
page_title: "google provider reference"
sidebar_current: "docs-google-provider-reference"
description: |-
  The Google provider is used to configure your GCP project, location, and creds
---

# `google` provider reference

-> We recently introduced the `google-beta` provider. See [Provider Versions](https://terraform.io/docs/providers/google/provider_versions.html)
for more details on how to use `google-beta`.

The `google` and `google-beta` provider blocks are used to configure default values for
your GCP project and location (`zone` and `region`), and add your credentials.

-> You can avoid using a provider block by using environment variables. Every
field of `google` and `google-beta` is inferred from your environment when it
has not been explicitly set. Even better - the GA and beta providers will both
share the same values.

## Example Usage - Basic provider blocks

```hcl
provider "google" {
  credentials = "${file("account.json")}"
  project     = "my-project-id"
  region      = "us-central1"
  zone        = "us-central1-c"
}
```

```hcl
provider "google-beta" {
  credentials = "${file("account.json")}"
  project     = "my-project-id"
  region      = "us-central1"
  zone        = "us-central1-c"
}
```

## Example Usage - Using beta features with `google-beta` 

To use Google Cloud Platform features that are in beta, explicitly set the provider for your
resource to `google-beta`. See [Provider Versions](https://terraform.io/docs/providers/google/provider_versions.html)
for a full reference on how to use different GCP versions with the Google provider.

```hcl
resource "google_compute_instance" "ga-instance" {
  provider = "google"

  # ...
}

resource "google_compute_instance" "beta-instance" {
  provider = "google-beta"

  # ...
}
```


## Configuration Reference

The following keys can be used to configure the provider. Both `google` and `google-beta`
share the same configuration.

* `credentials` - (Optional) The path or contents of a file that contains your
  service account private key in JSON format. You can download your existing
  [Google Cloud service account file] from the Google Cloud Console, or you can
  create a new one from the same page.

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

  -> [Service accounts][service accounts] are the recommended way
  to manage GCP credentials. [GCE metadata] is also acceptable, although it can
  only be used when running Terraform from within [certain GCP resources](https://cloud.google.com/docs/authentication/production#obtaining_credentials_on_compute_engine_kubernetes_engine_app_engine_flexible_environment_and_cloud_functions).
  Credentials obtained through `gcloud` are not guaranteed to work for all APIs.

* `project` - (Optional) The ID of the project to apply any resources to.  This
  can also be specified using any of the following environment variables (listed
  in order of precedence):

    * `GOOGLE_PROJECT`
    * `GOOGLE_CLOUD_PROJECT`
    * `GCLOUD_PROJECT`
    * `CLOUDSDK_CORE_PROJECT`

    -> `GOOGLE_PROJECT` is the recommended environment variable to use if
    you choose to add your project using environment variables.

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

[Google Cloud service account file]: https://console.cloud.google.com/apis/credentials/serviceaccountkey
[adc]: https://cloud.google.com/docs/authentication/production
[gce-service-account]: https://cloud.google.com/compute/docs/authentication
[gcloud adc]: https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login
[service accounts]: https://cloud.google.com/docs/authentication/getting-started
[GCE metadata]: https://cloud.google.com/docs/authentication/production#obtaining_credentials_on_compute_engine_kubernetes_engine_app_engine_flexible_environment_and_cloud_functions
