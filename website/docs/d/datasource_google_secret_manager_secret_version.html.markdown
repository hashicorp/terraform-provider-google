---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_secret_manager_secret_version"
sidebar_current: "docs-google-datasource-secret-manager-secret-version"
description: |-
  Get a Secret Manager secret's version.
---

# google\_secret\_manager\_secret\_version

Get a Secret Manager secret's version. For more information see the [official documentation](https://cloud.google.com/secret-manager/docs/) and [API](https://cloud.google.com/secret-manager/docs/reference/rest/v1beta1/projects.secrets.versions).

## Example Usage

```hcl
data "google_secret_manager_secret_version" "basic" {
  provider = google-beta
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