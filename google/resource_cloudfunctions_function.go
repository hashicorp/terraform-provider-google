package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/cloudfunctions/v1"

	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Min is 1 second, max is 9 minutes 540 sec
const functionTimeOutMax = 540
const functionTimeOutMin = 1
const functionDefaultTimeout = 60

var functionAllowedMemory = map[int]bool{
	128:  true,
	256:  true,
	512:  true,
	1024: true,
	2048: true,
}

// For now CloudFunctions are allowed only in the following locations.
// Please see https://cloud.google.com/about/locations/
var validCloudFunctionRegion = validation.StringInSlice([]string{"us-central1", "us-east1", "europe-west1", "asia-northeast1"}, true)

const functionDefaultAllowedMemoryMb = 256

type cloudFunctionId struct {
	Project string
	Region  string
	Name    string
}

func (s *cloudFunctionId) cloudFunctionId() string {
	return fmt.Sprintf("projects/%s/locations/%s/functions/%s", s.Project, s.Region, s.Name)
}

func (s *cloudFunctionId) locationId() string {
	return fmt.Sprintf("projects/%s/locations/%s", s.Project, s.Region)
}

func (s *cloudFunctionId) terraformId() string {
	return fmt.Sprintf("%s/%s/%s", s.Project, s.Region, s.Name)
}

func parseCloudFunctionId(id string, config *Config) (*cloudFunctionId, error) {
	parts := strings.Split(id, "/")

	cloudFuncIdRegex := regexp.MustCompile("^([a-z0-9-]+)/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")

	if cloudFuncIdRegex.MatchString(id) {
		return &cloudFunctionId{
			Project: parts[0],
			Region:  parts[1],
			Name:    parts[2],
		}, nil
	}

	return nil, fmt.Errorf("Invalid CloudFunction id format, expecting " +
		"`{projectId}/{regionId}/{cloudFunctionName}`")
}

func joinMapKeys(mapToJoin *map[int]bool) string {
	var keys []string
	for key := range *mapToJoin {
		keys = append(keys, strconv.Itoa(key))
	}
	return strings.Join(keys, ",")
}

func resourceCloudFunctionsFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFunctionsCreate,
		Read:   resourceCloudFunctionsRead,
		Update: resourceCloudFunctionsUpdate,
		Delete: resourceCloudFunctionsDestroy,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if len(value) > 48 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 48 characters", k))
					}
					if !regexp.MustCompile("^[a-zA-Z0-9-]+$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q can only contain letters, numbers and hyphens", k))
					}
					if !regexp.MustCompile("^[a-zA-Z]").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must start with a letter", k))
					}
					if !regexp.MustCompile("[a-zA-Z0-9]$").MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"%q must end with a number or a letter", k))
					}
					return
				},
			},

			"source_archive_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},

			"source_archive_object": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"available_memory_mb": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  functionDefaultAllowedMemoryMb,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					availableMemoryMB := v.(int)

					if functionAllowedMemory[availableMemoryMB] != true {
						errors = append(errors, fmt.Errorf("Allowed values for memory (in MB) are: %s . Got %d",
							joinMapKeys(&functionAllowedMemory), availableMemoryMB))
					}
					return
				},
			},

			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      functionDefaultTimeout,
				ValidateFunc: validation.IntBetween(functionTimeOutMin, functionTimeOutMax),
			},

			"entry_point": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"environment_variables": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"trigger_bucket": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Removed:       "This field is removed. Use `event_trigger` instead.",
				ConflictsWith: []string{"trigger_http", "trigger_topic"},
			},

			"trigger_http": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"trigger_bucket", "trigger_topic"},
			},

			"trigger_topic": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Removed:       "This field is removed. Use `event_trigger` instead.",
				ConflictsWith: []string{"trigger_http", "trigger_bucket"},
			},

			"event_trigger": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"trigger_http"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Required: true,
						},
						"failure_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retry": {
										Type: schema.TypeBool,
										// not strictly required, but this way an empty block can't be specified
										Required: true,
									},
								}},
						},
					},
				},
			},

			"https_trigger_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"retry_on_failure": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Removed:  "This field is removed. Use `event_trigger.failure_policy.retry` instead.",
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validCloudFunctionRegion,
			},
		},
	}
}

func resourceCloudFunctionsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}
	// We do this extra validation here since most regions are not valid, and the
	// error message that Cloud Functions has for "wrong region" is not specific.
	// Provider-level region fetching skips validation, because it's not possible
	// for the provider-level region to know about the field-level validator.
	_, errs := validCloudFunctionRegion(region, "region")
	if len(errs) > 0 {
		return errs[0]
	}

	cloudFuncId := &cloudFunctionId{
		Project: project,
		Region:  region,
		Name:    d.Get("name").(string),
	}

	function := &cloudfunctions.CloudFunction{
		Name:            cloudFuncId.cloudFunctionId(),
		Runtime:         d.Get("runtime").(string),
		ForceSendFields: []string{},
	}

	sourceArchiveBucket := d.Get("source_archive_bucket").(string)
	sourceArchiveObj := d.Get("source_archive_object").(string)
	function.SourceArchiveUrl = fmt.Sprintf("gs://%v/%v", sourceArchiveBucket, sourceArchiveObj)

	if v, ok := d.GetOk("available_memory_mb"); ok {
		availableMemoryMb := v.(int)
		function.AvailableMemoryMb = int64(availableMemoryMb)
	}

	if v, ok := d.GetOk("description"); ok {
		function.Description = v.(string)
	}

	if v, ok := d.GetOk("entry_point"); ok {
		function.EntryPoint = v.(string)
	}

	if v, ok := d.GetOk("timeout"); ok {
		function.Timeout = fmt.Sprintf("%vs", v.(int))
	}

	if v, ok := d.GetOk("event_trigger"); ok {
		function.EventTrigger = expandEventTrigger(v.([]interface{}), project)
	} else if v, ok := d.GetOk("trigger_http"); ok && v.(bool) {
		function.HttpsTrigger = &cloudfunctions.HttpsTrigger{}
	} else {
		return fmt.Errorf("One of `event_trigger` or `trigger_http` is required: " +
			"You must specify a trigger when deploying a new function.")
	}

	if _, ok := d.GetOk("labels"); ok {
		function.Labels = expandLabels(d)
	}

	if _, ok := d.GetOk("environment_variables"); ok {
		function.EnvironmentVariables = expandEnvironmentVariables(d)
	}

	log.Printf("[DEBUG] Creating cloud function: %s", function.Name)
	op, err := config.clientCloudFunctions.Projects.Locations.Functions.Create(
		cloudFuncId.locationId(), function).Do()
	if err != nil {
		return err
	}

	// Name of function should be unique
	d.SetId(cloudFuncId.terraformId())

	err = cloudFunctionsOperationWait(config.clientCloudFunctions, op, "Creating CloudFunctions Function")
	if err != nil {
		return err
	}

	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cloudFuncId, err := parseCloudFunctionId(d.Id(), config)
	if err != nil {
		return err
	}

	function, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(cloudFuncId.cloudFunctionId()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Target CloudFunctions Function %q", cloudFuncId.Name))
	}

	d.Set("name", cloudFuncId.Name)
	d.Set("description", function.Description)
	d.Set("entry_point", function.EntryPoint)
	d.Set("available_memory_mb", function.AvailableMemoryMb)
	sRemoved := strings.Replace(function.Timeout, "s", "", -1)
	timeout, err := strconv.Atoi(sRemoved)
	if err != nil {
		return err
	}
	d.Set("timeout", timeout)
	d.Set("labels", function.Labels)
	d.Set("runtime", function.Runtime)
	d.Set("environment_variables", function.EnvironmentVariables)
	if function.SourceArchiveUrl != "" {
		// sourceArchiveUrl should always be a Google Cloud Storage URL (e.g. gs://bucket/object)
		// https://cloud.google.com/functions/docs/reference/rest/v1/projects.locations.functions
		sourceURL, err := url.Parse(function.SourceArchiveUrl)
		if err != nil {
			return err
		}
		bucket := sourceURL.Host
		object := strings.TrimLeft(sourceURL.Path, "/")
		d.Set("source_archive_bucket", bucket)
		d.Set("source_archive_object", object)
	}

	if function.HttpsTrigger != nil {
		d.Set("trigger_http", true)
		d.Set("https_trigger_url", function.HttpsTrigger.Url)
	}

	d.Set("event_trigger", flattenEventTrigger(function.EventTrigger))

	d.Set("region", cloudFuncId.Region)
	d.Set("project", cloudFuncId.Project)

	return nil
}

func resourceCloudFunctionsUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Updating google_cloudfunctions_function")
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	cloudFuncId, err := parseCloudFunctionId(d.Id(), config)
	if err != nil {
		return err
	}

	d.Partial(true)

	function := cloudfunctions.CloudFunction{
		Name: cloudFuncId.cloudFunctionId(),
	}

	var updateMaskArr []string
	if d.HasChange("available_memory_mb") {
		availableMemoryMb := d.Get("available_memory_mb").(int)
		function.AvailableMemoryMb = int64(availableMemoryMb)
		updateMaskArr = append(updateMaskArr, "availableMemoryMb")
	}

	if d.HasChange("source_archive_bucket") || d.HasChange("source_archive_object") {
		sourceArchiveBucket := d.Get("source_archive_bucket").(string)
		sourceArchiveObj := d.Get("source_archive_object").(string)
		function.SourceArchiveUrl = fmt.Sprintf("gs://%v/%v", sourceArchiveBucket, sourceArchiveObj)
		updateMaskArr = append(updateMaskArr, "sourceArchiveUrl")
	}

	if d.HasChange("description") {
		function.Description = d.Get("description").(string)
		updateMaskArr = append(updateMaskArr, "description")
	}

	if d.HasChange("timeout") {
		function.Timeout = fmt.Sprintf("%vs", d.Get("timeout").(int))
		updateMaskArr = append(updateMaskArr, "timeout")
	}

	if d.HasChange("labels") {
		function.Labels = expandLabels(d)
		updateMaskArr = append(updateMaskArr, "labels")
	}

	if d.HasChange("runtime") {
		function.Runtime = d.Get("runtime").(string)
		updateMaskArr = append(updateMaskArr, "runtime")
	}

	if d.HasChange("environment_variables") {
		function.EnvironmentVariables = expandEnvironmentVariables(d)
		updateMaskArr = append(updateMaskArr, "environment_variables")
	}

	if d.HasChange("event_trigger") {
		function.EventTrigger = expandEventTrigger(d.Get("event_trigger").([]interface{}), project)
		updateMaskArr = append(updateMaskArr, "eventTrigger", "eventTrigger.failurePolicy.retry")
	}

	if len(updateMaskArr) > 0 {
		log.Printf("[DEBUG] Send Patch CloudFunction Configuration request: %#v", function)
		updateMask := strings.Join(updateMaskArr, ",")
		op, err := config.clientCloudFunctions.Projects.Locations.Functions.Patch(function.Name, &function).
			UpdateMask(updateMask).Do()

		if err != nil {
			return fmt.Errorf("Error while updating cloudfunction configuration: %s", err)
		}

		err = cloudFunctionsOperationWait(config.clientCloudFunctions, op,
			"Updating CloudFunctions Function")
		if err != nil {
			return err
		}
	}
	d.Partial(false)

	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cloudFuncId, err := parseCloudFunctionId(d.Id(), config)
	if err != nil {
		return err
	}

	op, err := config.clientCloudFunctions.Projects.Locations.Functions.Delete(cloudFuncId.cloudFunctionId()).Do()
	if err != nil {
		return err
	}
	err = cloudFunctionsOperationWait(config.clientCloudFunctions, op, "Deleting CloudFunctions Function")
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandEventTrigger(configured []interface{}, project string) *cloudfunctions.EventTrigger {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	eventType := data["event_type"].(string)
	shape := ""
	switch {
	case strings.HasPrefix(eventType, "providers/cloud.storage/eventTypes/"):
		shape = "projects/%s/buckets/%s"
	case strings.HasPrefix(eventType, "providers/cloud.pubsub/eventTypes/"):
		shape = "projects/%s/topics/%s"
	}

	return &cloudfunctions.EventTrigger{
		EventType:     eventType,
		Resource:      fmt.Sprintf(shape, project, data["resource"].(string)),
		FailurePolicy: expandFailurePolicy(data["failure_policy"].([]interface{})),
	}
}

func flattenEventTrigger(eventTrigger *cloudfunctions.EventTrigger) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if eventTrigger == nil {
		return result
	}

	result = append(result, map[string]interface{}{
		"event_type":     eventTrigger.EventType,
		"resource":       GetResourceNameFromSelfLink(eventTrigger.Resource),
		"failure_policy": flattenFailurePolicy(eventTrigger.FailurePolicy),
	})

	return result
}

func expandFailurePolicy(configured []interface{}) *cloudfunctions.FailurePolicy {
	if len(configured) == 0 || configured[0] == nil {
		return &cloudfunctions.FailurePolicy{}
	}

	if data := configured[0].(map[string]interface{}); data["retry"].(bool) {
		return &cloudfunctions.FailurePolicy{
			Retry: &cloudfunctions.Retry{},
		}
	}

	return nil
}

func flattenFailurePolicy(failurePolicy *cloudfunctions.FailurePolicy) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if failurePolicy == nil {
		return nil
	}

	result = append(result, map[string]interface{}{
		"retry": failurePolicy.Retry != nil,
	})

	return result
}
