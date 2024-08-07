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
subcategory: "Cloud Deploy"
description: |-
  The Cloud Deploy `Target` resource
---

# google_clouddeploy_target

The Cloud Deploy `Target` resource

## Example Usage - multi_target
tests creating and updating a multi-target
```hcl
resource "google_clouddeploy_target" "primary" {
  location          = "us-west1"
  name              = "target"
  deploy_parameters = {}
  description       = "multi-target description"

  execution_configs {
    usages            = ["RENDER", "DEPLOY"]
    execution_timeout = "3600s"
  }

  multi_target {
    target_ids = ["1", "2"]
  }

  project          = "my-project-name"
  require_approval = false

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }
  provider          = google-beta
}

```
## Example Usage - run_target
tests creating and updating a cloud run target
```hcl
resource "google_clouddeploy_target" "primary" {
  location          = "us-west1"
  name              = "target"
  deploy_parameters = {}
  description       = "basic description"

  execution_configs {
    usages            = ["RENDER", "DEPLOY"]
    execution_timeout = "3600s"
  }

  project          = "my-project-name"
  require_approval = false

  run {
    location = "projects/my-project-name/locations/us-west1"
  }

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }
  provider          = google-beta
}

```
## Example Usage - target
Creates a basic Cloud Deploy target
```hcl
resource "google_clouddeploy_target" "primary" {
  location = "us-west1"
  name     = "target"

  deploy_parameters = {
    deployParameterKey = "deployParameterValue"
  }

  description = "basic description"

  gke {
    cluster = "projects/my-project-name/locations/us-west1/clusters/example-cluster-name"
  }

  project          = "my-project-name"
  require_approval = false

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }
}


```

## Argument Reference

The following arguments are supported:

* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  Name of the `Target`. Format is `[a-z]([a-z0-9-]{0,61}[a-z0-9])?`.
  


- - -

* `annotations` -
  (Optional)
  Optional. User annotations. These attributes can only be set and used by the user, and not by Google Cloud Deploy. See https://google.aip.dev/128#annotations for more details such as format and size limitations.

**Note**: This field is non-authoritative, and will only manage the annotations present in your configuration.
Please refer to the field `effective_annotations` for all of the annotations present on the resource.
  
* `anthos_cluster` -
  (Optional)
  Information specifying an Anthos Cluster.
  
* `custom_target` -
  (Optional)
  Optional. Information specifying a Custom Target.
  
* `deploy_parameters` -
  (Optional)
  Optional. The deploy parameters to use for this target.
  
* `description` -
  (Optional)
  Optional. Description of the `Target`. Max length is 255 characters.
  
* `execution_configs` -
  (Optional)
  Configurations for all execution that relates to this `Target`. Each `ExecutionEnvironmentUsage` value may only be used in a single configuration; using the same value multiple times is an error. When one or more configurations are specified, they must include the `RENDER` and `DEPLOY` `ExecutionEnvironmentUsage` values. When no configurations are specified, execution will use the default specified in `DefaultPool`.
  
* `gke` -
  (Optional)
  Information specifying a GKE Cluster.
  
* `labels` -
  (Optional)
  Optional. Labels are attributes that can be set and used by both the user and by Google Cloud Deploy. Labels must meet the following constraints: * Keys and values can contain only lowercase letters, numeric characters, underscores, and dashes. * All characters must use UTF-8 encoding, and international characters are allowed. * Keys must start with a lowercase letter or international character. * Each resource is limited to a maximum of 64 labels. Both keys and values are additionally constrained to be <= 128 bytes.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field `effective_labels` for all of the labels present on the resource.
  
* `multi_target` -
  (Optional)
  Information specifying a multiTarget.
  
* `project` -
  (Optional)
  The project for the resource
  
* `require_approval` -
  (Optional)
  Optional. Whether or not the `Target` requires approval.
  
* `run` -
  (Optional)
  Information specifying a Cloud Run deployment target.
  


The `anthos_cluster` block supports:
    
* `membership` -
  (Optional)
  Membership of the GKE Hub-registered cluster to which to apply the Skaffold configuration. Format is `projects/{project}/locations/{location}/memberships/{membership_name}`.
    
The `custom_target` block supports:
    
* `custom_target_type` -
  (Required)
  Required. The name of the CustomTargetType. Format must be `projects/{project}/locations/{location}/customTargetTypes/{custom_target_type}`.
    
The `execution_configs` block supports:
    
* `artifact_storage` -
  (Optional)
  Optional. Cloud Storage location in which to store execution outputs. This can either be a bucket ("gs://my-bucket") or a path within a bucket ("gs://my-bucket/my-dir"). If unspecified, a default bucket located in the same region will be used.
    
* `execution_timeout` -
  (Optional)
  Optional. Execution timeout for a Cloud Build Execution. This must be between 10m and 24h in seconds format. If unspecified, a default timeout of 1h is used.
    
* `service_account` -
  (Optional)
  Optional. Google service account to use for execution. If unspecified, the project execution service account (-compute@developer.gserviceaccount.com) is used.
    
* `usages` -
  (Required)
  Required. Usages when this configuration should be applied.
    
* `verbose` -
  (Optional)
  Optional. If true, additional logging will be enabled when running builds in this execution environment.
    
* `worker_pool` -
  (Optional)
  Optional. The resource name of the `WorkerPool`, with the format `projects/{project}/locations/{location}/workerPools/{worker_pool}`. If this optional field is unspecified, the default Cloud Build pool will be used.
    
The `gke` block supports:
    
* `cluster` -
  (Optional)
  Information specifying a GKE Cluster. Format is `projects/{project_id}/locations/{location_id}/clusters/{cluster_id}.
    
* `internal_ip` -
  (Optional)
  Optional. If true, `cluster` is accessed using the private IP address of the control plane endpoint. Otherwise, the default IP address of the control plane endpoint is used. The default IP address is the private IP address for clusters with private control-plane endpoints and the public IP address otherwise. Only specify this option when `cluster` is a [private GKE cluster](https://cloud.google.com/kubernetes-engine/docs/concepts/private-cluster-concept).
    
* `proxy_url` -
  (Optional)
  Optional. If set, used to configure a [proxy](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#proxy) to the Kubernetes server.
    
The `multi_target` block supports:
    
* `target_ids` -
  (Required)
  Required. The target_ids of this multiTarget.
    
The `run` block supports:
    
* `location` -
  (Required)
  Required. The location where the Cloud Run Service should be located. Format is `projects/{project}/locations/{location}`.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/targets/{{name}}`

* `create_time` -
  Output only. Time at which the `Target` was created.
  
* `effective_annotations` -
  All of annotations (key/value pairs) present on the resource in GCP, including the annotations configured through Terraform, other clients and services.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `etag` -
  Optional. This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.
  
* `target_id` -
  Output only. Resource id of the `Target`.
  
* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.
  
* `uid` -
  Output only. Unique identifier of the `Target`.
  
* `update_time` -
  Output only. Most recent time at which the `Target` was updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Target can be imported using any of these accepted formats:
* `projects/{{project}}/locations/{{location}}/targets/{{name}}`
* `{{project}}/{{location}}/{{name}}`
* `{{location}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Target using one of the formats above. For example:


```tf
import {
  id = "projects/{{project}}/locations/{{location}}/targets/{{name}}"
  to = google_clouddeploy_target.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Target can be imported using one of the formats above. For example:

```
$ terraform import google_clouddeploy_target.default projects/{{project}}/locations/{{location}}/targets/{{name}}
$ terraform import google_clouddeploy_target.default {{project}}/{{location}}/{{name}}
$ terraform import google_clouddeploy_target.default {{location}}/{{name}}
```



