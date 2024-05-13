---
subcategory: "Firebase"
description: |-
  A Google Cloud Firebase Android application instance
---

# google_firebase_android_app

A Google Cloud Firebase Android application instance

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.


## Argument Reference

The following arguments are supported:


* `app_id` -
  (Required)
  The app_id of name of the Firebase androidApp.


- - -


* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{name}}`

* `name` -
  The fully qualified resource name of the AndroidApp, for example:
  projects/projectId/androidApps/appId

* `app_id` -
  Immutable. The globally unique, Firebase-assigned identifier of the AndroidApp.
  This identifier should be treated as an opaque token, as the data format is not specified.

* `display_name` -
  The user-assigned display name of the AndroidApp.

* `package_name` -
  The canonical package name of the Android app as would appear in the Google Play Developer Console.

* `sha1_hashes` -
  The SHA1 certificate hashes for the AndroidApp.

* `sha256_hashes` -
  The SHA256 certificate hashes for the AndroidApp.

* `etag` -
  This checksum is computed by the server based on the value of other fields, and it may be sent
  with update requests to ensure the client has an up-to-date value before proceeding.
