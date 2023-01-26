---
subcategory: "Tags"
description: |-
  Get a tag key within a GCP organization.
---

# google\_tags\_tag\_key

Get a tag key within a GCP org by `parent` and `short_name`.

## Example Usage

```tf
data "google_tags_tag_key" "environment_tag_key"{
  parent = "organizations/12345"
  short_name = "environment"
}
```

## Argument Reference

The following arguments are supported:

* `short_name` - (Required) The tag key's short_name.

* `parent` - (Required) The resource name of the parent organization in format `organizations/{org_id}`.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `tagKeys/{{name}}`

* `name` -
  The generated numeric id for the TagKey.

* `namespaced_name` -
  Namespaced name of the TagKey.

* `create_time` -
  Creation time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".

* `update_time` -
  Update time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".
