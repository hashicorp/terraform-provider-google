---
subcategory: "Cloud Pub/Sub"
page_title: "Google: google_pubsub_subscription"
description: |-
  Get information about a Google Cloud Pub/Sub Subscription.
---

# google\_pubsub\_subscription

Get information about a Google Cloud Pub/Sub Subscription. For more information see
the [official documentation](https://cloud.google.com/pubsub/docs/)
and [API](https://cloud.google.com/pubsub/docs/apis).

## Example Usage

```hcl
data "google_pubsub_subscription" "my-pubsub-subscription" {
  name = "my-pubsub-subscription"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Pub/Sub Subscription.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

See [google_pubsub_subscription](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/pubsub_subscription#argument-reference) resource for details of the available attributes.
