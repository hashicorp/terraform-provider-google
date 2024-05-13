---
subcategory: "Tags"
description: |-
  Get a tag key within a GCP organization or project.
---

# google_tags_tag_key

Get a tag key by org or project `parent` and `short_name`.

## Example Usage

```tf
data "google_tags_tag_key" "environment_tag_key"{
  parent = "organizations/12345"
  short_name = "environment"
}
```
```tf
data "google_tags_tag_key" "environment_tag_key"{
  parent = "projects/abc"
  short_name = "environment"
}
```

## Argument Reference

The following arguments are supported:

* `short_name` - (Required) The tag key's short_name.

* `parent` - (Required) The resource name of the parent organization or project. It can be in format `organizations/{org_id}` or `projects/{project_id_or_number}`.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `tagKeys/{{name}}`

* `name` -
  The generated numeric id for the TagKey.

* `namespaced_name` -
  Namespaced name of the TagKey which is in the format `{parentNamespace}/{shortName}`.

* `create_time` -
  Creation time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".

* `update_time` -
  Update time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".
