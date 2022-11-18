---
subcategory: "Cloud (Stackdriver) Logging"
page_title: "Google: google_logging_project_cmek_settings"
description: |-
  Describes the customer-managed encryption key (CMEK) settings associated with a project.
---

# google\_logging\_project\_cmek\_settings

Describes the customer-managed encryption key (CMEK) settings associated with a project.

To get more information about Service, see:

* [API documentation](https://cloud.google.com/logging/docs/reference/v2/rest/v2/projects/getCmekSettings)
* [Enable CMEK](https://cloud.google.com/logging/docs/routing/managed-encryption-storage#enable).

## Example Usage - Logging Project Cmek Settings Basic


```hcl
data "google_logging_project_cmek_settings" "cmek_settings" {
	project = "my-project-name"
}
```

## Argument Reference

The following arguments are supported:



- - -

* `project` - (Required) The ID of the project.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/cmekSettings`

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