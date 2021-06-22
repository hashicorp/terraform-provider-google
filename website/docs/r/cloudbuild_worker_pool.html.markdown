---
subcategory: "Cloud Build"
layout: "google"
page_title: "Google: google_cloudbuild_worker_pool"
sidebar_current: "docs-google-cloudbuild-worker-pool"
description: |-
  Configuration for custom WorkerPool to run builds
---

# google\_cloudbuild\_worker\_pool

Definition of custom Cloud Build WorkerPools for running jobs with custom configuration and custom networking.

-> This resource is not currently public, and requires allow-listing of projects prior to use.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.


## Example Usage

```hcl
resource "google_cloudbuild_worker_pool" "pool" {
  name = "my-pool"
  location = "europe-west1"
  worker_config {
    disk_size_gb = 100
    machine_type = "e2-standard-4"
    no_external_ip = false
  }
}
```

## Example Usage - Network Config

```hcl
resource "google_project_service" "servicenetworking" {
  service = "servicenetworking.googleapis.com"
  disable_on_destroy = false
}

resource "google_compute_network" "network" {
  name                    = "my-network"
  auto_create_subnetworks = false
  depends_on = [google_project_service.servicenetworking]
}

resource "google_compute_global_address" "worker_range" {
  name          = "worker-pool-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.network.id
}

resource "google_service_networking_connection" "worker_pool_conn" {
  network                 = google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.worker_range.name]
  depends_on              = [google_project_service.servicenetworking]
}

resource "google_cloudbuild_worker_pool" "pool" {
  name = "my-pool"
  location = "europe-west1"
  worker_config {
    disk_size_gb = 100
    machine_type = "e2-standard-4"
    no_external_ip = false
  }
  network_config {
    peered_network = google_compute_network.network.id
  }
  depends_on = [google_service_networking_connection.worker_pool_conn]
}
```

## Argument Reference

The following arguments are supported:

* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  User-defined name of the `WorkerPool`.
  

- - -

* `network_config` -
  (Optional)
  Network configuration for the `WorkerPool`.
  
* `project` -
  (Optional)
  The project for the resource
  
* `worker_config` -
  (Optional)
  Configuration to be used for a creating workers in the `WorkerPool`.
  


The `network_config` block supports:
    
* `peered_network` -
  (Required)
  Immutable. The network definition that the workers are peered to. If this section is left empty, the workers will be peered to `WorkerPool.project_id` on the service producer network. Must be in the format `projects/{project}/global/networks/{network}`, where `{project}` is a project number, such as `12345`, and `{network}` is the name of a VPC network in the project. See (https://cloud.google.com/cloud-build/docs/custom-workers/set-up-custom-worker-pool-environment#understanding_the_network_configuration_options)
    
The `worker_config` block supports:
    
* `disk_size_gb` -
  (Optional)
  Size of the disk attached to the worker, in GB. See (https://cloud.google.com/cloud-build/docs/custom-workers/worker-pool-config-file). Specify a value of up to 1000. If `0` is specified, Cloud Build will use a standard disk size.
    
* `machine_type` -
  (Optional)
  Machine type of a worker, such as `n1-standard-1`. See (https://cloud.google.com/cloud-build/docs/custom-workers/worker-pool-config-file). If left blank, Cloud Build will use `n1-standard-1`.
    
* `no_external_ip` -
  (Optional)
  If true, workers are created without any public address, which prevents network egress to public IPs.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/workerPools/{{name}}`

* `create_time` -
  Output only. Time at which the request to create the `WorkerPool` was received.
  
* `delete_time` -
  Output only. Time at which the request to delete the `WorkerPool` was received.
  
* `state` -
  Output only. WorkerPool state. Possible values: STATE_UNSPECIFIED, PENDING, APPROVED, REJECTED, CANCELLED
  
* `update_time` -
  Output only. Time at which the request to update the `WorkerPool` was received.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `update` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

WorkerPool can be imported using any of these accepted formats:

```
$ terraform import google_cloudbuild_worker_pool.default projects/{{project}}/locations/{{location}}/workerPools/{{name}}
$ terraform import google_cloudbuild_worker_pool.default {{project}}/{{location}}/{{name}}
$ terraform import google_cloudbuild_worker_pool.default {{location}}/{{name}}
```

