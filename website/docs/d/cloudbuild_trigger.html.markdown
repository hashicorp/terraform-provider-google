---
subcategory: "Cloud Build"
page_title: "Google: google_cloudbuild_trigger"
description: |-
  Get information about a Google CloudBuild Trigger.
---

# google\_cloudbuild\_trigger

To get more information about Cloudbuild Trigger, see:

* [API documentation](https://cloud.google.com/build/docs/api/reference/rest/v1/projects.triggers)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/build/docs/automating-builds/create-manage-triggers)

## Example Usage

```hcl
data "google_cloudbuild_trigger" "name" {
  project = "your-project-id"
  trigger_id = google_cloudbuild_trigger.filename-trigger.trigger_id
  location = "location of trigger build"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

* `trigger_id` - (Required) The unique identifier for the trigger..
    
* `location` - (Required) The Cloud Build location for the trigger.

- - -

## Attributes Reference

See [google_cloudbuild_trigger](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudbuild_trigger#project) resource for details of the available attributes.
