---
subcategory: "Secret Manager"
description: |-
  List the Secret Manager Secrets.
---

# google_secret_manager_secrets

Use this data source to list the Secret Manager Secrets

## Example Usage 


```hcl
data "google_secret_manager_secrets" "secrets" {
}
```

## Argument Reference

The following arguments are supported:

* `project` - (optional) The ID of the project.

* `filter` - (optional) Filter string, adhering to the rules in [List-operation filtering](https://cloud.google.com/secret-manager/docs/filtering). List only secrets matching the filter. If filter is empty, all secrets are listed.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `secrets` - A list of secrets matching the filter. Structure is [defined below](#nested_secrets).

<a name="nested_secrets"></a>The `secrets` block supports:

* `replication` -
  The replication policy of the secret data attached to the Secret.
  Structure is [documented below](#nested_replication).

* `labels` - The labels assigned to this Secret.

* `annotations` - Custom metadata about the secret.

* `version_aliases` - Mapping from version alias to version name.

* `topics` -
  A list of up to 10 Pub/Sub topics to which messages are published when control plane operations are called on the secret or its versions.
  Structure is [documented below](#nested_topics).

* `expire_time` - Timestamp in UTC when the Secret is scheduled to expire.

* `create_time` - The time at which the Secret was created.

* `rotation` -
  The rotation time and period for a Secret.
  Structure is [documented below](#nested_rotation).

* `project` - The ID of the project in which the resource belongs.


<a name="nested_replication"></a>The `replication` block supports:

* `auto` -
  The Secret will automatically be replicated without any restrictions.
  Structure is [documented below](#nested_auto).

* `user_managed` -
  The Secret will be replicated to the regions specified by the user.
  Structure is [documented below](#nested_user_managed).


<a name="nested_auto"></a>The `auto` block supports:

* `customer_managed_encryption` -
  The customer-managed encryption configuration of the Secret.
  Structure is [documented below](#nested_customer_managed_encryption).

<a name="nested_customer_managed_encryption"></a>The `customer_managed_encryption` block supports:

* `kms_key_name` -
  The resource name of the Cloud KMS CryptoKey used to encrypt secret payloads.

<a name="nested_user_managed"></a>The `user_managed` block supports:

* `replicas` -
  The list of Replicas for this Secret.
  Structure is [documented below](#nested_replicas).

<a name="nested_replicas"></a>The `replicas` block supports:

* `location` -
  The canonical IDs of the location to replicate data.

* `customer_managed_encryption` -
  Customer Managed Encryption for the secret.
  Structure is [documented below](#nested_customer_managed_encryption_user_managed).

<a name="nested_customer_managed_encryption_user_managed"></a>The `customer_managed_encryption` block supports:

* `kms_key_name` -
  Describes the Cloud KMS encryption key that will be used to protect destination secret.

<a name="nested_topics"></a>The `topics` block supports:

* `name` - The resource name of the Pub/Sub topic that will be published to.

<a name="nested_rotation"></a>The `rotation` block supports:

* `next_rotation_time` - Timestamp in UTC at which the Secret is scheduled to rotate.

* `rotation_period` - The Duration between rotation notifications.
