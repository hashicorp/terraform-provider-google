---
subcategory: "Cloud Pub/Sub"
page_title: "Google: google_pubsub_topic"
description: |-
  Get information about a Google Cloud Pub/Sub Topic.
---

# google\_pubsub\_topic

Get information about a Google Cloud Pub/Sub Topic. For more information see
the [official documentation](https://cloud.google.com/pubsub/docs/)
and [API](https://cloud.google.com/pubsub/docs/apis).

## Example Usage

```hcl
data "google_pubsub_topic" "my-pubsub-topic" {
  name = "my-pubsub-topic"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Pub/Sub Topic.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_pubsub_topic](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/pubsub_topic#argument-reference) resource for details of the available attributes.
