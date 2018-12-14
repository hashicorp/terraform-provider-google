package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/googleapi"
)

var dataflowTerminalStatesMap = map[string]struct{}{
	"JOB_STATE_DONE":       {},
	"JOB_STATE_FAILED":     {},
	"JOB_STATE_CANCELLED":  {},
	"JOB_STATE_UPDATED":    {},
	"JOB_STATE_DRAINING":   {},
	"JOB_STATE_DRAINED":    {},
	"JOB_STATE_CANCELLING": {},
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
				ForceNew: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
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

	env := dataflow.RuntimeEnvironment{
		TempLocation: d.Get("temp_gcs_location").(string),
		Zone:         zone,
		MaxWorkers:   int64(d.Get("max_workers").(int)),
	}

	request := dataflow.CreateJobFromTemplateRequest{
		JobName:     d.Get("name").(string),
		GcsPath:     d.Get("template_gcs_path").(string),
		Parameters:  params,
		Environment: &env,
	}

	job, err := createJob(config, project, region, &request)
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

	job, err := getJob(config, project, region, id)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataflow job %s", id))
	}

	d.Set("state", job.CurrentState)
	d.Set("name", job.Name)
	d.Set("project", project)

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
	requestedState, err := mapOnDelete(d.Get("on_delete").(string))
	if err != nil {
		return err
	}
	for _, ok := dataflowTerminalStatesMap[d.Get("state").(string)]; !ok; _, ok = dataflowTerminalStatesMap[d.Get("state").(string)] {
		job := &dataflow.Job{
			RequestedState: requestedState,
		}

		_, err = updateJob(config, project, region, id, job)
		if err != nil {
			if gerr, err_ok := err.(*googleapi.Error); !err_ok {
				// If we have an error and it's not a google-specific error, we should go ahead and return.
				return err
			} else if err_ok && strings.Contains(gerr.Message, "not yet ready for canceling") {
				// We'll sleep below to wait for the job to be ready to cancel.
			} else {
				return err
			}
		}

		err = resourceDataflowJobRead(d, meta)
		postReadState := d.Get("state").(string)
		log.Printf("[DEBUG] Job state: '%s'.", postReadState)
		if _, ok := dataflowTerminalStatesMap[postReadState]; !ok {
			// If we're not yet in a terminal state, we need to sleep a few seconds so we don't
			// exhaust our update quota with repeated attempts.
			time.Sleep(5 * time.Second)
		}
		if err != nil {
			return err
		}
	}

	// Only remove the job from state if it's actually successfully canceled.
	if _, ok := dataflowTerminalStatesMap[d.Get("state").(string)]; ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("There was a problem canceling the dataflow job '%s' - the final state was %s.", d.Id(), d.Get("state").(string))

}

func mapOnDelete(policy string) (string, error) {
	switch policy {
	case "cancel":
		return "JOB_STATE_CANCELLED", nil
	case "drain":
		return "JOB_STATE_DRAINING", nil
	default:
		return "", fmt.Errorf("Invalid `on_delete` policy: %s", policy)
	}
}

func createJob(config *Config, project string, region string, request *dataflow.CreateJobFromTemplateRequest) (*dataflow.Job, error) {
	if region == "" {
		return config.clientDataflow.Projects.Templates.Create(project, request).Do()
	}
	return config.clientDataflow.Projects.Locations.Templates.Create(project, region, request).Do()
}

func getJob(config *Config, project string, region string, id string) (*dataflow.Job, error) {
	if region == "" {
		return config.clientDataflow.Projects.Jobs.Get(project, id).Do()
	}
	return config.clientDataflow.Projects.Locations.Jobs.Get(project, region, id).Do()
}

func updateJob(config *Config, project string, region string, id string, job *dataflow.Job) (*dataflow.Job, error) {
	if region == "" {
		return config.clientDataflow.Projects.Jobs.Update(project, id, job).Do()
	}
	return config.clientDataflow.Projects.Locations.Jobs.Update(project, region, id, job).Do()
}
