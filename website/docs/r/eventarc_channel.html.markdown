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
page_title: "Google: google_eventarc_channel"
description: |-
  The Eventarc Channel resource
---

# google_eventarc_channel

The Eventarc Channel resource

## Example Usage - basic
```hcl
data "google_project" "test_project" {
	project_id  = "my-project-name"
}

data "google_kms_key_ring" "test_key_ring" {
	name     = "keyring"
	location = "us-west1"
}
  
data "google_kms_crypto_key" "key" {
	name     = "key"
	key_ring = data.google_kms_key_ring.test_key_ring.id
}

resource "google_kms_crypto_key_iam_binding" "key1_binding" {
    crypto_key_id = data.google_kms_crypto_key.key1.id
    role      = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  
    members = [
    "serviceAccount:service-${data.google_project.test_project.number}@gcp-sa-eventarc.iam.gserviceaccount.com",
    ]
}

resource "google_eventarc_channel" "primary" {
  location = "us-west1"
  name     = "channel"
  project  = "${data.google_project.test_project.project_id}"
  crypto_key_name =  "${data.google_kms_crypto_key.key1.id}"
  third_party_provider = "projects/${data.google_project.test_project.project_id}/locations/us-west1/providers/datadog"
  depends_on = [google_kms_crypto_key_iam_binding.key1_binding]
}
```

## Argument Reference

The following arguments are supported:

* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  Required. The resource name of the channel. Must be unique within the location on the project.
  


- - -

* `crypto_key_name` -
  (Optional)
  Optional. Resource name of a KMS crypto key (managed by the user) used to encrypt/decrypt their event data. It must match the pattern `projects/*/locations/*/keyRings/*/cryptoKeys/*`.
  
* `project` -
  (Optional)
  The project for the resource
  
* `third_party_provider` -
  (Optional)
  The name of the event provider (e.g. Eventarc SaaS partner) associated with the channel. This provider will be granted permissions to publish events to the channel. Format: `projects/{project}/locations/{location}/providers/{provider_id}`.
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/channels/{{name}}`

* `activation_token` -
  Output only. The activation token for the channel. The token must be used by the provider to register the channel for publishing.
  
* `create_time` -
  Output only. The creation time.
  
* `pubsub_topic` -
  Output only. The name of the Pub/Sub topic created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{project}/topics/{topic_id}`.
  
* `state` -
  Output only. The state of a Channel. Possible values: STATE_UNSPECIFIED, PENDING, ACTIVE, INACTIVE
  
* `uid` -
  Output only. Server assigned unique identifier for the channel. The value is a UUID4 string and guaranteed to remain unchanged until the resource is deleted.
  
* `update_time` -
  Output only. The last-modified time.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Channel can be imported using any of these accepted formats:

```
$ terraform import google_eventarc_channel.default projects/{{project}}/locations/{{location}}/channels/{{name}}
$ terraform import google_eventarc_channel.default {{project}}/{{location}}/{{name}}
$ terraform import google_eventarc_channel.default {{location}}/{{name}}
```



