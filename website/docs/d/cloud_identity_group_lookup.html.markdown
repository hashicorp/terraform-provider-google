---
subcategory: "Cloud Identity"
description: |-
  Look up a Cloud Identity Group using its email and namespace.
---

# google_cloud_identity_group_lookup

Use this data source to look up the resource name of a Cloud Identity Group by its [EntityKey](https://cloud.google.com/identity/docs/reference/rest/v1/EntityKey), i.e. the group's email.

https://cloud.google.com/identity/docs/concepts/overview#groups

## Example Usage

```tf
data "google_cloud_identity_group_lookup" "group" {
  group_key {
    id = "my-group@example.com"
  }
}
```

## Argument Reference

* `group_key` - (Required) The EntityKey of the Group to lookup. A unique identifier for an entity in the Cloud Identity Groups API.
An entity can represent either a group with an optional namespace or a user without a namespace.
The combination of id and namespace must be unique; however, the same id can be used with different namespaces. Structure is [documented below](#nested_group_key).

<a name="nested_group_key"></a>The `group_key` block supports:

* `id` -
  (Required) The ID of the entity.
  For Google-managed entities, the id is the email address of an existing group or user.
  For external-identity-mapped entities, the id is a string conforming
  to the Identity Source's requirements.

* `namespace` -
  (Optional) The namespace in which the entity exists.
  If not populated, the EntityKey represents a Google-managed entity
  such as a Google user or a Google Group.
  If populated, the EntityKey represents an external-identity-mapped group.
  The namespace must correspond to an identity source created in Admin Console
  and must be in the form of `identitysources/{identity_source_id}`.


## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `name` -
  Resource name of the Group in the format: groups/{group_id}, where `group_id` is the unique ID assigned to the Group.