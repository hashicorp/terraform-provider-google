---
subcategory: "Eventarc"
layout: "google"
page_title: "Google: google_eventarc_trigger"
sidebar_current: "docs-google-eventarc-trigger"
description: |-
  An event trigger sends messages to the event receiver service deployed on Cloud Run.
---

# google\_eventarc\_trigger

An event trigger sends messages to the event receiver service deployed on Cloud Run.

* [API documentation](https://cloud.google.com/eventarc/docs/reference/rest/v1/projects.locations.triggers)

## Example Usage

```hcl
resource "google_eventarc_trigger" "trigger" {
  name = "trigger"
  location = "us-central1"
  # matching_criteria is named event_filters in other tools
  matching_criteria {
    attribute = "type"
    value = "google.cloud.pubsub.topic.v1.messagePublished"
  }
  destination {
    # cloud_run_service is named cloud_run in other tools
    cloud_run_service {
      service = google_cloud_run_service.default.name
      region = "us-central1"
    }
  }
}

resource "google_cloud_run_service" "default" {
  name     = "service-eventarc"
  location = "us-central1"

  metadata {
    namespace = "my-project"
  }

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
        args  = ["arrgs"]
      }
      container_concurrency = 50
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `destination` -
  (Required)
  Required. Destination specifies where the events should be sent to.
  
* `location` -
  (Required)
  The location for the resource
  
* `matching_criteria` -
  (Required)
  Required. The criteria by which events are filtered. Only events that match with this criteria will be sent to the destination.

~> **NOTE:** `matching_criteria` is named `event_filters` in other tools.

* `name` -
  (Required)
  Required. The resource name of the trigger. Must be unique within the location on the project and must be in `projects/{project}/locations/{location}/triggers/{trigger}` format.
  

The `destination` block supports:
    
* `cloud_run_service` -
  (Optional)
  Cloud Run fully-managed service that receives the events. The service should be running in the same project as the trigger.

~> **NOTE:** `cloud_run_service` is named `cloud_run` in other tools.

The `matching_criteria` block supports:
    
* `attribute` -
  (Required)
  Required. The name of a CloudEvents attribute. Currently, only a subset of attributes can be specified. All triggers MUST provide a matching criteria for the 'type' attribute.
    
* `value` -
  (Required)
  Required. The value for the attribute.
    
- - -

* `labels` -
  (Optional)
  Optional. User labels attached to the triggers that can be used to group resources.
  
* `project` -
  (Optional)
  The project for the resource
  
* `service_account` -
  (Optional)
  Optional. The IAM service account email associated with the trigger. The service account represents the identity of the trigger. The principal who calls this API must have `iam.serviceAccounts.actAs` permission in the service account. See https://cloud.google.com/iam/docs/understanding-service-accounts?hl=en#sa\\\_common for more information. For Cloud Run destinations, this service account is used to generate identity tokens when invoking the service. See https://cloud.google.com/run/docs/triggering/pubsub-push#create-service-account for information on how to invoke authenticated Cloud Run services. In order to create Audit Log triggers, the service account should also have `roles/eventarc.eventReceiver` IAM role.
  
* `transport` -
  (Optional)
  Optional. In order to deliver messages, Eventarc may use other GCP products as transport intermediary. This field contains a reference to that transport intermediary. This information can be used for debugging purposes.
  

The `cloud_run_service` block supports:
    
* `service` -
  (Required)
  Required. The name of the Cloud run service being addressed (see https://cloud.google.com/run/docs/reference/rest/v1/namespaces.services). Only services located in the same project of the trigger object can be addressed.
    
* `path` -
  (Optional)
  Optional. The relative path on the Cloud Run service the events should be sent to. The value must conform to the definition of URI path segment (section 3.3 of RFC2396). Examples: "/route", "route", "route/subroute".
    
* `region` -
  (Optional)
  Required. The region the Cloud Run service is deployed in.
    
The `transport` block supports:
    
* `pubsub` -
  (Optional)
  The Pub/Sub topic and subscription used by Eventarc as delivery intermediary.
    The `pubsub` block supports:
    
* `topic` -
  (Optional)
  Optional. The name of the Pub/Sub topic created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{PROJECT_ID}/topics/{TOPIC_NAME}`. You may set an existing topic for triggers of the type `google.cloud.pubsub.topic.v1.messagePublished` only. The topic you provide here will not be deleted by Eventarc at trigger deletion.
    
* `subscription` -
  Output only. The name of the Pub/Sub subscription created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{PROJECT_ID}/subscriptions/{SUBSCRIPTION_NAME}`.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/triggers/{{name}}`

* `create_time` -
  Output only. The creation time.
  
* `etag` -
  Output only. This checksum is computed by the server based on the value of other fields, and may be sent only on create requests to ensure the client has an up-to-date value before proceeding.
  
* `update_time` -
  Output only. The last-modified time.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

Trigger can be imported using any of these accepted formats:

```
$ terraform import google_eventarc_trigger.default projects/{{project}}/locations/{{location}}/triggers/{{name}}
$ terraform import google_eventarc_trigger.default {{project}}/{{location}}/{{name}}
$ terraform import google_eventarc_trigger.default {{location}}/{{name}}
```