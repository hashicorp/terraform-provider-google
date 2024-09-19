---
subcategory: "Secret Manager"
description: |-
  List the Secret Manager Regional Secrets.
---

# google_secret_manager_regional_secrets

Use this data source to list the Secret Manager Regional Secrets.

## Example Usage 

```hcl
data "google_secret_manager_regional_secrets" "secrets" {
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (optional) The ID of the project.

* `filter` - (optional) Filter string, adhering to the rules in [List-operation filtering](https://cloud.google.com/secret-manager/docs/filtering). List only secrets matching the filter. If filter is empty, all regional secrets are listed from the specified location.

* `location` - (Required) The location of the regional secret.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `secrets` - A list of regional secrets present in the specified location and matching the filter. Structure is [defined below](#nested_secrets).

<a name="nested_secrets"></a>The `secrets` block supports:

* `labels` - The labels assigned to this regional secret.

* `annotations` - Custom metadata about the regional secret.

* `version_aliases` - Mapping from version alias to version name.

* `topics` -
  A list of up to 10 Pub/Sub topics to which messages are published when control plane operations are called on the regional secret or its versions.
  Structure is [documented below](#nested_topics).

* `expire_time` - Timestamp in UTC when the regional secret is scheduled to expire.

* `create_time` - The time at which the regional secret was created.

* `rotation` -
  The rotation time and period for a regional secret.
  Structure is [documented below](#nested_rotation).

* `project` - The ID of the project in which the resource belongs.

* `location` - The location in which the resource belongs.

* `secret_id` - The unique name of the resource.

* `name` - The resource name of the regional secret. Format: `projects/{{project}}/locations/{{location}}/secrets/{{secret_id}}`

* `version_destroy_ttl` - The version destroy ttl for the regional secret version.

* `customer_managed_encryption` -
  Customer Managed Encryption for the regional secret.
  Structure is [documented below](#nested_customer_managed_encryption_user_managed).

<a name="nested_topics"></a>The `topics` block supports:

* `name` - The resource name of the Pub/Sub topic that will be published to.

<a name="nested_rotation"></a>The `rotation` block supports:

* `next_rotation_time` - Timestamp in UTC at which the secret is scheduled to rotate.

* `rotation_period` - The Duration between rotation notifications.

<a name="nested_customer_managed_encryption_user_managed"></a>The `customer_managed_encryption` block supports:

* `kms_key_name` -
  Describes the Cloud KMS encryption key that will be used to protect destination secret.
