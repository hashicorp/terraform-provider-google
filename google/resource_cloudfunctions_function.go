package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/cloudfunctions/v1"

	"fmt"
	"log"
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

//For now CloudFunctions are allowed only in us-central1
//Please see https://cloud.google.com/about/locations/
var FUNCTION_ALLOWED_REGION = map[string]bool{
	"us-central1": true,
}

func resourceCloudFunctionsFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFunctionsCreate,
		Read:   resourceCloudFunctionsRead,
		Update: resourceCloudFunctionsUpdate,
		Delete: resourceCloudFunctionsDestroy,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
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
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"entry_point": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Default:  FUNCTION_DEFAULT_MEMORY,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"storage_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"storage_object": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"trigger_bucket": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"trigger_http", "trigger_topic"},
			},

			"trigger_http": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeBool},
				ConflictsWith: []string{"trigger_bucket", "trigger_topic"},
			},

			"trigger_topic": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"trigger_http", "trigger_bucket"},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceCloudFunctionsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := config.clientCloudFunctions

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	if FUNCTION_ALLOWED_REGION[region] != true {
		return fmt.Errorf("Invalid region. Now allowed only us-central1. See https://cloud.google.com/about/locations/")
	}

	funcName := d.Get("name").(string)

	function := &cloudfunctions.CloudFunction{
		Name: createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, funcName),
	}

	storageBucket := d.Get("storage_bucket").(string)
	storageObj := d.Get("storage_object").(string)
	function.SourceArchiveUrl = fmt.Sprintf("gs://%v/%v", storageBucket, storageObj)

	if v, ok := d.GetOk("memory"); ok {
		memory := v.(int)
		if FUNCTION_ALLOWED_MEMORY[memory] != true {
			return fmt.Errorf("Allowed values for memory are: 128MB, 256MB, 512MB, 1024MB, and 2048MB. "+
				"Got %d", memory)
		}
		function.AvailableMemoryMb = int64(memory)
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
		return fmt.Errorf("One of aguments [trigger_topic, trigger_bucket, trigger_http] is required: " +
			"You must specify a trigger when deploying a new function.")
	}

	if _, ok := d.GetOk("labels"); ok {
		function.Labels = expandLabels(d)
	}

	log.Printf("[DEBUG] Creating cloud function: %s", function.Name)
	op, err := service.Projects.Locations.Functions.Create(
		createCloudFunctionsPathString(CLOUDFUNCTIONS_REGION_ONLY, project, region, ""), function).Do()
	if err != nil {
		return err
	}
	err = cloudFunctionsOperationWait(config.clientCloudFunctions, op, "Creating CloudFunctions Function")
	if err != nil {
		return err
	}

	//Name of function should be unique
	d.SetId(d.Get("name").(string))

	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := config.clientCloudFunctions

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	getOpt, err := service.Projects.Locations.Functions.Get(
		createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
	if err != nil {
		return err
	}
	nameFromGet, err := getCloudFunctionName(getOpt.Name)
	if err != nil {
		return err
	}
	if name != nameFromGet {
		return fmt.Errorf("Name of Cloud Function is mismatched from what we get from Google Cloud %s != %s", name, nameFromGet)
	}
	funcRegion, err := getCloudFunctionRegion(getOpt.Name)
	if err != nil {
		return err
	}
	funcProject, err := getCloudFunctionProject(getOpt.Name)
	if err != nil {
		return err
	}

	d.Set("description", getOpt.Description)
	d.Set("entry_point", getOpt.EntryPoint)
	d.Set("memory", getOpt.AvailableMemoryMb)
	d.Set("region", funcRegion)
	timeout, err := readTimeout(getOpt.Timeout)
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
	d.Set("project", funcProject)

	return nil
}

func resourceCloudFunctionsUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := config.clientCloudFunctions

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	function := cloudfunctions.CloudFunction{
		Name: createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, d.Get("name").(string)),
	}

	_, err = service.Projects.Locations.Functions.Get(function.Name).Do()

	if err != nil {
		return fmt.Errorf("Function %s doesn't exists.", d.Get("name").(string))
	}

	if d.HasChange("labels") {
		function.Labels = expandLabels(d)

		op, err := service.Projects.Locations.Functions.Patch(function.Name, &function).
			UpdateMask("labels").Do()

		if err != nil {
			return fmt.Errorf("Error when updating labels: %s", err)
		}

		d.SetPartial("labels")
		err = cloudFunctionsOperationWait(config.clientCloudFunctions, op, "Updating CloudFunctions Function Labels")
		if err != nil {
			return err
		}
	}
	configUpdate := false
	var updateMaskArr []string
	var partialArr []string
	if d.HasChange("memory") {
		memory := d.Get("memory").(int)

		if FUNCTION_ALLOWED_MEMORY[memory] != true {
			return fmt.Errorf("Allowed values for memory are: 128MB, 256MB, 512MB, 1024MB, and 2048MB. "+
				"Got %d", memory)
		}
		function.AvailableMemoryMb = int64(memory)
		updateMaskArr = append(updateMaskArr, "availableMemoryMb")
		partialArr = append(partialArr, "memory")
		configUpdate = true
	}

	if d.HasChange("description") {
		function.Description = d.Get("description").(string)
		updateMaskArr = append(updateMaskArr, "description")
		partialArr = append(partialArr, "description")
		configUpdate = true
	}

	if d.HasChange("timeout") {
		function.Timeout = fmt.Sprintf("%vs", d.Get("timeout").(int))
		updateMaskArr = append(updateMaskArr, "timeout")
		partialArr = append(partialArr, "timeout")
		configUpdate = true
	}

	if configUpdate {
		log.Printf("[DEBUG] Send Patch CloudFunction Configuration request: %#v", function)
		updateMask := strings.Join(updateMaskArr, ",")
		op, err := service.Projects.Locations.Functions.Patch(function.Name, &function).
			UpdateMask(updateMask).Do()

		if err != nil {
			return fmt.Errorf("Error while updating cloudfunction configuration: %s", err)
		}
		for i := range partialArr {
			d.SetPartial(partialArr[i])
		}

		err = cloudFunctionsOperationWait(config.clientCloudFunctions, op, "Updating CloudFunctions Function")
		if err != nil {
			return err
		}
	}
	d.Partial(false)

	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	service := config.clientCloudFunctions

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	if len(name) == 0 {
		return fmt.Errorf("Error reading cloud function name %s.", name)
	}

	op, err := service.Projects.Locations.Functions.Delete(
		createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
	if err != nil {
		return err
	}
	err = cloudFunctionsOperationWait(config.clientCloudFunctions, op, "Deleting CloudFunctions Function")
	if err != nil {
		return err
	}

	return nil
}
