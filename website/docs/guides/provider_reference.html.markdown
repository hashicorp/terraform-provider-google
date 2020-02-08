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

See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html)
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

* `user_project_override` - (Optional) Defaults to false. If true, uses the
resource project for preconditions, quota, and billing, instead of the project
the credentials belong to. Not all resources support this- see the
documentation for each resource to learn whether it does.

* `{{service}}_custom_endpoint` - (Optional) The endpoint for a service's APIs,
such as `compute_custom_endpoint`. Defaults to the production GCP endpoint for
the service. This can be used to configure the Google provider to communicate
with GCP-like APIs such as [the Cloud Functions emulator](https://github.com/googlearchive/cloud-functions-emulator).
Values are expected to include the version of the service, such as
`https://www.googleapis.com/compute/v1/`.

* `batching` - (Optional) This block controls batching GCP calls for groups of specific resource types. Structure is documented below.
~>**NOTE**: Batching is not implemented for the majority or resources/request types and is bounded by two values. If you are running into issues with slow batches
resources, you may need to adjust one or both of 1) the core [`-parallelism`](https://www.terraform.io/docs/commands/apply.html#parallelism-n) flag, which controls how many concurrent resources are being operated on and 2) `send_after`, the time interval after which a batch is sent.

* `request_timeout` - (Optional) A duration string controlling the amount of time
the provider should wait for a single HTTP request.  This will not adjust the
amount of time the provider will wait for a logical operation - use the resource
timeout blocks for that.

The `batching` fields supports:

* `send_after` - (Optional) A duration string representing the amount of time
after which a request should be sent. Defaults to 3s. Note that if you increase
`parallelism` you should also increase this value.

* `enable_batching` - (Optional) Defaults to true. If false, disables batching
   so requests that have batching capabilities are instead is sent one by one.

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
    * https://www.googleapis.com/auth/userinfo.email

---

* `{{service}}_custom_endpoint` - (Optional) The endpoint for a service's APIs,
such as `compute_custom_endpoint`. Defaults to the production GCP endpoint for
the service. This can be used to configure the Google provider to communicate
with GCP-like APIs such as [the Cloud Functions emulator](https://github.com/googlearchive/cloud-functions-emulator).
Values are expected to include the version of the service, such as
`https://www.googleapis.com/compute/v1/`.

~> Support for custom endpoints is on a best-effort basis. The underlying
endpoint and default values for a resource can be changed at any time without
being considered a breaking change.

A full list of configurable keys, their default value (in the `google` provider
followed by `google-beta` if they differ), and an environment variable that can
be used for configuration are below:

* `access_context_manager_custom_endpoint` (`GOOGLE_ACCESS_CONTEXT_MANAGER_CUSTOM_ENDPOINT`) - `https://accesscontextmanager.googleapis.com/v1/`
* `app_engine_custom_endpoint` (`GOOGLE_APP_ENGINE_CUSTOM_ENDPOINT`) - `https://appengine.googleapis.com/v1/`
* `bigquery_custom_endpoint` (`GOOGLE_BIGQUERY_CUSTOM_ENDPOINT`) - `https://www.googleapis.com/bigquery/v2/`
* `bigtable_custom_endpoint` (`GOOGLE_BIGTABLE_CUSTOM_ENDPOINT`) - `https://bigtableadmin.googleapis.com/v2/`
* `binary_authorization_custom_endpoint` (`GOOGLE_BINARY_AUTHORIZATION_CUSTOM_ENDPOINT`) - `https://binaryauthorization.googleapis.com/v1/`
* `cloud_billing_custom_endpoint` (`GOOGLE_CLOUD_BILLING_CUSTOM_ENDPOINT`) - `https://cloudbilling.googleapis.com/v1/`
* `cloud_build_custom_endpoint` (`GOOGLE_CLOUD_BUILD_CUSTOM_ENDPOINT`) - `https://cloudbuild.googleapis.com/v1/`
* `cloud_functions_custom_endpoint` (`GOOGLE_CLOUD_FUNCTIONS_CUSTOM_ENDPOINT`) - `https://cloudfunctions.googleapis.com/v1/`
* `cloud_iot_custom_endpoint` (`GOOGLE_CLOUD_IOT_CUSTOM_ENDPOINT`) - `https://cloudiot.googleapis.com/v1/`
* `cloud_scheduler_custom_endpoint` (`GOOGLE_CLOUD_SCHEDULER_CUSTOM_ENDPOINT`) - `https://cloudscheduler.googleapis.com/v1/`
* `composer_custom_endpoint` (`GOOGLE_COMPOSER_CUSTOM_ENDPOINT`) - `https://composer.googleapis.com/v1beta1/`
* `compute_custom_endpoint` (`GOOGLE_COMPUTE_CUSTOM_ENDPOINT`) - `https://www.googleapis.com/compute/v1/` | `https://www.googleapis.com/compute/beta/`
* `compute_beta_custom_endpoint` (`GOOGLE_COMPUTE_BETA_CUSTOM_ENDPOINT`) - `https://www.googleapis.com/compute/beta/`
* `container_custom_endpoint` (`GOOGLE_CONTAINER_CUSTOM_ENDPOINT`) - `https://container.googleapis.com/v1/`
* `container_beta_custom_endpoint` (`GOOGLE_CONTAINER_BETA_CUSTOM_ENDPOINT`) - `https://container.googleapis.com/v1beta1/`
* `dataproc_custom_endpoint` (`GOOGLE_DATAPROC_CUSTOM_ENDPOINT`) - `https://dataproc.googleapis.com/v1/`
* `dataproc_beta_custom_endpoint` (`GOOGLE_DATAPROC_BETA_CUSTOM_ENDPOINT`) - `https://dataproc.googleapis.com/v1beta2/`
* `dataflow_custom_endpoint` (`GOOGLE_DATAFLOW_CUSTOM_ENDPOINT`) - `https://dataflow.googleapis.com/v1b3/`
* `dns_custom_endpoint` (`GOOGLE_DNS_CUSTOM_ENDPOINT`) - `https://www.googleapis.com/dns/v1/` | `https://www.googleapis.com/dns/v1beta2/`
* `dns_beta_custom_endpoint` (`GOOGLE_DNS_BETA_CUSTOM_ENDPOINT`) - `https://www.googleapis.com/dns/v1beta2/`
* `filestore_custom_endpoint` (`GOOGLE_FILESTORE_CUSTOM_ENDPOINT`) - `https://file.googleapis.com/v1/`
* `firestore_custom_endpoint` (`GOOGLE_FIRESTORE_CUSTOM_ENDPOINT`) - `https://firestore.googleapis.com/v1/`
* `iam_custom_endpoint` (`GOOGLE_IAM_CUSTOM_ENDPOINT`) - `https://iam.googleapis.com/v1/`
* `iam_credentials_custom_endpoint` (`GOOGLE_IAM_CREDENTIALS_CUSTOM_ENDPOINT`) - `https://iamcredentials.googleapis.com/v1/`
* `kms_custom_endpoint` (`GOOGLE_KMS_CUSTOM_ENDPOINT`) - `https://cloudkms.googleapis.com/v1/`
* `logging_custom_endpoint` (`GOOGLE_LOGGING_CUSTOM_ENDPOINT`) - `https://logging.googleapis.com/v2/`
* `monitoring_custom_endpoint` (`GOOGLE_MONITORING_CUSTOM_ENDPOINT`) - `https://monitoring.googleapis.com/v3/`
* `pubsub_custom_endpoint` (`GOOGLE_PUBSUB_CUSTOM_ENDPOINT`) - `https://pubsub.googleapis.com/v1/`
* `redis_custom_endpoint` (`GOOGLE_REDIS_CUSTOM_ENDPOINT`) - `https://redis.googleapis.com/v1/` | `https://redis.googleapis.com/v1beta1/`
* `resource_manager_custom_endpoint` (`GOOGLE_RESOURCE_MANAGER_CUSTOM_ENDPOINT`) - `https://cloudresourcemanager.googleapis.com/v1/`
* `resource_manager_v2beta1_custom_endpoint` (`GOOGLE_RESOURCE_MANAGER_V2BETA1_CUSTOM_ENDPOINT`) - `https://cloudresourcemanager.googleapis.com/v2beta1/`
* `runtimeconfig_custom_endpoint` (`GOOGLE_RUNTIMECONFIG_CUSTOM_ENDPOINT`) - `https://runtimeconfig.googleapis.com/v1beta1/`
* `security_center_custom_endpoints` (`GOOGLE_SECURITY_CENTER_CUSTOM_ENDPOINT`) - `https://securitycenter.googleapis.com/v1/`
* `service_management_custom_endpoint` (`GOOGLE_SERVICE_MANAGEMENT_CUSTOM_ENDPOINT`) - `https://servicemanagement.googleapis.com/v1/`
* `service_networking_custom_endpoint` (`GOOGLE_SERVICE_NETWORKING_CUSTOM_ENDPOINT`) - `https://servicenetworking.googleapis.com/v1/`
* `service_usage_custom_endpoint` (`GOOGLE_SERVICE_USAGE_CUSTOM_ENDPOINT`) - `https://serviceusage.googleapis.com/v1/`
* `source_repo_custom_endpoint` (`GOOGLE_SOURCE_REPO_CUSTOM_ENDPOINT`) - `https://sourcerepo.googleapis.com/v1/`
* `spanner_custom_endpoint` (`GOOGLE_SPANNER_CUSTOM_ENDPOINT`) - `https://spanner.googleapis.com/v1/`
* `sql_custom_endpoint` (`GOOGLE_SQL_CUSTOM_ENDPOINT`) - `https://www.googleapis.com/sql/v1beta4/`
* `storage_custom_endpoint` (`GOOGLE_STORAGE_CUSTOM_ENDPOINT`) - `https://www.googleapis.com/storage/v1/`
* `storage_transfer_custom_endpoint` (`GOOGLE_STORAGE_TRANSFER_CUSTOM_ENDPOINT`) - `https://storagetransfer.googleapis.com/v1/`
* `tpu_custom_endpoint` (`GOOGLE_TPU_CUSTOM_ENDPOINT`) - `https://tpu.googleapis.com/v1/`

The following keys are available exclusively in the `google-beta` provider:

* `container_analysis_custom_endpoint` (`GOOGLE_CONTAINER_ANALYSIS_CUSTOM_ENDPOINT`) - `https://containeranalysis.googleapis.com/v1beta1/`
* `iap_custom_endpoint` (`GOOGLE_IAP_CUSTOM_ENDPOINT`) - `https://iap.googleapis.com/v1beta1/`
* `monitoring_custom_endpoint` (`GOOGLE_MONITORING_CUSTOM_ENDPOINT`) - `https://monitoring.googleapis.com/v3/`
* `security_scanner_custom_endpoint` (`GOOGLE_SECURITY_SCANNER_CUSTOM_ENDPOINT`) - `https://websecurityscanner.googleapis.com/v1beta/`

-> Note that some endpoints are a versioned variant of another. These exist in
cases where the `google` provider uses multiple distinct endpoints, and both
need to be set. Additionally, in `google-beta`, they'll often use the same value
as their versioned counterpart but that won't necessarily always be the case.

[OAuth 2.0 access token]: https://developers.google.com/identity/protocols/OAuth2
[service account key file]: https://cloud.google.com/iam/docs/creating-managing-service-account-keys
[manage key files using the Cloud Console]: https://console.cloud.google.com/apis/credentials/serviceaccountkey
[adc]: https://cloud.google.com/docs/authentication/production
[gce-service-account]: https://cloud.google.com/compute/docs/authentication
[gcloud adc]: https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login
[service accounts]: https://cloud.google.com/docs/authentication/getting-started
[GCE metadata]: https://cloud.google.com/docs/authentication/production#obtaining_credentials_on_compute_engine_kubernetes_engine_app_engine_flexible_environment_and_cloud_functions
[scopes]: https://developers.google.com/identity/protocols/googlescopes

---

* `batching` - (Optional) Controls batching for specific GCP request types
  where users have encountered quota or speed issues using `count` with
  resources that affect the same GCP resource (e.g. `google_project_service`). 
  It is not used for every resource/request type and can only group parallel
  similar calls for nodes at a similar traversal time in the graph during
  `terraform apply` (e.g. resources created using `count` that affect a single
  `project`). Thus, it is also bounded by the `terraform` 
  [`-parallelism`](https://www.terraform.io/docs/commands/apply.html#parallelism-n) 
  flag, as reducing the number of parallel calls will reduce the number of
  simultaneous requests being added to a batcher.

  ~> **NOTE** Most resources/GCP request do not have batching implemented (see
  below for requests which use batching) Batching is really only needed for
  resources where several requests are made at the same time to an underlying
  GCP resource protected by a fairly low default quota, or with very slow
  operations with slower eventual propagation. If you're not completely sure
  what you are doing, avoid setting custom batching configuration.

**So far, batching is implemented for**:

* enabling project services using `google_project_service`.

The `batching` block supports the following fields.

* `send_after` - (Optional) A duration string representing the amount of time
after which a request should be sent. Defaults to 10s. Should be a non-negative
integer or float string with a unit suffix, such as "300ms", "1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

* `disable_batching` - (Optional) Defaults to false. If true, disables global
batching and each request is sent normally.

---
* `request_timeout` - (Optional) A duration string controlling the amount of time
the provider should wait for a single HTTP request.  This will not adjust the
amount of time the provider will wait for a logical operation - use the resource
timeout blocks for that.  This will adjust only the amount of time that a single
synchronous request will wait for a response.  The default is 30 seconds, and
that should be a suitable value in most cases.  Many GCP APIs will cancel a
request if no response is forthcoming within 30 seconds in any event.  In
limited cases, such as DNS record set creation, there is a synchronous request
to create the resource.  This may help in those cases.


---

* `user_project_override` - (Optional) Defaults to false. If true, uses the
resource project for preconditions, quota, and billing, instead of the project
the credentials belong to. Not all resources support this- see the
documentation for each resource to learn whether it does.

When set to false, the project the credentials belong to will be billed for the
request, and quota / API enablement checks will be done against that project.
For service account credentials, this is the project the service account was
created in. For credentials that come from the gcloud tool, this is a project
owned by Google. In order to properly use credentials that come from gcloud
with Terraform, it is recommended to set this property to true.

When set to true, the caller must have `serviceusage.services.use` permission
on the resource project.
