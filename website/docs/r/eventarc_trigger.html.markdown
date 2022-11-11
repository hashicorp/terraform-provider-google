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
subcategory: "Eventarc"
page_title: "Google: google_eventarc_trigger"
description: |-
  The Eventarc Trigger resource
---

# google_eventarc_trigger

The Eventarc Trigger resource

## Example Usage - basic
```hcl
resource "google_eventarc_trigger" "primary" {
	name = "name"
	location = "europe-west1"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default.name
			region = "europe-west1"
		}
	}
	labels = {
		foo = "bar"
	}
}

resource "google_pubsub_topic" "foo" {
	name = "topic"
}

resource "google_cloud_run_service" "default" {
	name     = "eventarc-service"
	location = "europe-west1"

	metadata {
		namespace = "my-project-name"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				ports {
					container_port = 8080
				}
			}
			container_concurrency = 50
			timeout_seconds = 100
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
  Required. null The list of filters that applies to event attributes. Only events that match all the provided filters will be sent to the destination.
  
* `name` -
  (Required)
  Required. The resource name of the trigger. Must be unique within the location on the project.
  


The `destination` block supports:
    
* `cloud_function` -
  (Optional)
  [WARNING] Configuring a Cloud Function in Trigger is not supported as of today. The Cloud Function resource name. Format: projects/{project}/locations/{location}/functions/{function}
    
* `cloud_run_service` -
  (Optional)
  Cloud Run fully-managed service that receives the events. The service should be running in the same project of the trigger.
    
* `gke` -
  (Optional)
  A GKE service capable of receiving events. The service should be running in the same project as the trigger.
    
* `workflow` -
  (Optional)
  The resource name of the Workflow whose Executions are triggered by the events. The Workflow resource should be deployed in the same project as the trigger. Format: `projects/{project}/locations/{location}/workflows/{workflow}`
    
The `matching_criteria` block supports:
    
* `attribute` -
  (Required)
  Required. The name of a CloudEvents attribute. Currently, only a subset of attributes are supported for filtering. All triggers MUST provide a filter for the 'type' attribute.
    
* `operator` -
  (Optional)
  Optional. The operator used for matching the events with the value of the filter. If not specified, only events that have an exact key-value pair specified in the filter are matched. The only allowed value is `match-path-pattern`.
    
* `value` -
  (Required)
  Required. The value for the attribute. See https://cloud.google.com/eventarc/docs/creating-triggers#trigger-gcloud for available values.
    
- - -

* `channel` -
  (Optional)
  Optional. The name of the channel associated with the trigger in `projects/{project}/locations/{location}/channels/{channel}` format. You must provide a channel to receive events from Eventarc SaaS partners.
  
* `labels` -
  (Optional)
  Optional. User labels attached to the triggers that can be used to group resources.
  
* `project` -
  (Optional)
  The project for the resource
  
* `service_account` -
  (Optional)
  Optional. The IAM service account email associated with the trigger. The service account represents the identity of the trigger. The principal who calls this API must have `iam.serviceAccounts.actAs` permission in the service account. See https://cloud.google.com/iam/docs/understanding-service-accounts#sa_common for more information. For Cloud Run destinations, this service account is used to generate identity tokens when invoking the service. See https://cloud.google.com/run/docs/triggering/pubsub-push#create-service-account for information on how to invoke authenticated Cloud Run services. In order to create Audit Log triggers, the service account should also have `roles/eventarc.eventReceiver` IAM role.
  
* `transport` -
  (Optional)
  Optional. In order to deliver messages, Eventarc may use other GCP products as transport intermediary. This field contains a reference to that transport intermediary. This information can be used for debugging purposes.
  


The `cloud_run_service` block supports:
    
* `path` -
  (Optional)
  Optional. The relative path on the Cloud Run service the events should be sent to. The value must conform to the definition of URI path segment (section 3.3 of RFC2396). Examples: "/route", "route", "route/subroute".
    
* `region` -
  (Optional)
  Required. The region the Cloud Run service is deployed in.
    
* `service` -
  (Required)
  Required. The name of the Cloud Run service being addressed. See https://cloud.google.com/run/docs/reference/rest/v1/namespaces.services. Only services located in the same project of the trigger object can be addressed.
    
The `gke` block supports:
    
* `cluster` -
  (Required)
  Required. The name of the cluster the GKE service is running in. The cluster must be running in the same project as the trigger being created.
    
* `location` -
  (Required)
  Required. The name of the Google Compute Engine in which the cluster resides, which can either be compute zone (for example, us-central1-a) for the zonal clusters or region (for example, us-central1) for regional clusters.
    
* `namespace` -
  (Required)
  Required. The namespace the GKE service is running in.
    
* `path` -
  (Optional)
  Optional. The relative path on the GKE service the events should be sent to. The value must conform to the definition of a URI path segment (section 3.3 of RFC2396). Examples: "/route", "route", "route/subroute".
    
* `service` -
  (Required)
  Required. Name of the GKE service.
    
The `transport` block supports:
    
* `pubsub` -
  (Optional)
  The Pub/Sub topic and subscription used by Eventarc as delivery intermediary.
    
The `pubsub` block supports:
    
* `subscription` -
  Output only. The name of the Pub/Sub subscription created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{PROJECT_ID}/subscriptions/{SUBSCRIPTION_NAME}`.
    
* `topic` -
  (Optional)
  Optional. The name of the Pub/Sub topic created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{PROJECT_ID}/topics/{TOPIC_NAME}. You may set an existing topic for triggers of the type google.cloud.pubsub.topic.v1.messagePublished` only. The topic you provide here will not be deleted by Eventarc at trigger deletion.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/triggers/{{name}}`

* `conditions` -
  Output only. The reason(s) why a trigger is in FAILED state.
  
* `create_time` -
  Output only. The creation time.
  
* `etag` -
  Output only. This checksum is computed by the server based on the value of other fields, and may be sent only on create requests to ensure the client has an up-to-date value before proceeding.
  
* `uid` -
  Output only. Server assigned unique identifier for the trigger. The value is a UUID4 string and guaranteed to remain unchanged until the resource is deleted.
  
* `update_time` -
  Output only. The last-modified time.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Trigger can be imported using any of these accepted formats:

```
$ terraform import google_eventarc_trigger.default projects/{{project}}/locations/{{location}}/triggers/{{name}}
$ terraform import google_eventarc_trigger.default {{project}}/{{location}}/{{name}}
$ terraform import google_eventarc_trigger.default {{location}}/{{name}}
```



