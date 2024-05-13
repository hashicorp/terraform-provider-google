---
subcategory: "Tags"
description: |-
  Get a tag value from the parent key and short_name.
---

# google_tags_tag_value

Get a tag value by `parent` key and `short_name`.

## Example Usage

```tf
data "google_tags_tag_value" "environment_prod_tag_value"{
  parent = "tagKeys/56789"
  short_name = "production"
}
```

## Argument Reference

The following arguments are supported:

* `short_name` - (Required) The tag value's short_name.

* `parent` - (Required) The resource name of the parent tagKey in format `tagKey/{name}`.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - an identifier for the resource with format `tagValues/{{name}}`

* `name` -
  The generated numeric id for the TagValue.

* `namespaced_name` -
  Namespaced name of the TagValue.

* `create_time` -
  Creation time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".

* `update_time` -
  Update time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".
