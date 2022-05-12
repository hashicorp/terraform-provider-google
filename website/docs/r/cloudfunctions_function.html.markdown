---
subcategory: "Cloud Functions"
layout: "google"
page_title: "Google: google_cloudfunctions_function"
sidebar_current: "docs-google-cloudfunctions-function"
description: |-
  Creates a new Cloud Function.
---

# google\_cloudfunctions\_function

Creates a new Cloud Function. For more information see:

* [API documentation](https://cloud.google.com/functions/docs/reference/rest/v1/projects.locations.functions)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/functions/docs)


~> **Warning:** As of November 1, 2019, newly created Functions are
private-by-default and will require [appropriate IAM permissions](https://cloud.google.com/functions/docs/reference/iam/roles)
to be invoked. See below examples for how to set up the appropriate permissions,
or view the [Cloud Functions IAM resources](/docs/providers/google/r/cloudfunctions_cloud_function_iam.html)
for Cloud Functions.

## Example Usage - Public Function

```hcl
resource "google_storage_bucket" "bucket" {
  name     = "test-bucket"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./path/to/zip/file/which/contains/code"
}

resource "google_cloudfunctions_function" "function" {
  name        = "function-test"
  description = "My function"
  runtime     = "nodejs16"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  entry_point           = "helloGET"
}

# IAM entry for all users to invoke the function
resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = google_cloudfunctions_function.function.project
  region         = google_cloudfunctions_function.function.region
  cloud_function = google_cloudfunctions_function.function.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}
```

## Example Usage - Single User

```hcl
resource "google_storage_bucket" "bucket" {
  name     = "test-bucket"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./path/to/zip/file/which/contains/code"
}

resource "google_cloudfunctions_function" "function" {
  name        = "function-test"
  description = "My function"
  runtime     = "nodejs16"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  timeout               = 60
  entry_point           = "helloGET"
  labels = {
    my-label = "my-label-value"
  }

  environment_variables = {
    MY_ENV_VAR = "my-env-var-value"
  }
}

# IAM entry for a single user to invoke the function
resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = google_cloudfunctions_function.function.project
  region         = google_cloudfunctions_function.function.region
  cloud_function = google_cloudfunctions_function.function.name

  role   = "roles/cloudfunctions.invoker"
  member = "user:myFunctionInvoker@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A user-defined name of the function. Function names must be unique globally.

* `runtime` - (Required) The runtime in which the function is going to run.
Eg. `"nodejs16"`, `"python39"`, `"dotnet3"`, `"go116"`, `"java11"`, `"ruby30"`, `"php74"`, etc. Check the [official doc](https://cloud.google.com/functions/docs/concepts/exec#runtimes) for the up-to-date list.

- - -

* `description` - (Optional) Description of the function.

* `available_memory_mb` - (Optional) Memory (in MB), available to the function. Default value is `256`. Possible values include `128`, `256`, `512`, `1024`, etc.

* `timeout` - (Optional) Timeout (in seconds) for the function. Default value is 60 seconds. Cannot be more than 540 seconds.

* `entry_point` - (Optional) Name of the function that will be executed when the Google Cloud Function is triggered.

* `event_trigger` - (Optional) A source that fires events in response to a condition in another service. Structure is [documented below](#nested_event_trigger). Cannot be used with `trigger_http`.

* `trigger_http` - (Optional) Boolean variable. Any HTTP request (of a supported type) to the endpoint will trigger function execution. Supported HTTP request types are: POST, PUT, GET, DELETE, and OPTIONS. Endpoint is returned as `https_trigger_url`. Cannot be used with `event_trigger`.

* `ingress_settings` - (Optional) String value that controls what traffic can reach the function. Allowed values are `ALLOW_ALL`, `ALLOW_INTERNAL_AND_GCLB` and `ALLOW_INTERNAL_ONLY`. Check [ingress documentation](https://cloud.google.com/functions/docs/networking/network-settings#ingress_settings) to see the impact of each settings value. Changes to this field will recreate the cloud function.

* `labels` - (Optional) A set of key/value label pairs to assign to the function. Label keys must follow the requirements at https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements.

* `service_account_email` - (Optional) If provided, the self-provided service account to run the function with.

* `environment_variables` - (Optional) A set of key/value environment variable pairs to assign to the function.

* `build_environment_variables` - (Optional) A set of key/value environment variable pairs available during build time.

* `vpc_connector` - (Optional) The VPC Network Connector that this cloud function can connect to. It should be set up as fully-qualified URI. The format of this field is `projects/*/locations/*/connectors/*`.

* `vpc_connector_egress_settings` - (Optional) The egress settings for the connector, controlling what traffic is diverted through it. Allowed values are `ALL_TRAFFIC` and `PRIVATE_RANGES_ONLY`. Defaults to `PRIVATE_RANGES_ONLY`. If unset, this field preserves the previously set value.

* `source_archive_bucket` - (Optional) The GCS bucket containing the zip archive which contains the function.

* `source_archive_object` - (Optional) The source archive object (file) in archive bucket.

* `source_repository` - (Optional) Represents parameters related to source repository where a function is hosted.
  Cannot be set alongside `source_archive_bucket` or `source_archive_object`. Structure is [documented below](#nested_source_repository). It must match the pattern `projects/{project}/locations/{location}/repositories/{repository}`.* 

* `docker_repository` - (Optional) User managed repository created in Artifact Registry optionally with a customer managed encryption key. If specified, deployments will use Artifact Registry. This is the repository to which the function docker image will be pushed after it is built by Cloud Build. If unspecified, Container Registry will be used by default, unless specified otherwise by other means.

* `kms_key_name` - (Optional) Resource name of a KMS crypto key (managed by the user) used to encrypt/decrypt function resources. It must match the pattern `projects/{project}/locations/{location}/keyRings/{key_ring}/cryptoKeys/{crypto_key}`.
  If specified, you must also provide an artifact registry repository using the `docker_repository` field that was created with the same KMS crypto key. Before deploying, please complete all pre-requisites described in https://cloud.google.com/functions/docs/securing/cmek#granting_service_accounts_access_to_the_key

* `max_instances` - (Optional) The limit on the maximum number of function instances that may coexist at a given time.

* `min_instances` - (Optional) The limit on the minimum number of function instances that may coexist at a given time.

* `secret_environment_variables` - (Optional) Secret environment variables configuration. Structure is [documented below](#nested_secret_environment_variables).

* `secret_volumes` - (Optional) Secret volumes configuration. Structure is [documented below](#nested_secret_volumes).

<a name="nested_event_trigger"></a>The `event_trigger` block supports:

* `event_type` - (Required) The type of event to observe. For example: `"google.storage.object.finalize"`.
See the documentation on [calling Cloud Functions](https://cloud.google.com/functions/docs/calling/) for a
full reference of accepted triggers.

* `resource` - (Required) Required. The name or partial URI of the resource from
which to observe events. For example, `"myBucket"` or `"projects/my-project/topics/my-topic"`

* `failure_policy` - (Optional) Specifies policy for failed executions. Structure is [documented below](#nested_failure_policy).

<a name="nested_failure_policy"></a>The `failure_policy` block supports:

* `retry` - (Required) Whether the function should be retried on failure. Defaults to `false`.

<a name="nested_source_repository"></a>The `source_repository` block supports:

* `url` - (Required) The URL pointing to the hosted repository where the function is defined. There are supported Cloud Source Repository URLs in the following formats:

    * To refer to a specific commit: `https://source.developers.google.com/projects/*/repos/*/revisions/*/paths/*`
    * To refer to a moveable alias (branch): `https://source.developers.google.com/projects/*/repos/*/moveable-aliases/*/paths/*`. To refer to HEAD, use the `master` moveable alias.
    * To refer to a specific fixed alias (tag): `https://source.developers.google.com/projects/*/repos/*/fixed-aliases/*/paths/*`

<a name="nested_secret_environment_variables"></a>The `secret_environment_variables` block supports:

* `key` - (Required) Name of the environment variable.

* `project_id` - (Optional) Project identifier (due to a known limitation, only project number is supported by this field) of the project that contains the secret. If not set, it will be populated with the function's project, assuming that the secret exists in the same project as of the function.

* `secret` - (Required) ID of the secret in secret manager (not the full resource name).

* `version` - (Required) Version of the secret (version number or the string "latest"). It is recommended to use a numeric version for secret environment variables as any updates to the secret value is not reflected until new clones start.

<a name="nested_secret_volumes"></a>The `secret_volumes` block supports:

* `mount_path` - (Required) The path within the container to mount the secret volume. For example, setting the mount_path as "/etc/secrets" would mount the secret value files under the "/etc/secrets" directory. This directory will also be completely shadowed and unavailable to mount any other secrets. Recommended mount paths: "/etc/secrets" Restricted mount paths: "/cloudsql", "/dev/log", "/pod", "/proc", "/var/log".

* `project_id` - (Optional) Project identifier (due to a known limitation, only project number is supported by this field) of the project that contains the secret. If not set, it will be populated with the function's project, assuming that the secret exists in the same project as of the function.

* `secret` - (Required) ID of the secret in secret manager (not the full resource name).

* `versions` - (Optional) List of secret versions to mount for this secret. If empty, the "latest" version of the secret will be made available in a file named after the secret under the mount point. Structure is [documented below](#nested_nested_versions).

<a name="nested_versions"></a>The `versions` block supports:

* `path` - (Required) Relative path of the file under the mount path where the secret value for this version will be fetched and made available. For example, setting the mount_path as "/etc/secrets" and path as "/secret_foo" would mount the secret value file at "/etc/secrets/secret_foo".

* `version` - (Required) Version of the secret (version number or the string "latest"). It is preferable to use "latest" version with secret volumes as secret value changes are reflected immediately.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `{{name}}`

* `https_trigger_url` - URL which triggers function execution. Returned only if `trigger_http` is used.

* `source_repository.0.deployed_url` - The URL pointing to the hosted repository where the function was defined at the time of deployment.

* `project` - Project of the function. If it is not provided, the provider project is used.

* `region` - Region of function. If it is not provided, the provider region is used.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 5 minutes.
- `update` - Default is 5 minutes.
- `delete` - Default is 5 minutes.

## Import

Functions can be imported using the `name` or `{{project}}/{{region}}/name`, e.g.

```
$ terraform import google_cloudfunctions_function.default function-test
$ terraform import google_cloudfunctions_function.default {{project}}/{{region}}/function-test
```
