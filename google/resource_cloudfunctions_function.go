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

const DEFAULT_FUNCTION_TIMEOUT_IN_SEC = 60

//Min is 1 second, max is 9 minutes 540 sec
const FUNCTION_TIMEOUT_MAX = 540
const FUNCTION_TIMEOUT_MIN = 1

//Allowed values are: 128MB, 256MB, 512MB, 1024MB, and 2048MB. By default, a new function is limited to 256MB of memory.
const FUNCTION_DEFAULT_MEMORY = 256

var FUNCTION_ALLOWED_MEMORY = map[int]bool{
	128:  true,
	256:  true,
	512:  true,
	1024: true,
	2048: true,
}

type cloudFunctionId struct {
	Project string
	Region  string
	Name    string
}

func (s *cloudFunctionId) cloudFunctionId() string {
	return fmt.Sprintf("projects/%s/locations/%s/functions/%s", s.Project, s.Region, s.Name)
}

func (s *cloudFunctionId) parentId() string {
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

			"storage_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_object": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"entry_point": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"available_memory_mb": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  FUNCTION_DEFAULT_MEMORY,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					availableMemoryMB := v.(int)

					if FUNCTION_ALLOWED_MEMORY[availableMemoryMB] != true {
						errors = append(errors, fmt.Errorf("Allowed values for memory (in MB) are: %s . Got %d",
							joinMapKeys(&FUNCTION_ALLOWED_MEMORY), availableMemoryMB))
					}
					return
				},
			},

			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      DEFAULT_FUNCTION_TIMEOUT_IN_SEC,
				ValidateFunc: validation.IntBetween(FUNCTION_TIMEOUT_MIN, FUNCTION_TIMEOUT_MAX),
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

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				//For now CloudFunctions are allowed only in us-central1
				//Please see https://cloud.google.com/about/locations/
				ValidateFunc: validation.StringInSlice([]string{"us-central1"}, true),
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

	cloudFuncId := &cloudFunctionId{
		Project: project,
		Region:  region,
		Name:    d.Get("name").(string),
	}

	function := &cloudfunctions.CloudFunction{
		Name: cloudFuncId.cloudFunctionId(),
	}

	storageBucket := d.Get("storage_bucket").(string)
	storageObj := d.Get("storage_object").(string)
	function.SourceArchiveUrl = fmt.Sprintf("gs://%v/%v", storageBucket, storageObj)

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

	if v, ok := d.GetOk("source"); ok {
		function.SourceArchiveUrl = v.(string)
	}

	if v, ok := d.GetOk("timeout"); ok {
		function.Timeout = fmt.Sprintf("%vs", v.(int))
	}

	v, triggHttpOk := d.GetOk("trigger_http")
	if triggHttpOk {
		if v.(bool) == true {
			function.HttpsTrigger = &cloudfunctions.HttpsTrigger{
				Url: "",
			}
		}
	}

	v, triggTopicOk := d.GetOk("trigger_topic")
	if triggTopicOk {
		//Make PubSub event publish as in https://cloud.google.com/functions/docs/calling/pubsub
		function.EventTrigger = &cloudfunctions.EventTrigger{
			//Other events are not supported
			EventType: "providers/cloud.pubsub/eventTypes/topic.publish",
			//Must be like projects/PROJECT_ID/topics/NAME
			//Topic must be in same project as function
			Resource: fmt.Sprintf("projects/%s/topics/%s", project, v.(string)),
		}
	}

	v, triggBucketOk := d.GetOk("trigger_bucket")
	if triggBucketOk {
		//Make Storage event as in https://cloud.google.com/functions/docs/calling/storage
		function.EventTrigger = &cloudfunctions.EventTrigger{
			EventType: "providers/cloud.storage/eventTypes/object.change",
			//Must be like projects/PROJECT_ID/buckets/NAME
			//Bucket must be in same project as function
			Resource: fmt.Sprintf("projects/%s/buckets/%s", project, v.(string)),
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
		cloudFuncId.parentId(), function).Do()
	if err != nil {
		return err
	}
	err = cloudFunctionsOperationWait(config.clientCloudFunctions, op, "Creating CloudFunctions Function")
	if err != nil {
		return err
	}

	//Name of function should be unique
	d.SetId(cloudFuncId.terraformId())

	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Reading google_cloudfunctions_function")
	config := meta.(*Config)

	cloudFuncId, err := parseCloudFunctionId(d.Id(), config)
	if err != nil {
		return err
	}

	getOpt, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(cloudFuncId.cloudFunctionId()).Do()
	if err != nil {
		return err
	}

	d.Set("name", cloudFuncId.Name)
	d.Set("description", getOpt.Description)
	d.Set("entry_point", getOpt.EntryPoint)
	d.Set("available_memory_mb", getOpt.AvailableMemoryMb)
	sRemoved := strings.Replace(getOpt.Timeout, "s", "", -1)
	timeout, err := strconv.Atoi(sRemoved)
	if err != nil {
		return err
	}
	d.Set("timeout", timeout)
	d.Set("labels", getOpt.Labels)
	if getOpt.SourceArchiveUrl != "" {
		sourceArr := strings.Split(getOpt.SourceArchiveUrl, "/")
		d.Set("storage_bucket", sourceArr[2])
		d.Set("storage_object", sourceArr[3])
	}

	if getOpt.HttpsTrigger != nil {
		d.Set("trigger_http", true)
	}
	if getOpt.EventTrigger != nil {
		switch getOpt.EventTrigger.EventType {
		//From https://github.com/google/google-api-go-client/blob/master/cloudfunctions/v1/cloudfunctions-gen.go#L335
		case "providers/cloud.pubsub/eventTypes/topic.publish":
			d.Set("trigger_topic", extractLastResourceFromUri(getOpt.EventTrigger.Resource))
		case "providers/cloud.storage/eventTypes/object.change":
			d.Set("trigger_bucket", extractLastResourceFromUri(getOpt.EventTrigger.Resource))
		}
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

	_, err = config.clientCloudFunctions.Projects.Locations.Functions.Get(function.Name).Do()

	if err != nil {
		return fmt.Errorf("Function %s doesn't exists.", d.Get("name").(string))
	}

	if d.HasChange("labels") {
		function.Labels = expandLabels(d)

		op, err := config.clientCloudFunctions.Projects.Locations.Functions.Patch(function.Name, &function).
			UpdateMask("labels").Do()

		if err != nil {
			return fmt.Errorf("Error when updating labels: %s", err)
		}

		d.SetPartial("labels")
		err = cloudFunctionsOperationWait(config.clientCloudFunctions, op,
			"Updating CloudFunctions Function Labels")
		if err != nil {
			return err
		}
	}
	configUpdate := false
	var updateMaskArr []string
	if d.HasChange("available_memory_mb") {
		availableMemoryMb := d.Get("available_memory_mb").(int)
		function.AvailableMemoryMb = int64(availableMemoryMb)
		updateMaskArr = append(updateMaskArr, "availableMemoryMb")
		configUpdate = true
	}

	if d.HasChange("description") {
		function.Description = d.Get("description").(string)
		updateMaskArr = append(updateMaskArr, "description")
		configUpdate = true
	}

	if d.HasChange("timeout") {
		function.Timeout = fmt.Sprintf("%vs", d.Get("timeout").(int))
		updateMaskArr = append(updateMaskArr, "timeout")
		configUpdate = true
	}

	if configUpdate {
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
	log.Printf("[DEBUG]: Destroying google_cloudfunctions_function")
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
