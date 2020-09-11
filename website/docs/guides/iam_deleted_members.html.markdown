---
layout: "google"
page_title: "Handling diffs with deleted IAM principals"
sidebar_current: "docs-google-provider-guides-iam-deleted"
description: |-
  Handling diffs with deleted IAM principals
---

# Handling diffs within IAM resources due to `deleted:` IAM principals

## What is a `deleted:` IAM principal?

Deleted IAM principals reference IAM principals (service accounts, users, groups) that were granted permission on an IAM policy when they were deleted. Deleting the IAM principal does not immediately purge associated permissions, but prefixes them with `deleted:` for \~30 days before deleting them. 

This results in IAM bindings existing for accounts that no longer exist. These bindings can be removed by removing the permission for the deleted account. Terraform will detect permissions that exist for deleted accounts, and in some cases try to remove them. This can potentially cause permanent diffs depending on which types of IAM resources are used within Terraform. This guide documents the different cases and how to handle the diffs in each case.

Validation will occur in the API preventing granting permissions on an IAM policy to a principal that matches a `deleted:` principal on the same policy. This is intended to prevent granting permission to a user that has the same name as a recently deleted user who had escalated permissions.

Prior to version 3.3.0 and 2.20.1 of the provider there was a bug with `deleted:` service accounts. It is strongly recommended to upgrade to 3.3.0+ to avoid issues.

For more information on deleted accounts, see the [official documentation](https://cloud.google.com/iam/docs/creating-managing-service-accounts) on service accounts.

## Intermediate phase

Deleted IAM principals are being introduced to GCP starting on September 14, 2020. Support for them will be fully rolled out over the few weeks following. This rollout period may present a period of time when the API will behave differently, requiring some extra action to be taken to resolve Terraform diffs. During this phase users will need to delete all references to `deleted:` members before adding bindings for a new principal with the same email. Attempting to add bindings for a new principal that shares an email with a `deleted:` principal on the same policy will result in the permissions being added to the `deleted:` principal, causing Terraform to attempt to recreate the binding. References to a `deleted:` principal must be fully removed from an IAM policy before permissions can be granted to the new principal.

After this intermediate phase, attempting to grant permissions to a principal that shares an email with a `deleted:` principal on the same policy will result in an error.

## Using `*_iam_policy` resources

`_iam_policy` allows you to declare the entire IAM policy from within Terraform. Users may see diffs on `deleted:` members in some cirtumstances, but applying the policy should succeed and resolve any issues. Specifying `deleted:` members is not allowed in Terraform, so any policy entirely managed by Terraform should automatically remove any deleted members when Terraform is run.

During the intermediate period it may be required to `taint` the `_iam_policy` resource to ensure any deleted principals are removed *before* the new principal is granted permission. This should only be necessary if you are continuing to see diffs after successful applies. For more information on using `taint` see the [official documentation](https://www.terraform.io/docs/commands/taint.html).

**Note:** Tainting an `iam_policy` resource will delete and recreate it. If the account that Terraform uses to provision GCP resources requires permissions granted by the `iam_policy` resource it may result in Terraform being unable to complete the apply. This is possible for `google_project_iam_policy` and `google_organization_iam_policy`. Special care should be taken when tainting these resources.

## `*_iam_binding` resources

`_iam_binding` resources handle all the members who are granted a specific role for an IAM policy. These resources may see diffs if a member they grant a role to is deleted and recreated. Due to these resources not controlling the entire IAM policy you may see issues around diffs not being resolved as requests can include both the deleted and non-deleted form of a principal.

Example diff caused by `deleted:` member on an IAM binding resource:
```hcl
  # google_secret_manager_secret_iam_binding.binding will be updated in-place
  ~ resource "google_secret_manager_secret_iam_binding" "binding" {
        id        = "projects/my-project/secrets/secret/roles/secretmanager.secretAccessor"
      ~ members   = [
          - "deleted:serviceAccount:myaccount@my-project.iam.gserviceaccount.com?uid=10231234122325702",
          + "serviceAccount:myaccount@my-project.iam.gserviceaccount.com",
        ]
        project   = "my-project"
        role      = "roles/secretmanager.secretAccessor"
        secret_id = "projects/my-project/secrets/secret"
    }
```

Tainting the `_iam_binding` resource using the [taint command](https://www.terraform.io/docs/commands/taint.html) may resolve this diff. If it does not, or results in an API error, you may need to manually remove any references to the deleted version of the principal causing the diff. This can be done through the Cloud Console UI for the resource. This could also be done via Terraform by specifying the entire IAM policy authoritatively using the `_iam_policy` resource.

## `*_iam_member` resources and inconsistent results

`_iam_member` resources manage a certain permission for a certain principal. This resource may see diffs if the principal that it grants permissions for is deleted and recreated. Similarly to the binding resources above these may see unresolved diffs on apply. Due to this resource only ensuring that a single role is present for a single principal, member resources will not remove permissions for `deleted:` members.

If applying the creation of an `_iam_member` resource results in an API error or the Terraform error: `Provider produced inconsistent result after apply` it may mean that the deleted version of the principal already exists on the IAM policy. Inspecting the policy either by [capturing the Terraform debug logs](https://www.terraform.io/docs/internals/debugging.html) or via the Cloud Console can verify that this is the cause. If the deleted version of the principal exists on the IAM policy, remove it using the Cloud Console or by using an authoritative IAM resource like `_iam_policy`.



