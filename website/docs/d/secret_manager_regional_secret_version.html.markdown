---
subcategory: "Secret Manager"
description: |-
  Get a Secret Manager regional secret's version.
---

# google_secret_manager_regional_secret_version

Get the value and metadata from a Secret Manager regional secret version. For more information see the [official documentation](https://cloud.google.com/secret-manager/docs/regional-secrets-overview) and [API](https://cloud.google.com/secret-manager/docs/reference/rest/v1/projects.secrets.versions). If you don't need the metadata (i.e., if you want to use a more limited role to access the regional secret version only), see also the [google_secret_manager_regional_secret_version_access](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/secret_manager_regional_secret_version_access) datasource.

## Example Usage

```hcl
data "google_secret_manager_regional_secret_version" "basic" {
  secret   = "my-secret"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project to get the secret version for. If it
    is not provided, the provider project is used.

* `secret` - (Required) The regional secret to get the secret version for.
    This can be either the reference of the regional secret as in `projects/{{project}}/locations/{{location}}/secrets/{{secret_id}}` or only the name of the regional secret as in `{{secret_id}}`. If only the name of the regional secret is provided, the location must also be provided.

* `location` - (Optional) Location of Secret Manager regional secret resource.
    It must be provided when the `secret` field provided consists of only the name of the regional secret.

* `version` - (Optional) The version of the regional secret to get. If it
    is not provided, the latest version is retrieved.

## Attributes Reference

The following attributes are exported:

* `secret_data` - The secret data. No larger than 64KiB.

* `name` - The resource name of the regional SecretVersion. Format:
  `projects/{{project}}/locations/{{location}}/secrets/{{secret_id}}/versions/{{version}}`

* `create_time` - The time at which the regional secret was created.

* `destroy_time` - The time at which the regional secret was destroyed. Only present if state is DESTROYED.

* `enabled` - True if the current state of the regional SecretVersion is enabled.

* `customer_managed_encryption` - The customer-managed encryption configuration of the regional secret. Structure is [documented below](#nested_customer_managed_encryption).

<a name="nested_customer_managed_encryption"></a>The `customer_managed_encryption` block contains:

* `kms_key_version_name` - The resource name of the Cloud KMS CryptoKey used to encrypt secret payloads.
