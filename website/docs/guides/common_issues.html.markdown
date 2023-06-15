---
page_title: "Common issues/FAQ"
description: |-
  Common issues and frequently asked questions when using the provider.
---

# Google Provider Common Issues/FAQ

## 403 Service API disabled

```
<service> API has not been used in project <project> before or it is disabled. Enable it by visiting https://console.developers.google.com/apis/api/<service>.googleapis.com/overview?project=<project> then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry.
```

Services must be [enabled in a project](https://cloud.google.com/service-usage/docs/enable-disable) before their service API can be used by the provider. The [`google_project_service` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_service) enables GCP service APIs with Terraform. 

For `google_project_service` guidance and troubleshooting, see the [advanced user guide](/docs/providers/google/guides/google_project_service.html).

## API not enabled in a different project than the resource

Quota projects refer to the project used in requests to GCP APIs for the purpose of preconditions, quota, and billing. By default, a resource’s quota project is determined by the API and may be the project associated with your credentials, or the resource project depending on the API. For most resources, `user_project_override` (and optionally `billing_project`) can control the quota project used in API requests. See the [provider_reference documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override) for more information.

## Handling default service accounts (user-managed)

Certain GCP services automatically create user-managed service accounts called default service accounts. These are granted a large set of permissions on project creation and are the responsibility of the user once they are created. See [Google’s guide on default service accounts](https://cloud.google.com/iam/docs/service-account-types#default).

Constraining the permissions or [replacing](https://github.com/terraform-google-modules/terraform-google-project-factory/blob/master/docs/FAQ.md#why-do-you-delete-the-default-service-account) the default service accounts entirely may be a suitable form of management. The Google provider offers the [google_project_default_service_accounts resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_project_default_service_accounts) as a way to manage default service accounts within Terraform.

## Handling Google-managed service accounts

Some services create service accounts that are fully managed by Google. These exist outside of user projects, so they do not appear when viewing a project’s service accounts. See Google’s information on [Google-managed service accounts](https://cloud.google.com/iam/docs/service-account-types#default).

The Google provider offers the [google_project_service_identity resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/project_service_identity), enabling access to the email address of Google-managed service accounts per service. 