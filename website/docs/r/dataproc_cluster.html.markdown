---
layout: "google"
page_title: "Google: google_dataproc_cluster"
sidebar_current: "docs-google-dataproc-cluster"
description: |-
  Manages a Cloud Dataproc cluster resource.
---

# google\_dataproc\_cluster

Manages a Cloud Dataproc cluster resource within GCP. For more information see
[the official dataproc documentation](https://cloud.google.com/dataproc/).


!> **Warning:** Due to limitations of the API, all arguments except
`labels`,`cluster_config.worker_config.num_instances` and `cluster_config.preemptible_worker_config.num_instances` are non-updateable. Changing others will cause recreation of the
whole cluster!

## Example Usage - Basic

```hcl
resource "google_dataproc_cluster" "simplecluster" {
    name       = "simplecluster"
    region     = "us-central1"
}
```

## Example Usage - Advanced

```hcl
resource "google_dataproc_cluster" "mycluster" {
    name       = "mycluster"
    region     = "us-central1"
    labels {
        foo = "bar"
    }

    cluster_config {
        staging_bucket        = "dataproc-staging-bucket"

        master_config {
            num_instances     = 1
            machine_type      = "n1-standard-1"
            disk_config {
                boot_disk_type = "pd-ssd"
                boot_disk_size_gb = 10
            }
        }

        worker_config {
            num_instances     = 2
            machine_type      = "n1-standard-1"
            disk_config {
                boot_disk_size_gb = 10
                num_local_ssds    = 1
            }
        }

        preemptible_worker_config {
            num_instances     = 0
        }

        # Override or set some custom properties
        software_config {
            image_version       = "1.3.7-deb9"
            override_properties = {
                "dataproc:dataproc.allow.zero.workers" = "true"
            }
        }

        gce_cluster_config {
            #network = "${google_compute_network.dataproc_network.name}"
            tags    = ["foo", "bar"]
        }

        # You can define multiple initialization_action blocks
        initialization_action {
            script      = "gs://dataproc-initialization-actions/stackdriver/stackdriver.sh"
            timeout_sec = 500
        }

    }
}
```

## Example Usage - Using a GPU accelerator

```hcl
resource "google_dataproc_cluster" "accelerated_cluster" {
    name   = "my-cluster-with-gpu"
    region = "us-central1"

    cluster_config {
        gce_cluster_config {
            zone = "us-central1-a"
        }

        master_config {
            accelerators {
                accelerator_type  = "nvidia-tesla-k80"
                accelerator_count = "1"
            }
        }
    }
}
```

## Argument Reference

* `name` - (Required) The name of the cluster, unique within the project and
	zone.

- - -

* `project` - (Optional) The ID of the project in which the `cluster` will exist. If it
	is not provided, the provider project is used.

* `region` - (Optional) The region in which the cluster and associated nodes will be created in.
   Defaults to `global`.

* `labels` - (Optional, Computed) The list of labels (key/value pairs) to be applied to
   instances in the cluster. GCP generates some itself including `goog-dataproc-cluster-name`
   which is the name of the cluster.

* `cluster_config` - (Optional) Allows you to configure various aspects of the cluster.
   Structure defined below.

- - -

The `cluster_config` block supports:

```hcl
    cluster_config {
        gce_cluster_config        { ... }
        master_config             { ... }
        worker_config             { ... }
        preemptible_worker_config { ... }
        software_config           { ... }

        # You can define multiple initialization_action blocks
        initialization_action     { ... }
    }
```

* `staging_bucket` - (Optional) The Cloud Storage staging bucket used to stage files,
   such as Hadoop jars, between client machines and the cluster.
   Note: If you don't explicitly specify a `staging_bucket`
   then GCP will auto create / assign one for you. However, you are not guaranteed
   an auto generated bucket which is solely dedicated to your cluster; it may be shared
   with other clusters in the same region/zone also choosing to use the auto generation
   option.

* `gce_cluster_config` (Optional) Common config settings for resources of Google Compute Engine cluster
   instances, applicable to all instances in the cluster. Structure defined below.

* `master_config` (Optional) The Google Compute Engine config settings for the master instances
   in a cluster.. Structure defined below.

* `worker_config` (Optional) The Google Compute Engine config settings for the worker instances
   in a cluster.. Structure defined below.

* `preemptible_worker_config` (Optional) The Google Compute Engine config settings for the additional (aka
   preemptible) instancesin a cluster. Structure defined below.

* `software_config` (Optional) The config settings for software inside the cluster.
   Structure defined below.

* `initialization_action` (Optional) Commands to execute on each node after config is completed.
   You can specify multiple versions of these. Structure defined below.

- - -

The `cluster_config.gce_cluster_config` block supports:

```hcl
    cluster_config {
        gce_cluster_config {

            zone = "us-central1-a"

            # One of the below to hook into a custom network / subnetwork
            network    = "${google_compute_network.dataproc_network.name}"
            subnetwork = "${google_compute_network.dataproc_subnetwork.name}"

            tags    = ["foo", "bar"]
        }
    }
```

* `zone` - (Optional, Computed) The GCP zone where your data is stored and used (i.e. where
	the master and the worker nodes will be created in). If `region` is set to 'global' (default)
	then `zone` is mandatory, otherwise GCP is able to make use of [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/auto-zone)
	to determine this automatically for you.
	Note: This setting additionally determines and restricts
	which computing resources are available for use with other configs such as
	`cluster_config.master_config.machine_type` and `cluster_config.worker_config.machine_type`.

* `network` - (Optional, Computed) The name or self_link of the Google Compute Engine
	network to the cluster will be part of. Conflicts with `subnetwork`.
	If neither is specified, this defaults to the "default" network.

* `subnetwork` - (Optional) The name or self_link of the Google Compute Engine
   subnetwork the cluster will be part of. Conflicts with `network`.

* `service_account` - (Optional) The service account to be used by the Node VMs.
	If not specified, the "default" service account is used.

* `service_account_scopes` - (Optional, Computed) The set of Google API scopes to be made available
	on all of the node VMs under the `service_account` specified. These can be
	either FQDNs, or scope aliases. The following scopes are necessary to ensure
	the correct functioning of the cluster:

  * `useraccounts-ro` (`https://www.googleapis.com/auth/cloud.useraccounts.readonly`)
  * `storage-rw`      (`https://www.googleapis.com/auth/devstorage.read_write`)
  * `logging-write`   (`https://www.googleapis.com/auth/logging.write`)

* `tags` - (Optional) The list of instance tags applied to instances in the cluster.
   Tags are used to identify valid sources or targets for network firewalls.

* `internal_ip_only` - (Optional) By default, clusters are not restricted to internal IP addresses, 
   and will have ephemeral external IP addresses assigned to each instance. If set to true, all 
   instances in the cluster will only have internal IP addresses. Note: Private Google Access 
   (also known as `privateIpGoogleAccess`) must be enabled on the subnetwork that the cluster 
   will be launched in.

* `metadata` - (Optional) A map of the Compute Engine metadata entries to add to all instances
   (see [Project and instance metadata](https://cloud.google.com/compute/docs/storing-retrieving-metadata#project_and_instance_metadata)).

- - -

The `cluster_config.master_config` block supports:

```hcl
    cluster_config {
        master_config {
            num_instances     = 1
            machine_type      = "n1-standard-1"
            disk_config {
                boot_disk_type    = "pd-ssd"
                boot_disk_size_gb = 10
                num_local_ssds    = 1
            }
        }
    }
```

* `num_instances`- (Optional, Computed) Specifies the number of master nodes to create.
   If not specified, GCP will default to a predetermined computed value (currently 1).

* `machine_type` - (Optional, Computed) The name of a Google Compute Engine machine type
   to create for the master. If not specified, GCP will default to a predetermined
   computed value (currently `n1-standard-4`).

* `image_uri` (Optional) The URI for the image to use for this worker.  See [the guide](https://cloud.google.com/dataproc/docs/guides/dataproc-images)
    for more information.

* `disk_config` (Optional) Disk Config

	* `boot_disk_type` - (Optional) The disk type of the primary disk attached to each node.
	One of `"pd-ssd"` or `"pd-standard"`. Defaults to `"pd-standard"`.

	* `boot_disk_size_gb` - (Optional, Computed) Size of the primary disk attached to each node, specified
	in GB. The primary disk contains the boot volume and system libraries, and the
	smallest allowed disk size is 10GB. GCP will default to a predetermined
	computed value if not set (currently 500GB). Note: If SSDs are not
	attached, it also contains the HDFS data blocks and Hadoop working directories.

	* `num_local_ssds` - (Optional) The amount of local SSD disks that will be
	attached to each master cluster node. Defaults to 0.

* `accelerators` (Optional) The Compute Engine accelerator (GPU) configuration for these instances. Can be specified multiple times.

    * `accelerator_type` - (Required) The short name of the accelerator type to expose to this instance. For example, `nvidia-tesla-k80`.

    * `accelerator_count` - (Required) The number of the accelerator cards of this type exposed to this instance. Often restricted to one of `1`, `2`, `4`, or `8`.

~> The Cloud Dataproc API can return unintuitive error messages when using accelerators; even when you have defined an accelerator, Auto Zone Placement does not exclusively select
zones that have that accelerator available. If you get a 400 error that the accelerator can't be found, this is a likely cause. Make sure you check [accelerator availability by zone](https://cloud.google.com/compute/docs/reference/rest/v1/acceleratorTypes/list)
if you are trying to use accelerators in a given zone.

- - -

The `cluster_config.worker_config` block supports:

```hcl
    cluster_config {
        worker_config {
            num_instances     = 3
            machine_type      = "n1-standard-1"
            disk_config {
                boot_disk_type    = "pd-standard"
                boot_disk_size_gb = 10
                num_local_ssds    = 1
            }
        }
    }
```

* `num_instances`- (Optional, Computed) Specifies the number of worker nodes to create.
   If not specified, GCP will default to a predetermined computed value (currently 2).
   There is currently a beta feature which allows you to run a
   [Single Node Cluster](https://cloud.google.com/dataproc/docs/concepts/single-node-clusters).
   In order to take advantage of this you need to set
   `"dataproc:dataproc.allow.zero.workers" = "true"` in
   `cluster_config.software_config.properties`

* `machine_type` - (Optional, Computed) The name of a Google Compute Engine machine type
   to create for the worker nodes. If not specified, GCP will default to a predetermined
   computed value (currently `n1-standard-4`).

* `disk_config` (Optional) Disk Config

    * `boot_disk_type` - (Optional) The disk type of the primary disk attached to each node.
	One of `"pd-ssd"` or `"pd-standard"`. Defaults to `"pd-standard"`.

    * `boot_disk_size_gb` - (Optional, Computed) Size of the primary disk attached to each worker node, specified
    in GB. The smallest allowed disk size is 10GB. GCP will default to a predetermined
    computed value if not set (currently 500GB). Note: If SSDs are not
	attached, it also contains the HDFS data blocks and Hadoop working directories.

    * `num_local_ssds` - (Optional) The amount of local SSD disks that will be
	attached to each worker cluster node. Defaults to 0.

* `image_uri` (Optional) The URI for the image to use for this worker.  See [the guide](https://cloud.google.com/dataproc/docs/guides/dataproc-images)
    for more information.

* `accelerators` (Optional) The Compute Engine accelerator configuration for these instances. Can be specified multiple times.

    * `accelerator_type` - (Required) The short name of the accelerator type to expose to this instance. For example, `nvidia-tesla-k80`.

    * `accelerator_count` - (Required) The number of the accelerator cards of this type exposed to this instance. Often restricted to one of `1`, `2`, `4`, or `8`.

~> The Cloud Dataproc API can return unintuitive error messages when using accelerators; even when you have defined an accelerator, Auto Zone Placement does not exclusively select
zones that have that accelerator available. If you get a 400 error that the accelerator can't be found, this is a likely cause. Make sure you check [accelerator availability by zone](https://cloud.google.com/compute/docs/reference/rest/v1/acceleratorTypes/list)
if you are trying to use accelerators in a given zone.

- - -

The `cluster_config.preemptible_worker_config` block supports:

```hcl
    cluster_config {
        preemptible_worker_config {
            num_instances     = 1
            disk_config {
                boot_disk_size_gb = 10
            }
        }
    }
```

Note: Unlike `worker_config`, you cannot set the `machine_type` value directly. This
will be set for you based on whatever was set for the `worker_config.machine_type` value.

* `num_instances`- (Optional) Specifies the number of preemptible nodes to create.
   Defaults to 0.

* `disk_config` (Optional) Disk Config

    * `boot_disk_size_gb` - (Optional, Computed) Size of the primary disk attached to each preemptible worker node, specified
    in GB. The smallest allowed disk size is 10GB. GCP will default to a predetermined
    computed value if not set (currently 500GB). Note: If SSDs are not
	attached, it also contains the HDFS data blocks and Hadoop working directories.

- - -

The `cluster_config.software_config` block supports:

```hcl
    cluster_config {
        # Override or set some custom properties
        software_config {
            image_version       = "1.3.7-deb9"
            override_properties = {
                "dataproc:dataproc.allow.zero.workers" = "true"
            }
        }
    }
```

* `image_version` - (Optional, Computed) The Cloud Dataproc image version to use
   for the cluster - this controls the sets of software versions
   installed onto the nodes when you create clusters. If not specified, defaults to the
   latest version. For a list of valid versions see
   [Cloud Dataproc versions](https://cloud.google.com/dataproc/docs/concepts/dataproc-versions)

* `override_properties` - (Optional) A list of override and additional properties (key/value pairs)
   used to modify various aspects of the common configuration files used when creating
   a cluster. For a list of valid properties please see
  [Cluster properties](https://cloud.google.com/dataproc/docs/concepts/cluster-properties)

- - -

The `initialization_action` block (Optional) can be specified multiple times and supports:

```hcl
    cluster_config {
        # You can define multiple initialization_action blocks
        initialization_action {
            script      = "gs://dataproc-initialization-actions/stackdriver/stackdriver.sh"
            timeout_sec = 500
        }
    }
```

* `script`- (Required) The script to be executed during initialization of the cluster.
   The script must be a GCS file with a gs:// prefix.

* `timeout_sec` - (Optional, Computed) The maximum duration (in seconds) which `script` is
   allowed to take to execute its action. GCP will default to a predetermined
   computed value if not set (currently 300).

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `cluster_config.master_config.instance_names` - List of master instance names which
   have been assigned to the cluster.

* `cluster_config.worker_config.instance_names` - List of worker instance names which have been assigned
	to the cluster.

* `cluster_config.preemptible_worker_config.instance_names` - List of preemptible instance names which have been assigned
	to the cluster.

* `cluster_config.bucket` - The name of the cloud storage bucket ultimately used to house the staging data
   for the cluster. If `staging_bucket` is specified, it will contain this value, otherwise
   it will be the auto generated name.

* `cluster_config.software_config.properties` - A list of the properties used to set the daemon config files.
   This will include any values supplied by the user via `cluster_config.software_config.override_properties`

## Timeouts

`google_dataproc_cluster` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating clusters.
- `update` - (Default `5 minutes`) Used for updating clusters
- `delete` - (Default `5 minutes`) Used for destroying clusters.
