---
subcategory: "Firebase"
page_title: "Google: google_firebase_android_app"
description: |-
  A Google Cloud Firebase Android application instance
---

# google\_firebase\_android\_app

A Google Cloud Firebase Android application instance

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.


## Argument Reference

The following arguments are supported:


* `app_id` -
  (Required)
  The app_ip of name of the Firebase androidApp.


- - -


* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{name}}`

* `name` -
  The fully qualified resource name of the App, for example:
  projects/projectId/androidApps/appId

* `app_id` -
  Immutable. The globally unique, Firebase-assigned identifier of the App.
  This identifier should be treated as an opaque token, as the data format is not specified.

* `display_name` -
  The user-assigned display name of the App.

* `package_name` -
  The canonical package name of the Android app as would appear in the Google Play Developer Console.
