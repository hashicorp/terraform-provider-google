---
subcategory: "Cloud Bigtable"
layout: "google"
page_title: "Google: google_bigtable_instance"
sidebar_current: "docs-google-bigtable-instance"
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


## Example Usage - Production Instance

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


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name (also called Instance Id in the Cloud Console) of the Cloud Bigtable instance.

* `cluster` - (Required) A block of cluster configuration options. This can be specified at least once, and up to 4 times.
See [structure below](#nested_cluster).

-----

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `instance_type` - (Optional, Deprecated) The instance type to create. One of `"DEVELOPMENT"` or `"PRODUCTION"`. Defaults to `"PRODUCTION"`.
    It is recommended to leave this field unspecified since the distinction between `"DEVELOPMENT"` and `"PRODUCTION"` instances is going away,
    and all instances will become `"PRODUCTION"` instances. This means that new and existing `"DEVELOPMENT"` instances will be converted to
    `"PRODUCTION"` instances. It is recommended for users to use `"PRODUCTION"` instances in any case, since a 1-node `"PRODUCTION"` instance
    is functionally identical to a `"DEVELOPMENT"` instance, but without the accompanying restrictions.

* `display_name` - (Optional) The human-readable display name of the Bigtable instance. Defaults to the instance `name`.

* `deletion_protection` - (Optional) Whether or not to allow Terraform to destroy the instance. Unless this field is set to false
in Terraform state, a `terraform destroy` or `terraform apply` that would delete the instance will fail.

* `labels` - (Optional) A set of key/value label pairs to assign to the resource. Label keys must follow the requirements at https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements.


-----

<a name="nested_cluster"></a>The `cluster` block supports the following arguments:

* `cluster_id` - (Required) The ID of the Cloud Bigtable cluster.

* `zone` - (Optional) The zone to create the Cloud Bigtable cluster in. If it not
specified, the provider zone is used. Each cluster must have a different zone in the same region. Zones that support
Bigtable instances are noted on the [Cloud Bigtable locations page](https://cloud.google.com/bigtable/docs/locations).

* `num_nodes` - (Optional) The number of nodes in your Cloud Bigtable cluster.
Required, with a minimum of `1` for a `PRODUCTION` instance. Must be left unset
for a `DEVELOPMENT` instance.

* `storage_type` - (Optional) The storage type to use. One of `"SSD"` or
`"HDD"`. Defaults to `"SSD"`.

* `kms_key_name` - (Optional) Describes the Cloud KMS encryption key that will be used to protect the destination Bigtable cluster. The requirements for this key are: 1) The Cloud Bigtable service account associated with the project that contains this cluster must be granted the `cloudkms.cryptoKeyEncrypterDecrypter` role on the CMEK key. 2) Only regional keys can be used and the region of the CMEK key must match the region of the cluster. 3) All clusters within an instance must use the same CMEK key. Values are of the form `projects/{project}/locations/{location}/keyRings/{keyring}/cryptoKeys/{key}`

!> **Warning**: Modifying this field will cause Terraform to delete/recreate the entire resource. 

-> **Note**: To remove this field once it is set, set the value to an empty string. Removing the field entirely from the config will cause the provider to default to the backend value.

!> **Warning:** Modifying the `storage_type` or `zone` of an existing cluster (by
`cluster_id`) will cause Terraform to delete/recreate the entire
`google_bigtable_instance` resource. If these values are changing, use a new
`cluster_id`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/instances/{{name}}`

## Import

Bigtable Instances can be imported using any of these accepted formats:

```
$ terraform import google_bigtable_instance.default projects/{{project}}/instances/{{name}}
$ terraform import google_bigtable_instance.default {{project}}/{{name}}
$ terraform import google_bigtable_instance.default {{name}}
```
