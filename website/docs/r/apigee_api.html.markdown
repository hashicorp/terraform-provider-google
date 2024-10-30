---
subcategory: "Apigee"
page_title: "Google: google_apigee_api"
description: |-
    An Apigee API proxy is essentially a layer that sits in front of your backend APIs. It acts as an intermediary between your API consumers (like mobile apps or websites) and your backend services.   

    Think of it like a gatekeeper or a middleman:

    * Decoupling: It decouples the app-facing API from your backend services. This means you can make changes to your backend systems without affecting the apps that use your API, as long as the API proxy interface remains consistent.   
    * Abstraction: It hides the complexities of your backend systems, presenting a simplified and consistent interface to your API consumers.
    * Control: It gives you fine-grained control over how your APIs are accessed and used, allowing you to enforce security policies, rate limits, and other controls.
---

# google_apigee_api

To get more information about API proxies see, see:

* [API documentation](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.apis)
* How-to Guides
  * [API proxies](https://cloud.google.com/apigee/docs/resources)


## Example Usage

```hcl
data "archive_file" "bundle" {
  type             = "zip"
  source_dir       = "${path.module}/bundle"
  output_path      = "${path.module}/bundle.zip"
  output_file_mode = "0644"
}

resource "google_apigee_sharedflow" "sharedflow" {
  name          = "shareflow1"
  org_id        = var.org_id
  config_bundle = data.archive_file.bundle.output_path
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The ID of the API proxy.

* `org_id` -
  (Required)
  The Apigee Organization name associated with the Apigee instance.

* `config_bundle` -
  (Required)
  Path to the config zip bundle.

- - -

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `organizations/{{org_id}}/apis/{{name}}`

* `meta_data` -
  Metadata describing the API proxy.
  Structure is [documented below](#nested_meta_data).

* `revision` -
  A list of revisions of this API proxy.

* `latest_revision_id` -
  The id of the most recently created revision for this API proxy.

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

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

An API proxy can be imported using any of these accepted formats:

* `{{org_id}}/apis/{{name}}`
* `{{org_id}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import API proxy using one of the formats above. For example:

```tf
import {
  id = "{{org_id}}/apis/{{name}}"
  to = google_apigee_api.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), API proxy can be imported using one of the formats above. For example:

```
terraform import google_apigee_api.default {{org_id}}/apis/{{name}}
terraform import google_apigee_api.default {{org_id}}/{{name}}
```
