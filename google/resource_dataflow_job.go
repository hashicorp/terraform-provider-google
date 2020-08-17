package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	dataflow "google.golang.org/api/dataflow/v1b3"
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
		Update: resourceDataflowJobUpdateByReplacement,
		Delete: resourceDataflowJobDelete,
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		CustomizeDiff: customdiff.All(
			resourceDataflowJobTypeCustomizeDiff,
		),
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
				Description: `The GCS path to the Dataflow job template.`,
			},

			"temp_gcs_location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `A writeable location on GCS for the Dataflow job to dump its temporary data.`,
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
				DiffSuppressFunc: resourceDataflowJobLabelDiffSuppress,
				Description:      `User labels to be specified for the job. Keys and values should follow the restrictions specified in the labeling restrictions page. NOTE: Google-provided Dataflow templates often provide default labels that begin with goog-dataflow-provided. Unless explicitly set in config, these labels will be ignored to prevent diffs on re-apply.`,
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
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The network to which VMs will be assigned. If it is not provided, "default" will be used.`,
			},

			"subnetwork": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The subnetwork to which VMs will be assigned. Should be of the form "regions/REGION/subnetworks/SUBNETWORK".`,
			},

			"machine_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The machine type to use for the job.`,
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
		},
	}
}

func resourceDataflowJobTypeCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	// All non-virtual fields are ForceNew for batch jobs
	if d.Get("type") == "JOB_TYPE_BATCH" {
		resourceSchema := resourceDataflowJob().Schema
		for field := range resourceSchema {
			if field == "on_delete" {
				continue
			}
			// Labels map will likely have suppressed changes, so we check each key instead of the parent field
			if field == "labels" {
				resourceDataflowJobIterateMapForceNew(field, d)
			} else if d.HasChange(field) {
				d.ForceNew(field)
			}
		}
	}

	return nil
}

func resourceDataflowJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	params := expandStringMap(d, "parameters")

	env, err := resourceDataflowJobSetupEnv(d, config)
	if err != nil {
		return err
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

	d.Set("job_id", job.Id)
	d.Set("state", job.CurrentState)
	d.Set("name", job.Name)
	d.Set("type", job.Type)
	d.Set("project", project)
	d.Set("labels", job.Labels)

	sdkPipelineOptions, err := ConvertToMap(job.Environment.SdkPipelineOptions)
	if err != nil {
		return err
	}
	optionsMap := sdkPipelineOptions["options"].(map[string]interface{})
	d.Set("template_gcs_path", optionsMap["templateLocation"])
	d.Set("temp_gcs_location", optionsMap["tempLocation"])

	if _, ok := dataflowTerminalStatesMap[job.CurrentState]; ok {
		log.Printf("[DEBUG] Removing resource '%s' because it is in state %s.\n", job.Name, job.CurrentState)
		d.SetId("")
		return nil
	}
	d.SetId(job.Id)

	return nil
}

// Stream update method. Batch job changes should have been set to ForceNew via custom diff
func resourceDataflowJobUpdateByReplacement(d *schema.ResourceData, meta interface{}) error {
	// Don't send an update request if only virtual fields have changes
	if resourceDataflowJobIsVirtualUpdate(d) {
		return nil
	}

	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	params := expandStringMap(d, "parameters")

	env, err := resourceDataflowJobSetupEnv(d, config)
	if err != nil {
		return err
	}

	request := dataflow.LaunchTemplateParameters{
		JobName:     d.Get("name").(string),
		Parameters:  params,
		Environment: &env,
		Update:      true,
	}

	var response *dataflow.LaunchTemplateResponse
	err = retryTimeDuration(func() (updateErr error) {
		response, updateErr = resourceDataflowJobLaunchTemplate(config, project, region, d.Get("template_gcs_path").(string), &request)
		return updateErr
	}, time.Minute*time.Duration(5), isDataflowJobUpdateRetryableError)
	if err != nil {
		return err
	}

	if err := waitForDataflowJobToBeUpdated(d, config, response.Job.Id, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("Error updating job with job ID %q: %v", d.Id(), err)
	}

	d.SetId(response.Job.Id)

	return resourceDataflowJobRead(d, meta)
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
			gerr, isGoogleErr := updateErr.(*googleapi.Error)
			if !isGoogleErr {
				// If we have an error and it's not a google-specific error, we should go ahead and return.
				return resource.NonRetryableError(updateErr)
			}

			if strings.Contains(gerr.Message, "not yet ready for canceling") {
				// Retry cancelling job if it's not ready.
				// Sleep to avoid hitting update quota with repeated attempts.
				time.Sleep(5 * time.Second)
				return resource.RetryableError(updateErr)
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
		return config.clientDataflow.Projects.Jobs.Get(project, id).View("JOB_VIEW_ALL").Do()
	}
	return config.clientDataflow.Projects.Locations.Jobs.Get(project, region, id).View("JOB_VIEW_ALL").Do()
}

func resourceDataflowJobUpdateJob(config *Config, project string, region string, id string, job *dataflow.Job) (*dataflow.Job, error) {
	if region == "" {
		return config.clientDataflow.Projects.Jobs.Update(project, id, job).Do()
	}
	return config.clientDataflow.Projects.Locations.Jobs.Update(project, region, id, job).Do()
}

func resourceDataflowJobLaunchTemplate(config *Config, project string, region string, gcsPath string, request *dataflow.LaunchTemplateParameters) (*dataflow.LaunchTemplateResponse, error) {
	if region == "" {
		return config.clientDataflow.Projects.Templates.Launch(project, request).GcsPath(gcsPath).Do()
	}
	return config.clientDataflow.Projects.Locations.Templates.Launch(project, region, request).GcsPath(gcsPath).Do()
}

func resourceDataflowJobSetupEnv(d *schema.ResourceData, config *Config) (dataflow.RuntimeEnvironment, error) {
	zone, err := getZone(d, config)
	if err != nil {
		return dataflow.RuntimeEnvironment{}, err
	}

	labels := expandStringMap(d, "labels")

	additionalExperiments := convertStringSet(d.Get("additional_experiments").(*schema.Set))

	env := dataflow.RuntimeEnvironment{
		MaxWorkers:            int64(d.Get("max_workers").(int)),
		Network:               d.Get("network").(string),
		ServiceAccountEmail:   d.Get("service_account_email").(string),
		Subnetwork:            d.Get("subnetwork").(string),
		TempLocation:          d.Get("temp_gcs_location").(string),
		MachineType:           d.Get("machine_type").(string),
		IpConfiguration:       d.Get("ip_configuration").(string),
		AdditionalUserLabels:  labels,
		Zone:                  zone,
		AdditionalExperiments: additionalExperiments,
	}
	return env, nil
}

func resourceDataflowJobIterateMapForceNew(mapKey string, d *schema.ResourceDiff) {
	obj := d.Get(mapKey).(map[string]interface{})
	for k := range obj {
		entrySchemaKey := mapKey + "." + k
		if d.HasChange(entrySchemaKey) {
			// ForceNew must be called on the parent map to trigger
			d.ForceNew(mapKey)
			break
		}
	}
}

func resourceDataflowJobIterateMapHasChange(mapKey string, d *schema.ResourceData) bool {
	obj := d.Get(mapKey).(map[string]interface{})
	for k := range obj {
		entrySchemaKey := mapKey + "." + k
		if d.HasChange(entrySchemaKey) {
			return true
		}
	}
	return false
}

func resourceDataflowJobIsVirtualUpdate(d *schema.ResourceData) bool {
	// on_delete is the only virtual field
	if d.HasChange("on_delete") {
		// Check if other fields have changes, which would require an actual update request
		resourceSchema := resourceDataflowJob().Schema
		for field := range resourceSchema {
			if field == "on_delete" {
				continue
			}
			// Labels map will likely have suppressed changes, so we check each key instead of the parent field
			if (field == "labels" && resourceDataflowJobIterateMapHasChange(field, d)) ||
				(field != "labels" && d.HasChange(field)) {
				return false
			}
		}
		// on_delete is changing, but nothing else
		return true
	}

	return false
}

func waitForDataflowJobToBeUpdated(d *schema.ResourceData, config *Config, replacementJobID string, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		project, err := getProject(d, config)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		region, err := getRegion(d, config)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		replacementJob, err := resourceDataflowJobGetJob(config, project, region, replacementJobID)
		if err != nil {
			if isRetryableError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		state := replacementJob.CurrentState
		switch state {
		case "", "JOB_STATE_PENDING":
			return resource.RetryableError(fmt.Errorf("the replacement job with ID %q has pending state %q.", replacementJobID, state))
		case "JOB_STATE_FAILED":
			return resource.NonRetryableError(fmt.Errorf("the replacement job with ID %q failed with state %q.", replacementJobID, state))
		default:
			log.Printf("[DEBUG] the replacement job with ID %q has state %q.", replacementJobID, state)
			return nil
		}
	})
}
