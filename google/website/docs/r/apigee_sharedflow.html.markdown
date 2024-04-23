---
subcategory: "Apigee"
page_title: "Google: google_apigee_sharedflow"
description: |-
  You can combine policies and resources into a shared flow that you can consume from multiple API proxies, and even from other shared flows.
---

# google\_apigee\_shared\_flow

You can combine policies and resources into a shared flow that you can consume from multiple API proxies, and even from other shared flows. Although it's like a proxy, a shared flow has no endpoint. It can be used only from an API proxy or shared flow that's in the same organization as the shared flow itself.


To get more information about SharedFlow, see:

* [API documentation](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.sharedflows)
* How-to Guides
    * [Sharedflows](https://cloud.google.com/apigee/docs/resources)


## Argument Reference

The following arguments are supported:


* `name` -
  (Required)
  The ID of the shared flow.

* `org_id` -
  (Required)
  The Apigee Organization name associated with the Apigee instance.

* `config_bundle` -
  (Required)
  Path to the config zip bundle.

- - -



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `organizations/{{org_id}}/sharedflows/{{name}}`

* `meta_data` -
  Metadata describing the shared flow.
  Structure is [documented below](#nested_meta_data).

* `revision` -
  A list of revisions of this shared flow.

* `latest_revision_id` -
  The id of the most recently created revision for this shared flow.

* `md5hash` -
  (Computed) Base 64 MD5 hash of the uploaded data. It is speculative as remote does not return hash of the bundle. Remote changes are detected using returned last_modified timestamp.

* `detect_md5hash` - 
  (Optional) Detect changes to local config bundle file or changes made outside of Terraform. MD5 hash of the data, encoded using base64. Hash is automatically computed without need for user input.


<a name="nested_meta_data"></a>The `meta_data` block contains:

* `created_at` -
  (Optional)
  Time at which the API proxy was created, in milliseconds since epoch.

* `last_modified_at` -
  (Optional)
  Time at which the API proxy was most recently modified, in milliseconds since epoch.

* `sub_type` -
  (Optional)
  The type of entity described

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


SharedFlow can be imported using any of these accepted formats:

* `{{org_id}}/sharedflows/{{name}}`
* `{{org_id}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import SharedFlow using one of the formats above. For example:

```tf
import {
  id = "{{org_id}}/sharedflows/{{name}}"
  to = google_apigee_sharedflow.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), SharedFlow can be imported using one of the formats above. For example:

```
$ terraform import google_apigee_sharedflow.default {{org_id}}/sharedflows/{{name}}
$ terraform import google_apigee_sharedflow.default {{org_id}}/{{name}}
```
