---
subcategory: "Secret Manager"
description: |-
  Get a Secret Manager secret's version.
---

# google_secret_manager_secret_version

Get the value and metadata from a Secret Manager secret version. For more information see the [official documentation](https://cloud.google.com/secret-manager/docs/) and [API](https://cloud.google.com/secret-manager/docs/reference/rest/v1/projects.secrets.versions). If you don't need the metadata (i.e., if you want to use a more limited role to access the secret version only), see also the [google_secret_manager_secret_version_access](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/secret_manager_secret_version_access) datasource.

## Example Usage

```hcl
data "google_secret_manager_secret_version" "basic" {
  secret = "my-secret"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project to get the secret version for. If it
    is not provided, the provider project is used.

* `secret` - (Required) The secret to get the secret version for.

* `version` - (Optional) The version of the secret to get. If it
    is not provided, the latest version is retrieved.


## Attributes Reference

The following attributes are exported:

* `secret_data` - The secret data. No larger than 64KiB.

* `name` - The resource name of the SecretVersion. Format:
  `projects/{{project}}/secrets/{{secret_id}}/versions/{{version}}`

* `create_time` - The time at which the Secret was created.

* `destroy_time` - The time at which the Secret was destroyed. Only present if state is DESTROYED.

* `enabled` - True if the current state of the SecretVersion is enabled.
