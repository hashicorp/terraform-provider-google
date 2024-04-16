---
subcategory: "Cloud Platform"
description: |-
  Ensures that a member:role pairing does not exist in a project's IAM policy.
---

# google\_project\_iam\member\_remove

Ensures that a member:role pairing does not exist in a project's IAM policy. 

On create, this resource will modify the policy to remove the `member` from the
`role`. If the membership is ever re-added, the next refresh will clear this
resource from state, proposing re-adding it to correct the membership. Import is
not supported- this resource will acquire the current policy and modify it as
part of creating the resource.

This resource will conflict with `google_project_iam_policy` and
`google_project_iam_binding` resources that share a role, as well as
`google_project_iam_member` resources that target the same membership. When
multiple resources conflict the final state is not guaranteed to include or omit
the membership. Subsequent `terraform apply` calls will always show a diff
until the configuration is corrected.

For more information see
[the official documentation](https://cloud.google.com/iam/docs/granting-changing-revoking-access)
and
[API reference](https://cloud.google.com/resource-manager/reference/rest/v1/projects/setIamPolicy).

## Example Usage

```hcl
data "google_project" "target_project {}

resource "google_project_iam_member_remove" "foo" {
  role     = "roles/editor"
  project  = google_project.target_project.project_id
  member  = "serviceAccount:${google_project.target_project.number}-compute@developer.gserviceaccount.com"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The project id of the target project.

* `role` - (Required) The target role that should be removed. 

* `member` - (Required) The IAM principal that should not have the target role.
  Each entry can have one of the following values:
  * **user:{emailid}**: An email address that represents a specific Google account. For example, alice@gmail.com or joe@example.com.
  * **serviceAccount:{emailid}**: An email address that represents a service account. For example, my-other-app@appspot.gserviceaccount.com.
  * **group:{emailid}**: An email address that represents a Google group. For example, admins@example.com.
  * **domain:{domain}**: A G Suite domain (primary, instead of alias) name that represents all the users of that domain. For example, google.com or example.com.

