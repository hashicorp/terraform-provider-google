---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_node_types"
sidebar_current: "docs-google-datasource-compute-node-types"
description: |-
  Provides list of available Google Compute Engine node types for
  sole-tenant nodes.
---

# google\_compute\_node\_types

Provides available node types for Compute Engine sole-tenant nodes in a zone
for a given project. For more information, see [the official documentation](https://cloud.google.com/compute/docs/nodes/#types) and [API](https://cloud.google.com/compute/docs/reference/rest/v1/nodeTypes).

## Example Usage

```hcl
data "google_compute_node_types" "central1b" {
  zone = "us-central1-b"
}

resource "google_compute_node_template" "tmpl" {
  name      = "terraform-test-tmpl"
  region    = "us-central1"
  node_type = data.google_compute_node_types.types.names[0]
}
```

## Argument Reference

The following arguments are supported:

* `zone` (Optional) - The zone to list node types for. Should be in zone of intended node groups and region of referencing node template. If `zone` is not specified, the provider-level zone must be set and is used
instead.

* `project` (Optional) - ID of the project to list available node types for.
Should match the project the nodes of this type will be deployed to.
Defaults to the project that the provider is authenticated with.

## Attributes Reference

The following attributes are exported:

* `names` - A list of node types available in the given zone and project.
