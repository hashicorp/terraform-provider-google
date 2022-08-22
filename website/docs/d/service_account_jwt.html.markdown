---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_service_account_jwt"
sidebar_current: "docs-google-service-account-jwt"
description: |-
  Produces an arbitrary self-signed JWT for service accounts
---

# google\_service\_account\_jwt

This data source provides a [self-signed JWT](https://cloud.google.com/iam/docs/create-short-lived-credentials-direct#sa-credentials-jwt).  Tokens issued from this data source are typically used to call external services that accept JWTs for authentication.

## Example Usage

Note: in order to use the following, the caller must have _at least_ `roles/iam.serviceAccountTokenCreator` on the `target_service_account`.

```hcl
data "google_service_account_jwt" "foo" {
  target_service_account = "impersonated-account@project.iam.gserviceaccount.com"

  payload = jsonencode({
    foo: "bar",
    sub: "subject",
  })
}

output "jwt" {
  value = data.google_service_account_jwt.foo.jwt
}
```

## Argument Reference

The following arguments are supported:

* `target_service_account` (Required) - The email of the service account that will sign the JWT.
* `payload` (Required) - The JSON-encoded JWT claims set to include in the self-signed JWT.
* `delegates` (Optional) - Delegate chain of approvals needed to perform full impersonation. Specify the fully qualified service account name.

## Attributes Reference

The following attribute is exported:

* `jwt` - The signed JWT containing the JWT Claims Set from the `payload`.
