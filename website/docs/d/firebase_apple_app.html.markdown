---
subcategory: "Firebase"
page_title: "Google: google_firebase_apple_app"
description: |-
  A Google Cloud Firebase Apple application instance
---

# google\_firebase\_apple\_app

A Google Cloud Firebase Apple application instance

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.


## Argument Reference

The following arguments are supported:


* `app_id` -
  (Required)
  The app_id of name of the Firebase iosApp.


- - -


* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{name}}`

* `name` -
  The fully qualified resource name of the App, for example:
  projects/projectId/iosApps/appId

* `app_id` -
  Immutable. The globally unique, Firebase-assigned identifier of the App.
  This identifier should be treated as an opaque token, as the data format is not specified.

* `display_name` -
  The user-assigned display name of the App.

* `bundle_id` -
  The canonical bundle ID of the Apple app as it would appear in the Apple AppStore.

* `app_store_id` -
  The automatically generated Apple ID assigned to the Apple app by Apple in the Apple App Store.

* `team_id` -
  The Apple Developer Team ID associated with the App in the App Store.
