---
subcategory: "Firebase"
page_title: "Google: google_firebase_hosting_channel"
description: |-
  A Google Cloud Firebase Hosting Channel instance
---

# google_firebase_hosting_channel

A Google Cloud Firebase Hosting Channel instance

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.


## Argument Reference

The following arguments are supported:

* `site_id` - 
  (Required)
  The ID of the site this channel belongs to.

* `channel_id` - 
  (Required)
  The ID of the channel. Use `channel_id = "live"` for the default channel of a site.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - An identifier for the resource with format `sites/{{site_id}}/channels/{{channel_id}}`. Same as `name`

* `name` - The fully-qualified resource name for the channel, in the format: `sites/{{site_id}}/channels/{{channel_id}}`.