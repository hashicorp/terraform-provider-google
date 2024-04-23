---
subcategory: "Tags"
description: |-
  Get tag values from the parent key.
---

# google\_tags\_tag\_values

Get tag values from a `parent` key.

## Example Usage

```tf
data "google_tags_tag_values" "environment_tag_values"{
  parent = "tagKeys/56789"
}
```

## Argument Reference

The following arguments are supported:


* `parent` - (Required) The resource name of the parent tagKey in format `tagKey/{name}`.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `name` - an identifier for the resource with format `tagValues/{{name}}`

* `namespaced_name` -
  Namespaced name of the TagValue.

* `short_name` -
  User-assigned short name for TagValue. The short name should be unique for TagValues within the same parent TagKey.

* `parent` -
  The resource name of the new TagValue's parent TagKey. Must be of the form tagKeys/{tag_key_id}.

* `create_time` -
  Creation time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".

* `update_time` -
  Update time.
  A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".

* `description` - 
  User-assigned description of the TagValue.
