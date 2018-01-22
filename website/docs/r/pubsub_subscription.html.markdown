---
layout: "google"
page_title: "Google: google_pubsub_subscription"
sidebar_current: "docs-google-pubsub-subscription"
description: |-
  Creates a subscription in Google's pubsub  queueing system
---

# google\_pubsub\_subscription

Creates a subscription in Google's pubsub queueing system. For more information see
[the official documentation](https://cloud.google.com/pubsub/docs) and
[API](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.subscriptions).


## Example Usage

```hcl
resource "google_pubsub_topic" "default-topic" {
  name = "default-topic"
}

resource "google_pubsub_subscription" "default" {
  name  = "default-subscription"
  topic = "${google_pubsub_topic.default-topic.name}"

  ack_deadline_seconds = 20

  push_config {
    push_endpoint = "https://example.com/push"

    attributes {
      x-goog-version = "v1"
    }
  }
}
```

If the subscription has a topic in a different project:

```hcl
resource "google_pubsub_topic" "topic-different-project" {
  project = "another-project"
  name = "topic-different-project"
}

resource "google_pubsub_subscription" "default" {
  name  = "default-subscription"
  topic = "${google_pubsub_topic.topic-different-project.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by pubsub.
    Changing this forces a new resource to be created.

* `topic` - (Required) The topic name or id to bind this subscription to, required by pubsub.
    Changing this forces a new resource to be created.

- - -

* `ack_deadline_seconds` - (Optional) The maximum number of seconds a
    subscriber has to acknowledge a received message, otherwise the message is
    redelivered. Changing this forces a new resource to be created.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `push_config` - (Optional) Block configuration for push options. More
    configuration options are detailed below.

The optional `push_config` block supports:

* `push_endpoint` - (Required) The URL of the endpoint to which messages should
    be pushed. Changing this forces a new resource to be created.

* `attributes` - (Optional) Key-value pairs of API supported attributes used
    to control aspects of the message delivery. Currently, only
    `x-goog-version` is supported, which controls the format of the data
    delivery. For more information, read [the API docs
    here](https://cloud.google.com/pubsub/reference/rest/v1/projects.subscriptions#PushConfig.FIELDS.attributes).
    Changing this forces a new resource to be created.

## Attributes Reference

* `path` - Path of the subscription in the format `projects/{project}/subscriptions/{sub}`

## Import

Pubsub subscription can be imported using the `name`, e.g.

```
$ terraform import google_pubsub_subscription.default default-subscription
```

