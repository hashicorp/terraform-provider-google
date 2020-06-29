---
subcategory: "Cloud Identity"
layout: "google"
page_title: "Google: google_cloud_identity_group_memberships"
sidebar_current: "docs-google-datasource-cloud-identity-group-memberships"
description: |-
  Get list of the Cloud Identity Group Memberships within a Group.
---

# google_cloud_identity_group_memberships

Use this data source to get list of the Cloud Identity Group Memberships within a given Group.

https://cloud.google.com/identity/docs/concepts/overview#memberships

## Example Usage

```tf
data "google_cloud_identity_group_memberships" "members" {
  group = "groups/123eab45c6defghi"
}
```

## Argument Reference

* `group` - The parent Group resource under which to lookup the Membership names. Must be of the form groups/{group_id}.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `memberships` - The list of memberships under the given group. Structure is documented below.

The `memberships` block contains:

* `name` -
  The resource name of the Membership, of the form groups/{group_id}/memberships/{membership_id}.

* `roles` - The MembershipRoles that apply to the Membership. Structure is documented below.

* `member_key` -
  (Optional)
  EntityKey of the member.  Structure is documented below.

* `preferred_member_key` -
  (Optional)
  EntityKey of the member.  Structure is documented below.

The `roles` block supports:

* `name` - The name of the MembershipRole. One of OWNER, MANAGER, MEMBER.


The `member_key` block supports:

* `id` - The ID of the entity. For Google-managed entities, the id is the email address of an existing
  group or user. For external-identity-mapped entities, the id is a string conforming
  to the Identity Source's requirements.

* `namespace` - The namespace in which the entity exists.
  If not populated, the EntityKey represents a Google-managed entity
  such as a Google user or a Google Group.
  If populated, the EntityKey represents an external-identity-mapped group.

The `preferred_member_key` block supports:

* `id` - The ID of the entity. For Google-managed entities, the id is the email address of an existing
  group or user. For external-identity-mapped entities, the id is a string conforming
  to the Identity Source's requirements.

* `namespace` - The namespace in which the entity exists.
  If not populated, the EntityKey represents a Google-managed entity
  such as a Google user or a Google Group.
  If populated, the EntityKey represents an external-identity-mapped group.
