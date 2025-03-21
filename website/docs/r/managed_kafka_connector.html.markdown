---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
#
# ----------------------------------------------------------------------------
#
#     This code is generated by Magic Modules using the following:
#
#     Configuration: https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/managedkafka/Connector.yaml
#     Template:      https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.html.markdown.tmpl
#
#     DO NOT EDIT this file directly. Any changes made to this file will be
#     overwritten during the next generation cycle.
#
# ----------------------------------------------------------------------------
subcategory: "Managed Kafka"
description: |-
  A Managed Service for Kafka Connect Connectors.
---

# google_managed_kafka_connector

A Managed Service for Kafka Connect Connectors.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.


<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=managedkafka_connector_basic&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Managedkafka Connector Basic


```hcl
resource "google_compute_network" "mkc_network" {
  name                    = "my-network-0"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "mkc_subnet" {
  name          = "my-subnetwork-0"
  ip_cidr_range = "10.4.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.mkc_network.id
}

resource "google_compute_subnetwork" "mkc_additional_subnet" {
  name          = "my-additional-subnetwork-0"
  ip_cidr_range = "10.5.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.mkc_network.id
}

resource "google_pubsub_topic" "cps_topic" {
  name = "my-cps-topic"
  
  message_retention_duration = "86600s"
}

resource "google_managed_kafka_cluster" "gmk_cluster" {
  cluster_id = "my-cluster"
  location = "us-central1"
  capacity_config {
    vcpu_count = 3
    memory_bytes = 3221225472
  }
  gcp_config {
    access_config {
      network_configs {
        subnet = "projects/${data.google_project.project.project_id}/regions/us-central1/subnetworks/${google_compute_subnetwork.mkc_subnet.id}"
      }
    }
  }
}

resource "google_managed_kafka_topic" "gmk_topic" {
  topic_id = "my-topic"
  cluster = google_managed_kafka_cluster.gmk_cluster.cluster_id
  location = "us-central1"
  partition_count = 2
  replication_factor = 3
}

resource "google_managed_kafka_connect_cluster" "mkc_cluster" {
  connect_cluster_id = "my-connect-cluster"
  kafka_cluster = "projects/${data.google_project.project.project_id}/locations/us-central1/clusters/${google_managed_kafka_cluster.gmk_cluster.cluster_id}"
  location = "us-central1"
  capacity_config {
    vcpu_count = 12
    memory_bytes = 21474836480
  }
  gcp_config {
    access_config {
      network_configs {
        primary_subnet = "projects/${data.google_project.project.project_id}/regions/us-central1/subnetworks/${google_compute_subnetwork.mkc_subnet.id}"
        additional_subnets = ["${google_compute_subnetwork.mkc_additional_subnet.id}"]
        dns_domain_names = ["${google_managed_kafka_cluster.gmk_cluster.cluster_id}.us-central1.managedkafka-staging.${data.google_project.project.project_id}.cloud-staging.goog"]
      }
    }
  }
  labels = {
    key = "value"
  }
}

resource "google_managed_kafka_connector" "example" {
  connector_id = "my-connector"
  connect_cluster = google_managed_kafka_connect_cluster.mkc_cluster.connect_cluster_id
  location = "us-central1"
  configs = {
    "connector.class" = "com.google.pubsub.kafka.sink.CloudPubSubSinkConnector"
    "name" = "my-connector"
    "tasks.max" = "1"
    "topics" = "${google_managed_kafka_topic.gmk_topic.topic_id}"
    "cps.topic" = "${google_pubsub_topic.cps_topic.name}"
    "cps.project" = "${data.google_project.project.project_id}"
    "value.converter" = "org.apache.kafka.connect.storage.StringConverter"
    "key.converter" = "org.apache.kafka.connect.storage.StringConverter"
  }
  task_restart_policy {
    minimum_backoff = "60s"
    maximum_backoff = "1800s"
  }
}

data "google_project" "project" {
}
```

## Argument Reference

The following arguments are supported:


* `location` -
  (Required)
  ID of the location of the Kafka Connect resource. See https://cloud.google.com/managed-kafka/docs/locations for a list of supported locations.

* `connect_cluster` -
  (Required)
  The connect cluster name.

* `connector_id` -
  (Required)
  The ID to use for the connector, which will become the final component of the connector's name. This value is structured like: `my-connector-id`.


- - -


* `configs` -
  (Optional)
  Connector config as keys/values. The keys of the map are connector property names, for example: `connector.class`, `tasks.max`, `key.converter`.

* `task_restart_policy` -
  (Optional)
  A policy that specifies how to restart the failed connectors/tasks in a Cluster resource. If not set, the failed connectors/tasks won't be restarted.
  Structure is [documented below](#nested_task_restart_policy).

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


<a name="nested_task_restart_policy"></a>The `task_restart_policy` block supports:

* `minimum_backoff` -
  (Optional)
  The minimum amount of time to wait before retrying a failed task. This sets a lower bound for the backoff delay.
  A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".

* `maximum_backoff` -
  (Optional)
  The maximum amount of time to wait before retrying a failed task. This sets an upper bound for the backoff delay.
  A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/connectClusters/{{connect_cluster}}/connectors/{{connector_id}}`

* `name` -
  The name of the connector. The `connector` segment is used when connecting directly to the connect cluster. Structured like: `projects/PROJECT_ID/locations/LOCATION/connectClusters/CONNECT_CLUSTER/connectors/CONNECTOR_ID`.

* `state` -
  The current state of the connect. Possible values: `STATE_UNSPECIFIED`, `UNASSIGNED`, `RUNNING`, `PAUSED`, `FAILED`, `RESTARTING`, and `STOPPED`.


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 60 minutes.
- `update` - Default is 30 minutes.
- `delete` - Default is 30 minutes.

## Import


Connector can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{location}}/connectClusters/{{connect_cluster}}/connectors/{{connector_id}}`
* `{{project}}/{{location}}/{{connect_cluster}}/{{connector_id}}`
* `{{location}}/{{connect_cluster}}/{{connector_id}}`


In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Connector using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{location}}/connectClusters/{{connect_cluster}}/connectors/{{connector_id}}"
  to = google_managed_kafka_connector.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Connector can be imported using one of the formats above. For example:

```
$ terraform import google_managed_kafka_connector.default projects/{{project}}/locations/{{location}}/connectClusters/{{connect_cluster}}/connectors/{{connector_id}}
$ terraform import google_managed_kafka_connector.default {{project}}/{{location}}/{{connect_cluster}}/{{connector_id}}
$ terraform import google_managed_kafka_connector.default {{location}}/{{connect_cluster}}/{{connector_id}}
```

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
