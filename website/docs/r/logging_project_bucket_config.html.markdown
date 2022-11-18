---
subcategory: "Cloud (Stackdriver) Logging"
page_title: "Google: google_logging_project_bucket_config"
description: |-
  Manages a project-level logging bucket config.
---

# google\_logging\_project\_bucket\_config

Manages a project-level logging bucket config. For more information see
[the official logging documentation](https://cloud.google.com/logging/docs/) and
[Storing Logs](https://cloud.google.com/logging/docs/storage).

~> **Note:** Logging buckets are automatically created for a given folder, project, organization, billingAccount and cannot be deleted. Creating a resource of this type will acquire and update the resource that already exists at the desired location. These buckets cannot be removed so deleting this resource will remove the bucket config from your terraform state but will leave the logging bucket unchanged. The buckets that are currently automatically created are "_Default" and "_Required".

## Example Usage

```hcl
resource "google_project" "default" {
	project_id = "your-project-id"
	name       = "your-project-id"
	org_id     = "123456789"
}

resource "google_logging_project_bucket_config" "basic" {
	project    = google_project.default.id
	location  = "global"
	retention_days = 30
	bucket_id = "_Default"
}
```

Create logging bucket with customId

```hcl
resource "google_logging_project_bucket_config" "basic" {
	project    = "project_id"
	location  = "global"
	retention_days = 30
	bucket_id = "custom-bucket"
}
```

Create logging bucket with customId and cmekSettings

```hcl
data "google_logging_project_cmek_settings" "cmek_settings" {
	project = "project_id"
}

resource "google_kms_key_ring" "keyring" {
	name     = "keyring-example"
	location = "us-central1"
}

resource "google_kms_crypto_key" "key" {
	name            = "crypto-key-example"
	key_ring        = google_kms_key_ring.keyring.id
	rotation_period = "100000s"
}

resource "google_kms_crypto_key_iam_binding" "crypto_key_binding" {
	crypto_key_id = google_kms_crypto_key.key.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	
	members = [
		"serviceAccount:${data.google_logging_project_cmek_settings.cmek_settings.service_account_id}",
	]
}

resource "google_logging_project_bucket_config" "example-project-bucket-cmek-settings" {
	project        = "project_id"
	location       = "us-central1"
	retention_days = 30
	bucket_id      = "custom-bucket"

	cmek_settings {
		kms_key_name = google_kms_crypto_key.key.id
	}

	depends_on   = [google_kms_crypto_key_iam_binding.crypto_key_binding]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The parent resource that contains the logging bucket.

* `location` - (Required) The location of the bucket.

* `bucket_id` - (Required) The name of the logging bucket. Logging automatically creates two log buckets: `_Required` and `_Default`.

* `description` - (Optional) Describes this bucket.

* `retention_days` - (Optional) Logs will be retained by default for this amount of time, after which they will automatically be deleted. The minimum retention period is 1 day. If this value is set to zero at bucket creation time, the default time of 30 days will be used.

* `cmek_settings` - (Optional) The CMEK settings of the log bucket. If present, new log entries written to this log bucket are encrypted using the CMEK key provided in this configuration. If a log bucket has CMEK settings, the CMEK settings cannot be disabled later by updating the log bucket. Changing the KMS key is allowed. Structure is [documented below](#nested_cmek_settings).


<a name="nested_cmek_settings"></a>The `cmek_settings` block supports:

* `name` - The resource name of the CMEK settings.

* `kms_key_name` - The resource name for the configured Cloud KMS key.
KMS key name format:
`'projects/[PROJECT_ID]/locations/[LOCATION]/keyRings/[KEYRING]/cryptoKeys/[KEY]'`
To enable CMEK for the bucket, set this field to a valid kmsKeyName for which the associated service account has the required cloudkms.cryptoKeyEncrypterDecrypter roles assigned for the key.
The Cloud KMS key used by the bucket can be updated by changing the kmsKeyName to a new valid key name. Encryption operations that are in progress will be completed with the key that was in use when they started. Decryption operations will be completed using the key that was used at the time of encryption unless access to that key has been revoked.
See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.

* `kms_key_version_name` - The CryptoKeyVersion resource name for the configured Cloud KMS key.
KMS key name format:
`'projects/[PROJECT_ID]/locations/[LOCATION]/keyRings/[KEYRING]/cryptoKeys/[KEY]/cryptoKeyVersions/[VERSION]'`
For example:
"projects/my-project/locations/us-central1/keyRings/my-ring/cryptoKeys/my-key/cryptoKeyVersions/1"
This is a read-only field used to convey the specific configured CryptoKeyVersion of kms_key that has been configured. It will be populated in cases where the CMEK settings are bound to a single key version.

* `service_account_id` - The service account associated with a project for which CMEK will apply.
Before enabling CMEK for a logging bucket, you must first assign the cloudkms.cryptoKeyEncrypterDecrypter role to the service account associated with the project for which CMEK will apply. Use [v2.getCmekSettings](https://cloud.google.com/logging/docs/reference/v2/rest/v2/TopLevel/getCmekSettings#google.logging.v2.ConfigServiceV2.GetCmekSettings) to obtain the service account ID.
See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/buckets/{{bucket_id}}`

* `name` -  The resource name of the bucket. For example: "projects/my-project-id/locations/my-location/buckets/my-bucket-id"

* `lifecycle_state` -  The bucket's lifecycle such as active or deleted. See [LifecycleState](https://cloud.google.com/logging/docs/reference/v2/rest/v2/billingAccounts.buckets#LogBucket.LifecycleState).

## Import

This resource can be imported using the following format:

```
$ terraform import google_logging_project_bucket_config.default projects/{{project}}/locations/{{location}}/buckets/{{bucket_id}}
```
