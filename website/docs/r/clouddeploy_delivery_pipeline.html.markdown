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
  The Cloud Deploy `DeliveryPipeline` resource
---

# google_clouddeploy_delivery_pipeline

The Cloud Deploy `DeliveryPipeline` resource

## Example Usage - canary_delivery_pipeline
Creates a basic Cloud Deploy delivery pipeline
```hcl
resource "google_clouddeploy_delivery_pipeline" "primary" {
  location = "us-west1"
  name     = "pipeline"

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  description = "basic description"

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }

  project = "my-project-name"

  serial_pipeline {
    stages {
      deploy_parameters {
        values = {
          deployParameterKey = "deployParameterValue"
        }

        match_target_labels = {}
      }

      profiles  = ["example-profile-one", "example-profile-two"]
      target_id = "example-target-one"
    }

    stages {
      profiles  = []
      target_id = "example-target-two"
    }
  }
  provider = google-beta
}

```
## Example Usage - canary_service_networking_delivery_pipeline
Creates a basic Cloud Deploy delivery pipeline
```hcl
resource "google_clouddeploy_delivery_pipeline" "primary" {
  location = "us-west1"
  name     = "pipeline"

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  description = "basic description"

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }

  project = "my-project-name"

  serial_pipeline {
    stages {
      deploy_parameters {
        values = {
          deployParameterKey = "deployParameterValue"
        }

        match_target_labels = {}
      }

      profiles  = ["example-profile-one", "example-profile-two"]
      target_id = "example-target-one"
    }

    stages {
      profiles  = []
      target_id = "example-target-two"
    }
  }
  provider = google-beta
}

```
## Example Usage - canaryrun_delivery_pipeline
Creates a basic Cloud Deploy delivery pipeline
```hcl
resource "google_clouddeploy_delivery_pipeline" "primary" {
  location = "us-west1"
  name     = "pipeline"

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  description = "basic description"

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }

  project = "my-project-name"

  serial_pipeline {
    stages {
      deploy_parameters {
        values = {
          deployParameterKey = "deployParameterValue"
        }

        match_target_labels = {}
      }

      profiles  = ["example-profile-one", "example-profile-two"]
      target_id = "example-target-one"
    }

    stages {
      profiles  = []
      target_id = "example-target-two"
    }
  }
  provider = google-beta
}

```
## Example Usage - delivery_pipeline
Creates a basic Cloud Deploy delivery pipeline
```hcl
resource "google_clouddeploy_delivery_pipeline" "primary" {
  location = "us-west1"
  name     = "pipeline"

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  description = "basic description"

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }

  project = "my-project-name"

  serial_pipeline {
    stages {
      deploy_parameters {
        values = {
          deployParameterKey = "deployParameterValue"
        }

        match_target_labels = {}
      }

      profiles  = ["example-profile-one", "example-profile-two"]
      target_id = "example-target-one"
    }

    stages {
      profiles  = []
      target_id = "example-target-two"
    }
  }
}


```
## Example Usage - verify_delivery_pipeline
tests creating and updating a delivery pipeline with deployment verification strategy
```hcl
resource "google_clouddeploy_delivery_pipeline" "primary" {
  location = "us-west1"
  name     = "pipeline"

  annotations = {
    my_first_annotation = "example-annotation-1"

    my_second_annotation = "example-annotation-2"
  }

  description = "basic description"

  labels = {
    my_first_label = "example-label-1"

    my_second_label = "example-label-2"
  }

  project = "my-project-name"

  serial_pipeline {
    stages {
      deploy_parameters {
        values = {
          deployParameterKey = "deployParameterValue"
        }

        match_target_labels = {}
      }

      profiles  = ["example-profile-one", "example-profile-two"]
      target_id = "example-target-one"
    }

    stages {
      profiles  = []
      target_id = "example-target-two"
    }
  }
  provider = google-beta
}

```

## Argument Reference

The following arguments are supported:

* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  Name of the `DeliveryPipeline`. Format is [a-z][a-z0-9\-]{0,62}.
  


The `phase_configs` block supports:
    
* `percentage` -
  (Required)
  Required. Percentage deployment for the phase.
    
* `phase_id` -
  (Required)
  Required. The ID to assign to the `Rollout` phase. This value must consist of lower-case letters, numbers, and hyphens, start with a letter and end with a letter or a number, and have a max length of 63 characters. In other words, it must match the following regex: `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`.
    
* `profiles` -
  (Optional)
  Skaffold profiles to use when rendering the manifest for this phase. These are in addition to the profiles list specified in the `DeliveryPipeline` stage.
    
* `verify` -
  (Optional)
  Whether to run verify tests after the deployment.
    
- - -

* `annotations` -
  (Optional)
  User annotations. These attributes can only be set and used by the user, and not by Google Cloud Deploy. See https://google.aip.dev/128#annotations for more details such as format and size limitations.
  
* `description` -
  (Optional)
  Description of the `DeliveryPipeline`. Max length is 255 characters.
  
* `labels` -
  (Optional)
  Labels are attributes that can be set and used by both the user and by Google Cloud Deploy. Labels must meet the following constraints: * Keys and values can contain only lowercase letters, numeric characters, underscores, and dashes. * All characters must use UTF-8 encoding, and international characters are allowed. * Keys must start with a lowercase letter or international character. * Each resource is limited to a maximum of 64 labels. Both keys and values are additionally constrained to be <= 128 bytes.
  
* `project` -
  (Optional)
  The project for the resource
  
* `serial_pipeline` -
  (Optional)
  SerialPipeline defines a sequential set of stages for a `DeliveryPipeline`.
  
* `suspended` -
  (Optional)
  When suspended, no new releases or rollouts can be created, but in-progress ones will complete.
  


The `serial_pipeline` block supports:
    
* `stages` -
  (Optional)
  Each stage specifies configuration for a `Target`. The ordering of this list defines the promotion flow.
    
The `stages` block supports:
    
* `deploy_parameters` -
  (Optional)
  Optional. The deploy parameters to use for the target in this stage.
    
* `profiles` -
  (Optional)
  Skaffold profiles to use when rendering the manifest for this stage's `Target`.
    
* `strategy` -
  (Optional)
  Optional. The strategy to use for a `Rollout` to this stage.
    
* `target_id` -
  (Optional)
  The target_id to which this stage points. This field refers exclusively to the last segment of a target name. For example, this field would just be `my-target` (rather than `projects/project/locations/location/targets/my-target`). The location of the `Target` is inferred to be the same as the location of the `DeliveryPipeline` that contains this `Stage`.
    
The `deploy_parameters` block supports:
    
* `match_target_labels` -
  (Optional)
  Optional. Deploy parameters are applied to targets with match labels. If unspecified, deploy parameters are applied to all targets (including child targets of a multi-target).
    
* `values` -
  (Required)
  Required. Values are deploy parameters in key-value pairs.
    
The `strategy` block supports:
    
* `canary` -
  (Optional)
  Canary deployment strategy provides progressive percentage based deployments to a Target.
    
* `standard` -
  (Optional)
  Standard deployment strategy executes a single deploy and allows verifying the deployment.
    
The `canary` block supports:
    
* `canary_deployment` -
  (Optional)
  Configures the progressive based deployment for a Target.
    
* `custom_canary_deployment` -
  (Optional)
  Configures the progressive based deployment for a Target, but allows customizing at the phase level where a phase represents each of the percentage deployments.
    
* `runtime_config` -
  (Optional)
  Optional. Runtime specific configurations for the deployment strategy. The runtime configuration is used to determine how Cloud Deploy will split traffic to enable a progressive deployment.
    
The `canary_deployment` block supports:
    
* `percentages` -
  (Required)
  Required. The percentage based deployments that will occur as a part of a `Rollout`. List is expected in ascending order and each integer n is 0 <= n < 100.
    
* `verify` -
  (Optional)
  Whether to run verify tests after each percentage deployment.
    
The `custom_canary_deployment` block supports:
    
* `phase_configs` -
  (Required)
  Required. Configuration for each phase in the canary deployment in the order executed.
    
The `runtime_config` block supports:
    
* `cloud_run` -
  (Optional)
  Cloud Run runtime configuration.
    
* `kubernetes` -
  (Optional)
  Kubernetes runtime configuration.
    
The `cloud_run` block supports:
    
* `automatic_traffic_control` -
  (Optional)
  Whether Cloud Deploy should update the traffic stanza in a Cloud Run Service on the user's behalf to facilitate traffic splitting. This is required to be true for CanaryDeployments, but optional for CustomCanaryDeployments.
    
The `kubernetes` block supports:
    
* `gateway_service_mesh` -
  (Optional)
  Kubernetes Gateway API service mesh configuration.
    
* `service_networking` -
  (Optional)
  Kubernetes Service networking configuration.
    
The `gateway_service_mesh` block supports:
    
* `deployment` -
  (Required)
  Required. Name of the Kubernetes Deployment whose traffic is managed by the specified HTTPRoute and Service.
    
* `http_route` -
  (Required)
  Required. Name of the Gateway API HTTPRoute.
    
* `service` -
  (Required)
  Required. Name of the Kubernetes Service.
    
The `service_networking` block supports:
    
* `deployment` -
  (Required)
  Required. Name of the Kubernetes Deployment whose traffic is managed by the specified Service.
    
* `disable_pod_overprovisioning` -
  (Optional)
  Optional. Whether to disable Pod overprovisioning. If Pod overprovisioning is disabled then Cloud Deploy will limit the number of total Pods used for the deployment strategy to the number of Pods the Deployment has on the cluster.
    
* `service` -
  (Required)
  Required. Name of the Kubernetes Service.
    
The `standard` block supports:
    
* `verify` -
  (Optional)
  Whether to verify a deployment.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}`

* `condition` -
  Output only. Information around the state of the Delivery Pipeline.
  
* `create_time` -
  Output only. Time at which the pipeline was created.
  
* `etag` -
  This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.
  
* `uid` -
  Output only. Unique identifier of the `DeliveryPipeline`.
  
* `update_time` -
  Output only. Most recent time at which the pipeline was updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

DeliveryPipeline can be imported using any of these accepted formats:

```
$ terraform import google_clouddeploy_delivery_pipeline.default projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}
$ terraform import google_clouddeploy_delivery_pipeline.default {{project}}/{{location}}/{{name}}
$ terraform import google_clouddeploy_delivery_pipeline.default {{location}}/{{name}}
```



