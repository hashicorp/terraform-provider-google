---
layout: "google"
page_title: "Google: google_dataproc_cluster"
sidebar_current: "docs-google-dataproc-cluster"
description: |-
  Manages a Cloud Dataproc cluster resource.
---

# google\_dataproc\_cluster

Manages a Cloud Dataproc cluster resource within GCE. For more information see
[the official dataproc documentation](https://cloud.google.com/dataproc/).


!> **Warning:** Due to limitations of the API, all arguments except
`labels`,`worker_config.num_workers` and `worker_config.preemptible_num_workers` are non-updateable. Changing any will cause recreation of the
whole cluster!

## Example usage

```hcl
resource "google_dataproc_cluster" "mycluster" {
    name   = "dproc-cluster-unique-name"
    region = "us-central1"

    master_config {
        num_masters       = 1
        machine_type      = "n1-standard-1"
        boot_disk_size_gb = 10
        num_local_ssds    = 1
    }

    worker_config {
    	num_workers             = 2
        machine_type            = "n1-standard-1"
        boot_disk_size_gb       = 10
        num_local_ssds          = 1

        preemptible_num_workers       = 1
        preemptible_boot_disk_size_gb = 10
    }

    initialization_action_timeout_sec = 500
    initialization_actions = [
       "gs://dataproc-initialization-actions/stackdriver/stackdriver.sh"
    ]
    
    labels {
      foo = "bar"
    }
    
    tags = ["foo", "bar"]

  }
}
```

## Argument Reference

* `name` - (Required) The name of the cluster, unique within the project and
    zone.

- - -

* `region` - (Optional) The region that the cluster and associated nodes will be created in.
   Defaults to `global`.

* `zone` - (Optional) The GCP zone where your data is stored and used (i.e. where
    the master and the worker nodes will be created in). If region is set to 'global'
    then `zone` IS mandatory, otherwise GCP is able to make use of [Auto Zone Placement](https://cloud.google.com/dataproc/docs/concepts/auto-zone)
    to determine this automatically for you.
    Note: This setting additionally determines and restricts
    which computing resources are available for use with other config such as
    `master_config.machine_type` and `worker_config.machine_type`.

* `staging_bucket` - (Optional) The Cloud Storage staging bucket used to stage files,
   such as Hadoop jars, between client machines and the cluster.
   Note: If you don't explicitly specify a `staging_bucket`
   then GCP will auto create / assign one for you, however you are NOT guaranteed to get
   an auto generated bucket which is solely dedicated to your cluster; it may be shared
   with other clusters in the same region/zone also choosing to use the auto generation
   option.

* `delete_autogen_bucket` (Optional) If this is set to true, upon destroying the cluster,
   if no explicit `staging_bucket` was specified (i.e. an auto generated bucket was relied
   upon) then this auto generated bucket will also be deleted as part of the cluster destroy.
   By default this is set to false.

* `image_version` - (Optional) The Cloud Dataproc image version to use
   for the cluster - this essentially controls the sets of software versions
   installed onto the nodes when you create clusters. If not specified, defaults to the
   latest version. For a list of valid versions see
   [Cloud Dataproc versions](https://cloud.google.com/dataproc/docs/concepts/dataproc-versions)

* `network` - (Optional) The name or self_link of the Google Compute Engine
    network to the cluster will be part of. Conflicts with `subnetwork`.
    If neither is specified, this defaults to the "default" network.

* `subnetwork` - (Optional) The name or self_link of the Google Compute Engine
   subnetwork the cluster will be part of. Conflicts with `network`.

* `initialization_actions` - (Optional) A list of scripts to be executed during
   initialisation of the cluster. Each must be a GCS file with a gs:// prefix.

* `initialization_action_timeout_sec` - (Optional) If `initialization_actions` is set,
   then this value specifies the maximum duration (in seconds) of each initialization
   action.

* `service_account` - (Optional) The service account to be used by the Node VMs.
    If not specified, the "default" service account is used.

* `service_scopes` - (Optional) The set of Google API scopes to be made available
    on all of the node VMs under the `service_account` specified. These can be
    either FQDNs, or scope aliases. The following scopes are necessary to ensure
    the correct functioning of the cluster:

  * `useraccounts-ro` (`https://www.googleapis.com/auth/cloud.useraccounts.readonly`)
  * `storage-rw`      (`https://www.googleapis.com/auth/devstorage.read_write`)
  * `logging-write`   (`https://www.googleapis.com/auth/logging.write`)

* `metadata` - (Optional) The metadata key/value pairs assigned to instances in
    the cluster.

* `labels` - (Optional) The list of labels (key/value pairs) to be applied to
   instances in the cluster.

* `properties` - (Optional) A list of override and additional properties (key/value pairs)
   used to modify various aspects of the common configuration files used when creating
   a cluster. For a list of valid please see
  [Cluster properties](https://cloud.google.com/dataproc/docs/concepts/cluster-properties)

* `tags` - (Optional) The list of instance tags applied to instances in the cluster.
   Tags are used to identify valid sources or targets for network firewalls.

The **master_config** supports:

* `num_masters`- (Optional) Specifies the number of master nodes to create.
   Defaults to 1.

* `machine_type` - (Optional) The name of a Google Compute Engine machine type.
    Defaults to `n1-standard-4`.

* `boot_disk_size_gb` - (Optional) Size of the primary disk attached to each node, specified
    in GB. The primary disk contains the boot volume and system libraries, and the
    smallest allowed disk size is 10GB, but defaults to 500GB. Note: If SSDs are not
    attached, it also contains the HDFS data blocks and Hadoop working directories.

* `num_local_ssds` - (Optional) The amount of local SSD disks that will be
    attached to each cluster node. Defaults to 0.

The **worker_config** supports:

* `num_worker`- (Optional) Specifies the number of worker nodes to create.
   Defaults to 2 which is the minimum. There is currently a beta feature which allows you to run a
   [Single Node Cluster](https://cloud.google.com/dataproc/docs/concepts/single-node-clusters).
   In order to take advantage of this you need to set the property `"dataproc:dataproc.allow.zero.workers" = "true"`

* `machine_type` - (Optional) The name of a Google Compute Engine machine type.
    Defaults to `n1-standard-4`.

* `boot_disk_size_gb` - (Optional) Size of the disk attached to each node, specified
    in GB. The smallest allowed disk size is 10GB. Defaults to 500GB.

* `num_local_ssds` - (Optional) The amount of local SSD disks that will be
    attached to each cluster node. Defaults to 0.

* `preemptible_num_workers`- (Optional) Specifies the number of master nodes to create.
   Defaults to 0. Note, the Machine type used is whatever is specified for the
   `worker_config.machine_type`

* `preemptible_boot_disk_size_gb`- (Optional) (Optional) Size of the disk attached to each
   preemptible node, specified in GB. The smallest allowed disk size is 10GB.
   Defaults to 500GB.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `master_instance_names` - List of master instance names which have been assigned
    to the cluster.

* `worker_instance_names` - List of worker instance names which have been assigned
    to the cluster.

* `preemptible_instance_names` - List of preemptible instance names which have been assigned
    to the cluster.

* `bucket` - The name of the cloud storage bucket ultimately used to house the staging data
   for the cluster. If `staging_bucket` is specified, it will contain this value, otherwise
   it will be the auto generated name.

<a id="timeouts"></a>
## Timeouts

`google_dataproc_cluster` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating clusters.
- `update` - (Default `5 minutes`) Used for updates to clusters
- `delete` - (Default `5 minutes`) Used for destroying clusters.
