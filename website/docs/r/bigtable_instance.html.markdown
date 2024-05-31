---
subcategory: "Cloud Bigtable"
description: |-
  Creates a Google Bigtable instance.
---

# google_bigtable_instance

Creates a Google Bigtable instance. For more information see:

* [API documentation](https://cloud.google.com/bigtable/docs/reference/admin/rest/v2/projects.instances.clusters)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/bigtable/docs)


-> **Note**: It is strongly recommended to set `lifecycle { prevent_destroy = true }`
on instances in order to prevent accidental data loss. See
[Terraform docs](https://www.terraform.io/docs/configuration/resources.html#prevent_destroy)
for more information on lifecycle parameters.

-> **Note**: On newer versions of the provider, you must explicitly set `deletion_protection=false`
(and run `terraform apply` to write the field to state) in order to destroy an instance.
It is recommended to not set this field (or set it to true) until you're ready to destroy.


## Example Usage - Simple Instance

```hcl
resource "google_bigtable_instance" "production-instance" {
  name = "tf-instance"

  cluster {
    cluster_id   = "tf-instance-cluster"
    num_nodes    = 1
    storage_type = "HDD"
  }

  labels = {
    my-label = "prod-label"
  }
}
```

## Example Usage - Replicated Instance

```hcl
resource "google_bigtable_instance" "production-instance" {
  name = "tf-instance"

  # A cluster with fixed number of nodes.
  cluster {
    cluster_id   = "tf-instance-cluster1"
    num_nodes    = 1
    storage_type = "HDD"
    zone    = "us-central1-c"
  }

  # a cluster with auto scaling.
  cluster {
    cluster_id   = "tf-instance-cluster2"
    storage_type = "HDD"
    zone    = "us-central1-b"
    autoscaling_config {
      min_nodes = 1
      max_nodes = 3
      cpu_target = 50
    }
  }

  labels = {
    my-label = "prod-label"
  }
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name (also called Instance Id in the Cloud Console) of the Cloud Bigtable instance. Must be 6-33 characters and must only contain hyphens, lowercase letters and numbers.

* `cluster` - (Required) A block of cluster configuration options. This can be specified at least once, and up 
to as many as possible within 8 cloud regions. Removing the field entirely from the config will cause the provider
to default to the backend value. See [structure below](#nested_cluster).

-----

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `instance_type` - (Optional, Deprecated) The instance type to create. One of `"DEVELOPMENT"` or `"PRODUCTION"`. Defaults to `"PRODUCTION"`.
    It is recommended to leave this field unspecified since the distinction between `"DEVELOPMENT"` and `"PRODUCTION"` instances is going away,
    and all instances will become `"PRODUCTION"` instances. This means that new and existing `"DEVELOPMENT"` instances will be converted to
    `"PRODUCTION"` instances. It is recommended for users to use `"PRODUCTION"` instances in any case, since a 1-node `"PRODUCTION"` instance
    is functionally identical to a `"DEVELOPMENT"` instance, but without the accompanying restrictions.

* `display_name` - (Optional) The human-readable display name of the Bigtable instance. Defaults to the instance `name`.

* `force_destroy` - (Optional) Deleting a BigTable instance can be blocked if any backups are present in the instance. When `force_destroy` is set to true, Terraform will delete all backups found in the BigTable instance before attempting to delete the instance itself. Defaults to false.

* `deletion_protection` - (Optional) Whether or not to allow Terraform to destroy the instance. Unless this field is set to false
in Terraform state, a `terraform destroy` or `terraform apply` that would delete the instance will fail. Defaults to true.

* `labels` - (Optional) A set of key/value label pairs to assign to the resource. Label keys must follow the requirements at https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements.

  **Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
  Please refer to the field 'effective_labels' for all of the labels present on the resource.

* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.

* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.

-----

<a name="nested_cluster"></a>The `cluster` block supports the following arguments:

* `cluster_id` - (Required) The ID of the Cloud Bigtable cluster. Must be 6-30 characters and must only contain hyphens, lowercase letters and numbers.

* `zone` - (Optional) The zone to create the Cloud Bigtable cluster in. If it not
specified, the provider zone is used. Each cluster must have a different zone in the same region. Zones that support
Bigtable instances are noted on the [Cloud Bigtable locations page](https://cloud.google.com/bigtable/docs/locations).

* `num_nodes` - (Optional) The number of nodes in the cluster.
If no value is set, Cloud Bigtable automatically allocates nodes based on your data footprint and optimized for 50% storage utilization.

* `autoscaling_config` - (Optional) [Autoscaling](https://cloud.google.com/bigtable/docs/autoscaling#parameters) config for the cluster, contains the following arguments:

  * `min_nodes` - (Required) The minimum number of nodes for autoscaling.
  * `max_nodes` - (Required) The maximum number of nodes for autoscaling.
  * `cpu_target` - (Required) The target CPU utilization for autoscaling, in percentage. Must be between 10 and 80.
  * `storage_target` - The target storage utilization for autoscaling, in GB, for each node in a cluster. This number is limited between 2560 (2.5TiB) and 5120 (5TiB) for a SSD cluster and between 8192 (8TiB) and 16384 (16 TiB) for an HDD cluster. If not set, whatever is already set for the cluster will not change, or if the cluster is just being created, it will use the default value of 2560 for SSD clusters and 8192 for HDD clusters.

!> **Warning**: Only one of `autoscaling_config` or `num_nodes` should be set for a cluster. If both are set, `num_nodes` is ignored. If none is set, autoscaling will be disabled and sized to the current node count.

* `storage_type` - (Optional) The storage type to use. One of `"SSD"` or
`"HDD"`. Defaults to `"SSD"`.

* `kms_key_name` - (Optional) Describes the Cloud KMS encryption key that will be used to protect the destination Bigtable cluster. The requirements for this key are: 1) The Cloud Bigtable service account associated with the project that contains this cluster must be granted the `cloudkms.cryptoKeyEncrypterDecrypter` role on the CMEK key. 2) Only regional keys can be used and the region of the CMEK key must match the region of the cluster.

-> **Note**: Removing the field entirely from the config will cause the provider to default to the backend value.

!> **Warning:** Modifying the `storage_type`, `zone` or `kms_key_name` of an existing cluster (by
`cluster_id`) will cause Terraform to delete/recreate the entire
`google_bigtable_instance` resource. If these values are changing, use a new
`cluster_id`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/instances/{{name}}`
* `cluster.0.state` - describes the current state of the cluster.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 60 minutes.
- `update` - Default is 60 minutes.
- `read` - Default is 60 minutes.

Adding clusters to existing instances can take a long time. Consider setting a higher value to timeouts if you plan on doing that.

## Import

Bigtable Instances can be imported using any of these accepted formats:

* `projects/{{project}}/instances/{{name}}`
* `{{project}}/{{name}}`
* `{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to Bigtable Instances using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/instances/{{name}}"
  to = google_bigtable_instance.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Bigtable Instances can be imported using one of the formats above. For example:

```
$ terraform import google_bigtable_instance.default projects/{{project}}/instances/{{name}}
$ terraform import google_bigtable_instance.default {{project}}/{{name}}
$ terraform import google_bigtable_instance.default {{name}}
```
