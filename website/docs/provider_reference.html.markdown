---
layout: "google"
page_title: "Google Provider Configuration Reference"
sidebar_current: "docs-google-provider-reference"
description: |-
  Configuration reference for the Google provider for Terraform.
---

# Google Provider Configuration Reference

-> Try out Terraform 0.12 with the Google provider! `google` and `google-beta` are 0.12-compatible from `2.5.0` onwards.

The `google` and `google-beta` provider blocks are used to configure the
credentials you use to authenticate with GCP, as well as a default project and
location (`zone` and/or `region`) for your resources.

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

To use Google Cloud Platform features that are in beta, you need to both:

* Explicitly define a `google-beta` provider block

* explicitly set the provider for your resource to `google-beta`.

See [Provider Versions](https://terraform.io/docs/providers/google/provider_versions.html)
for a full reference on how to use features from different GCP API versions in
the Google provider.

```hcl
resource "google_compute_instance" "ga-instance" {
  provider = "google"

  # ...
}

resource "google_compute_instance" "beta-instance" {
  provider = "google-beta"

  # ...
}

provider "google-beta" {}
```


## Configuration Reference

The following attributes can be used to configure the provider. The quick
reference should be sufficient for most use cases, but see the full reference
if you're interested in more details. Both `google` and `google-beta` share the
same configuration. 

### Quick Reference

* `credentials` - (Optional) Either the path to or the contents of a
[service account key file] in JSON format. You can
[manage key files using the Cloud Console].

* `project` - (Optional) The default project to manage resources in. If another
project is specified on a resource, it will take precedence.

* `region` - (Optional) The default region to manage resources in. If another
region is specified on a regional resource, it will take precedence.

* `zone` - (Optional) The default zone to manage resources in. Generally, this
zone should be within the default region you specified. If another zone is
specified on a zonal resource, it will take precedence.

---

* `scopes` - (Optional) The list of OAuth 2.0 [scopes] requested when generating
an access token using the service account key specified in `credentials`.

* `access_token` - (Optional) A temporary [OAuth 2.0 access token] obtained from
the Google Authorization server, i.e. the `Authorization: Bearer` token used to
authenticate HTTP requests to GCP APIs. This is an alternative to `credentials`,
and ignores the `scopes` field. If both are specified, `access_token` will be
used over the `credentials` field.

### Full Reference

* `credentials` - (Optional) Either the path to or the contents of a
[service account key file] in JSON format. You can
[manage key files using the Cloud Console]. Your service account key file is
used to complete a two-legged OAuth 2.0 flow to obtain access tokens to
authenticate with the GCP API as needed; Terraform will use it to reauthenticate
automatically when tokens expire. Alternatively, this can be specified using the
`GOOGLE_CREDENTIALS` environment variable or any of the following ordered
by precedence.

    * GOOGLE_CREDENTIALS
    * GOOGLE_CLOUD_KEYFILE_JSON
    * GCLOUD_KEYFILE_JSON

    Using Terraform-specific [service accounts] to authenticate with GCP is the
    recommended practice when using Terraform. If no Terraform-specific
    credentials are specified, the provider will fall back to using
    [Google Application Default Credentials][adc]. To use them, you can enter
    the path of your service account key file in the
    `GOOGLE_APPLICATION_CREDENTIALS` environment variable, or configure
    authentication through one of the following;

* If you're running Terraform from a GCE instance, default credentials
are automatically available. See
[Creating and Enabling Service Accounts for Instances][gce-service-account]
for more details.

* On your computer, you can make your Google identity available by
running [`gcloud auth application-default login`][gcloud adc]. This
approach isn't recommended- some APIs are not compatible with
credentials obtained through `gcloud`.

---

* `project` - (Optional) The default project to manage resources in. If another
project is specified on a resource, it will take precedence. This can also be
specified using the `GOOGLE_PROJECT` environment variable, or any of the
following ordered by precedence.

    * GOOGLE_PROJECT
    * GOOGLE_CLOUD_PROJECT
    * GCLOUD_PROJECT
    * CLOUDSDK_CORE_PROJECT

---

* `region` - (Optional) The default region to manage resources in. If another
region is specified on a regional resource, it will take precedence.
Alternatively, this can be specified using the `GOOGLE_REGION` environment
variable or any of the following ordered by precedence.

    * GOOGLE_REGION
    * GCLOUD_REGION
    * CLOUDSDK_COMPUTE_REGION

---

* `zone` - (Optional) The default zone to manage resources in. Generally, this
zone should be within the default region you specified. If another zone is
specified on a zonal resource, it will take precedence. Alternatively, this can
be specified using the `GOOGLE_ZONE` environment variable or any of the
following ordered by precedence.

    * GOOGLE_ZONE
    * GCLOUD_ZONE
    * CLOUDSDK_COMPUTE_ZONE

---

* `access_token` - (Optional) A temporary [OAuth 2.0 access token] obtained from
the Google Authorization server, i.e. the `Authorization: Bearer` token used to
authenticate HTTP requests to GCP APIs. If both are specified, `access_token` will be
used over the `credentials` field. This is an alternative to `credentials`,
and ignores the `scopes` field. Alternatively, this can be specified using the
`GOOGLE_OAUTH_ACCESS_TOKEN` environment variable.

    -> These access tokens cannot be renewed by Terraform and thus will only
    work until they expire. If you anticipate Terraform needing access for
    longer than a token's lifetime (default `1 hour`), please use a service
    account key with `credentials` instead.

---

* `scopes` - (Optional) The list of OAuth 2.0 [scopes] requested when generating
an access token using the service account key specified in `credentials`.

    By default, the following scopes are configured:

    * https://www.googleapis.com/auth/compute
    * https://www.googleapis.com/auth/cloud-platform
    * https://www.googleapis.com/auth/ndev.clouddns.readwrite
    * https://www.googleapis.com/auth/devstorage.full_control

[OAuth 2.0 access token]: https://developers.google.com/identity/protocols/OAuth2
[service account key file]: https://cloud.google.com/iam/docs/creating-managing-service-account-keys
[manage key files using the Cloud Console]: https://console.cloud.google.com/apis/credentials/serviceaccountkey
[adc]: https://cloud.google.com/docs/authentication/production
[gce-service-account]: https://cloud.google.com/compute/docs/authentication
[gcloud adc]: https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login
[service accounts]: https://cloud.google.com/docs/authentication/getting-started
[GCE metadata]: https://cloud.google.com/docs/authentication/production#obtaining_credentials_on_compute_engine_kubernetes_engine_app_engine_flexible_environment_and_cloud_functions
[scopes]: https://developers.google.com/identity/protocols/googlescopes
