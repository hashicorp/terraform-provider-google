// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package clouddeploy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	clouddeploy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/clouddeploy"

	"github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceClouddeployDeliveryPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceClouddeployDeliveryPipelineCreate,
		Read:   resourceClouddeployDeliveryPipelineRead,
		Update: resourceClouddeployDeliveryPipelineUpdate,
		Delete: resourceClouddeployDeliveryPipelineDelete,

		Importer: &schema.ResourceImporter{
			State: resourceClouddeployDeliveryPipelineImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
			tpgresource.SetLabelsDiff,
			tpgresource.SetAnnotationsDiff,
		),

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location for the resource",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the `DeliveryPipeline`. Format is `[a-z]([a-z0-9-]{0,61}[a-z0-9])?`.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the `DeliveryPipeline`. Max length is 255 characters.",
			},

			"effective_annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "All of annotations (key/value pairs) present on the resource in GCP, including the annotations configured through Terraform, other clients and services.",
			},

			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.",
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"serial_pipeline": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "SerialPipeline defines a sequential set of stages for a `DeliveryPipeline`.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineSchema(),
			},

			"suspended": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When suspended, no new releases or rollouts can be created, but in-progress ones will complete.",
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "User annotations. These attributes can only be set and used by the user, and not by Google Cloud Deploy. See https://google.aip.dev/128#annotations for more details such as format and size limitations.\n\n**Note**: This field is non-authoritative, and will only manage the annotations present in your configuration.\nPlease refer to the field `effective_annotations` for all of the annotations present on the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output only. Information around the state of the Delivery Pipeline.",
				Elem:        ClouddeployDeliveryPipelineConditionSchema(),
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time at which the pipeline was created.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels are attributes that can be set and used by both the user and by Google Cloud Deploy. Labels must meet the following constraints: * Keys and values can contain only lowercase letters, numeric characters, underscores, and dashes. * All characters must use UTF-8 encoding, and international characters are allowed. * Keys must start with a lowercase letter or international character. * Each resource is limited to a maximum of 64 labels. Both keys and values are additionally constrained to be <= 128 bytes.\n\n**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.\nPlease refer to the field `effective_labels` for all of the labels present on the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"terraform_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The combination of labels configured directly on the resource and default labels configured on the provider.",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Unique identifier of the `DeliveryPipeline`.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Most recent time at which the pipeline was updated.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"stages": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Each stage specifies configuration for a `Target`. The ordering of this list defines the promotion flow.",
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"deploy_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. The deploy parameters to use for the target in this stage.",
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesDeployParametersSchema(),
			},

			"profiles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Skaffold profiles to use when rendering the manifest for this stage's `Target`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. The strategy to use for a `Rollout` to this stage.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategySchema(),
			},

			"target_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The target_id to which this stage points. This field refers exclusively to the last segment of a target name. For example, this field would just be `my-target` (rather than `projects/project/locations/location/targets/my-target`). The location of the `Target` is inferred to be the same as the location of the `DeliveryPipeline` that contains this `Stage`.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesDeployParametersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"values": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Required. Values are deploy parameters in key-value pairs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"match_target_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. Deploy parameters are applied to targets with match labels. If unspecified, deploy parameters are applied to all targets (including child targets of a multi-target).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"canary": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Canary deployment strategy provides progressive percentage based deployments to a Target.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanarySchema(),
			},

			"standard": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Standard deployment strategy executes a single deploy and allows verifying the deployment.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanarySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"canary_deployment": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configures the progressive based deployment for a Target.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentSchema(),
			},

			"custom_canary_deployment": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configures the progressive based deployment for a Target, but allows customizing at the phase level where a phase represents each of the percentage deployments.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentSchema(),
			},

			"runtime_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Runtime specific configurations for the deployment strategy. The runtime configuration is used to determine how Cloud Deploy will split traffic to enable a progressive deployment.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"percentages": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The percentage based deployments that will occur as a part of a `Rollout`. List is expected in ascending order and each integer n is 0 <= n < 100.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},

			"postdeploy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Configuration for the postdeploy job of the last phase. If this is not configured, postdeploy job will not be present.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploySchema(),
			},

			"predeploy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Configuration for the predeploy job of the first phase. If this is not configured, predeploy job will not be present.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploySchema(),
			},

			"verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to run verify tests after each percentage deployment.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the postdeploy job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the predeploy job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"phase_configs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Configuration for each phase in the canary deployment in the order executed.",
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"percentage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Required. Percentage deployment for the phase.",
			},

			"phase_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The ID to assign to the `Rollout` phase. This value must consist of lower-case letters, numbers, and hyphens, start with a letter and end with a letter or a number, and have a max length of 63 characters. In other words, it must match the following regex: `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`.",
			},

			"postdeploy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Configuration for the postdeploy job of this phase. If this is not configured, postdeploy job will not be present for this phase.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploySchema(),
			},

			"predeploy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Configuration for the predeploy job of this phase. If this is not configured, predeploy job will not be present for this phase.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploySchema(),
			},

			"profiles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Skaffold profiles to use when rendering the manifest for this phase. These are in addition to the profiles list specified in the `DeliveryPipeline` stage.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to run verify tests after the deployment.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the postdeploy job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the predeploy job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cloud_run": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cloud Run runtime configuration.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRunSchema(),
			},

			"kubernetes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Kubernetes runtime configuration.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRunSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"automatic_traffic_control": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether Cloud Deploy should update the traffic stanza in a Cloud Run Service on the user's behalf to facilitate traffic splitting. This is required to be true for CanaryDeployments, but optional for CustomCanaryDeployments.",
			},

			"canary_revision_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A list of tags that are added to the canary revision while the canary phase is in progress.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"prior_revision_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A list of tags that are added to the prior revision while the canary phase is in progress.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"stable_revision_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A list of tags that are added to the final stable revision when the stable phase is applied.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gateway_service_mesh": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Kubernetes Gateway API service mesh configuration.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshSchema(),
			},

			"service_networking": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Kubernetes Service networking configuration.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworkingSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"deployment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Kubernetes Deployment whose traffic is managed by the specified HTTPRoute and Service.",
			},

			"http_route": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Gateway API HTTPRoute.",
			},

			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Kubernetes Service.",
			},

			"pod_selector_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The label to use when selecting Pods for the Deployment and Service resources. This label must already be present in both resources.",
			},

			"route_destinations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Route destinations allow configuring the Gateway API HTTPRoute to be deployed to additional clusters. This option is available for multi-cluster service mesh set ups that require the route to exist in the clusters that call the service. If unspecified, the HTTPRoute will only be deployed to the Target cluster.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinationsSchema(),
			},

			"route_update_wait_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The time to wait for route updates to propagate. The maximum configurable time is 3 hours, in seconds format. If unspecified, there is no wait time.",
			},

			"stable_cutback_duration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The amount of time to migrate traffic back from the canary Service to the original Service during the stable phase deployment. If specified, must be between 15s and 3600s. If unspecified, there is no cutback time.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The clusters where the Gateway API HTTPRoute resource will be deployed to. Valid entries include the associated entities IDs configured in the Target resource and \"@self\" to include the Target cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"propagate_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. Whether to propagate the Kubernetes Service to the route destination clusters. The Service will always be deployed to the Target cluster even if the HTTPRoute is not. This option may be used to facilitiate successful DNS lookup in the route destination clusters. Can only be set to true if destinations are specified.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworkingSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"deployment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Kubernetes Deployment whose traffic is managed by the specified Service.",
			},

			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. Name of the Kubernetes Service.",
			},

			"disable_pod_overprovisioning": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. Whether to disable Pod overprovisioning. If Pod overprovisioning is disabled then Cloud Deploy will limit the number of total Pods used for the deployment strategy to the number of Pods the Deployment has on the cluster.",
			},

			"pod_selector_label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The label to use when selecting Pods for the Deployment resource. This label must already be present in the Deployment.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"postdeploy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Configuration for the postdeploy job. If this is not configured, postdeploy job will not be present.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploySchema(),
			},

			"predeploy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Configuration for the predeploy job. If this is not configured, predeploy job will not be present.",
				MaxItems:    1,
				Elem:        ClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPredeploySchema(),
			},

			"verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to verify a deployment.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the postdeploy job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPredeploySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the predeploy job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployDeliveryPipelineConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pipeline_ready_condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details around the Pipeline's overall status.",
				Elem:        ClouddeployDeliveryPipelineConditionPipelineReadyConditionSchema(),
			},

			"targets_present_condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details around targets enumerated in the pipeline.",
				Elem:        ClouddeployDeliveryPipelineConditionTargetsPresentConditionSchema(),
			},

			"targets_type_condition": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details on the whether the targets enumerated in the pipeline are of the same type.",
				Elem:        ClouddeployDeliveryPipelineConditionTargetsTypeConditionSchema(),
			},
		},
	}
}

func ClouddeployDeliveryPipelineConditionPipelineReadyConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the Pipeline is in a valid state. Otherwise at least one condition in `PipelineCondition` is in an invalid state. Iterate over those conditions and see which condition(s) has status = false to find out what is wrong with the Pipeline.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last time the condition was updated.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineConditionTargetsPresentConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"missing_targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of Target names that are missing. For example, projects/{project_id}/locations/{location_name}/targets/{target_name}.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if there aren't any missing Targets.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last time the condition was updated.",
			},
		},
	}
}

func ClouddeployDeliveryPipelineConditionTargetsTypeConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"error_details": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human readable error message.",
			},

			"status": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the targets are all a comparable type. For example this is true if all targets are GKE clusters. This is false if some targets are Cloud Run targets and others are GKE clusters.",
			},
		},
	}
}

func resourceClouddeployDeliveryPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.DeliveryPipeline{
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		Description:    dcl.String(d.Get("description").(string)),
		Annotations:    tpgresource.CheckStringMap(d.Get("effective_annotations")),
		Labels:         tpgresource.CheckStringMap(d.Get("effective_labels")),
		Project:        dcl.String(project),
		SerialPipeline: expandClouddeployDeliveryPipelineSerialPipeline(d.Get("serial_pipeline")),
		Suspended:      dcl.Bool(d.Get("suspended").(bool)),
	}

	id, err := obj.ID()
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	d.SetId(id)
	directive := tpgdclresource.CreateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLClouddeployClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyDeliveryPipeline(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating DeliveryPipeline: %s", err)
	}

	log.Printf("[DEBUG] Finished creating DeliveryPipeline %q: %#v", d.Id(), res)

	return resourceClouddeployDeliveryPipelineRead(d, meta)
}

func resourceClouddeployDeliveryPipelineRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.DeliveryPipeline{
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		Description:    dcl.String(d.Get("description").(string)),
		Annotations:    tpgresource.CheckStringMap(d.Get("effective_annotations")),
		Labels:         tpgresource.CheckStringMap(d.Get("effective_labels")),
		Project:        dcl.String(project),
		SerialPipeline: expandClouddeployDeliveryPipelineSerialPipeline(d.Get("serial_pipeline")),
		Suspended:      dcl.Bool(d.Get("suspended").(bool)),
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLClouddeployClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutRead))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.GetDeliveryPipeline(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ClouddeployDeliveryPipeline %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("description", res.Description); err != nil {
		return fmt.Errorf("error setting description in state: %s", err)
	}
	if err = d.Set("effective_annotations", res.Annotations); err != nil {
		return fmt.Errorf("error setting effective_annotations in state: %s", err)
	}
	if err = d.Set("effective_labels", res.Labels); err != nil {
		return fmt.Errorf("error setting effective_labels in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("serial_pipeline", flattenClouddeployDeliveryPipelineSerialPipeline(res.SerialPipeline)); err != nil {
		return fmt.Errorf("error setting serial_pipeline in state: %s", err)
	}
	if err = d.Set("suspended", res.Suspended); err != nil {
		return fmt.Errorf("error setting suspended in state: %s", err)
	}
	if err = d.Set("annotations", flattenClouddeployDeliveryPipelineAnnotations(res.Annotations, d)); err != nil {
		return fmt.Errorf("error setting annotations in state: %s", err)
	}
	if err = d.Set("condition", flattenClouddeployDeliveryPipelineCondition(res.Condition)); err != nil {
		return fmt.Errorf("error setting condition in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("error setting etag in state: %s", err)
	}
	if err = d.Set("labels", flattenClouddeployDeliveryPipelineLabels(res.Labels, d)); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("terraform_labels", flattenClouddeployDeliveryPipelineTerraformLabels(res.Labels, d)); err != nil {
		return fmt.Errorf("error setting terraform_labels in state: %s", err)
	}
	if err = d.Set("uid", res.Uid); err != nil {
		return fmt.Errorf("error setting uid in state: %s", err)
	}
	if err = d.Set("update_time", res.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time in state: %s", err)
	}

	return nil
}
func resourceClouddeployDeliveryPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.DeliveryPipeline{
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		Description:    dcl.String(d.Get("description").(string)),
		Annotations:    tpgresource.CheckStringMap(d.Get("effective_annotations")),
		Labels:         tpgresource.CheckStringMap(d.Get("effective_labels")),
		Project:        dcl.String(project),
		SerialPipeline: expandClouddeployDeliveryPipelineSerialPipeline(d.Get("serial_pipeline")),
		Suspended:      dcl.Bool(d.Get("suspended").(bool)),
	}
	directive := tpgdclresource.UpdateDirective
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLClouddeployClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutUpdate))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	res, err := client.ApplyDeliveryPipeline(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating DeliveryPipeline: %s", err)
	}

	log.Printf("[DEBUG] Finished creating DeliveryPipeline %q: %#v", d.Id(), res)

	return resourceClouddeployDeliveryPipelineRead(d, meta)
}

func resourceClouddeployDeliveryPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.DeliveryPipeline{
		Location:       dcl.String(d.Get("location").(string)),
		Name:           dcl.String(d.Get("name").(string)),
		Description:    dcl.String(d.Get("description").(string)),
		Annotations:    tpgresource.CheckStringMap(d.Get("effective_annotations")),
		Labels:         tpgresource.CheckStringMap(d.Get("effective_labels")),
		Project:        dcl.String(project),
		SerialPipeline: expandClouddeployDeliveryPipelineSerialPipeline(d.Get("serial_pipeline")),
		Suspended:      dcl.Bool(d.Get("suspended").(bool)),
	}

	log.Printf("[DEBUG] Deleting DeliveryPipeline %q", d.Id())
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	client := transport_tpg.NewDCLClouddeployClient(config, userAgent, billingProject, d.Timeout(schema.TimeoutDelete))
	if bp, err := tpgresource.ReplaceVars(d, config, client.Config.BasePath); err != nil {
		d.SetId("")
		return fmt.Errorf("Could not format %q: %w", client.Config.BasePath, err)
	} else {
		client.Config.BasePath = bp
	}
	if err := client.DeleteDeliveryPipeline(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting DeliveryPipeline: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting DeliveryPipeline %q", d.Id())
	return nil
}

func resourceClouddeployDeliveryPipelineImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/deliveryPipelines/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandClouddeployDeliveryPipelineSerialPipeline(o interface{}) *clouddeploy.DeliveryPipelineSerialPipeline {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipeline
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipeline
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipeline{
		Stages: expandClouddeployDeliveryPipelineSerialPipelineStagesArray(obj["stages"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipeline(obj *clouddeploy.DeliveryPipelineSerialPipeline) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"stages": flattenClouddeployDeliveryPipelineSerialPipelineStagesArray(obj.Stages),
	}

	return []interface{}{transformed}

}
func expandClouddeployDeliveryPipelineSerialPipelineStagesArray(o interface{}) []clouddeploy.DeliveryPipelineSerialPipelineStages {
	if o == nil {
		return make([]clouddeploy.DeliveryPipelineSerialPipelineStages, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]clouddeploy.DeliveryPipelineSerialPipelineStages, 0)
	}

	items := make([]clouddeploy.DeliveryPipelineSerialPipelineStages, 0, len(objs))
	for _, item := range objs {
		i := expandClouddeployDeliveryPipelineSerialPipelineStages(item)
		items = append(items, *i)
	}

	return items
}

func expandClouddeployDeliveryPipelineSerialPipelineStages(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStages {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStages
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStages{
		DeployParameters: expandClouddeployDeliveryPipelineSerialPipelineStagesDeployParametersArray(obj["deploy_parameters"]),
		Profiles:         tpgdclresource.ExpandStringArray(obj["profiles"]),
		Strategy:         expandClouddeployDeliveryPipelineSerialPipelineStagesStrategy(obj["strategy"]),
		TargetId:         dcl.String(obj["target_id"].(string)),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesArray(objs []clouddeploy.DeliveryPipelineSerialPipelineStages) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenClouddeployDeliveryPipelineSerialPipelineStages(&item)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployDeliveryPipelineSerialPipelineStages(obj *clouddeploy.DeliveryPipelineSerialPipelineStages) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"deploy_parameters": flattenClouddeployDeliveryPipelineSerialPipelineStagesDeployParametersArray(obj.DeployParameters),
		"profiles":          obj.Profiles,
		"strategy":          flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategy(obj.Strategy),
		"target_id":         obj.TargetId,
	}

	return transformed

}
func expandClouddeployDeliveryPipelineSerialPipelineStagesDeployParametersArray(o interface{}) []clouddeploy.DeliveryPipelineSerialPipelineStagesDeployParameters {
	if o == nil {
		return make([]clouddeploy.DeliveryPipelineSerialPipelineStagesDeployParameters, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]clouddeploy.DeliveryPipelineSerialPipelineStagesDeployParameters, 0)
	}

	items := make([]clouddeploy.DeliveryPipelineSerialPipelineStagesDeployParameters, 0, len(objs))
	for _, item := range objs {
		i := expandClouddeployDeliveryPipelineSerialPipelineStagesDeployParameters(item)
		items = append(items, *i)
	}

	return items
}

func expandClouddeployDeliveryPipelineSerialPipelineStagesDeployParameters(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesDeployParameters {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesDeployParameters
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesDeployParameters{
		Values:            tpgresource.CheckStringMap(obj["values"]),
		MatchTargetLabels: tpgresource.CheckStringMap(obj["match_target_labels"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesDeployParametersArray(objs []clouddeploy.DeliveryPipelineSerialPipelineStagesDeployParameters) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenClouddeployDeliveryPipelineSerialPipelineStagesDeployParameters(&item)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesDeployParameters(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesDeployParameters) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"values":              obj.Values,
		"match_target_labels": obj.MatchTargetLabels,
	}

	return transformed

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategy(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategy {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategy
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategy
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategy{
		Canary:   expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanary(obj["canary"]),
		Standard: expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandard(obj["standard"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategy(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategy) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"canary":   flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanary(obj.Canary),
		"standard": flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandard(obj.Standard),
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanary(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanary {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanary
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanary
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanary{
		CanaryDeployment:       expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment(obj["canary_deployment"]),
		CustomCanaryDeployment: expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment(obj["custom_canary_deployment"]),
		RuntimeConfig:          expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig(obj["runtime_config"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanary(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanary) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"canary_deployment":        flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment(obj.CanaryDeployment),
		"custom_canary_deployment": flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment(obj.CustomCanaryDeployment),
		"runtime_config":           flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig(obj.RuntimeConfig),
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment{
		Percentages: tpgdclresource.ExpandIntegerArray(obj["percentages"]),
		Postdeploy:  expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy(obj["postdeploy"]),
		Predeploy:   expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy(obj["predeploy"]),
		Verify:      dcl.Bool(obj["verify"].(bool)),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"percentages": obj.Percentages,
		"postdeploy":  flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy(obj.Postdeploy),
		"predeploy":   flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy(obj.Predeploy),
		"verify":      obj.Verify,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy{
		Actions: tpgdclresource.ExpandStringArray(obj["actions"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"actions": obj.Actions,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy{
		Actions: tpgdclresource.ExpandStringArray(obj["actions"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"actions": obj.Actions,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment{
		PhaseConfigs: expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsArray(obj["phase_configs"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"phase_configs": flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsArray(obj.PhaseConfigs),
	}

	return []interface{}{transformed}

}
func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsArray(o interface{}) []clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs {
	if o == nil {
		return make([]clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs, 0)
	}

	items := make([]clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs, 0, len(objs))
	for _, item := range objs {
		i := expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs(item)
		items = append(items, *i)
	}

	return items
}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs{
		Percentage: dcl.Int64(int64(obj["percentage"].(int))),
		PhaseId:    dcl.String(obj["phase_id"].(string)),
		Postdeploy: expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy(obj["postdeploy"]),
		Predeploy:  expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy(obj["predeploy"]),
		Profiles:   tpgdclresource.ExpandStringArray(obj["profiles"]),
		Verify:     dcl.Bool(obj["verify"].(bool)),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsArray(objs []clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs(&item)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"percentage": obj.Percentage,
		"phase_id":   obj.PhaseId,
		"postdeploy": flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy(obj.Postdeploy),
		"predeploy":  flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy(obj.Predeploy),
		"profiles":   obj.Profiles,
		"verify":     obj.Verify,
	}

	return transformed

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy{
		Actions: tpgdclresource.ExpandStringArray(obj["actions"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"actions": obj.Actions,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy{
		Actions: tpgdclresource.ExpandStringArray(obj["actions"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"actions": obj.Actions,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig{
		CloudRun:   expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun(obj["cloud_run"]),
		Kubernetes: expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes(obj["kubernetes"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cloud_run":  flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun(obj.CloudRun),
		"kubernetes": flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes(obj.Kubernetes),
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun{
		AutomaticTrafficControl: dcl.Bool(obj["automatic_traffic_control"].(bool)),
		CanaryRevisionTags:      tpgdclresource.ExpandStringArray(obj["canary_revision_tags"]),
		PriorRevisionTags:       tpgdclresource.ExpandStringArray(obj["prior_revision_tags"]),
		StableRevisionTags:      tpgdclresource.ExpandStringArray(obj["stable_revision_tags"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"automatic_traffic_control": obj.AutomaticTrafficControl,
		"canary_revision_tags":      obj.CanaryRevisionTags,
		"prior_revision_tags":       obj.PriorRevisionTags,
		"stable_revision_tags":      obj.StableRevisionTags,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes{
		GatewayServiceMesh: expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh(obj["gateway_service_mesh"]),
		ServiceNetworking:  expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking(obj["service_networking"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"gateway_service_mesh": flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh(obj.GatewayServiceMesh),
		"service_networking":   flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking(obj.ServiceNetworking),
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh{
		Deployment:            dcl.String(obj["deployment"].(string)),
		HttpRoute:             dcl.String(obj["http_route"].(string)),
		Service:               dcl.String(obj["service"].(string)),
		PodSelectorLabel:      dcl.String(obj["pod_selector_label"].(string)),
		RouteDestinations:     expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations(obj["route_destinations"]),
		RouteUpdateWaitTime:   dcl.String(obj["route_update_wait_time"].(string)),
		StableCutbackDuration: dcl.String(obj["stable_cutback_duration"].(string)),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"deployment":              obj.Deployment,
		"http_route":              obj.HttpRoute,
		"service":                 obj.Service,
		"pod_selector_label":      obj.PodSelectorLabel,
		"route_destinations":      flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations(obj.RouteDestinations),
		"route_update_wait_time":  obj.RouteUpdateWaitTime,
		"stable_cutback_duration": obj.StableCutbackDuration,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations{
		DestinationIds:   tpgdclresource.ExpandStringArray(obj["destination_ids"]),
		PropagateService: dcl.Bool(obj["propagate_service"].(bool)),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMeshRouteDestinations) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"destination_ids":   obj.DestinationIds,
		"propagate_service": obj.PropagateService,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking{
		Deployment:                 dcl.String(obj["deployment"].(string)),
		Service:                    dcl.String(obj["service"].(string)),
		DisablePodOverprovisioning: dcl.Bool(obj["disable_pod_overprovisioning"].(bool)),
		PodSelectorLabel:           dcl.String(obj["pod_selector_label"].(string)),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"deployment":                   obj.Deployment,
		"service":                      obj.Service,
		"disable_pod_overprovisioning": obj.DisablePodOverprovisioning,
		"pod_selector_label":           obj.PodSelectorLabel,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandard(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandard {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyStandard
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyStandard
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandard{
		Postdeploy: expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy(obj["postdeploy"]),
		Predeploy:  expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy(obj["predeploy"]),
		Verify:     dcl.Bool(obj["verify"].(bool)),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandard(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandard) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"postdeploy": flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy(obj.Postdeploy),
		"predeploy":  flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy(obj.Predeploy),
		"verify":     obj.Verify,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy{
		Actions: tpgdclresource.ExpandStringArray(obj["actions"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"actions": obj.Actions,
	}

	return []interface{}{transformed}

}

func expandClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy(o interface{}) *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy {
	if o == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyDeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy{
		Actions: tpgdclresource.ExpandStringArray(obj["actions"]),
	}
}

func flattenClouddeployDeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy(obj *clouddeploy.DeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"actions": obj.Actions,
	}

	return []interface{}{transformed}

}

func flattenClouddeployDeliveryPipelineCondition(obj *clouddeploy.DeliveryPipelineCondition) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"pipeline_ready_condition":  flattenClouddeployDeliveryPipelineConditionPipelineReadyCondition(obj.PipelineReadyCondition),
		"targets_present_condition": flattenClouddeployDeliveryPipelineConditionTargetsPresentCondition(obj.TargetsPresentCondition),
		"targets_type_condition":    flattenClouddeployDeliveryPipelineConditionTargetsTypeCondition(obj.TargetsTypeCondition),
	}

	return []interface{}{transformed}

}

func flattenClouddeployDeliveryPipelineConditionPipelineReadyCondition(obj *clouddeploy.DeliveryPipelineConditionPipelineReadyCondition) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"status":      obj.Status,
		"update_time": obj.UpdateTime,
	}

	return []interface{}{transformed}

}

func flattenClouddeployDeliveryPipelineConditionTargetsPresentCondition(obj *clouddeploy.DeliveryPipelineConditionTargetsPresentCondition) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"missing_targets": obj.MissingTargets,
		"status":          obj.Status,
		"update_time":     obj.UpdateTime,
	}

	return []interface{}{transformed}

}

func flattenClouddeployDeliveryPipelineConditionTargetsTypeCondition(obj *clouddeploy.DeliveryPipelineConditionTargetsTypeCondition) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"error_details": obj.ErrorDetails,
		"status":        obj.Status,
	}

	return []interface{}{transformed}

}

func flattenClouddeployDeliveryPipelineLabels(v map[string]string, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	if l, ok := d.Get("labels").(map[string]interface{}); ok {
		for k := range l {
			transformed[k] = v[k]
		}
	}

	return transformed
}

func flattenClouddeployDeliveryPipelineTerraformLabels(v map[string]string, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	if l, ok := d.Get("terraform_labels").(map[string]interface{}); ok {
		for k := range l {
			transformed[k] = v[k]
		}
	}

	return transformed
}

func flattenClouddeployDeliveryPipelineAnnotations(v map[string]string, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}

	transformed := make(map[string]interface{})
	if l, ok := d.Get("annotations").(map[string]interface{}); ok {
		for k := range l {
			transformed[k] = v[k]
		}
	}

	return transformed
}
