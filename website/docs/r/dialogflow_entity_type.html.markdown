---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
#
# ----------------------------------------------------------------------------
#
#     This code is generated by Magic Modules using the following:
#
#     Configuration: https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/dialogflow/EntityType.yaml
#     Template:      https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.html.markdown.tmpl
#
#     DO NOT EDIT this file directly. Any changes made to this file will be
#     overwritten during the next generation cycle.
#
# ----------------------------------------------------------------------------
subcategory: "Dialogflow"
description: |-
  Represents an entity type.
---

# google_dialogflow_entity_type

Represents an entity type. Entity types serve as a tool for extracting parameter values from natural language queries.


To get more information about EntityType, see:

* [API documentation](https://cloud.google.com/dialogflow/docs/reference/rest/v2/projects.agent.entityTypes)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/dialogflow/docs/)

## Example Usage - Dialogflow Entity Type Basic


```hcl
resource "google_dialogflow_agent" "basic_agent" {
  display_name = "example_agent"
  default_language_code = "en"
  time_zone = "America/New_York"
}

resource "google_dialogflow_entity_type" "basic_entity_type" {
  depends_on = [google_dialogflow_agent.basic_agent]
  display_name = "basic-entity-type"
  kind = "KIND_MAP"
  entities {
    value = "value1"
    synonyms = ["synonym1","synonym2"]
  }
  entities {
    value = "value2"
    synonyms = ["synonym3","synonym4"]
  }
}
```

## Argument Reference

The following arguments are supported:


* `display_name` -
  (Required)
  The name of this entity type to be displayed on the console.

* `kind` -
  (Required)
  Indicates the kind of entity type.
  * KIND_MAP: Map entity types allow mapping of a group of synonyms to a reference value.
  * KIND_LIST: List entity types contain a set of entries that do not map to reference values. However, list entity
  types can contain references to other entity types (with or without aliases).
  * KIND_REGEXP: Regexp entity types allow to specify regular expressions in entries values.
  Possible values are: `KIND_MAP`, `KIND_LIST`, `KIND_REGEXP`.


* `enable_fuzzy_extraction` -
  (Optional)
  Enables fuzzy entity extraction during classification.

* `entities` -
  (Optional)
  The collection of entity entries associated with the entity type.
  Structure is [documented below](#nested_entities).

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.



<a name="nested_entities"></a>The `entities` block supports:

* `value` -
  (Required)
  The primary value associated with this entity entry. For example, if the entity type is vegetable, the value
  could be scallions.
  For KIND_MAP entity types:
  * A reference value to be used in place of synonyms.
  For KIND_LIST entity types:
  * A string that can contain references to other entity types (with or without aliases).

* `synonyms` -
  (Required)
  A collection of value synonyms. For example, if the entity type is vegetable, and value is scallions, a synonym
  could be green onions.
  For KIND_LIST entity types:
  * This collection must contain exactly one synonym equal to value.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{name}}`

* `name` -
  The unique identifier of the entity type.
  Format: projects/<Project ID>/agent/entityTypes/<Entity type ID>.


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


EntityType can be imported using any of these accepted formats:

* `{{name}}`


In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import EntityType using one of the formats above. For example:

```tf
import {
  id = "{{name}}"
  to = google_dialogflow_entity_type.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), EntityType can be imported using one of the formats above. For example:

```
$ terraform import google_dialogflow_entity_type.default {{name}}
```

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
