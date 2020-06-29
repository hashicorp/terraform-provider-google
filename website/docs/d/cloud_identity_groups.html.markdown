---
subcategory: "Cloud Identity"
layout: "google"
page_title: "Google: google_cloud_identity_groups"
sidebar_current: "docs-google-datasource-cloud-identity-groups"
description: |-
  Get list of the Cloud Identity Groups under a customer or namespace.
---

# google_cloud_identity_groups

Use this data source to get list of the Cloud Identity Groups under a customer or namespace.

https://cloud.google.com/identity/docs/concepts/overview#groups

## Example Usage

```tf
data "google_cloud_identity_groups" "groups" {
  parent = "customers/A01b123xz"
}
```

## Argument Reference

* `parent` - The parent resource under which to list all Groups. Must be of the form identitysources/{identity_source_id} for external- identity-mapped groups or customers/{customer_id} for Google Groups.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `groups` - The list of groups under the provided customer or namespace. Structure is documented below.

The `groups` block contains:

* `name` -
  Resource name of the Group in the format: groups/{group_id}, where `group_id` is the unique ID assigned to the Group.

* `group_key` -
  EntityKey of the Group.  Structure is documented below.

* `display_name` -
  The display name of the Group.

* `description` -
  An extended description to help users determine the purpose of a Group.

* `labels` -The labels that apply to the Group.
  Contains 'cloudidentity.googleapis.com/groups.discussion_forum': '' if the Group is a Google Group or
  'system/groups/external': '' if the Group is an external-identity-mapped group.

The `group_key` block supports:

* `id` -
  The ID of the entity.
  For Google-managed entities, the id is the email address of an existing group or user.
  For external-identity-mapped entities, the id is a string conforming
  to the Identity Source's requirements.

* `namespace` -
  The namespace in which the entity exists.
  If not populated, the EntityKey represents a Google-managed entity
  such as a Google user or a Google Group.
  If populated, the EntityKey represents an external-identity-mapped group.
  The namespace must correspond to an identity source created in Admin Console
  and must be in the form of `identitysources/{identity_source_id}`.