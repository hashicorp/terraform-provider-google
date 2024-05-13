---
subcategory: "Secret Manager"
page_title: "Google: google_secret_manager_secret_version_access"
description: |-
  Get a payload of Secret Manager secret's version.
---

# google_secret_manager_secret_version_access

Get the value from a Secret Manager secret version. This is similar to the [google_secret_manager_secret_version](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/secret_manager_secret_version) datasource, but it only requires the [Secret Manager Secret Accessor](https://cloud.google.com/secret-manager/docs/access-control#secretmanager.secretAccessor) role. For more information see the [official documentation](https://cloud.google.com/secret-manager/docs/) and [API](https://cloud.google.com/secret-manager/docs/reference/rest/v1/projects.secrets.versions/access).

## Example Usage

```hcl
data "google_secret_manager_secret_version_access" "basic" {
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
