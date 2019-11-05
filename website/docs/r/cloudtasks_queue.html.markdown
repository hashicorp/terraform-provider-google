---
layout: "google"
page_title: "Google: google_cloudtasks_queue"
sidebar_current: "docs-google-cloudtasks-queue"
description: |-
 Configures task queues for the project
---

 # google\_cloudtasks\_queue

 Configures [Cloud Cloud Tasks Queues](https://cloud.google.com/tasks/docs/creating-queues)
for a project.

 For more information, see,
[the Project API documentation](https://cloud.google.com/tasks/docs/reference/rest/v2/projects.locations.queues).

 ## Example Usage

 ```hcl
resource "google_cloudtasks_queue" "foo" {
  name     = "foo"
  location = "us-central1"
  rate_limits {
    max_dispatches_per_second = 500
    max_concurrent_dispatches = 15
  }
  retry {
    max_attempts = 15
  }
}
```

 ## Argument Reference

 The following arguments are supported:

 * `name` - (Required) The name of the queue.

 * `project` - (Optional) The ID of the project in which the resource belongs. If it is not provided, the provider project is used.

 * `location` - (Required) The region where the queue will be created.

  - - -

 * `app_engine_routing_override` - (Optional) Overrides for task-level AppEngine routing.

 * `rate_limits` - (Optional) Rate limits for the task dispatches.

 * `retry` - (Optional) Settings that determine the retry behavior.

  - - -

 The `app_engine_routing_override` block supports:

 * `instance` - (Optional) App instance to which the task is sent.

 * `service` - (Optional) Service to which the task is sent.

 * `version` - (Optional) Version to which the task is sent.

 The `rate_limits` block supports:

 * `max_dispatches_per_second` - (Optional) The maximum rate at which tasks are dispatched from this queue. For App Engine queues, the maximum allowed value is 500.

 * `max_concurrent_dispatches` - (Optional) The maximum number of concurrent tasks that Cloud Tasks allows to be dispatched for this queue.

 The `rate_limits` block supports: 

  * `max_attempts` - (Optional) Number of attempts per task.

  * `max_doublings` - (Optional) The time between retries will double `max_doublings` times.

  * `max_backoff` - (Optional) A task will be scheduled for retry between `min_backoff` and `max_backoff` duration after it fails if the queue's RetryConfig specifies that the task should be retried. Defined in seconds (string ending in "s").

  * `min_backoff` - (Optional) A task will be scheduled for retry between `min_backoff` and `max_backoff` duration after it fails if the queue's RetryConfig specifies that the task should be retried. Defined in seconds (string ending in "s").

 ## Attributes Reference

  * `max_burst_size` - The max burst size that limits how fast tasks in queue are processwed. Cloud Tasks will pick the value based on the value of `max_dispatches_per_second`.

  * `app_engine_routing_override.0.host` - The host that the task is sent to. Constructed from the project ID, service, version, and instance in the format <app-id>.appspot.com.

 ## Import

 This resource can be imported using the project ID:

 `terraform import google_cloudtasks_queue.foo project-id/location-name/queue-name`
