package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/cloudfunctions/v1"

	"fmt"
	"log"
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

// For now CloudFunctions are allowed only in us-central1
// Please see https://cloud.google.com/about/locations/
var validCloudFunctionRegion = validation.StringInSlice([]string{"us-central1"}, true)

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
				ForceNew: true,
			},

			"source_archive_object": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"trigger_bucket": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
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
				ForceNew:      true,
				ConflictsWith: []string{"trigger_http", "trigger_bucket"},
			},

			"https_trigger_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"retry_on_failure": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"trigger_http"},
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
		Name: cloudFuncId.cloudFunctionId(),
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

	v, triggHttpOk := d.GetOk("trigger_http")
	if triggHttpOk && v.(bool) {
		function.HttpsTrigger = &cloudfunctions.HttpsTrigger{}
	}

	v, triggTopicOk := d.GetOk("trigger_topic")
	if triggTopicOk {
		// Make PubSub event publish as in https://cloud.google.com/functions/docs/calling/pubsub
		function.EventTrigger = &cloudfunctions.EventTrigger{
			// Other events are not supported
			EventType: "google.pubsub.topic.publish",
			// Must be like projects/PROJECT_ID/topics/NAME
			// Topic must be in same project as function
			Resource: fmt.Sprintf("projects/%s/topics/%s", project, v.(string)),
		}
		if d.Get("retry_on_failure").(bool) {
			function.EventTrigger.FailurePolicy = &cloudfunctions.FailurePolicy{
				Retry: &cloudfunctions.Retry{},
			}
		}
	}

	v, triggBucketOk := d.GetOk("trigger_bucket")
	if triggBucketOk {
		// Make Storage event as in https://cloud.google.com/functions/docs/calling/storage
		function.EventTrigger = &cloudfunctions.EventTrigger{
			EventType: "providers/cloud.storage/eventTypes/object.change",
			// Must be like projects/PROJECT_ID/buckets/NAME
			// Bucket must be in same project as function
			Resource: fmt.Sprintf("projects/%s/buckets/%s", project, v.(string)),
		}
		if d.Get("retry_on_failure").(bool) {
			function.EventTrigger.FailurePolicy = &cloudfunctions.FailurePolicy{
				Retry: &cloudfunctions.Retry{},
			}
		}
	}

	if !triggHttpOk && !triggTopicOk && !triggBucketOk {
		return fmt.Errorf("One of arguments [trigger_topic, trigger_bucket, trigger_http] is required: " +
			"You must specify a trigger when deploying a new function.")
	}

	if _, ok := d.GetOk("labels"); ok {
		function.Labels = expandLabels(d)
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
	if function.SourceArchiveUrl != "" {
		sourceArr := strings.Split(function.SourceArchiveUrl, "/")
		d.Set("source_archive_bucket", sourceArr[2])
		d.Set("source_archive_object", sourceArr[3])
	}

	if function.HttpsTrigger != nil {
		d.Set("trigger_http", true)
		d.Set("https_trigger_url", function.HttpsTrigger.Url)
	}

	if function.EventTrigger != nil {
		switch function.EventTrigger.EventType {
		// From https://github.com/google/google-api-go-client/blob/master/cloudfunctions/v1/cloudfunctions-gen.go#L335
		case "google.pubsub.topic.publish":
			d.Set("trigger_topic", GetResourceNameFromSelfLink(function.EventTrigger.Resource))
		case "providers/cloud.storage/eventTypes/object.change":
			d.Set("trigger_bucket", GetResourceNameFromSelfLink(function.EventTrigger.Resource))
		}
		retry := function.EventTrigger.FailurePolicy != nil && function.EventTrigger.FailurePolicy.Retry != nil
		d.Set("retry_on_failure", retry)
	}
	d.Set("region", cloudFuncId.Region)
	d.Set("project", cloudFuncId.Project)

	return nil
}

func resourceCloudFunctionsUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Updating google_cloudfunctions_function")
	config := meta.(*Config)

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

	if d.HasChange("retry_on_failure") {
		if d.Get("retry_on_failure").(bool) {
			if function.EventTrigger == nil {
				function.EventTrigger = &cloudfunctions.EventTrigger{}
			}
			function.EventTrigger.FailurePolicy = &cloudfunctions.FailurePolicy{
				Retry: &cloudfunctions.Retry{},
			}
		}
		updateMaskArr = append(updateMaskArr, "eventTrigger.failurePolicy.retry")
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
