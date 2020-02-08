package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/googleapi"
)

const resourceDataflowJobGoogleProvidedLabelPrefix = "labels.goog-dataflow-provided"

var dataflowTerminalStatesMap = map[string]struct{}{
	"JOB_STATE_DONE":      {},
	"JOB_STATE_FAILED":    {},
	"JOB_STATE_CANCELLED": {},
	"JOB_STATE_UPDATED":   {},
	"JOB_STATE_DRAINED":   {},
}

func resourceDataflowJobLabelDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Example Diff: "labels.goog-dataflow-provided-template-version": "word_count" => ""
	if strings.HasPrefix(k, resourceDataflowJobGoogleProvidedLabelPrefix) && new == "" {
		// Suppress diff if field is a Google Dataflow-provided label key and has no explicitly set value in Config.
		return true
	}

	// Let diff be determined by labels (above)
	if strings.HasPrefix(k, "labels.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}

func resourceDataflowJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataflowJobCreate,
		Read:   resourceDataflowJobRead,
		Delete: resourceDataflowJobDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"template_gcs_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"temp_gcs_location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"max_workers": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"labels": {
				Type:             schema.TypeMap,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: resourceDataflowJobLabelDiffSuppress,
			},

			"on_delete": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"cancel", "drain"}, false),
				Optional:     true,
				Default:      "drain",
				ForceNew:     true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_account_email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"network": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"machine_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ip_configuration": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"WORKER_IP_PUBLIC", "WORKER_IP_PRIVATE", ""}, false),
			},
		},
	}
}

func resourceDataflowJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	params := expandStringMap(d, "parameters")
	labels := expandStringMap(d, "labels")

	env := dataflow.RuntimeEnvironment{
		MaxWorkers:           int64(d.Get("max_workers").(int)),
		Network:              d.Get("network").(string),
		ServiceAccountEmail:  d.Get("service_account_email").(string),
		Subnetwork:           d.Get("subnetwork").(string),
		TempLocation:         d.Get("temp_gcs_location").(string),
		MachineType:          d.Get("machine_type").(string),
		IpConfiguration:      d.Get("ip_configuration").(string),
		AdditionalUserLabels: labels,
		Zone:                 zone,
	}

	request := dataflow.CreateJobFromTemplateRequest{
		JobName:     d.Get("name").(string),
		GcsPath:     d.Get("template_gcs_path").(string),
		Parameters:  params,
		Environment: &env,
	}

	job, err := resourceDataflowJobCreateJob(config, project, region, &request)
	if err != nil {
		return err
	}
	d.SetId(job.Id)

	return resourceDataflowJobRead(d, meta)
}

func resourceDataflowJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	id := d.Id()

	job, err := resourceDataflowJobGetJob(config, project, region, id)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataflow job %s", id))
	}

	d.Set("state", job.CurrentState)
	d.Set("name", job.Name)
	d.Set("project", project)
	d.Set("labels", job.Labels)

	if _, ok := dataflowTerminalStatesMap[job.CurrentState]; ok {
		log.Printf("[DEBUG] Removing resource '%s' because it is in state %s.\n", job.Name, job.CurrentState)
		d.SetId("")
		return nil
	}
	d.SetId(job.Id)

	return nil
}

func resourceDataflowJobDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	id := d.Id()

	requestedState, err := resourceDataflowJobMapRequestedState(d.Get("on_delete").(string))
	if err != nil {
		return err
	}

	// Retry updating the state while the job is not ready to be canceled/drained.
	err = resource.Retry(time.Minute*time.Duration(15), func() *resource.RetryError {
		// To terminate a dataflow job, we update the job with a requested
		// terminal state.
		job := &dataflow.Job{
			RequestedState: requestedState,
		}

		_, updateErr := resourceDataflowJobUpdateJob(config, project, region, id, job)
		if updateErr != nil {
			gerr, isGoogleErr := err.(*googleapi.Error)
			if !isGoogleErr {
				// If we have an error and it's not a google-specific error, we should go ahead and return.
				return resource.NonRetryableError(err)
			}

			if strings.Contains(gerr.Message, "not yet ready for canceling") {
				// Retry cancelling job if it's not ready.
				// Sleep to avoid hitting update quota with repeated attempts.
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}

			if strings.Contains(gerr.Message, "Job has terminated") {
				// Job has already been terminated, skip.
				return nil
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Wait for state to reach terminal state (canceled/drained/done)
	_, ok := dataflowTerminalStatesMap[d.Get("state").(string)]
	for !ok {
		log.Printf("[DEBUG] Waiting for job with job state %q to terminate...", d.Get("state").(string))
		time.Sleep(5 * time.Second)

		err = resourceDataflowJobRead(d, meta)
		if err != nil {
			return fmt.Errorf("Error while reading job to see if it was properly terminated: %v", err)
		}
		_, ok = dataflowTerminalStatesMap[d.Get("state").(string)]
	}

	// Only remove the job from state if it's actually successfully canceled.
	if _, ok := dataflowTerminalStatesMap[d.Get("state").(string)]; ok {
		log.Printf("[DEBUG] Removing dataflow job with final state %q", d.Get("state").(string))
		d.SetId("")
		return nil
	}
	return fmt.Errorf("Unable to cancel the dataflow job '%s' - final state was %q.", d.Id(), d.Get("state").(string))
}

func resourceDataflowJobMapRequestedState(policy string) (string, error) {
	switch policy {
	case "cancel":
		return "JOB_STATE_CANCELLED", nil
	case "drain":
		return "JOB_STATE_DRAINING", nil
	default:
		return "", fmt.Errorf("Invalid `on_delete` policy: %s", policy)
	}
}

func resourceDataflowJobCreateJob(config *Config, project string, region string, request *dataflow.CreateJobFromTemplateRequest) (*dataflow.Job, error) {
	if region == "" {
		return config.clientDataflow.Projects.Templates.Create(project, request).Do()
	}
	return config.clientDataflow.Projects.Locations.Templates.Create(project, region, request).Do()
}

func resourceDataflowJobGetJob(config *Config, project string, region string, id string) (*dataflow.Job, error) {
	if region == "" {
		return config.clientDataflow.Projects.Jobs.Get(project, id).Do()
	}
	return config.clientDataflow.Projects.Locations.Jobs.Get(project, region, id).Do()
}

func resourceDataflowJobUpdateJob(config *Config, project string, region string, id string, job *dataflow.Job) (*dataflow.Job, error) {
	if region == "" {
		return config.clientDataflow.Projects.Jobs.Update(project, id, job).Do()
	}
	return config.clientDataflow.Projects.Locations.Jobs.Update(project, region, id, job).Do()
}
