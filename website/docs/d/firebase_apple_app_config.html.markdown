---
subcategory: "Firebase"
description: |-
  A Google Cloud Firebase Apple application configuration
---

# google\_firebase\_apple\_app\_config

A Google Cloud Firebase Apple application configuration

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

To get more information about iosApp, see:

* [API documentation](https://firebase.google.com/docs/projects/api/reference/rest/v1beta1/projects.iosApps)
* How-to Guides
    * [Official Documentation](https://firebase.google.com/)


## Argument Reference
The following arguments are supported:

* `app_id` - (Required) The id of the Firebase iOS App.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `config_filename` -
  The filename that the configuration artifact for the IosApp is typically saved as.

* `config_file_contents` -
  The content of the XML configuration file as a base64-encoded string.
