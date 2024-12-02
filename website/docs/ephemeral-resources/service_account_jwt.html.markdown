---
subcategory: "Cloud Platform"
description: |-
  Produces an arbitrary self-signed JWT for service accounts
---

# google_service_account_jwt

This ephemeral resource provides a [self-signed JWT](https://cloud.google.com/iam/docs/create-short-lived-credentials-direct#sa-credentials-jwt).  Tokens issued from this ephemeral resource are typically used to call external services that accept JWTs for authentication.

## Example Usage

-> **Note:** If you run this example configuration you will be able to see ephemeral.google_service_account_jwt.foo in terraform plan and apply terminal output but you will not see it in state, as ephemeral resources are excluded from state. In future, when write-only attributes are added to resources in the Google provider, ephemeral resources such as google_service_account_jwt could be used to set field values when creating managed resources.

Note: in order to use the following, the caller must have _at least_ `roles/iam.serviceAccountTokenCreator` on the `target_service_account`.

```hcl
ephemeral "google_service_account_jwt" "foo" {
  target_service_account = "impersonated-account@project.iam.gserviceaccount.com"

  payload = jsonencode({
    foo: "bar",
    sub: "subject",
  })

  expires_in = 60
}
```

## Argument Reference

The following arguments are supported:

* `target_service_account` (Required) - The email of the service account that will sign the JWT.
* `payload` (Required) - The JSON-encoded JWT claims set to include in the self-signed JWT.
* `expires_in` (Optional) - Number of seconds until the JWT expires. If set and non-zero an `exp` claim will be added to the payload derived from the current timestamp plus expires_in seconds.
* `delegates` (Optional) - Delegate chain of approvals needed to perform full impersonation. Specify the fully qualified service account name.

## Attributes Reference

The following attribute is exported:

* `jwt` - The signed JWT containing the JWT Claims Set from the `payload`.
