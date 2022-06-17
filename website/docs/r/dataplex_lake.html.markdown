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
layout: "google"
page_title: "Google: google_dataplex_lake"
sidebar_current: "docs-google-dataplex-lake"
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

  labels = {
    my-lake = "exists"
  }

  project = "my-project-name"
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
  
* `metastore_status` -
  Output only. Metastore status of the lake.
  
* `service_account` -
  Output only. Service account associated with this lake. This service account must be authorized to access or operate on resources managed by the lake.
  
* `state` -
  Output only. Current state of the lake. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED
  
* `uid` -
  Output only. System generated globally unique ID for the lake. This ID will be different if the lake is deleted and re-created with the same name.
  
* `update_time` -
  Output only. The time when the lake was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Lake can be imported using any of these accepted formats:

```
$ terraform import google_dataplex_lake.default projects/{{project}}/locations/{{location}}/lakes/{{name}}
$ terraform import google_dataplex_lake.default {{project}}/{{location}}/{{name}}
$ terraform import google_dataplex_lake.default {{location}}/{{name}}
```



