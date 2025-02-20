---
page_title: "Use write-only attributes in the Google Cloud provider"
description: |-
  How to use write-only attributes in the Google Cloud provider
---

# Write-only attributes in the Google Cloud provider

Write-only attributes allow users to access and use data in their configurations without that data being stored in Terraform state.

Write-only attributes are available in Terraform v1.11 and later. For more information, see the [official HashiCorp documentation for Write-only Attributes](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/write-only-arguments).

To mark the launch of write-only attributes, the Google Cloud provider has added the following write-only attributes:
- [`google_compute_disk: disk_encryption_key.raw_key_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_disk.html#raw_key_wo)
- [`google_compute_disk: disk_encryption_key.rsa_encrypted_key_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_disk.html#rsa_encryption_key_wo)
- [`google_compute_region_disk: disk_encryption_key.raw_key_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_region_disk.html#raw_key_wo)
- [`google_sql_user: password_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_user#password-1)
- [`google_secret_manager_secret_version: secret_data_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/secret_manager_secret_version#secret_data_wo)
- [`google_bigquery_data_transfer_config: sensitive_params.secret_access_key_wo`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/bigquery_data_transfer_config#secret_access_key_wo)

These were chosen due to them being marked as sensitive already in the provider. Although sensitive attributes do not appear in `terraform plan`, they are still stored in the Terraform state. Write-only attributes allow users to access and use data in their configurations without that data being stored in Terraform state.

## Use the Google Cloud provider's new write-only attributes

The following sections show how to use the new write-only attributes in the Google Cloud provider.

### Applying a write-only attribute

The following example shows how to apply a write-only attribute. All write-only attributes are marked with the `wo` suffix and can not be used with the `sensitive` attribute.

```hcl
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret-version"

  labels = {
    label = "my-label"
  }

  replication {
    auto {}
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data_wo = "secret-data"
  secret_data_wo_version = 1
  enabled = true
}
```

During `terraform plan` you will see that the write-only attribute is marked appropriately:

```
  # google_secret_manager_secret_version.secret-version-basic will be created
  + resource "google_secret_manager_secret_version" "secret-version-basic" {
      + create_time            = (known after apply)
      + deletion_policy        = "DELETE"
      + destroy_time           = (known after apply)
      + enabled                = true
      + id                     = (known after apply)
      + is_secret_data_base64  = false
      + name                   = (known after apply)
      + secret                 = (known after apply)
      + secret_data_wo         = (write-only attribute)
      + secret_data_wo_version = 1
      + version                = (known after apply)
    }
```

Upon `terrform apply` you will see in `terraform.tfstate` that the write-only attribute from the configuration is not reflected in the state:

```hcl
...
    {
      "mode": "managed",
      "type": "google_secret_manager_secret_version",
      "name": "secret-version-basic",
      "provider": "provider[\"registry.terraform.io/hashicorp/google\"]",
      "instances": [
        {
          "schema_version": 0,
          "attributes": {
            "create_time": "2025-02-19T05:56:20.942452Z",
            "deletion_policy": "DELETE",
            "destroy_time": "",
            "enabled": true,
            "id": "projects/871647908372/secrets/secret-version/versions/1",
            "is_secret_data_base64": false,
            "name": "projects/871647908372/secrets/secret-version/versions/1",
            "secret": "projects/hc-terraform-testing/secrets/secret-version",
            "secret_data": null,
            "secret_data_wo": null,
            "secret_data_wo_version": 1,
            "timeouts": null,
            "version": "1"
          },
          ...
        }
      ]
    }
```

Any value that is set for a write-only attribute is nulled out before the RPC response is sent to Terraform.

### Updating write-only attributes

Since write-only attributes are not stored in the Terraform state, they cannot be updated by just changing the value in the configuration due to the attribute being nulled out.

In order to update a write-only attribute we must change the write-only attribute's version.

```hcl
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "secret-version"

  labels = {
    label = "my-label"
  }

  replication {
    auto {}
  }
}


resource "google_secret_manager_secret_version" "secret-version-basic" {
  secret = google_secret_manager_secret.secret-basic.id
  secret_data_wo = "new-secret-data"
  secret_data_wo_version = 2
  enabled = true
}
```

A `terraform apply` of this configuration will allow you to update the write-only attribute as it will destroy and recreate the resource.

```hcl
Terraform used the selected providers to generate the following execution plan. Resource actions
are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # google_secret_manager_secret_version.secret-version-basic must be replaced
-/+ resource "google_secret_manager_secret_version" "secret-version-basic" {
      ~ create_time            = "2025-02-19T05:56:20.942452Z" -> (known after apply)
      + destroy_time           = (known after apply)
      ~ id                     = "projects/871647908372/secrets/secret-version/versions/1" -> (known after apply)
      ~ name                   = "projects/871647908372/secrets/secret-version/versions/1" -> (known after apply)
      ~ secret_data_wo_version = 1 -> 2 # forces replacement
      ~ version                = "1" -> (known after apply)
        # (5 unchanged attributes hidden)
    }

Plan: 1 to add, 0 to change, 1 to destroy.
```
