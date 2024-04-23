---
subcategory: "Tags"
description: |-
  Get tag keys within a GCP organization or project.
---

# google\_tags\_tag\_keys

Get tag keys by org or project `parent`.

## Example Usage

```tf
data "google_tags_tag_keys" "environment_tag_key"{
  parent = "organizations/12345"
}
```
```tf
data "google_tags_tag_keys" "environment_tag_key"{
  parent = "projects/abc"
}
```

## Argument Reference

The following arguments are supported:

* `parent` - (Required) The resource name of the parent organization or project. It can be in format `organizations/{org_id}` or `projects/{project_id_or_number}`.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `name` - an identifier for the resource with format `tagKeys/{{name}}`

* `namespaced_name` -
  Namespaced name of the TagKey which is in the format `{parentNamespace}/{shortName}`.

* `create_time` -
  Creation time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".

* `update_time` -
  Update time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".

* `short_name` -
  The user friendly name for a TagKey. The short name should be unique for TagKeys wihting the same tag namespace.

* `parent` -
  The resource name of the TagKey's parent. A TagKey can be parented by an Orgination or a Project.

* `description` -
  User-assigned description of the TagKey.

* `purpose` -
  A purpose denotes that this Tag is intended for use in policies of a specific policy engine, and will involve that policy engine in management operations involving this Tag. A purpose does not grant a policy engine exclusive rights to the Tag, and it may be referenced by other policy engines.

* `purpose_data` - 
  Purpose data corresponds to the policy system that the tag is intended for. See documentation for Purpose for formatting of this field.

