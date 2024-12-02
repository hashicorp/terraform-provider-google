---
subcategory: "Cloud Platform"
description: |-
  Produces OpenID Connect token for service accounts
---

# google_service_account_id_token

This ephemeral resource provides a Google OpenID Connect (`oidc`) `id_token`.  Tokens issued from this ephemeral resource are typically used to call external services that accept OIDC tokens for authentication (e.g. [Google Cloud Run](https://cloud.google.com/run/docs/authenticating/service-to-service)).

For more information see
[OpenID Connect](https://openid.net/specs/openid-connect-core-1_0.html#IDToken).

## Example Usage - ServiceAccount JSON credential file.

-> **Note:** If you run this example configuration you will be able to see ephemeral.google_service_account_id_token.oidc in terraform plan and apply terminal output but you will not see it in state, as ephemeral resources are excluded from state. In future, when write-only attributes are added to resources in the Google provider, ephemeral resources such as google_service_account_id_token could be used to set field values when creating managed resources.

  `google_service_account_id_token` will use the configured [provider credentials](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#credentials-1)

  ```hcl
  ephemeral "google_service_account_id_token" "oidc" {
    target_audience = "https://foo.bar/"
  }
  ```

## Example Usage - Service Account Impersonation.

-> **Note:** If you run this example configuration you will be able to see ephemeral.google_service_account_id_token.oidc in terraform plan and apply terminal output but you will not see it in state, as ephemeral resources are excluded from state. In future, when write-only attributes are added to resources in the Google provider, ephemeral resources such as google_service_account_id_token could be used to set field values when creating managed resources.

  Ephemeral resource `google_service_account_id_token` will use background impersonated credentials provided by [google_service_account_access_token](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/service_account_access_token).

  Note: to use the following, you must grant `target_service_account` the
  `roles/iam.serviceAccountTokenCreator` role on itself.

  ```hcl
  data "google_service_account_access_token" "impersonated" {
    provider = google
    target_service_account = "impersonated-account@project.iam.gserviceaccount.com"
    delegates = []
    scopes = ["userinfo-email", "cloud-platform"]
    lifetime = "300s"
  }

  provider "google" {
    alias  = "impersonated"
    access_token = data.google_service_account_access_token.impersonated.access_token
  }

  ephemeral "google_service_account_id_token" "oidc" {
    provider = google.impersonated
    target_service_account = "impersonated-account@project.iam.gserviceaccount.com"
    delegates = []
    include_email = true
    target_audience = "https://foo.bar/"
  }

  ```

## Argument Reference

The following arguments are supported:

* `target_audience` (Required) - The audience claim for the `id_token`.
* `target_service_account` (Optional) - The email of the service account being impersonated.  Used only when using impersonation mode.
* `delegates` (Optional) - Delegate chain of approvals needed to perform full impersonation. Specify the fully qualified service account name.   Used only when using impersonation mode.
* `include_email` (Optional) Include the verified email in the claim. Used only when using impersonation mode.

## Attributes Reference

The following attribute is exported:

* `id_token` - The `id_token` representing the new generated identity.
