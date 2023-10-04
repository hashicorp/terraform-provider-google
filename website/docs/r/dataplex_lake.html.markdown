---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
#
# ----------------------------------------------------------------------------
#
#     This file is managed by Magic Modules (https:#github.com/GoogleCloudPlatform/magic-modules)
#     and is based on the DCL (https:#github.com/GoogleCloudPlatform/declarative-resource-client-library).
#     Changes will need to be made to the DCL or Magic Modules instead of here.
#
#     We are not currently able to accept contributions to this file. If changes
#     are required, please file an issue at https:#github.com/hashicorp/terraform-provider-google/issues/new/choose
#
# ----------------------------------------------------------------------------
subcategory: "Dataplex"
description: |-
  The Dataplex Lake resource
---

# google_dataplex_lake

The Dataplex Lake resource

## Example Usage - basic_lake
A basic example of a dataplex lake
```hcl
resource "google_dataplex_lake" "primary" {
  location     = "us-west1"
  name         = "lake"
  description  = "Lake for DCL"
  display_name = "Lake for DCL"
  project      = "my-project-name"

  labels = {
    my-lake = "exists"
  }
}


```

## Argument Reference

The following arguments are supported:

* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  The name of the lake.
  


- - -

* `description` -
  (Optional)
  Optional. Description of the lake.
  
* `display_name` -
  (Optional)
  Optional. User friendly display name.
  
* `labels` -
  (Optional)
  Optional. User-defined labels for the lake.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field `effective_labels` for all of the labels present on the resource.
  
* `metastore` -
  (Optional)
  Optional. Settings to manage lake and Dataproc Metastore service instance association.
  
* `project` -
  (Optional)
  The project for the resource
  


The `metastore` block supports:
    
* `service` -
  (Optional)
  Optional. A relative reference to the Dataproc Metastore (https://cloud.google.com/dataproc-metastore/docs) service associated with the lake: `projects/{project_id}/locations/{location_id}/services/{service_id}`
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/lakes/{{name}}`

* `asset_status` -
  Output only. Aggregated status of the underlying assets of the lake.
  
* `create_time` -
  Output only. The time when the lake was created.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `metastore_status` -
  Output only. Metastore status of the lake.
  
* `service_account` -
  Output only. Service account associated with this lake. This service account must be authorized to access or operate on resources managed by the lake.
  
* `state` -
  Output only. Current state of the lake. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED
  
* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.
  
* `uid` -
  Output only. System generated globally unique ID for the lake. This ID will be different if the lake is deleted and re-created with the same name.
  
* `update_time` -
  Output only. The time when the lake was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Lake can be imported using any of these accepted formats:
* `projects/{{project}}/locations/{{location}}/lakes/{{name}}`
* `{{project}}/{{location}}/{{name}}`
* `{{location}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Lake using one of the formats above. For example:


```tf
import {
  id = "projects/{{project}}/locations/{{location}}/lakes/{{name}}"
  to = google_dataplex_lake.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Lake can be imported using one of the formats above. For example:

```
$ terraform import google_dataplex_lake.default projects/{{project}}/locations/{{location}}/lakes/{{name}}
$ terraform import google_dataplex_lake.default {{project}}/{{location}}/{{name}}
$ terraform import google_dataplex_lake.default {{location}}/{{name}}
```



