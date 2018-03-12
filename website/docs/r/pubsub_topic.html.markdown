---
layout: "google"
page_title: "Google: google_pubsub_topic"
sidebar_current: "docs-google-pubsub-topic-x"
description: |-
  Creates a topic in Google's pubsub queueing system
---

# google\_pubsub\_topic

Creates a topic in Google's pubsub queueing system. For more information see
[the official documentation](https://cloud.google.com/pubsub/docs) and
[API](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics).


## Example Usage

```hcl
resource "google_pubsub_topic" "mytopic" {
  name = "default-topic"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the pubsub topic.
    Changing this forces a new resource to be created.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

Only the arguments listed above are exposed as attributes.

## Import

Pubsub topics can be imported using the `name` or full topic id, e.g.

```
$ terraform import google_pubsub_topic.mytopic default-topic
```
```
$ terraform import google_pubsub_topic.mytopic projects/my-gcp-project/topics/default-topic
```
When importing using only the name, the provider project must be set.
