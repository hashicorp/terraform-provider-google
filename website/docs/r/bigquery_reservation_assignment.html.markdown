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
subcategory: "BigQuery Reservation"
page_title: "Google: google_bigquery_reservation_assignment"
description: |-
  The BigqueryReservation Assignment resource
---

# google_bigquery_reservation_assignment

The BigqueryReservation Assignment resource

## Example Usage - basic
```hcl
resource "google_bigquery_reservation" "basic" {
  name  = "tf-test-my-reservation%{random_suffix}"
  project = "my-project-name"
  location = "us-central1"
  slot_capacity = 0
  ignore_idle_slots = false
}

resource "google_bigquery_reservation_assignment" "primary" {
  assignee  = "projects/my-project-name"
  job_type = "PIPELINE"
  reservation = google_bigquery_reservation.basic.id
}
```

## Argument Reference

The following arguments are supported:

* `assignee` -
  (Required)
  The resource which will use the reservation. E.g. projects/myproject, folders/123, organizations/456.
  
* `job_type` -
  (Required)
  Types of job, which could be specified when using the reservation. Possible values: JOB_TYPE_UNSPECIFIED, PIPELINE, QUERY
  
* `reservation` -
  (Required)
  The reservation for the resource
  


- - -

* `location` -
  (Optional)
  The location for the resource
  
* `project` -
  (Optional)
  The project for the resource
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments/{{name}}`

* `name` -
  Output only. The resource name of the assignment.
  
* `state` -
  Assignment will remain in PENDING state if no active capacity commitment is present. It will become ACTIVE when some capacity commitment becomes active. Possible values: STATE_UNSPECIFIED, PENDING, ACTIVE
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Assignment can be imported using any of these accepted formats:

```
$ terraform import google_bigquery_reservation_assignment.default projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments/{{name}}
$ terraform import google_bigquery_reservation_assignment.default {{project}}/{{location}}/{{reservation}}/{{name}}
$ terraform import google_bigquery_reservation_assignment.default {{location}}/{{reservation}}/{{name}}
```



