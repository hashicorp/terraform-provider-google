---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_folder_iam_audit_config"
sidebar_current: "docs-google-folder-iam-audit-config"
description: |-
 Allows management of audit logging config for a given service for a Google Cloud Platform folder.
---

## google\_folder\_iam\_audit\_config

Allows management of audit logging config for a given service for a Google Cloud Platform folder.

```hcl
resource "google_folder_iam_audit_config" "config" {
  folder = "folders/{folder_id}"
  service = "allServices"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:joebloggs@hashicorp.com",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `folder` - (Required) The resource name of the folder in which you want to manage the audit logging config. Its format is folders/{folder_id}.

* `service` - (Required) Service which will be enabled for audit logging.  The special value `allServices` covers all services.  Note that if there are google\_folder\_iam\_audit\_config resources covering both `allServices` and a specific service then the union of the two AuditConfigs is used for that service: the `log_types` specified in each `audit_log_config` are enabled, and the `exempted_members` in each `audit_log_config` are exempted.

* `audit_log_config` - (Required) The configuration for logging of each type of permission.  This can be specified multiple times.  Structure is documented below.

---

The `audit_log_config` block supports:

* `log_type` - (Required) Permission type for which logging is to be configured.  Must be one of `DATA_READ`, `DATA_WRITE`, or `ADMIN_READ`.

* `exempted_members` - (Optional) Identities that do not cause logging for this type of permission.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

## Import
IAM audit config imports use the identifier of the resource in question and the service, e.g.

```
terraform import google_folder_iam_audit_config.config "{{folder_id}} foo.googleapis.com"
```
