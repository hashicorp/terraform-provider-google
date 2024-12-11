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

func ResourceClouddeployTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceClouddeployTargetCreate,
		Read:   resourceClouddeployTargetRead,
		Update: resourceClouddeployTargetUpdate,
		Delete: resourceClouddeployTargetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceClouddeployTargetImport,
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
				Description: "Name of the `Target`. Format is `[a-z]([a-z0-9-]{0,61}[a-z0-9])?`.",
			},

			"anthos_cluster": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Information specifying an Anthos Cluster.",
				MaxItems:      1,
				Elem:          ClouddeployTargetAnthosClusterSchema(),
				ConflictsWith: []string{"gke", "run", "multi_target", "custom_target"},
			},

			"associated_entities": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Optional. Map of entity IDs to their associated entities. Associated entities allows specifying places other than the deployment target for specific features. For example, the Gateway API canary can be configured to deploy the HTTPRoute to a different cluster(s) than the deployment cluster using associated entities. An entity ID must consist of lower-case letters, numbers, and hyphens, start with a letter and end with a letter or a number, and have a max length of 63 characters. In other words, it must match the following regex: `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`.",
				Elem:        ClouddeployTargetAssociatedEntitiesSchema(),
				Set:         schema.HashResource(ClouddeployTargetAssociatedEntitiesSchema()),
			},

			"custom_target": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Optional. Information specifying a Custom Target.",
				MaxItems:      1,
				Elem:          ClouddeployTargetCustomTargetSchema(),
				ConflictsWith: []string{"gke", "anthos_cluster", "run", "multi_target"},
			},

			"deploy_parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. The deploy parameters to use for this target.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. Description of the `Target`. Max length is 255 characters.",
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

			"execution_configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Configurations for all execution that relates to this `Target`. Each `ExecutionEnvironmentUsage` value may only be used in a single configuration; using the same value multiple times is an error. When one or more configurations are specified, they must include the `RENDER` and `DEPLOY` `ExecutionEnvironmentUsage` values. When no configurations are specified, execution will use the default specified in `DefaultPool`.",
				Elem:        ClouddeployTargetExecutionConfigsSchema(),
			},

			"gke": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Information specifying a GKE Cluster.",
				MaxItems:      1,
				Elem:          ClouddeployTargetGkeSchema(),
				ConflictsWith: []string{"anthos_cluster", "run", "multi_target", "custom_target"},
			},

			"multi_target": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Information specifying a multiTarget.",
				MaxItems:      1,
				Elem:          ClouddeployTargetMultiTargetSchema(),
				ConflictsWith: []string{"gke", "anthos_cluster", "run", "custom_target"},
			},

			"project": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "The project for the resource",
			},

			"require_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. Whether or not the `Target` requires approval.",
			},

			"run": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Information specifying a Cloud Run deployment target.",
				MaxItems:      1,
				Elem:          ClouddeployTargetRunSchema(),
				ConflictsWith: []string{"gke", "anthos_cluster", "multi_target", "custom_target"},
			},

			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. User annotations. These attributes can only be set and used by the user, and not by Google Cloud Deploy. See https://google.aip.dev/128#annotations for more details such as format and size limitations.\n\n**Note**: This field is non-authoritative, and will only manage the annotations present in your configuration.\nPlease refer to the field `effective_annotations` for all of the annotations present on the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Time at which the `Target` was created.",
			},

			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional. This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional. Labels are attributes that can be set and used by both the user and by Google Cloud Deploy. Labels must meet the following constraints: * Keys and values can contain only lowercase letters, numeric characters, underscores, and dashes. * All characters must use UTF-8 encoding, and international characters are allowed. * Keys must start with a lowercase letter or international character. * Each resource is limited to a maximum of 64 labels. Both keys and values are additionally constrained to be <= 128 bytes.\n\n**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.\nPlease refer to the field `effective_labels` for all of the labels present on the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Resource id of the `Target`.",
			},

			"terraform_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The combination of labels configured directly on the resource and default labels configured on the provider.",
			},

			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Unique identifier of the `Target`.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Output only. Most recent time at which the `Target` was updated.",
			},
		},
	}
}

func ClouddeployTargetAnthosClusterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"membership": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Membership of the GKE Hub-registered cluster to which to apply the Skaffold configuration. Format is `projects/{project}/locations/{location}/memberships/{membership_name}`.",
			},
		},
	}
}

func ClouddeployTargetAssociatedEntitiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the key in the map for which this object is mapped to in the API",
			},

			"anthos_clusters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Information specifying Anthos clusters as associated entities.",
				Elem:        ClouddeployTargetAssociatedEntitiesAnthosClustersSchema(),
			},

			"gke_clusters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional. Information specifying GKE clusters as associated entities.",
				Elem:        ClouddeployTargetAssociatedEntitiesGkeClustersSchema(),
			},
		},
	}
}

func ClouddeployTargetAssociatedEntitiesAnthosClustersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"membership": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. Membership of the GKE Hub-registered cluster to which to apply the Skaffold configuration. Format is `projects/{project}/locations/{location}/memberships/{membership_name}`.",
			},
		},
	}
}

func ClouddeployTargetAssociatedEntitiesGkeClustersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. Information specifying a GKE Cluster. Format is `projects/{project_id}/locations/{location_id}/clusters/{cluster_id}`.",
			},

			"internal_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. If true, `cluster` is accessed using the private IP address of the control plane endpoint. Otherwise, the default IP address of the control plane endpoint is used. The default IP address is the private IP address for clusters with private control-plane endpoints and the public IP address otherwise. Only specify this option when `cluster` is a [private GKE cluster](https://cloud.google.com/kubernetes-engine/docs/concepts/private-cluster-concept).",
			},

			"proxy_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. If set, used to configure a [proxy](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#proxy) to the Kubernetes server.",
			},
		},
	}
}

func ClouddeployTargetCustomTargetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_target_type": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Required. The name of the CustomTargetType. Format must be `projects/{project}/locations/{location}/customTargetTypes/{custom_target_type}`.",
			},
		},
	}
}

func ClouddeployTargetExecutionConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"usages": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. Usages when this configuration should be applied.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"artifact_storage": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Cloud Storage location in which to store execution outputs. This can either be a bucket (\"gs://my-bucket\") or a path within a bucket (\"gs://my-bucket/my-dir\"). If unspecified, a default bucket located in the same region will be used.",
			},

			"execution_timeout": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Execution timeout for a Cloud Build Execution. This must be between 10m and 24h in seconds format. If unspecified, a default timeout of 1h is used.",
			},

			"service_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Optional. Google service account to use for execution. If unspecified, the project execution service account (-compute@developer.gserviceaccount.com) is used.",
			},

			"verbose": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. If true, additional logging will be enabled when running builds in this execution environment.",
			},

			"worker_pool": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Optional. The resource name of the `WorkerPool`, with the format `projects/{project}/locations/{location}/workerPools/{worker_pool}`. If this optional field is unspecified, the default Cloud Build pool will be used.",
			},
		},
	}
}

func ClouddeployTargetGkeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      "Information specifying a GKE Cluster. Format is `projects/{project_id}/locations/{location_id}/clusters/{cluster_id}.",
			},

			"internal_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Optional. If true, `cluster` is accessed using the private IP address of the control plane endpoint. Otherwise, the default IP address of the control plane endpoint is used. The default IP address is the private IP address for clusters with private control-plane endpoints and the public IP address otherwise. Only specify this option when `cluster` is a [private GKE cluster](https://cloud.google.com/kubernetes-engine/docs/concepts/private-cluster-concept).",
			},

			"proxy_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. If set, used to configure a [proxy](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#proxy) to the Kubernetes server.",
			},
		},
	}
}

func ClouddeployTargetMultiTargetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Required. The target_ids of this multiTarget.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ClouddeployTargetRunSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Required. The location where the Cloud Run Service should be located. Format is `projects/{project}/locations/{location}`.",
			},
		},
	}
}

func resourceClouddeployTargetCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.Target{
		Location:           dcl.String(d.Get("location").(string)),
		Name:               dcl.String(d.Get("name").(string)),
		AnthosCluster:      expandClouddeployTargetAnthosCluster(d.Get("anthos_cluster")),
		AssociatedEntities: expandClouddeployTargetAssociatedEntitiesMap(d.Get("associated_entities")),
		CustomTarget:       expandClouddeployTargetCustomTarget(d.Get("custom_target")),
		DeployParameters:   tpgresource.CheckStringMap(d.Get("deploy_parameters")),
		Description:        dcl.String(d.Get("description").(string)),
		Annotations:        tpgresource.CheckStringMap(d.Get("effective_annotations")),
		Labels:             tpgresource.CheckStringMap(d.Get("effective_labels")),
		ExecutionConfigs:   expandClouddeployTargetExecutionConfigsArray(d.Get("execution_configs")),
		Gke:                expandClouddeployTargetGke(d.Get("gke")),
		MultiTarget:        expandClouddeployTargetMultiTarget(d.Get("multi_target")),
		Project:            dcl.String(project),
		RequireApproval:    dcl.Bool(d.Get("require_approval").(bool)),
		Run:                expandClouddeployTargetRun(d.Get("run")),
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
	res, err := client.ApplyTarget(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error creating Target: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Target %q: %#v", d.Id(), res)

	return resourceClouddeployTargetRead(d, meta)
}

func resourceClouddeployTargetRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.Target{
		Location:           dcl.String(d.Get("location").(string)),
		Name:               dcl.String(d.Get("name").(string)),
		AnthosCluster:      expandClouddeployTargetAnthosCluster(d.Get("anthos_cluster")),
		AssociatedEntities: expandClouddeployTargetAssociatedEntitiesMap(d.Get("associated_entities")),
		CustomTarget:       expandClouddeployTargetCustomTarget(d.Get("custom_target")),
		DeployParameters:   tpgresource.CheckStringMap(d.Get("deploy_parameters")),
		Description:        dcl.String(d.Get("description").(string)),
		Annotations:        tpgresource.CheckStringMap(d.Get("effective_annotations")),
		Labels:             tpgresource.CheckStringMap(d.Get("effective_labels")),
		ExecutionConfigs:   expandClouddeployTargetExecutionConfigsArray(d.Get("execution_configs")),
		Gke:                expandClouddeployTargetGke(d.Get("gke")),
		MultiTarget:        expandClouddeployTargetMultiTarget(d.Get("multi_target")),
		Project:            dcl.String(project),
		RequireApproval:    dcl.Bool(d.Get("require_approval").(bool)),
		Run:                expandClouddeployTargetRun(d.Get("run")),
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
	res, err := client.GetTarget(context.Background(), obj)
	if err != nil {
		resourceName := fmt.Sprintf("ClouddeployTarget %q", d.Id())
		return tpgdclresource.HandleNotFoundDCLError(err, d, resourceName)
	}

	if err = d.Set("location", res.Location); err != nil {
		return fmt.Errorf("error setting location in state: %s", err)
	}
	if err = d.Set("name", res.Name); err != nil {
		return fmt.Errorf("error setting name in state: %s", err)
	}
	if err = d.Set("anthos_cluster", flattenClouddeployTargetAnthosCluster(res.AnthosCluster)); err != nil {
		return fmt.Errorf("error setting anthos_cluster in state: %s", err)
	}
	if err = d.Set("associated_entities", flattenClouddeployTargetAssociatedEntitiesMap(res.AssociatedEntities)); err != nil {
		return fmt.Errorf("error setting associated_entities in state: %s", err)
	}
	if err = d.Set("custom_target", flattenClouddeployTargetCustomTarget(res.CustomTarget)); err != nil {
		return fmt.Errorf("error setting custom_target in state: %s", err)
	}
	if err = d.Set("deploy_parameters", res.DeployParameters); err != nil {
		return fmt.Errorf("error setting deploy_parameters in state: %s", err)
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
	if err = d.Set("execution_configs", flattenClouddeployTargetExecutionConfigsArray(res.ExecutionConfigs)); err != nil {
		return fmt.Errorf("error setting execution_configs in state: %s", err)
	}
	if err = d.Set("gke", flattenClouddeployTargetGke(res.Gke)); err != nil {
		return fmt.Errorf("error setting gke in state: %s", err)
	}
	if err = d.Set("multi_target", flattenClouddeployTargetMultiTarget(res.MultiTarget)); err != nil {
		return fmt.Errorf("error setting multi_target in state: %s", err)
	}
	if err = d.Set("project", res.Project); err != nil {
		return fmt.Errorf("error setting project in state: %s", err)
	}
	if err = d.Set("require_approval", res.RequireApproval); err != nil {
		return fmt.Errorf("error setting require_approval in state: %s", err)
	}
	if err = d.Set("run", flattenClouddeployTargetRun(res.Run)); err != nil {
		return fmt.Errorf("error setting run in state: %s", err)
	}
	if err = d.Set("annotations", flattenClouddeployTargetAnnotations(res.Annotations, d)); err != nil {
		return fmt.Errorf("error setting annotations in state: %s", err)
	}
	if err = d.Set("create_time", res.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time in state: %s", err)
	}
	if err = d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("error setting etag in state: %s", err)
	}
	if err = d.Set("labels", flattenClouddeployTargetLabels(res.Labels, d)); err != nil {
		return fmt.Errorf("error setting labels in state: %s", err)
	}
	if err = d.Set("target_id", res.TargetId); err != nil {
		return fmt.Errorf("error setting target_id in state: %s", err)
	}
	if err = d.Set("terraform_labels", flattenClouddeployTargetTerraformLabels(res.Labels, d)); err != nil {
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
func resourceClouddeployTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.Target{
		Location:           dcl.String(d.Get("location").(string)),
		Name:               dcl.String(d.Get("name").(string)),
		AnthosCluster:      expandClouddeployTargetAnthosCluster(d.Get("anthos_cluster")),
		AssociatedEntities: expandClouddeployTargetAssociatedEntitiesMap(d.Get("associated_entities")),
		CustomTarget:       expandClouddeployTargetCustomTarget(d.Get("custom_target")),
		DeployParameters:   tpgresource.CheckStringMap(d.Get("deploy_parameters")),
		Description:        dcl.String(d.Get("description").(string)),
		Annotations:        tpgresource.CheckStringMap(d.Get("effective_annotations")),
		Labels:             tpgresource.CheckStringMap(d.Get("effective_labels")),
		ExecutionConfigs:   expandClouddeployTargetExecutionConfigsArray(d.Get("execution_configs")),
		Gke:                expandClouddeployTargetGke(d.Get("gke")),
		MultiTarget:        expandClouddeployTargetMultiTarget(d.Get("multi_target")),
		Project:            dcl.String(project),
		RequireApproval:    dcl.Bool(d.Get("require_approval").(bool)),
		Run:                expandClouddeployTargetRun(d.Get("run")),
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
	res, err := client.ApplyTarget(context.Background(), obj, directive...)

	if _, ok := err.(dcl.DiffAfterApplyError); ok {
		log.Printf("[DEBUG] Diff after apply returned from the DCL: %s", err)
	} else if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error updating Target: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Target %q: %#v", d.Id(), res)

	return resourceClouddeployTargetRead(d, meta)
}

func resourceClouddeployTargetDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	obj := &clouddeploy.Target{
		Location:           dcl.String(d.Get("location").(string)),
		Name:               dcl.String(d.Get("name").(string)),
		AnthosCluster:      expandClouddeployTargetAnthosCluster(d.Get("anthos_cluster")),
		AssociatedEntities: expandClouddeployTargetAssociatedEntitiesMap(d.Get("associated_entities")),
		CustomTarget:       expandClouddeployTargetCustomTarget(d.Get("custom_target")),
		DeployParameters:   tpgresource.CheckStringMap(d.Get("deploy_parameters")),
		Description:        dcl.String(d.Get("description").(string)),
		Annotations:        tpgresource.CheckStringMap(d.Get("effective_annotations")),
		Labels:             tpgresource.CheckStringMap(d.Get("effective_labels")),
		ExecutionConfigs:   expandClouddeployTargetExecutionConfigsArray(d.Get("execution_configs")),
		Gke:                expandClouddeployTargetGke(d.Get("gke")),
		MultiTarget:        expandClouddeployTargetMultiTarget(d.Get("multi_target")),
		Project:            dcl.String(project),
		RequireApproval:    dcl.Bool(d.Get("require_approval").(bool)),
		Run:                expandClouddeployTargetRun(d.Get("run")),
	}

	log.Printf("[DEBUG] Deleting Target %q", d.Id())
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
	if err := client.DeleteTarget(context.Background(), obj); err != nil {
		return fmt.Errorf("Error deleting Target: %s", err)
	}

	log.Printf("[DEBUG] Finished deleting Target %q", d.Id())
	return nil
}

func resourceClouddeployTargetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/targets/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/targets/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandClouddeployTargetAnthosCluster(o interface{}) *clouddeploy.TargetAnthosCluster {
	if o == nil {
		return clouddeploy.EmptyTargetAnthosCluster
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyTargetAnthosCluster
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.TargetAnthosCluster{
		Membership: dcl.String(obj["membership"].(string)),
	}
}

func flattenClouddeployTargetAnthosCluster(obj *clouddeploy.TargetAnthosCluster) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"membership": obj.Membership,
	}

	return []interface{}{transformed}

}

func expandClouddeployTargetAssociatedEntitiesMap(o interface{}) map[string]clouddeploy.TargetAssociatedEntities {
	if o == nil {
		return make(map[string]clouddeploy.TargetAssociatedEntities)
	}

	o = o.(*schema.Set).List()

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make(map[string]clouddeploy.TargetAssociatedEntities)
	}

	items := make(map[string]clouddeploy.TargetAssociatedEntities)
	for _, item := range objs {
		i := expandClouddeployTargetAssociatedEntities(item)
		if item != nil {
			items[item.(map[string]interface{})["entity_id"].(string)] = *i
		}
	}

	return items
}

func expandClouddeployTargetAssociatedEntities(o interface{}) *clouddeploy.TargetAssociatedEntities {
	if o == nil {
		return clouddeploy.EmptyTargetAssociatedEntities
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.TargetAssociatedEntities{
		AnthosClusters: expandClouddeployTargetAssociatedEntitiesAnthosClustersArray(obj["anthos_clusters"]),
		GkeClusters:    expandClouddeployTargetAssociatedEntitiesGkeClustersArray(obj["gke_clusters"]),
	}
}

func flattenClouddeployTargetAssociatedEntitiesMap(objs map[string]clouddeploy.TargetAssociatedEntities) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for name, item := range objs {
		i := flattenClouddeployTargetAssociatedEntities(&item, name)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployTargetAssociatedEntities(obj *clouddeploy.TargetAssociatedEntities, name string) interface{} {
	if obj == nil {
		return nil
	}
	transformed := map[string]interface{}{
		"anthos_clusters": flattenClouddeployTargetAssociatedEntitiesAnthosClustersArray(obj.AnthosClusters),
		"gke_clusters":    flattenClouddeployTargetAssociatedEntitiesGkeClustersArray(obj.GkeClusters),
	}

	transformed["entity_id"] = name

	return transformed

}
func expandClouddeployTargetAssociatedEntitiesAnthosClustersArray(o interface{}) []clouddeploy.TargetAssociatedEntitiesAnthosClusters {
	if o == nil {
		return make([]clouddeploy.TargetAssociatedEntitiesAnthosClusters, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]clouddeploy.TargetAssociatedEntitiesAnthosClusters, 0)
	}

	items := make([]clouddeploy.TargetAssociatedEntitiesAnthosClusters, 0, len(objs))
	for _, item := range objs {
		i := expandClouddeployTargetAssociatedEntitiesAnthosClusters(item)
		items = append(items, *i)
	}

	return items
}

func expandClouddeployTargetAssociatedEntitiesAnthosClusters(o interface{}) *clouddeploy.TargetAssociatedEntitiesAnthosClusters {
	if o == nil {
		return clouddeploy.EmptyTargetAssociatedEntitiesAnthosClusters
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.TargetAssociatedEntitiesAnthosClusters{
		Membership: dcl.String(obj["membership"].(string)),
	}
}

func flattenClouddeployTargetAssociatedEntitiesAnthosClustersArray(objs []clouddeploy.TargetAssociatedEntitiesAnthosClusters) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenClouddeployTargetAssociatedEntitiesAnthosClusters(&item)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployTargetAssociatedEntitiesAnthosClusters(obj *clouddeploy.TargetAssociatedEntitiesAnthosClusters) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"membership": obj.Membership,
	}

	return transformed

}
func expandClouddeployTargetAssociatedEntitiesGkeClustersArray(o interface{}) []clouddeploy.TargetAssociatedEntitiesGkeClusters {
	if o == nil {
		return make([]clouddeploy.TargetAssociatedEntitiesGkeClusters, 0)
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return make([]clouddeploy.TargetAssociatedEntitiesGkeClusters, 0)
	}

	items := make([]clouddeploy.TargetAssociatedEntitiesGkeClusters, 0, len(objs))
	for _, item := range objs {
		i := expandClouddeployTargetAssociatedEntitiesGkeClusters(item)
		items = append(items, *i)
	}

	return items
}

func expandClouddeployTargetAssociatedEntitiesGkeClusters(o interface{}) *clouddeploy.TargetAssociatedEntitiesGkeClusters {
	if o == nil {
		return clouddeploy.EmptyTargetAssociatedEntitiesGkeClusters
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.TargetAssociatedEntitiesGkeClusters{
		Cluster:    dcl.String(obj["cluster"].(string)),
		InternalIP: dcl.Bool(obj["internal_ip"].(bool)),
		ProxyUrl:   dcl.String(obj["proxy_url"].(string)),
	}
}

func flattenClouddeployTargetAssociatedEntitiesGkeClustersArray(objs []clouddeploy.TargetAssociatedEntitiesGkeClusters) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenClouddeployTargetAssociatedEntitiesGkeClusters(&item)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployTargetAssociatedEntitiesGkeClusters(obj *clouddeploy.TargetAssociatedEntitiesGkeClusters) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cluster":     obj.Cluster,
		"internal_ip": obj.InternalIP,
		"proxy_url":   obj.ProxyUrl,
	}

	return transformed

}

func expandClouddeployTargetCustomTarget(o interface{}) *clouddeploy.TargetCustomTarget {
	if o == nil {
		return clouddeploy.EmptyTargetCustomTarget
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyTargetCustomTarget
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.TargetCustomTarget{
		CustomTargetType: dcl.String(obj["custom_target_type"].(string)),
	}
}

func flattenClouddeployTargetCustomTarget(obj *clouddeploy.TargetCustomTarget) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"custom_target_type": obj.CustomTargetType,
	}

	return []interface{}{transformed}

}
func expandClouddeployTargetExecutionConfigsArray(o interface{}) []clouddeploy.TargetExecutionConfigs {
	if o == nil {
		return nil
	}

	objs := o.([]interface{})
	if len(objs) == 0 || objs[0] == nil {
		return nil
	}

	items := make([]clouddeploy.TargetExecutionConfigs, 0, len(objs))
	for _, item := range objs {
		i := expandClouddeployTargetExecutionConfigs(item)
		items = append(items, *i)
	}

	return items
}

func expandClouddeployTargetExecutionConfigs(o interface{}) *clouddeploy.TargetExecutionConfigs {
	if o == nil {
		return nil
	}

	obj := o.(map[string]interface{})
	return &clouddeploy.TargetExecutionConfigs{
		Usages:           expandClouddeployTargetExecutionConfigsUsagesArray(obj["usages"]),
		ArtifactStorage:  dcl.StringOrNil(obj["artifact_storage"].(string)),
		ExecutionTimeout: dcl.StringOrNil(obj["execution_timeout"].(string)),
		ServiceAccount:   dcl.StringOrNil(obj["service_account"].(string)),
		Verbose:          dcl.Bool(obj["verbose"].(bool)),
		WorkerPool:       dcl.String(obj["worker_pool"].(string)),
	}
}

func flattenClouddeployTargetExecutionConfigsArray(objs []clouddeploy.TargetExecutionConfigs) []interface{} {
	if objs == nil {
		return nil
	}

	items := []interface{}{}
	for _, item := range objs {
		i := flattenClouddeployTargetExecutionConfigs(&item)
		items = append(items, i)
	}

	return items
}

func flattenClouddeployTargetExecutionConfigs(obj *clouddeploy.TargetExecutionConfigs) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"usages":            flattenClouddeployTargetExecutionConfigsUsagesArray(obj.Usages),
		"artifact_storage":  obj.ArtifactStorage,
		"execution_timeout": obj.ExecutionTimeout,
		"service_account":   obj.ServiceAccount,
		"verbose":           obj.Verbose,
		"worker_pool":       obj.WorkerPool,
	}

	return transformed

}

func expandClouddeployTargetGke(o interface{}) *clouddeploy.TargetGke {
	if o == nil {
		return clouddeploy.EmptyTargetGke
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyTargetGke
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.TargetGke{
		Cluster:    dcl.String(obj["cluster"].(string)),
		InternalIP: dcl.Bool(obj["internal_ip"].(bool)),
		ProxyUrl:   dcl.String(obj["proxy_url"].(string)),
	}
}

func flattenClouddeployTargetGke(obj *clouddeploy.TargetGke) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"cluster":     obj.Cluster,
		"internal_ip": obj.InternalIP,
		"proxy_url":   obj.ProxyUrl,
	}

	return []interface{}{transformed}

}

func expandClouddeployTargetMultiTarget(o interface{}) *clouddeploy.TargetMultiTarget {
	if o == nil {
		return clouddeploy.EmptyTargetMultiTarget
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyTargetMultiTarget
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.TargetMultiTarget{
		TargetIds: tpgdclresource.ExpandStringArray(obj["target_ids"]),
	}
}

func flattenClouddeployTargetMultiTarget(obj *clouddeploy.TargetMultiTarget) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"target_ids": obj.TargetIds,
	}

	return []interface{}{transformed}

}

func expandClouddeployTargetRun(o interface{}) *clouddeploy.TargetRun {
	if o == nil {
		return clouddeploy.EmptyTargetRun
	}
	objArr := o.([]interface{})
	if len(objArr) == 0 || objArr[0] == nil {
		return clouddeploy.EmptyTargetRun
	}
	obj := objArr[0].(map[string]interface{})
	return &clouddeploy.TargetRun{
		Location: dcl.String(obj["location"].(string)),
	}
}

func flattenClouddeployTargetRun(obj *clouddeploy.TargetRun) interface{} {
	if obj == nil || obj.Empty() {
		return nil
	}
	transformed := map[string]interface{}{
		"location": obj.Location,
	}

	return []interface{}{transformed}

}

func flattenClouddeployTargetLabels(v map[string]string, d *schema.ResourceData) interface{} {
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

func flattenClouddeployTargetTerraformLabels(v map[string]string, d *schema.ResourceData) interface{} {
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

func flattenClouddeployTargetAnnotations(v map[string]string, d *schema.ResourceData) interface{} {
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

func flattenClouddeployTargetExecutionConfigsUsagesArray(obj []clouddeploy.TargetExecutionConfigsUsagesEnum) interface{} {
	if obj == nil {
		return nil
	}
	items := []string{}
	for _, item := range obj {
		items = append(items, string(item))
	}
	return items
}
func expandClouddeployTargetExecutionConfigsUsagesArray(o interface{}) []clouddeploy.TargetExecutionConfigsUsagesEnum {
	objs := o.([]interface{})
	items := make([]clouddeploy.TargetExecutionConfigsUsagesEnum, 0, len(objs))
	for _, item := range objs {
		i := clouddeploy.TargetExecutionConfigsUsagesEnumRef(item.(string))
		items = append(items, *i)
	}
	return items
}
