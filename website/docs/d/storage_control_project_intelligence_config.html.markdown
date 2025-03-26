---
subcategory: "Cloud Storage Control"
description: |-
  Gets a Project Storage Intelligence config.
---

# google_storage_control_project_intelligence_config

Use this data source to get information about a Project Storage Intelligence config resource.
See [the official documentation](https://cloud.google.com/storage/docs/storage-intelligence/overview#resource)
and
[API](https://cloud.google.com/storage/docs/json_api/v1/intelligenceConfig).


## Example Usage

```hcl
data "google_storage_control_project_intelligence_config" "sample-config" {
  name = "my-project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name or number of the GCP project.


## Attributes Reference

See [google_storage_control_project_intelligence_config](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/storage_control_project_intelligence_config#argument-reference) resource for details of the available attributes.
