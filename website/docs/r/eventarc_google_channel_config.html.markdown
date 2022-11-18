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
page_title: "Google: google_eventarc_google_channel_config"
description: |-
  The Eventarc GoogleChannelConfig resource
---

# google_eventarc_google_channel_config

The Eventarc GoogleChannelConfig resource

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

resource "google_eventarc_google_channel_config" "primary" {
  location = "us-west1"
  name     = "channel"
  project  = "${data.google_project.test_project.project_id}"
  crypto_key_name =  "${data.google_kms_crypto_key.key1.id}"
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
  Required. The resource name of the config. Must be in the format of, `projects/{project}/locations/{location}/googleChannelConfig`.
  


- - -

* `crypto_key_name` -
  (Optional)
  Optional. Resource name of a KMS crypto key (managed by the user) used to encrypt/decrypt their event data. It must match the pattern `projects/*/locations/*/keyRings/*/cryptoKeys/*`.
  
* `project` -
  (Optional)
  The project for the resource
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/googleChannelConfig`

* `update_time` -
  Output only. The last-modified time.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

GoogleChannelConfig can be imported using any of these accepted formats:

```
$ terraform import google_eventarc_google_channel_config.default projects/{{project}}/locations/{{location}}/googleChannelConfig
$ terraform import google_eventarc_google_channel_config.default {{project}}/{{location}}
$ terraform import google_eventarc_google_channel_config.default {{location}}
```



