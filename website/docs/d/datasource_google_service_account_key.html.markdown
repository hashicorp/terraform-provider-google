---
layout: "google"
page_title: "Google: google_service_account_key"
sidebar_current: "docs-google-datasource-service-account-key"
description: |-
  Get a Google Cloud Platform service account Public Key
---

# google\_service\_account\_key

Get service account public key. For more information, see [the official documentation](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) and [API](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts.keys/get).


## Example Usage

```hcl
data "google_service_account" "myaccount" {
  account_id = "myaccount"
}

data "google_service_account_key" "mykey" {
  service_account_id = "${data.google_service_account.myaccount.name}"
  public_key_type = "TYPE_X509_PEM_FILE"
}

output "mykey_public_key" {
  value = "${data.google_service_account_key.mykey.public_key}"
}
```

## Argument Reference

The following arguments are supported:

* `service_account_id` - (Required) The Service account id of the Key Pair. This can be a string in the format
`{ACCOUNT}` or `projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}`, where `{ACCOUNT}` is the email address or
unique id of the service account. If the `{ACCOUNT}` syntax is used, the project will be inferred from the account.

* `project` - (Optional) The ID of the project that the service account will be created in.
    Defaults to the provider project configuration.

* `public_key_type` (Optional) The output format of the public key requested. X509_PEM is the default output format.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `name` - The name used for this key pair

* `public_key` - The public key, base64 encoded
