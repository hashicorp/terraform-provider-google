---
subcategory: "Cloud (Stackdriver) Logging"
description: |-
  Describes the settings associated with a folder.
---

# google_logging_folder_settings

Describes the settings associated with a folder.

To get more information about LoggingFolderSettings, see:

* [API documentation](https://cloud.google.com/logging/docs/reference/v2/rest/v2/folders/getSettings)
* [Configure default settings for organizations and folders](https://cloud.google.com/logging/docs/default-settings).

## Example Usage - Logging Folder Settings Basic

```hcl
data "google_logging_folder_settings" "settings" {
  folder = "my-folder-name"
}
```

## Argument Reference

The following arguments are supported:

- - -

* `folder` - (Required) The ID of the folder for which to retrieve settings.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `folders/{{folder}}/settings`

* `name` - The resource name of the settings.

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

* `kms_service_account_id` - The service account associated with a project for which CMEK will apply.
Before enabling CMEK for a logging bucket, you must first assign the cloudkms.cryptoKeyEncrypterDecrypter role to the service account associated with the project for which CMEK will apply. See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.

* `logging_service_account_id` - The service account for the given container. Sinks use this service account as their writerIdentity if no custom service account is provided.

* `disable_default_sink` -  If set to true, the _Default sink in newly created projects and folders will created in a disabled state. This can be used to automatically disable log storage if there is already an aggregated sink configured in the hierarchy. The _Default sink can be re-enabled manually if needed.

* `storage_location` -  The storage location that Cloud Logging will use to create new resources when a location is needed but not explicitly provided.
