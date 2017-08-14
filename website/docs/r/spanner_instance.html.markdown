---
layout: "google"
page_title: "Google: google_spanner_instance"
sidebar_current: "docs-google-spanner-instance"
description: |-
  Creates and manages a Google Spanner Instance.
---

# google\_spanner\_instance

Creates and manages a Google Spanner Instance. For more information, see the [official documentation](https://cloud.google.com/spanner/), or the [JSON API](https://cloud.google.com/spanner/docs/reference/rest/v1/projects.instances).

## Example Usage

Example creating a Spanner instance.

```hcl
resource "google_spanner_instance" "main" {
  config       = "regional-europe-west1"
  display_name = "main-instance"
  name         = "main-instance"
  num_nodes    = 1
}
```

## Argument Reference

The following arguments are supported:

* `config` - (Required) The name of the instance's configuration (similar but not
   quite the same as a region) which defines defines the geographic placement and
   replication of your databases in this instance. It determines where your data
   is stored. Values are typically of the form `regional-europe-west1` , `us-central` etc.
   In order to obtain a valid list please consult the
   [Configuration section of the docs](https://cloud.google.com/spanner/docs/instances).

* `display_name` - (Required) The descriptive name for this instance as it appears
   in UIs. Can be updated, however should be kept globally unique to avoid confusion.

- - -

* `name` - (Optional, Computed) The unique name (ID) of the instance. If the name is left
    blank, Terraform will randomly generate one when the instance is first
    created.

* `num_nodes` - (Optional, Computed) The number of nodes allocated to this instance.
   Defaults to `1`. This can be updated after creation.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `labels` - (Optional) A mapping (key/value pairs) of labels to assign to the instance.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `state` - The current state of the instance.

## Import

Instances can be imported using their `name` and optionally
the `project` in which it is defined (Often used when the project is different
to that defined in the provider), The format is thus either `{instanceId}` or
`{projectId}/{instanceId}`. e.g.

```
$ terraform import google_spanner_instance.master instance123

$ terraform import google_spanner_instance.master project123/instance456

```
