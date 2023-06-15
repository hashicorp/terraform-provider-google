---
subcategory: "Tags"
description: |-
  A LocationTagBinding represents a connection between a TagValue and a Regional cloud resources.
---

# google\_tags\_location\_tag\_binding

A TagBinding represents a connection between a TagValue and a Regional cloud resource (currently project, folder, or organization). Once a TagBinding is created, the TagValue is applied to all the descendants of the cloud resource.


To get more information about TagBinding, see:

* [API documentation](https://cloud.google.com/resource-manager/reference/rest/v3/tagBindings)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/resource-manager/docs/tags/tags-creating-and-managing)

## Example Usage

To bind a tag to a Cloud Run instance:

```hcl
resource "google_project" "project" {
	project_id = "project_id"
	name       = "project_id"
	org_id     = "123456789"
}

resource "google_tags_tag_key" "key" {
	parent      = "organizations/123456789"
	short_name  = "keyname"
	description = "For keyname resources."
}

resource "google_tags_tag_value" "value" {
	parent      = "tagKeys/${google_tags_tag_key.key.name}"
	short_name  = "valuename"
	description = "For valuename resources."
}

resource "google_tags_location_tag_binding" "binding" {
	parent    = "//run.googleapis.com/projects/${data.google_project.project.number}/locations/${google_cloud_run_service.default.location}/services/${google_cloud_run_service.default.name}"
	tag_value = "tagValues/${google_tags_tag_value.value.name}"
	location  = "us-central1"
}
```

To bind a (firewall) tag to compute instance:

```hcl
resource "google_project" "project" {
	project_id = "project_id"
	name       = "project_id"
	org_id     = "123456789"
}

resource "google_tags_tag_key" "key" {
	parent      = "organizations/123456789"
	short_name  = "keyname"
	description = "For keyname resources."
}

resource "google_tags_tag_value" "value" {
	parent      = "tagKeys/${google_tags_tag_key.key.name}"
	short_name  = "valuename"
	description = "For valuename resources."
}

resource "google_tags_location_tag_binding" "binding" {
	parent    = "//compute.googleapis.com/projects/${google_project.project.number}/zones/us-central1-a/instances/${google_compute_instance.instance.instance_id}"
	tag_value = "tagValues/${google_tags_tag_value.value.name}"
	location  = "us-central1"
}
```

## Argument Reference

The following arguments are supported:


* `parent` -
  (Required)
  The full resource name of the resource the TagValue is bound to. E.g. //cloudresourcemanager.googleapis.com/projects/123

* `tag_value` -
  (Required)
  The TagValue of the TagBinding. Must be of the form tagValues/456.

* `location` -
  (Required)
  Location of the resource.

- - -



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{location}}/{{name}}`

* `name` -
  The generated id for the TagBinding. This is a string of the form: `tagBindings/{parent}/{tag-value-name}`


## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


TagBinding can be imported using any of these accepted formats:

```
$ terraform import google_tags_location_tag_binding.default {{location}}/{{name}}
```
