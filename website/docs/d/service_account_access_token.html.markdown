---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_service_account_access_token"
sidebar_current: "docs-google-service-account-access-token"
description: |-
  Produces access_token for impersonated service accounts
---

# google\_service\_account\_access\_token

This data source provides a google `oauth2` `access_token` for a different service account than the one initially running the script.

For more information see
[the official documentation](https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials) as well as [iamcredentials.generateAccessToken()](https://cloud.google.com/iam/credentials/reference/rest/v1/projects.serviceAccounts/generateAccessToken)

## Example Usage

To allow `service_A` to impersonate `service_B`, grant the [Service Account Token Creator](https://cloud.google.com/iam/docs/service-accounts#the_service_account_token_creator_role) on B to A. 

In the IAM policy below, `service_A` is given the Token Creator role impersonate `service_B`

```sh
resource "google_service_account_iam_binding" "token-creator-iam" {
	service_account_id = "projects/-/serviceAccounts/service_B@projectB.iam.gserviceaccount.com"
	role               = "roles/iam.serviceAccountTokenCreator"
	members = [
		"serviceAccount:service_A@projectA.iam.gserviceaccount.com",
	]
}
```

Once the IAM permissions are set, you can apply the new token to a provider bootstrapped with it.  Any resources that references the aliased provider will run as the new identity.

In the example below, `google_project` will run as `service_B`.

```hcl
provider "google" {
}

data "google_client_config" "default" {
  provider = google
}

data "google_service_account_access_token" "default" {
  provider               = google
  target_service_account = "service_B@projectB.iam.gserviceaccount.com"
  scopes                 = ["userinfo-email", "cloud-platform"]
  lifetime               = "300s"
}

provider "google" {
  alias        = "impersonated"
  access_token = data.google_service_account_access_token.default.access_token
}

data "google_client_openid_userinfo" "me" {
  provider = google.impersonated
}

output "target-email" {
  value = data.google_client_openid_userinfo.me.email
}
```

> *Note*: the generated token is non-refreshable and can have a maximum `lifetime` of `3600` seconds.

## Argument Reference

The following arguments are supported:

* `target_service_account` (Required) - The service account _to_ impersonate (e.g. `service_B@your-project-id.iam.gserviceaccount.com`)
* `scopes` (Required) - The scopes the new credential should have (e.g. `["storage-ro", "cloud-platform"]`)
* `delegates` (Optional) - Delegate chain of approvals needed to perform full impersonation. Specify the fully qualified service account name.  (e.g. `["projects/-/serviceAccounts/delegate-svc-account@project-id.iam.gserviceaccount.com"]`)
* `lifetime` (Optional) Lifetime of the impersonated token (defaults to its max: `3600s`).

## Attributes Reference

The following attribute is exported:

* `access_token` - The `access_token` representing the new generated identity.
