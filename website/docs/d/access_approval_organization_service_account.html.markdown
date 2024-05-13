---
subcategory: "Access Approval"
description: |-
  Get the email address of an organization's Access Approval service account.
---

# google_access_approval_organization_service_account

Get the email address of an organization's Access Approval service account.

Each Google Cloud organization has a unique service account used by Access Approval.
When using Access Approval with a
[custom signing key](https://cloud.google.com/cloud-provider-access-management/access-approval/docs/review-approve-access-requests-custom-keys),
this account needs to be granted the `cloudkms.signerVerifier` IAM role on the
Cloud KMS key used to sign approvals.

## Example Usage

```hcl
data "google_access_approval_organization_service_account" "service_account" {
  organization_id = "my-organization"
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.signerVerifier"
  member        = "serviceAccount:${data.google_access_approval_organization_service_account.service_account.account_email}"
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required) The organization ID the service account was created for.

## Attributes Reference

The following attributes are exported:

* `name` - The Access Approval service account resource name. Format is "organizations/{organization_id}/serviceAccount".

* `account_email` - The email address of the service account. This value is
often used to refer to the service account in order to grant IAM permissions.
