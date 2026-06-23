---
subcategory: "Cloud Platform"
description: |-
  List IAM member bindings for a Google Cloud folder for use with terraform query
  and .tfquery.hcl files.
---

# google_folder_iam_member (list)

Lists IAM **member bindings** for a Google Cloud folder for use with
[`terraform query`](https://developer.hashicorp.com/terraform/cli/commands/query) and
**`.tfquery.hcl`** files. Results correspond to existing
[`google_folder_iam_member`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_folder_iam)
managed resources.

For how list resources work in this provider, file layout, Terraform version requirements, and
shared `list` block arguments, refer to the guide
[Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).

## Example

```hcl
list "google_folder_iam_member" "all" {
  provider = google

  config {
    folder = "folders/123456789"
  }
}
```

Run `terraform query` from the directory that contains the `.tfquery.hcl` file.

## Configuration (`config` block)

* `folder` - (Required) Folder ID to list IAM member from.
  For example, `folders/123456789`

* `role` - (Optional) If set, only bindings with this exact role are returned.
  For example, `roles/editor`. If unset, bindings for all roles are returned.

* `member` - (Optional) If set, only bindings where this principle is a member
  are returned. For example, `user:jane@example.com`. If unset, bindings for 
  all roles are returned.


## Results

By default each result includes **resource identity** for `google_folder_iam_member` (see
[Resource identity](https://developer.hashicorp.com/terraform/language/resources/identities)):

* `folder` - Folder ID the binding belongs to.
* `role` - The Iam role, e.g. `roles/editor`.
* `member` The principal, e.g. `user:jane@example.com`.

With `include_resource = true` on the `list` block, results also include the full resource-style
attributes documented for the managed
[`google_folder_iam_member` resource](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_folder_iam#attributes-reference)
(for example `etag`and `condition` where present in state).