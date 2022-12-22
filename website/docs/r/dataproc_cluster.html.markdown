---
subcategory: "Dataproc"
page_title: "Google: google_dataproc_cluster"
description: |-
  Manages a Cloud Dataproc cluster resource.
---

# google\_dataproc\_cluster

Manages a Cloud Dataproc cluster resource within GCP.

* [API documentation](https://cloud.google.com/dataproc/docs/reference/rest/v1/projects.regions.clusters)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/dataproc/docs)


!> **Warning:** Due to limitations of the API, all arguments except
`labels`,`cluster_config.worker_config.num_instances` and `cluster_config.preemptible_worker_config.num_instances` are non-updatable. Changing others will cause recreation of the
whole cluster!

## Example Usage - Basic

```hcl
resource "google_dataproc_cluster" "simplecluster" {
  name   = "simplecluster"
  region = "us-central1"
}
```

## Example Usage - Advanced

```hcl
resource "google_service_account" "default" {
  account_id   = "service-account-id"
  display_name = "Service Account"
}

resource "google_dataproc_cluster" "mycluster" {
  name     = "mycluster"
  region   = "us-central1"
  graceful_decommission_timeout = "120s"
  labels = {
    foo = "bar"
  }

  cluster_config {
    staging_bucket = "dataproc-staging-bucket"

    master_config {
      num_instances = 1
      machine_type  = "e2-medium"
      disk_config {
        boot_disk_type    = "pd-ssd"
        boot_disk_size_gb = 30
      }
    }

    worker_config {
      num_instances    = 2
      machine_type     = "e2-medium"
      min_cpu_platform = "Intel Skylake"
      disk_config {
        boot_disk_size_gb = 30
        num_local_ssds    = 1
      }
    }

    preemptible_worker_config {
      num_instances = 0
    }

    # Override or set some custom properties
    software_config {
      image_version = "2.0.35-debian10"
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    gce_cluster_config {
      tags = ["foo", "bar"]
      # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
      service_account = google_service_account.default.email
      service_account_scopes = [
        "cloud-platform"
      ]
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

* `virtual_cluster_config` - (Optional) Allows you to configure a virtual Dataproc on GKE cluster.
   Structure [defined below](#nested_virtual_cluster_config).

* `cluster_config` - (Optional) Allows you to configure various aspects of the cluster.
   Structure [defined below](#nested_cluster_config).

* `graceful_decommission_timeout` - (Optional) Allows graceful decomissioning when you change the number of worker nodes directly through a terraform apply.
      Does not affect auto scaling decomissioning from an autoscaling policy.
      Graceful decommissioning allows removing nodes from the cluster without interrupting jobs in progress.
      Timeout specifies how long to wait for jobs in progress to finish before forcefully removing nodes (and potentially interrupting jobs).
      Default timeout is 0 (for forceful decommission), and the maximum allowed timeout is 1 day. (see JSON representation of
      [Duration](https://developers.google.com/protocol-buffers/docs/proto3#json)).
      Only supported on Dataproc image versions 1.2 and higher.
      For more context see the [docs](https://cloud.google.com/dataproc/docs/reference/rest/v1/projects.regions.clusters/patch#query-parameters)
- - -

<a name="nested_virtual_cluster_config"></a>The `virtual_cluster_config` block supports:

```hcl
    virtual_cluster_config {
        auxiliary_services_config { ... }
        kubernetes_cluster_config { ... }
    }
```

* `staging_bucket` - (Optional) The Cloud Storage staging bucket used to stage files,
   such as Hadoop jars, between client machines and the cluster.
   Note: If you don't explicitly specify a `staging_bucket`
   then GCP will auto create / assign one for you. However, you are not guaranteed
   an auto generated bucket which is solely dedicated to your cluster; it may be shared
   with other clusters in the same region/zone also choosing to use the auto generation
   option.

* `auxiliary_services_config` (Optional) Configuration of auxiliary services used by this cluster. 
   Structure [defined below](#nested_auxiliary_services_config).

* `kubernetes_cluster_config` (Required) The configuration for running the Dataproc cluster on Kubernetes.
   Structure [defined below](#nested_kubernetes_cluster_config).
- - -

<a name="nested_auxiliary_services_config"></a>The `auxiliary_services_config` block supports:

```hcl
    virtual_cluster_config {
      auxiliary_services_config {
        metastore_config {
          dataproc_metastore_service = google_dataproc_metastore_service.metastore_service.id
        }

        spark_history_server_config {
          dataproc_cluster = google_dataproc_cluster.dataproc_cluster.id
        }
      }
    }
```

* `metastore_config` (Optional) The Hive Metastore configuration for this workload. 

  * `dataproc_metastore_service` (Required) Resource name of an existing Dataproc Metastore service.

* `spark_history_server_config` (Optional) The Spark History Server configuration for the workload.

  * `dataproc_cluster` (Optional) Resource name of an existing Dataproc Cluster to act as a Spark History Server for the workload.
- - -

<a name="nested_kubernetes_cluster_config"></a>The `kubernetes_cluster_config` block supports:

```hcl
    virtual_cluster_config {
      kubernetes_cluster_config {
        kubernetes_namespace = "foobar"

        kubernetes_software_config {
          component_version = {
            "SPARK" : "3.1-dataproc-7"
          }

          properties = {
            "spark:spark.eventLog.enabled": "true"
          }
        }

        gke_cluster_config {
          gke_cluster_target = google_container_cluster.primary.id

          node_pool_target {
            node_pool = "dpgke"
            roles = ["DEFAULT"]

            node_pool_config {
              autoscaling {
                min_node_count = 1
                max_node_count = 6
              }
              
              config {
                machine_type      = "n1-standard-4"
                preemptible       = true
                local_ssd_count   = 1
                min_cpu_platform  = "Intel Sandy Bridge"
              }

              locations = ["us-central1-c"]
            }
          }
        }
      }
    }
```

* `kubernetes_namespace` (Optional) A namespace within the Kubernetes cluster to deploy into. 
   If this namespace does not exist, it is created. 
   If it  exists, Dataproc verifies that another Dataproc VirtualCluster is not installed into it. 
   If not specified, the name of the Dataproc Cluster is used.

* `kubernetes_software_config` (Required) The software configuration for this Dataproc cluster running on Kubernetes.

  * `component_version` (Required) The components that should be installed in this Dataproc cluster. The key must be a string from the   
     KubernetesComponent enumeration. The value is the version of the software to be installed. At least one entry must be specified.
    * **NOTE** : `component_version[SPARK]` is mandatory to set, or the creation of the cluster will fail.

  * `properties` (Optional) The properties to set on daemon config files. Property keys are specified in prefix:property format, 
     for example spark:spark.kubernetes.container.image.

* `gke_cluster_config` (Required) The configuration for running the Dataproc cluster on GKE.

  * `gke_cluster_target` (Optional) A target GKE cluster to deploy to. It must be in the same project and region as the Dataproc cluster 
     (the GKE cluster can be zonal or regional)

  * `node_pool_target` (Optional) GKE node pools where workloads will be scheduled. At least one node pool must be assigned the `DEFAULT` 
     GkeNodePoolTarget.Role. If a GkeNodePoolTarget is not specified, Dataproc constructs a `DEFAULT` GkeNodePoolTarget. 
     Each role can be given to only one GkeNodePoolTarget. All node pools must have the same location settings.

    * `node_pool` (Required) The target GKE node pool.

    * `roles` (Required) The roles associated with the GKE node pool. 
       One of `"DEFAULT"`, `"CONTROLLER"`, `"SPARK_DRIVER"` or `"SPARK_EXECUTOR"`.

    * `node_pool_config` (Input only) The configuration for the GKE node pool. 
       If specified, Dataproc attempts to create a node pool with the specified shape. 
       If one with the same name already exists, it is verified against all specified fields. 
       If a field differs, the virtual cluster creation will fail.

      * `autoscaling` (Optional) The autoscaler configuration for this node pool. 
         The autoscaler is enabled only when a valid configuration is present.

        * `min_node_count` (Optional) The minimum number of nodes in the node pool. Must be >= 0 and <= maxNodeCount.

        * `max_node_count` (Optional) The maximum number of nodes in the node pool. Must be >= minNodeCount, and must be > 0.

      * `config` (Optional) The node pool configuration.

        * `machine_type` (Optional) The name of a Compute Engine machine type.

        * `local_ssd_count` (Optional) The number of local SSD disks to attach to the node, 
           which is limited by the maximum number of disks allowable per zone.

        * `preemptible` (Optional) Whether the nodes are created as preemptible VM instances. 
           Preemptible nodes cannot be used in a node pool with the CONTROLLER role or in the DEFAULT node pool if the 
           CONTROLLER role is not assigned (the DEFAULT node pool will assume the CONTROLLER role).

        * `min_cpu_platform` (Optional) Minimum CPU platform to be used by this instance. 
           The instance may be scheduled on the specified or a newer CPU platform. 
           Specify the friendly names of CPU platforms, such as "Intel Haswell" or "Intel Sandy Bridge".

        * `spot` (Optional) Spot flag for enabling Spot VM, which is a rebrand of the existing preemptible flag.

      * `locations` (Optional) The list of Compute Engine zones where node pool nodes associated 
         with a Dataproc on GKE virtual cluster will be located.
- - -

<a name="nested_cluster_config"></a>The `cluster_config` block supports:

```hcl
    cluster_config {
        gce_cluster_config        { ... }
        master_config             { ... }
        worker_config             { ... }
        preemptible_worker_config { ... }
        software_config           { ... }

        # You can define multiple initialization_action blocks
        initialization_action     { ... }
        encryption_config         { ... }
        endpoint_config           { ... }
        metastore_config          { ... }
    }
```

* `staging_bucket` - (Optional) The Cloud Storage staging bucket used to stage files,
   such as Hadoop jars, between client machines and the cluster.
   Note: If you don't explicitly specify a `staging_bucket`
   then GCP will auto create / assign one for you. However, you are not guaranteed
   an auto generated bucket which is solely dedicated to your cluster; it may be shared
   with other clusters in the same region/zone also choosing to use the auto generation
   option.

* `temp_bucket` - (Optional) The Cloud Storage temp bucket used to store ephemeral cluster
   and jobs data, such as Spark and MapReduce history files.
   Note: If you don't explicitly specify a `temp_bucket` then GCP will auto create / assign one for you.

* `gce_cluster_config` (Optional) Common config settings for resources of Google Compute Engine cluster
   instances, applicable to all instances in the cluster. Structure [defined below](#nested_gce_cluster_config).

* `master_config` (Optional) The Google Compute Engine config settings for the master instances
   in a cluster. Structure [defined below](#nested_master_config).

* `worker_config` (Optional) The Google Compute Engine config settings for the worker instances
   in a cluster. Structure [defined below](#nested_worker_config).

* `preemptible_worker_config` (Optional) The Google Compute Engine config settings for the additional
   instances in a cluster. Structure [defined below](#nested_preemptible_worker_config).
  * **NOTE** : `preemptible_worker_config` is
   an alias for the api's [secondaryWorkerConfig](https://cloud.google.com/dataproc/docs/reference/rest/v1/ClusterConfig#InstanceGroupConfig). The name doesn't necessarily mean it is preemptible and is named as
   such for legacy/compatibility reasons.

* `software_config` (Optional) The config settings for software inside the cluster.
   Structure [defined below](#nested_software_config).

* `security_config` (Optional) Security related configuration. Structure [defined below](#nested_security_config).

* `autoscaling_config` (Optional)  The autoscaling policy config associated with the cluster.
   Note that once set, if `autoscaling_config` is the only field set in `cluster_config`, it can
   only be removed by setting `policy_uri = ""`, rather than removing the whole block.
   Structure [defined below](#nested_autoscaling_config).

* `initialization_action` (Optional) Commands to execute on each node after config is completed.
   You can specify multiple versions of these. Structure [defined below](#nested_initialization_action).

* `encryption_config` (Optional) The Customer managed encryption keys settings for the cluster.
   Structure [defined below](#nested_encryption_config).

* `lifecycle_config` (Optional) The settings for auto deletion cluster schedule.
   Structure [defined below](#nested_lifecycle_config).

* `endpoint_config` (Optional) The config settings for port access on the cluster.
   Structure [defined below](#nested_endpoint_config).

* `metastore_config` (Optional) The config setting for metastore service with the cluster.
   Structure [defined below](#nested_metastore_config).
- - -

<a name="nested_gce_cluster_config"></a>The `cluster_config.gce_cluster_config` block supports:

```hcl
  cluster_config {
    gce_cluster_config {
      zone = "us-central1-a"

      # One of the below to hook into a custom network / subnetwork
      network    = google_compute_network.dataproc_network.name
      subnetwork = google_compute_network.dataproc_subnetwork.name

      tags = ["foo", "bar"]
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

* `service_account_scopes` - (Optional, Computed) The set of Google API scopes
    to be made available on all of the node VMs under the `service_account`
    specified. Both OAuth2 URLs and gcloud
    short names are supported. To allow full access to all Cloud APIs, use the
    `cloud-platform` scope. See a complete list of scopes [here](https://cloud.google.com/sdk/gcloud/reference/alpha/compute/instances/set-scopes#--scopes).

* `tags` - (Optional) The list of instance tags applied to instances in the cluster.
   Tags are used to identify valid sources or targets for network firewalls.

* `internal_ip_only` - (Optional) By default, clusters are not restricted to internal IP addresses,
   and will have ephemeral external IP addresses assigned to each instance. If set to true, all
   instances in the cluster will only have internal IP addresses. Note: Private Google Access
   (also known as `privateIpGoogleAccess`) must be enabled on the subnetwork that the cluster
   will be launched in.

* `metadata` - (Optional) A map of the Compute Engine metadata entries to add to all instances
   (see [Project and instance metadata](https://cloud.google.com/compute/docs/storing-retrieving-metadata#project_and_instance_metadata)).

* `shielded_instance_config` (Optional) Shielded Instance Config for clusters using [Compute Engine Shielded VMs](https://cloud.google.com/security/shielded-cloud/shielded-vm).

- - -


The `cluster_config.gce_cluster_config.shielded_instance_config` block supports:

```hcl
cluster_config{
  gce_cluster_config{
    shielded_instance_config{
      enable_secure_boot          = true
      enable_vtpm                 = true
      enable_integrity_monitoring = true
    }
  }
}
```

* `enable_secure_boot` - (Optional) Defines whether instances have Secure Boot enabled.

* `enable_vtpm` - (Optional) Defines whether instances have the [vTPM](https://cloud.google.com/security/shielded-cloud/shielded-vm#vtpm) enabled.

* `enable_integrity_monitoring` - (Optional) Defines whether instances have integrity monitoring enabled.

- - -

<a name="nested_master_config"></a>The `cluster_config.master_config` block supports:

```hcl
cluster_config {
  master_config {
    num_instances    = 1
    machine_type     = "e2-medium"
    min_cpu_platform = "Intel Skylake"

    disk_config {
      boot_disk_type    = "pd-ssd"
      boot_disk_size_gb = 30
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

* `min_cpu_platform` - (Optional, Computed) The name of a minimum generation of CPU family
   for the master. If not specified, GCP will default to a predetermined computed value
   for each zone. See [the guide](https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform)
   for details about which CPU families are available (and defaulted) for each zone.

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

<a name="nested_worker_config"></a>The `cluster_config.worker_config` block supports:

```hcl
cluster_config {
  worker_config {
    num_instances    = 3
    machine_type     = "e2-medium"
    min_cpu_platform = "Intel Skylake"

    disk_config {
      boot_disk_type    = "pd-standard"
      boot_disk_size_gb = 30
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

* `min_cpu_platform` - (Optional, Computed) The name of a minimum generation of CPU family
   for the master. If not specified, GCP will default to a predetermined computed value
   for each zone. See [the guide](https://cloud.google.com/compute/docs/instances/specify-min-cpu-platform)
   for details about which CPU families are available (and defaulted) for each zone.

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

<a name="nested_preemptible_worker_config"></a>The `cluster_config.preemptible_worker_config` block supports:

```hcl
cluster_config {
  preemptible_worker_config {
    num_instances = 1

    disk_config {
      boot_disk_type    = "pd-standard"
      boot_disk_size_gb = 30
      num_local_ssds    = 1
    }
  }
}
```

Note: Unlike `worker_config`, you cannot set the `machine_type` value directly. This
will be set for you based on whatever was set for the `worker_config.machine_type` value.

* `num_instances`- (Optional) Specifies the number of preemptible nodes to create.
   Defaults to 0.

* `preemptibility`- (Optional) Specifies the preemptibility of the secondary workers. The default value is `PREEMPTIBLE`
  Accepted values are:
  * PREEMPTIBILITY_UNSPECIFIED
  * NON_PREEMPTIBLE
  * PREEMPTIBLE
  * SPOT

* `disk_config` (Optional) Disk Config

    * `boot_disk_type` - (Optional) The disk type of the primary disk attached to each preemptible worker node.
	One of `"pd-ssd"` or `"pd-standard"`. Defaults to `"pd-standard"`.

    * `boot_disk_size_gb` - (Optional, Computed) Size of the primary disk attached to each preemptible worker node, specified
    in GB. The smallest allowed disk size is 10GB. GCP will default to a predetermined
    computed value if not set (currently 500GB). Note: If SSDs are not
	attached, it also contains the HDFS data blocks and Hadoop working directories.

	* `num_local_ssds` - (Optional) The amount of local SSD disks that will be
	attached to each preemptible worker node. Defaults to 0.

- - -

<a name="nested_software_config"></a>The `cluster_config.software_config` block supports:

```hcl
cluster_config {
  # Override or set some custom properties
  software_config {
    image_version = "2.0.35-debian10"

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

* `optional_components` - (Optional) The set of optional components to activate on the cluster.
    Accepted values are:
    * ANACONDA
    * DRUID
    * FLINK
    * HBASE
    * HIVE_WEBHCAT
    * JUPYTER
    * PRESTO
    * RANGER
    * SOLR
    * ZEPPELIN
    * ZOOKEEPER

- - -

<a name="nested_security_config"></a>The `cluster_config.security_config` block supports:

```hcl
cluster_config {
  # Override or set some custom properties
  security_config {
    kerberos_config {
      kms_key_uri = "projects/projectId/locations/locationId/keyRings/keyRingId/cryptoKeys/keyId"
      root_principal_password_uri = "bucketId/o/objectId"
    }
  }
}
```

* `kerberos_config` (Required) Kerberos Configuration

    * `cross_realm_trust_admin_server` - (Optional) The admin server (IP or hostname) for the
       remote trusted realm in a cross realm trust relationship.

    * `cross_realm_trust_kdc` - (Optional) The KDC (IP or hostname) for the
       remote trusted realm in a cross realm trust relationship.

    * `cross_realm_trust_realm` - (Optional) The remote realm the Dataproc on-cluster KDC will
       trust, should the user enable cross realm trust.

    * `cross_realm_trust_shared_password_uri` - (Optional) The Cloud Storage URI of a KMS
       encrypted file containing the shared password between the on-cluster Kerberos realm
       and the remote trusted realm, in a cross realm trust relationship.

    * `enable_kerberos` - (Optional) Flag to indicate whether to Kerberize the cluster.

    * `kdc_db_key_uri` - (Optional) The Cloud Storage URI of a KMS encrypted file containing
       the master key of the KDC database.

    * `key_password_uri` - (Optional) The Cloud Storage URI of a KMS encrypted file containing
       the password to the user provided key. For the self-signed certificate, this password
       is generated by Dataproc.

    * `keystore_uri` - (Optional) The Cloud Storage URI of the keystore file used for SSL encryption.
       If not provided, Dataproc will provide a self-signed certificate.

    * `keystore_password_uri` - (Optional) The Cloud Storage URI of a KMS encrypted file containing
       the password to the user provided keystore. For the self-signed certificated, the password
       is generated by Dataproc.

    * `kms_key_uri` - (Required) The URI of the KMS key used to encrypt various sensitive files.

    * `realm` - (Optional) The name of the on-cluster Kerberos realm. If not specified, the
       uppercased domain of hostnames will be the realm.

    * `root_principal_password_uri` - (Required) The Cloud Storage URI of a KMS encrypted file
       containing the root principal password.

    * `tgt_lifetime_hours` - (Optional) The lifetime of the ticket granting ticket, in hours.

    * `truststore_password_uri` - (Optional) The Cloud Storage URI of a KMS encrypted file
       containing the password to the user provided truststore. For the self-signed
       certificate, this password is generated by Dataproc.

    * `truststore_uri` - (Optional) The Cloud Storage URI of the truststore file used for
       SSL encryption. If not provided, Dataproc will provide a self-signed certificate.

- - -

<a name="nested_autoscaling_config"></a>The `cluster_config.autoscaling_config` block supports:

```hcl
cluster_config {
  # Override or set some custom properties
  autoscaling_config {
    policy_uri = "projects/projectId/locations/region/autoscalingPolicies/policyId"
  }
}
```

* `policy_uri` - (Required) The autoscaling policy used by the cluster.

Only resource names including projectid and location (region) are valid. Examples:

`https://www.googleapis.com/compute/v1/projects/[projectId]/locations/[dataproc_region]/autoscalingPolicies/[policy_id]`
`projects/[projectId]/locations/[dataproc_region]/autoscalingPolicies/[policy_id]`
Note that the policy must be in the same project and Cloud Dataproc region.

- - -

<a name="nested_initialization_action"></a>The `initialization_action` block (Optional) can be specified multiple times and supports:

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

- - -

<a name="nested_encryption_config"></a>The `encryption_config` block supports:

```hcl
cluster_config {
  encryption_config {
    kms_key_name = "projects/projectId/locations/region/keyRings/keyRingName/cryptoKeys/keyName"
  }
}
```

* `kms_key_name` - (Required) The Cloud KMS key name to use for PD disk encryption for
   all instances in the cluster.

- - -

<a name="nested_lifecycle_config"></a>The `lifecycle_config` block supports:

```hcl
cluster_config {
  lifecycle_config {
    idle_delete_ttl = "10m"
    auto_delete_time = "2120-01-01T12:00:00.01Z"
  }
}
```

* `idle_delete_ttl` - (Optional) The duration to keep the cluster alive while idling
  (no jobs running). After this TTL, the cluster will be deleted. Valid range: [10m, 14d].

* `auto_delete_time` - (Optional) The time when cluster will be auto-deleted.
  A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds.
  Example: "2014-10-02T15:01:23.045123456Z".

- - -

<a name="nested_endpoint_config"></a>The `endpoint_config` block (Optional, Computed, Beta) supports:

```hcl
cluster_config {
  endpoint_config {
    enable_http_port_access = "true"
  }
}
```

* `enable_http_port_access` - (Optional) The flag to enable http access to specific ports
  on the cluster from external sources (aka Component Gateway). Defaults to false.


<a name="nested_metastore_config"></a>The `metastore_config` block (Optional, Computed, Beta) supports:

```hcl
cluster_config {
  metastore_config {
    dataproc_metastore_service = "projects/projectId/locations/region/services/serviceName"
  }
}
```

* `dataproc_metastore_service` - (Required) Resource name of an existing Dataproc Metastore service.

Only resource names including projectid and location (region) are valid. Examples:

`projects/[projectId]/locations/[dataproc_region]/services/[service-name]`

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `cluster_config.0.master_config.0.instance_names` - List of master instance names which
   have been assigned to the cluster.

* `cluster_config.0.worker_config.0.instance_names` - List of worker instance names which have been assigned
	to the cluster.

* `cluster_config.0.preemptible_worker_config.0.instance_names` - List of preemptible instance names which have been assigned
	to the cluster.

* `cluster_config.0.bucket` - The name of the cloud storage bucket ultimately used to house the staging data
   for the cluster. If `staging_bucket` is specified, it will contain this value, otherwise
   it will be the auto generated name.

* `cluster_config.0.software_config.0.properties` - A list of the properties used to set the daemon config files.
   This will include any values supplied by the user via `cluster_config.software_config.override_properties`

* `cluster_config.0.lifecycle_config.0.idle_start_time` - Time when the cluster became idle
  (most recent job finished) and became eligible for deletion due to idleness.

* `cluster_config.0.endpoint_config.0.http_ports` - The map of port descriptions to URLs. Will only be populated if
  `enable_http_port_access` is true.

## Import

This resource does not support import.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 45 minutes.
- `update` - Default is 45 minutes.
- `delete` - Default is 45 minutes.
