// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataflow

import (
	"context"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDataflowJobResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// ForceNew applies to both stream and batch jobs
				ForceNew:    true,
				Description: `A unique name for the resource, required by Dataflow.`,
			},

			"template_gcs_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The Google Cloud Storage path to the Dataflow job template.`,
			},

			"temp_gcs_location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `A writeable location on Google Cloud Storage for the Dataflow job to dump its temporary data.`,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew applies to both stream and batch jobs
				ForceNew:    true,
				Description: `The zone in which the created job should run. If it is not provided, the provider zone is used.`,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew applies to both stream and batch jobs
				ForceNew:    true,
				Description: `The region in which the created job should run.`,
			},

			"max_workers": {
				Type:     schema.TypeInt,
				Optional: true,
				// ForceNew applies to both stream and batch jobs
				ForceNew:    true,
				Description: `The number of workers permitted to work on the job. More workers may improve processing speed at additional cost.`,
			},

			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `Key/Value pairs to be passed to the Dataflow job (as used in the template).`,
			},

			"labels": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: resourceDataflowJobLabelDiffSuppress,
				Description:      `User labels to be specified for the job. Keys and values should follow the restrictions specified in the labeling restrictions page. NOTE: Google-provided Dataflow templates often provide default labels that begin with goog-dataflow-provided. Unless explicitly set in config, these labels will be ignored to prevent diffs on re-apply.`,
			},

			"transform_name_mapping": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `Only applicable when updating a pipeline. Map of transform name prefixes of the job to be replaced with the corresponding name prefixes of the new job.`,
			},

			"on_delete": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"cancel", "drain"}, false),
				Optional:     true,
				Default:      "drain",
				Description:  `One of "drain" or "cancel". Specifies behavior of deletion during terraform destroy.`,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// ForceNew applies to both stream and batch jobs
				ForceNew:    true,
				Description: `The project in which the resource belongs.`,
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current state of the resource, selected from the JobState enum.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of this job, selected from the JobType enum.`,
			},
			"service_account_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Service Account email used to create the job.`,
			},

			"network": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The network to which VMs will be assigned. If it is not provided, "default" will be used.`,
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The subnetwork to which VMs will be assigned. Should be of the form "regions/REGION/subnetworks/SUBNETWORK".`,
			},

			"machine_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The machine type to use for the job.`,
			},

			"kms_key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name for the Cloud KMS key for the job. Key format is: projects/PROJECT_ID/locations/LOCATION/keyRings/KEY_RING/cryptoKeys/KEY`,
			},

			"ip_configuration": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  `The configuration for VM IPs. Options are "WORKER_IP_PUBLIC" or "WORKER_IP_PRIVATE".`,
				ValidateFunc: validation.StringInSlice([]string{"WORKER_IP_PUBLIC", "WORKER_IP_PRIVATE", ""}, false),
			},

			"additional_experiments": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: `List of experiments that should be used by the job. An example value is ["enable_stackdriver_agent_metrics"].`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique ID of this job.`,
			},

			"enable_streaming_engine": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Indicates if the job should use the streaming engine feature.`,
			},

			"skip_wait_on_job_termination": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `If true, treat DRAINING and CANCELLING as terminal job states and do not wait for further changes before removing from terraform state and moving on. WARNING: this will lead to job name conflicts if you do not ensure that the job names are different, e.g. by embedding a release ID or by using a random_id.`,
			},
		},
		UseJSONNumber: true,
	}
}

func ResourceDataflowJobStateUpgradeV0(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return tpgresource.LabelsStateUpgrade(rawState, resourceDataflowJobGoogleLabelPrefix)
}
