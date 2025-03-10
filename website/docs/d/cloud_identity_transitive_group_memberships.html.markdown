---
subcategory: "Cloud Identity"
description: |-
  Get a list of direct and indirect Cloud Identity Group Memberships within a Group.
---

# google_cloud_identity_group_transitive_memberships

Use this data source to get list of the Cloud Identity Group Memberships within a given Group. Whereas `google_cloud_identity_group_memberships` returns details of only direct members of the group, `google_cloud_identity_group_transitive_memberships` will return details about both direct and indirect members. For example, a user is an indirect member of Group A if the user is a direct member of Group B and Group B is a direct member of Group A.

To get more information about TransitiveGroupMembership, see:

* [API documentation](https://cloud.google.com/identity/docs/reference/rest/v1/groups.memberships/searchTransitiveMemberships)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/identity/docs/how-to/memberships-google-groups)

## Example Usage

```tf
data "google_cloud_identity_group_transitive_memberships" "members" {
  group = "groups/123eab45c6defghi"
}
```

## Argument Reference

* `group` - (Required) The parent Group resource to search transitive memberships in. Must be of the form groups/{group_id}.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `memberships` - The list of memberships under the given group. Structure is [documented below](#nested_memberships).

<a name="nested_memberships"></a>The `memberships` block contains:

* `roles` - The TransitiveMembershipRoles that apply to the Membership. Structure is [documented below](#nested_roles).

* `member` - EntityKey of the member.  This value will be either a userKey in the format `users/000000000000000000000` with a numerical id or a groupKey in the format `groups/000ab0000ab0000` with a hexadecimal id.

* `relation_type` - The relation between the group and the transitive member. The value can be DIRECT, INDIRECT, or DIRECT_AND_INDIRECT.

* `preferred_member_key` -
  (Optional)
  EntityKey of the member.  Structure is [documented below](#nested_preferred_member_key).

<a name="nested_roles"></a>The `roles` block supports:

* `role` - The name of the TransitiveMembershipRole. One of OWNER, MANAGER, MEMBER.

<a name="nested_preferred_member_key"></a>The `preferred_member_key` block supports:

* `id` - The ID of the entity. For Google-managed entities, the id is the email address of an existing
  group or user. For external-identity-mapped entities, the id is a string conforming
  to the Identity Source's requirements.

* `namespace` - The namespace in which the entity exists.
  If not populated, the EntityKey represents a Google-managed entity
  such as a Google user or a Google Group.
  If populated, the EntityKey represents an external-identity-mapped group.
