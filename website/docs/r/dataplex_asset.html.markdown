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
  The Dataplex Asset resource
---

# google_dataplex_asset

The Dataplex Asset resource

## Example Usage - basic_asset
```hcl
resource "google_storage_bucket" "basic_bucket" {
  name          = "bucket"
  location      = "us-west1"
  uniform_bucket_level_access = true
  lifecycle {
    ignore_changes = [
      labels
    ]
  }
 
  project = "my-project-name"
}
 
resource "google_dataplex_lake" "basic_lake" {
  name         = "lake"
  location     = "us-west1"
  project = "my-project-name"
}
 
 
resource "google_dataplex_zone" "basic_zone" {
  name         = "zone"
  location     = "us-west1"
  lake = google_dataplex_lake.basic_lake.name
  type = "RAW"
 
  discovery_spec {
    enabled = false
  }
 
 
  resource_spec {
    location_type = "SINGLE_REGION"
  }
 
  project = "my-project-name"
}
 
 
resource "google_dataplex_asset" "primary" {
  name          = "asset"
  location      = "us-west1"
 
  lake = google_dataplex_lake.basic_lake.name
  dataplex_zone = google_dataplex_zone.basic_zone.name
 
  discovery_spec {
    enabled = false
  }
 
  resource_spec {
    name = "projects/my-project-name/buckets/bucket"
    type = "STORAGE_BUCKET"
  }

  labels = {
    env     = "foo"
    my-asset = "exists"
  }

 
  project = "my-project-name"
  depends_on = [
    google_storage_bucket.basic_bucket
  ]
}
```

## Argument Reference

The following arguments are supported:

* `dataplex_zone` -
  (Required)
  The zone for the resource
  
* `discovery_spec` -
  (Required)
  Required. Specification of the discovery feature applied to data referenced by this asset. When this spec is left unset, the asset will use the spec set on the parent zone.
  
* `lake` -
  (Required)
  The lake for the resource
  
* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  The name of the asset.
  
* `resource_spec` -
  (Required)
  Required. Immutable. Specification of the resource that is referenced by this asset.
  


The `discovery_spec` block supports:
    
* `csv_options` -
  (Optional)
  Optional. Configuration for CSV data.
    
* `enabled` -
  (Required)
  Required. Whether discovery is enabled.
    
* `exclude_patterns` -
  (Optional)
  Optional. The list of patterns to apply for selecting data to exclude during discovery. For Cloud Storage bucket assets, these are interpreted as glob patterns used to match object names. For BigQuery dataset assets, these are interpreted as patterns to match table names.
    
* `include_patterns` -
  (Optional)
  Optional. The list of patterns to apply for selecting data to include during discovery if only a subset of the data should considered. For Cloud Storage bucket assets, these are interpreted as glob patterns used to match object names. For BigQuery dataset assets, these are interpreted as patterns to match table names.
    
* `json_options` -
  (Optional)
  Optional. Configuration for Json data.
    
* `schedule` -
  (Optional)
  Optional. Cron schedule (https://en.wikipedia.org/wiki/Cron) for running discovery periodically. Successive discovery runs must be scheduled at least 60 minutes apart. The default value is to run discovery every 60 minutes. To explicitly set a timezone to the cron tab, apply a prefix in the cron tab: "CRON_TZ=${IANA_TIME_ZONE}" or TZ=${IANA_TIME_ZONE}". The ${IANA_TIME_ZONE} may only be a valid string from IANA time zone database. For example, "CRON_TZ=America/New_York 1 * * * *", or "TZ=America/New_York 1 * * * *".
    
The `resource_spec` block supports:
    
* `name` -
  (Optional)
  Immutable. Relative name of the cloud resource that contains the data that is being managed within a lake. For example: `projects/{project_number}/buckets/{bucket_id}` `projects/{project_number}/datasets/{dataset_id}`
    
* `read_access_mode` -
  (Optional)
  Optional. Determines how read permissions are handled for each asset and their associated tables. Only available to storage buckets assets. Possible values: DIRECT, MANAGED
    
* `type` -
  (Required)
  Required. Immutable. Type of resource. Possible values: STORAGE_BUCKET, BIGQUERY_DATASET
    
- - -

* `description` -
  (Optional)
  Optional. Description of the asset.
  
* `display_name` -
  (Optional)
  Optional. User friendly display name.
  
* `labels` -
  (Optional)
  Optional. User defined labels for the asset.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration. Please refer to the field `effective_labels` for all of the labels present on the resource.
  
* `project` -
  (Optional)
  The project for the resource
  


The `csv_options` block supports:
    
* `delimiter` -
  (Optional)
  Optional. The delimiter being used to separate values. This defaults to ','.
    
* `disable_type_inference` -
  (Optional)
  Optional. Whether to disable the inference of data type for CSV data. If true, all columns will be registered as strings.
    
* `encoding` -
  (Optional)
  Optional. The character encoding of the data. The default is UTF-8.
    
* `header_rows` -
  (Optional)
  Optional. The number of rows to interpret as header rows that should be skipped when reading data rows.
    
The `json_options` block supports:
    
* `disable_type_inference` -
  (Optional)
  Optional. Whether to disable the inference of data type for Json data. If true, all columns will be registered as their primitive types (strings, number or boolean).
    
* `encoding` -
  (Optional)
  Optional. The character encoding of the data. The default is UTF-8.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{dataplex_zone}}/assets/{{name}}`

* `create_time` -
  Output only. The time when the asset was created.
  
* `discovery_status` -
  Output only. Status of the discovery feature applied to data referenced by this asset.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `resource_status` -
  Output only. Status of the resource referenced by this asset.
  
* `security_status` -
  Output only. Status of the security policy applied to resource referenced by this asset.
  
* `state` -
  Output only. Current state of the asset. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED
  
* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.
  
* `uid` -
  Output only. System generated globally unique ID for the asset. This ID will be different if the asset is deleted and re-created with the same name.
  
* `update_time` -
  Output only. The time when the asset was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Asset can be imported using any of these accepted formats:

```
$ terraform import google_dataplex_asset.default projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{dataplex_zone}}/assets/{{name}}
$ terraform import google_dataplex_asset.default {{project}}/{{location}}/{{lake}}/{{dataplex_zone}}/{{name}}
$ terraform import google_dataplex_asset.default {{location}}/{{lake}}/{{dataplex_zone}}/{{name}}
```



