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
  The Dataplex Zone resource
---

# google_dataplex_zone

The Dataplex Zone resource

## Example Usage - basic_zone
A basic example of a dataplex zone
```hcl
resource "google_dataplex_zone" "primary" {
  discovery_spec {
    enabled = false
  }

  lake     = google_dataplex_lake.basic.name
  location = "us-west1"
  name     = "zone"

  resource_spec {
    location_type = "MULTI_REGION"
  }

  type         = "RAW"
  description  = "Zone for DCL"
  display_name = "Zone for DCL"
  project      = "my-project-name"
  labels       = {}
}

resource "google_dataplex_lake" "basic" {
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

* `discovery_spec` -
  (Required)
  Required. Specification of the discovery feature applied to data in this zone.
  
* `lake` -
  (Required)
  The lake for the resource
  
* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  The name of the zone.
  
* `resource_spec` -
  (Required)
  Required. Immutable. Specification of the resources that are referenced by the assets within this zone.
  
* `type` -
  (Required)
  Required. Immutable. The type of the zone. Possible values: TYPE_UNSPECIFIED, RAW, CURATED
  


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
    
* `location_type` -
  (Required)
  Required. Immutable. The location type of the resources that are allowed to be attached to the assets within this zone. Possible values: LOCATION_TYPE_UNSPECIFIED, SINGLE_REGION, MULTI_REGION
    
- - -

* `description` -
  (Optional)
  Optional. Description of the zone.
  
* `display_name` -
  (Optional)
  Optional. User friendly display name.
  
* `labels` -
  (Optional)
  Optional. User defined labels for the zone.

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

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{name}}`

* `asset_status` -
  Output only. Aggregated status of the underlying assets of the zone.
  
* `create_time` -
  Output only. The time when the zone was created.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `state` -
  Output only. Current state of the zone. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED
  
* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.
  
* `uid` -
  Output only. System generated globally unique ID for the zone. This ID will be different if the zone is deleted and re-created with the same name.
  
* `update_time` -
  Output only. The time when the zone was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Zone can be imported using any of these accepted formats:

```
$ terraform import google_dataplex_zone.default projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{name}}
$ terraform import google_dataplex_zone.default {{project}}/{{location}}/{{lake}}/{{name}}
$ terraform import google_dataplex_zone.default {{location}}/{{lake}}/{{name}}
```



