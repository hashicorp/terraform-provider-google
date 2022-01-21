---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
#
# ----------------------------------------------------------------------------
#
#     This file is managed by Magic Modules (https:#github.com/GoogleCloudPlatform/magic-modules)
#     and is based on the DCL (https:#github.com/GoogleCloudPlatform/declarative-resource-client-library).
#     Changes will need to be made to the DCL or Magic Modules instead of here.
#
#     We are not currently able to accept contributions to this file. If changes
#     are required, please file an issue at https:#github.com/hashicorp/terraform-provider-google/issues/new/choose
#
# ----------------------------------------------------------------------------
subcategory: "ContainerAws"
layout: "google"
page_title: "Google: google_container_aws_node_pool"
description: |-
An Anthos node pool running on AWS.
---

# google_container_aws_node_pool

An Anthos node pool running on AWS.

For more information, see:
* [Multicloud overview](https://cloud.google.com/anthos/clusters/docs/multi-cloud)
## Example Usage - basic_aws_cluster
A basic example of a containeraws node pool
```hcl
data "google_container_aws_versions" "versions" {
  project = "my-project-name"
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  authorization {
    admin_users {
      username = "emailAddress:my@service-account.com"
    }
  }

  aws_region = "my-aws-region"

  control_plane {
    aws_services_authentication {
      role_arn          = "arn:aws:iam::012345678910:role/my--1p-dev-oneplatform"
      role_session_name = "my--1p-dev-session"
    }

    config_encryption {
      kms_key_arn = "arn:aws:kms:my-aws-region:012345678910:key/12345678-1234-1234-1234-123456789111"
    }

    database_encryption {
      kms_key_arn = "arn:aws:kms:my-aws-region:012345678910:key/12345678-1234-1234-1234-123456789111"
    }

    iam_instance_profile = "my--1p-dev-controlplane"
    subnet_ids           = ["subnet-00000000000000000"]
    version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"
    instance_type        = "t3.medium"

    main_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:my-aws-region:012345678910:key/12345678-1234-1234-1234-123456789111"
      size_gib    = 10
      volume_type = "GP3"
    }

    proxy_config {
      secret_arn     = "arn:aws:secretsmanager:us-west-2:126285863215:secret:proxy_config20210824150329476300000001-ABCDEF"
      secret_version = "12345678-ABCD-EFGH-IJKL-987654321098"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:my-aws-region:012345678910:key/12345678-1234-1234-1234-123456789111"
      size_gib    = 10
      volume_type = "GP3"
    }

    security_group_ids = ["sg-00000000000000000"]

    ssh_config {
      ec2_key_pair = "my--1p-dev-ssh"
    }

    tags = {
      owner = "emailAddress:my@service-account.com"
    }
  }

  fleet {
    project = "my-project-number"
  }

  location = "us-west1"
  name     = "name"

  networking {
    pod_address_cidr_blocks     = ["10.2.0.0/16"]
    service_address_cidr_blocks = ["10.1.0.0/16"]
    vpc_id                      = "vpc-00000000000000000"
  }

  annotations = {
    label-one = "value-one"
  }

  description = "A sample aws cluster"
  project     = "my-project-name"
}


resource "google_container_aws_node_pool" "primary" {
  autoscaling {
    max_node_count = 5
    min_node_count = 1
  }

  cluster = google_container_aws_cluster.primary.name

  config {
    config_encryption {
      kms_key_arn = "arn:aws:kms:my-aws-region:012345678910:key/12345678-1234-1234-1234-123456789111"
    }

    iam_instance_profile = "my--1p-dev-nodepool"
    instance_type        = "t3.medium"

    labels = {
      label-one = "value-one"
    }

    root_volume {
      iops        = 3000
      kms_key_arn = "arn:aws:kms:my-aws-region:012345678910:key/12345678-1234-1234-1234-123456789111"
      size_gib    = 10
      volume_type = "GP3"
    }

    security_group_ids = ["sg-00000000000000000"]

    ssh_config {
      ec2_key_pair = "my--1p-dev-ssh"
    }

    tags = {
      tag-one = "value-one"
    }

    taints {
      effect = "PREFER_NO_SCHEDULE"
      key    = "taint-key"
      value  = "taint-value"
    }
  }

  location = "us-west1"

  max_pods_constraint {
    max_pods_per_node = 110
  }

  name      = "node-pool-name"
  subnet_id = "subnet-00000000000000000"
  version   = "${data.google_container_aws_versions.versions.valid_versions[0]}"

  annotations = {
    label-one = "value-one"
  }

  project = "my-project-name"
}


```

## Argument Reference

The following arguments are supported:

* `autoscaling` -
  (Required)
  Required. Autoscaler configuration for this node pool.
  
* `cluster` -
  (Required)
  The awsCluster for the resource
  
* `config` -
  (Required)
  Required. The configuration of the node pool.
  
* `location` -
  (Required)
  The location for the resource
  
* `max_pods_constraint` -
  (Required)
  Required. The constraint on the maximum number of pods that can be run simultaneously on a node in the node pool.
  
* `name` -
  (Required)
  The name of this resource.
  
* `subnet_id` -
  (Required)
  Required. The subnet where the node pool node run.
  
* `version` -
  (Required)
  Required. The Kubernetes version to run on this node pool (e.g. `1.19.10-gke.1000`). You can list all supported versions on a given Google Cloud region by calling GetAwsServerConfig.
  


The `autoscaling` block supports:
    
* `max_node_count` -
  (Required)
  Required. Maximum number of nodes in the NodePool. Must be >= min_node_count.
    
* `min_node_count` -
  (Required)
  Required. Minimum number of nodes in the NodePool. Must be >= 1 and <= max_node_count.
    
The `config` block supports:
    
* `config_encryption` -
  (Required)
  Required. The ARN of the AWS KMS key used to encrypt node pool configuration.
    
* `iam_instance_profile` -
  (Required)
  Required. The name of the AWS IAM role assigned to nodes in the pool.
    
* `instance_type` -
  (Optional)
  Optional. The AWS instance type. When unspecified, it defaults to `t3.medium`.
    
* `labels` -
  (Optional)
  Optional. The initial labels assigned to nodes of this node pool. An object containing a list of "key": value pairs. Example: { "name": "wrench", "mass": "1.3kg", "count": "3" }.
    
* `root_volume` -
  (Optional)
  Optional. Template for the root volume provisioned for node pool nodes. Volumes will be provisioned in the availability zone assigned to the node pool subnet. When unspecified, it defaults to 32 GiB with the GP2 volume type.
    
* `security_group_ids` -
  (Optional)
  Optional. The IDs of additional security groups to add to nodes in this pool. The manager will automatically create security groups with minimum rules needed for a functioning cluster.
    
* `ssh_config` -
  (Optional)
  Optional. The SSH configuration.
    
* `tags` -
  (Optional)
  Optional. Key/value metadata to assign to each underlying AWS resource. Specify at most 50 pairs containing alphanumerics, spaces, and symbols (.+-=_:@/). Keys can be up to 127 Unicode characters. Values can be up to 255 Unicode characters.
    
* `taints` -
  (Optional)
  Optional. The initial taints assigned to nodes of this node pool.
    
The `config_encryption` block supports:
    
* `kms_key_arn` -
  (Required)
  Required. The ARN of the AWS KMS key used to encrypt node pool configuration.
    
The `max_pods_constraint` block supports:
    
* `max_pods_per_node` -
  (Required)
  Required. The maximum number of pods to schedule on a single node.
    
- - -

* `annotations` -
  (Optional)
  Optional. Annotations on the node pool. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Key can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.
  
* `project` -
  (Optional)
  The project for the resource
  


The `root_volume` block supports:
    
* `iops` -
  (Optional)
  Optional. The number of I/O operations per second (IOPS) to provision for GP3 volume.
    
* `kms_key_arn` -
  (Optional)
  Optional. The Amazon Resource Name (ARN) of the Customer Managed Key (CMK) used to encrypt AWS EBS volumes. If not specified, the default Amazon managed key associated to the AWS region where this cluster runs will be used.
    
* `size_gib` -
  (Optional)
  Optional. The size of the volume, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.
    
* `volume_type` -
  (Optional)
  Optional. Type of the EBS volume. When unspecified, it defaults to GP2 volume. Possible values: VOLUME_TYPE_UNSPECIFIED, GP2, GP3
    
The `ssh_config` block supports:
    
* `ec2_key_pair` -
  (Required)
  Required. The name of the EC2 key pair used to login into cluster machines.
    
The `taints` block supports:
    
* `effect` -
  (Required)
  Required. The taint effect. Possible values: EFFECT_UNSPECIFIED, NO_SCHEDULE, PREFER_NO_SCHEDULE, NO_EXECUTE
    
* `key` -
  (Required)
  Required. Key for the taint.
    
* `value` -
  (Required)
  Required. Value for the taint.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools/{{name}}`

* `create_time` -
  Output only. The time at which this node pool was created.
  
* `etag` -
  Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.
  
* `reconciling` -
  Output only. If set, there are currently changes in flight to the node pool.
  
* `state` -
  Output only. The lifecycle state of the node pool. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED
  
* `uid` -
  Output only. A globally unique identifier for the node pool.
  
* `update_time` -
  Output only. The time at which this node pool was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

NodePool can be imported using any of these accepted formats:

```
$ terraform import google_container_aws_node_pool.default projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools/{{name}}
$ terraform import google_container_aws_node_pool.default {{project}}/{{location}}/{{cluster}}/{{name}}
$ terraform import google_container_aws_node_pool.default {{location}}/{{cluster}}/{{name}}
```



