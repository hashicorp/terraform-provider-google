---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_service_account_id_token"
sidebar_current: "docs-google-service-account-id-token"
description: |-
  Produces OpenID Connect token for service accounts
---

# google\_service\_account\id\_token

This data source provides a google OpenID Connect (`oidc`) `id_token` 

For more information see
[OpenID Connect](https://openid.net/specs/openid-connect-core-1_0.html#IDToken) as well as [Authenticating using Google OpenID Connect Tokens](https://github.com/salrashid123/google_id_token)

## Example Usage

There are several ways to use this data source to get an `id_token`:

* Using ServiceAccount JSON credential file.
  In this mode, a the datasource will the service account provided by `GOOGLE_APPLICATION_CREDENTIALS` or `GOOGLE_CLOUD_KEYFILE_JSON` environment variable. 
  
  Sample configuration:
  ```yaml
  provider google {}

  data "google_service_account_id_token" oidc {
    target_audience = "https://foo.bar/"
  }

  output "oidc_token" {
    value = data.google_service_account_id_token.oidc.id_token
  }
  ```

* Using `Application Default Credentials` on ComputeEngine or Kubernetes Engine
  In this mode, the datastource will use the the [instance identity document](https://cloud.google.com/compute/docs/instances/verifying-instance-identity) to provide the `id_token.  The sample usage is the same as the mode above except no environment variables are needed.

* Using Service Account Impersonation.
  In this mode, the service account that runs terraform will attempt to impersonate another service account and then acquire the second service accounts `id_token`.  This mode is similar to `access_token` that is provided by [google_service_account_access_token](https://www.terraform.io/docs/providers/google/d/datasource_google_service_account_access_token.html).  Utilizing this mechanism is requires IAM `TokenCreator` role to be granted to the the origin service account _on_ the target account.

  There are further two variations on how this mode is used depending on how the source token is acquired.

  1. If the source token is provided directly as a variable using the requestors token, then grant the requestor the `TokenCreator` role on the target service account.  In the following case, the identity in context with `gcloud` should have the TokenCreator role on `impersonated-account@project.iam.gserviceaccount.com`


  Then run the terraform apply and specify the token directly

     ```terraform apply --var static_access_token=`gcloud auth print-access-token` ```

  ```yaml
  provider google {}

  variable "static_access_token" {
    type = string
  }

  data "google_service_account_id_token" oidc {

    // The following requires the owner of static_access_token to have
    // ServiceAccountTokenCreator permissions on  
    // impersonated-account@some_project.iam.gserviceaccount.com 
    access_token = var.static_access_token
    
    // for any of the static or impersonated types provided
    target_service_account = "impersonated-account@project.iam.gserviceaccount.com"
    // delegates = []
    include_email = true

    target_audience = "https://foo.bar/"

  }

  output "oidc_token" {
    value = data.google_service_account_id_token.oidc.id_token
  }  
  ```

2. If the source token is provided by [google_service_account_access_token](https://www.terraform.io/docs/providers/google/d/datasource_google_service_account_access_token.html), then the target service account must have the `TokenCreator` role _on itself`.

  >> Note: this variation should be extremely rare since you should be able to directly use the source system's service account JSON file directly.

  ```yaml
  provider google {}

  data "google_service_account_access_token" "default" {
  provider = google
  target_service_account = "impersonated-account@project.iam.gserviceaccount.com"
  delegates = []
  scopes = ["userinfo-email", "cloud-platform"]
  lifetime = "300s"
  }

  provider google {
    alias  = "impersonated"
    access_token = data.google_service_account_access_token.default.access_token
  }

  data "google_service_account_id_token" oidc {
    provider = google.impersonated
    
    target_service_account = "impersonated-account@project.iam.gserviceaccount.com"
    include_email = true

    target_audience = "https://foo.bar/"
  }

  output "oidc_token" {
    value = data.google_service_account_id_token.oidc.id_token
  }
  ```

## Argument Reference

The following arguments are supported:

* `target_audience` (Required) - The audience claim value to to issued the `id_token` for.
* `access_token` (Optional) - Raw access_token to use as the source credential.  Specifying this value will use IAMCredentials API even if any other variable is set.
* `target_service_account` (Optional) - The email of the service account bing impersonated. Used only when using impersonation mode.
* `delegates` (Optional) - Delegate chain of approvals needed to perform full impersonation. Specify the fully qualified service account name.   Used only when using impersonation mode.
* `include_email` (Optional) Include the verified email in the claim. Used only when using impersonation mode.

## Attributes Reference

The following attribute is exported:

* `id_token` - The `id_token` representing the new generated identity.
