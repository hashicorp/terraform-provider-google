---
subcategory: "Cloud Functions (2nd gen)"
page_title: "Google: google_cloudfunctions2_function"
description: |-
  Get information about a Google Cloud Function (2nd gen).
---

# google\_cloudfunctions2\_function

Get information about a Google Cloud Function (2nd gen). For more information see:

* [API documentation](https://cloud.google.com/functions/docs/reference/rest/v2beta/projects.locations.functions).

## Example Usage

```hcl
data "google_cloudfunctions2_function" "my-function" {
  name = "function"
  location = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Cloud Function (2nd gen).

* `location` - (Required) The location in which the resource belongs.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_cloudfunctions2_function](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudfunctions2_function) resource for details of all the available attributes.