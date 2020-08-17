---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_iam_policy"
sidebar_current: "docs-google-datasource-iam-policy"
description: |-
  Generates an IAM policy that can be referenced by other resources, applying
  the policy to them.
---

# google\_iam\_policy

Generates an IAM policy document that may be referenced by and applied to
other Google Cloud Platform resources, such as the `google_project` resource.

**Note:** Several restrictions apply when setting IAM policies through this API.
See the [setIamPolicy docs](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setIamPolicy)
for a list of these restrictions.

```hcl
data "google_iam_policy" "admin" {
  binding {
    role = "roles/compute.instanceAdmin"

    members = [
      "serviceAccount:your-custom-sa@your-project.iam.gserviceaccount.com",
    ]
  }

  binding {
    role = "roles/storage.objectViewer"

    members = [
      "user:alice@gmail.com",
    ]
  }

  audit_config {
    service = "cloudkms.googleapis.com"
    audit_log_configs {
      log_type = "DATA_READ",
      exempted_members = ["user:you@domain.com"]
    }

    audit_log_configs {
      log_type = "DATA_WRITE",
    }

    audit_log_configs {
      log_type = "ADMIN_READ",
    }
  }
}
```

This data source is used to define IAM policies to apply to other resources.
Currently, defining a policy through a datasource and referencing that policy
from another resource is the only way to apply an IAM policy to a resource.

## Argument Reference

The following arguments are supported:

* `binding` (Required) - A nested configuration block (described below)
  defining a binding to be included in the policy document. Multiple
  `binding` arguments are supported.

Each document configuration must have one or more `binding` blocks, which
each accept the following arguments:

* `role` (Required) - The role/permission that will be granted to the members.
  See the [IAM Roles](https://cloud.google.com/compute/docs/access/iam) documentation for a complete list of roles.
  Note that custom roles must be of the format `[projects|organizations]/{parent-name}/roles/{role-name}`.

* `members` (Required) - An array of identities that will be granted the privilege in the `role`. For more details on format and restrictions see https://cloud.google.com/billing/reference/rest/v1/Policy#Binding
  Each entry can have one of the following values:
  * **allUsers**: A special identifier that represents anyone who is on the internet; with or without a Google account. It **can't** be used with the `google_project` resource.
  * **allAuthenticatedUsers**: A special identifier that represents anyone who is authenticated with a Google account or a service account. It **can't** be used with the `google_project` resource.
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

* `audit_config` (Optional) - A nested configuration block that defines logging additional configuration for your project.
  * `service` (Required) Defines a service that will be enabled for audit logging. For example, `storage.googleapis.com`, `cloudsql.googleapis.com`. `allServices` is a special value that covers all services.
  * `audit_log_configs` (Required) A nested block that defines the operations you'd like to log.
    * `log_type` (Required) Defines the logging level. `DATA_READ`, `DATA_WRITE` and `ADMIN_READ` capture different types of events. See [the audit configuration documentation](https://cloud.google.com/resource-manager/reference/rest/Shared.Types/AuditConfig) for more details.
    * `exempted_members` (Optional) Specifies the identities that are exempt from these types of logging operations. Follows the same format of the `members` array for `binding`.

* `condition` - (Optional) An [IAM Condition](https://cloud.google.com/iam/docs/conditions-overview) for a given binding. Structure is documented below.

The `condition` block supports:

* `expression` - (Required) Textual representation of an expression in Common Expression Language syntax.

* `title` - (Required) A title for the expression, i.e. a short string describing its purpose.

* `description` - (Optional) An optional description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.

## Attributes Reference

The following attribute is exported:

* `policy_data` - The above bindings serialized in a format suitable for
  referencing from a resource that supports IAM.
