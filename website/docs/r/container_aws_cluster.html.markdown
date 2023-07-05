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
description: |-
  An Anthos cluster running on AWS.
---

# google_container_aws_cluster

An Anthos cluster running on AWS.

For more information, see:
* [Multicloud overview](https://cloud.google.com/anthos/clusters/docs/multi-cloud)
## Example Usage - basic_aws_cluster
A basic example of a containeraws cluster
```hcl
data "google_container_aws_versions" "versions" {
  project = "my-project-name"
  location = "us-west1"
}

resource "google_container_aws_cluster" "primary" {
  authorization {
    admin_users {
      username = "my@service-account.com"
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
      owner = "my@service-account.com"
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


```

## Argument Reference

The following arguments are supported:

* `authorization` -
  (Required)
  Configuration related to the cluster RBAC settings.
  
* `aws_region` -
  (Required)
  The AWS region where the cluster runs. Each Google Cloud region supports a subset of nearby AWS regions. You can call to list all supported AWS regions within a given Google Cloud region.
  
* `control_plane` -
  (Required)
  Configuration related to the cluster control plane.
  
* `fleet` -
  (Required)
  Fleet configuration.
  
* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  The name of this resource.
  
* `networking` -
  (Required)
  Cluster-wide networking configuration.
  


The `authorization` block supports:
    
* `admin_users` -
  (Required)
  Users to perform operations as a cluster admin. A managed ClusterRoleBinding will be created to grant the `cluster-admin` ClusterRole to the users. Up to ten admin users can be provided. For more info on RBAC, see https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles
    
The `admin_users` block supports:
    
* `username` -
  (Required)
  The name of the user, e.g. `my-gcp-id@gmail.com`.
    
The `control_plane` block supports:
    
* `aws_services_authentication` -
  (Required)
  Authentication configuration for management of AWS resources.
    
* `config_encryption` -
  (Required)
  The ARN of the AWS KMS key used to encrypt cluster configuration.
    
* `database_encryption` -
  (Required)
  The ARN of the AWS KMS key used to encrypt cluster secrets.
    
* `iam_instance_profile` -
  (Required)
  The name of the AWS IAM instance pofile to assign to each control plane replica.
    
* `instance_placement` -
  (Optional)
  (Beta only) Details of placement information for an instance.
    
* `instance_type` -
  (Optional)
  Optional. The AWS instance type. When unspecified, it defaults to `m5.large`.
    
* `main_volume` -
  (Optional)
  Optional. Configuration related to the main volume provisioned for each control plane replica. The main volume is in charge of storing all of the cluster's etcd state. Volumes will be provisioned in the availability zone associated with the corresponding subnet. When unspecified, it defaults to 8 GiB with the GP2 volume type.
    
* `proxy_config` -
  (Optional)
  Proxy configuration for outbound HTTP(S) traffic.
    
* `root_volume` -
  (Optional)
  Optional. Configuration related to the root volume provisioned for each control plane replica. Volumes will be provisioned in the availability zone associated with the corresponding subnet. When unspecified, it defaults to 32 GiB with the GP2 volume type.
    
* `security_group_ids` -
  (Optional)
  Optional. The IDs of additional security groups to add to control plane replicas. The Anthos Multi-Cloud API will automatically create and manage security groups with the minimum rules needed for a functioning cluster.
    
* `ssh_config` -
  (Optional)
  Optional. SSH configuration for how to access the underlying control plane machines.
    
* `subnet_ids` -
  (Required)
  The list of subnets where control plane replicas will run. A replica will be provisioned on each subnet and up to three values can be provided. Each subnet must be in a different AWS Availability Zone (AZ).
    
* `tags` -
  (Optional)
  Optional. A set of AWS resource tags to propagate to all underlying managed AWS resources. Specify at most 50 pairs containing alphanumerics, spaces, and symbols (.+-=_:@/). Keys can be up to 127 Unicode characters. Values can be up to 255 Unicode characters.
    
* `version` -
  (Required)
  The Kubernetes version to run on control plane replicas (e.g. `1.19.10-gke.1000`). You can list all supported versions on a given Google Cloud region by calling .
    
The `aws_services_authentication` block supports:
    
* `role_arn` -
  (Required)
  The Amazon Resource Name (ARN) of the role that the Anthos Multi-Cloud API will assume when managing AWS resources on your account.
    
* `role_session_name` -
  (Optional)
  Optional. An identifier for the assumed role session. When unspecified, it defaults to `multicloud-service-agent`.
    
The `config_encryption` block supports:
    
* `kms_key_arn` -
  (Required)
  The ARN of the AWS KMS key used to encrypt cluster configuration.
    
The `database_encryption` block supports:
    
* `kms_key_arn` -
  (Required)
  The ARN of the AWS KMS key used to encrypt cluster secrets.
    
The `fleet` block supports:
    
* `membership` -
  The name of the managed Hub Membership resource associated to this cluster. Membership names are formatted as projects/<project-number>/locations/global/membership/<cluster-id>.
    
* `project` -
  (Optional)
  The number of the Fleet host project where this cluster will be registered.
    
The `networking` block supports:
    
* `per_node_pool_sg_rules_disabled` -
  (Optional)
  Disable the per node pool subnet security group rules on the control plane security group. When set to true, you must also provide one or more security groups that ensure node pools are able to send requests to the control plane on TCP/443 and TCP/8132. Failure to do so may result in unavailable node pools.
    
* `pod_address_cidr_blocks` -
  (Required)
  All pods in the cluster are assigned an RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.
    
* `service_address_cidr_blocks` -
  (Required)
  All services in the cluster are assigned an RFC1918 IPv4 address from these ranges. Only a single range is supported. This field cannot be changed after creation.
    
* `vpc_id` -
  (Required)
  The VPC associated with the cluster. All component clusters (i.e. control plane and node pools) run on a single VPC. This field cannot be changed after creation.
    
- - -

* `annotations` -
  (Optional)
  Optional. Annotations on the cluster. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Key can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.
  
* `description` -
  (Optional)
  Optional. A human readable description of this cluster. Cannot be longer than 255 UTF-8 encoded bytes.
  
* `logging_config` -
  (Optional)
  (Beta only) Logging configuration.
  
* `project` -
  (Optional)
  The project for the resource
  


The `instance_placement` block supports:
    
* `tenancy` -
  (Optional)
  The tenancy for the instance. Possible values: TENANCY_UNSPECIFIED, DEFAULT, DEDICATED, HOST
    
The `main_volume` block supports:
    
* `iops` -
  (Optional)
  Optional. The number of I/O operations per second (IOPS) to provision for GP3 volume.
    
* `kms_key_arn` -
  (Optional)
  Optional. The Amazon Resource Name (ARN) of the Customer Managed Key (CMK) used to encrypt AWS EBS volumes. If not specified, the default Amazon managed key associated to the AWS region where this cluster runs will be used.
    
* `size_gib` -
  (Optional)
  Optional. The size of the volume, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.
    
* `throughput` -
  (Optional)
  Optional. The throughput to provision for the volume, in MiB/s. Only valid if the volume type is GP3.
    
* `volume_type` -
  (Optional)
  Optional. Type of the EBS volume. When unspecified, it defaults to GP2 volume. Possible values: VOLUME_TYPE_UNSPECIFIED, GP2, GP3
    
The `proxy_config` block supports:
    
* `secret_arn` -
  (Required)
  The ARN of the AWS Secret Manager secret that contains the HTTP(S) proxy configuration.
    
* `secret_version` -
  (Required)
  The version string of the AWS Secret Manager secret that contains the HTTP(S) proxy configuration.
    
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
    
* `throughput` -
  (Optional)
  Optional. The throughput to provision for the volume, in MiB/s. Only valid if the volume type is GP3.
    
* `volume_type` -
  (Optional)
  Optional. Type of the EBS volume. When unspecified, it defaults to GP2 volume. Possible values: VOLUME_TYPE_UNSPECIFIED, GP2, GP3
    
The `ssh_config` block supports:
    
* `ec2_key_pair` -
  (Required)
  The name of the EC2 key pair used to login into cluster machines.
    
The `logging_config` block supports:
    
* `component_config` -
  (Optional)
  Configuration of the logging components.
    
The `component_config` block supports:
    
* `enable_components` -
  (Optional)
  Components of the logging configuration to be enabled.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/awsClusters/{{name}}`

* `create_time` -
  Output only. The time at which this cluster was created.
  
* `endpoint` -
  Output only. The endpoint of the cluster's API server.
  
* `etag` -
  Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.
  
* `reconciling` -
  Output only. If set, there are currently changes in flight to the cluster.
  
* `state` -
  Output only. The current state of the cluster. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED
  
* `uid` -
  Output only. A globally unique identifier for the cluster.
  
* `update_time` -
  Output only. The time at which this cluster was last updated.
  
* `workload_identity_config` -
  Output only. Workload Identity settings.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Cluster can be imported using any of these accepted formats:

```
$ terraform import google_container_aws_cluster.default projects/{{project}}/locations/{{location}}/awsClusters/{{name}}
$ terraform import google_container_aws_cluster.default {{project}}/{{location}}/{{name}}
$ terraform import google_container_aws_cluster.default {{location}}/{{name}}
```



