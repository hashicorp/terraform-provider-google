package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/cloudfunctions/v1"

	"fmt"
	"log"
	"strconv"
	"time"
	"strings"
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

func resourceCloudFunctionsFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFunctionsCreate,
		Read:   resourceCloudFunctionsRead,
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

			"entry_point": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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

			"remove_labels": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"retry": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      DEFAULT_FUNCTION_TIMEOUT_IN_SEC,
				ValidateFunc: validation.IntBetween(FUNCTION_TIMEOUT_MIN, FUNCTION_TIMEOUT_MAX),
			},

			"update_labels": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"stage_bucket": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"trigger_bucket": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"trigger_http": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},

			"trigger_topic": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	funcName := d.Get("name").(string)

	var memory int
	if v, ok := d.GetOk("memory"); ok {
		memory = v.(int)
		if FUNCTION_ALLOWED_MEMORY[memory] != true {
			return fmt.Errorf("Allowed values for memory are: 128MB, 256MB, 512MB, 1024MB, and 2048MB. "+
				"Got %d", memory)
		}
	}

	var entryPoint string
	if v, ok := d.GetOk("entry_point"); ok {
		entryPoint = v.(string)
	}

	var timeout int
	if v, ok := d.GetOk("timeout"); ok {
		timeout = v.(int)
	}

	var triggerHttp *cloudfunctions.HttpsTrigger
	triggerHttp = nil
	var triggerTopicOrBucket *cloudfunctions.EventTrigger
	if v, ok := d.GetOk("trigger_http"); ok {
		if v.(bool) == true {
			triggerHttp = &cloudfunctions.HttpsTrigger{
				Url: "",
			}
		}
	}
	triggerTopicOrBucket = nil
	if v, ok := d.GetOk("trigger_topic"); ok {
		//Make PubSub event publish as in https://cloud.google.com/functions/docs/calling/pubsub
		triggerTopicOrBucket = &cloudfunctions.EventTrigger{
			EventType: "providers/cloud.pubsub/eventTypes/topic.publish",
			Resource:  v.(string),
		}
	}
	if v, ok := d.GetOk("trigger_bucket"); ok {
		if triggerTopicOrBucket != nil {
			//It was previously initialized by trigger_topic - can't do both
			return fmt.Errorf("One of aguments only [trigger_bucket, trigger_http] must be used.")
		}
		//Make Storage event as in https://cloud.google.com/functions/docs/calling/storage
		triggerTopicOrBucket = &cloudfunctions.EventTrigger{
			EventType: "providers/cloud.storage/eventTypes/object.change",
			Resource:  v.(string),
		}
	}

	if triggerHttp == nil && triggerTopicOrBucket == nil {
		//It's bad when no trigger is specified
		return fmt.Errorf(
			"One of aguments [trigger_topic, trigger_bucket, trigger_http] is required: " +
				"You must specify a trigger when deploying a new function.")
	}
	if triggerHttp != nil && triggerTopicOrBucket != nil {
		//Also bad when too many triggers specified
		return fmt.Errorf(
			"Only one of aguments [trigger_topic, trigger_bucket, trigger_http] is allowed.")
	}

	function := &cloudfunctions.CloudFunction{
		AvailableMemoryMb: int64(memory),
		EntryPoint:        entryPoint,
		HttpsTrigger:      triggerHttp,
		EventTrigger:      triggerTopicOrBucket,
		Timeout:           fmt.Sprintf("%vs", timeout),
		Name:              createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, funcName),
		SourceArchiveUrl:  "gs://test-cloudfunctions-sk/index.zip",
	}

	log.Printf("[DEBUG] Creating cloud function: %s", function.Name)
	createOpt, err := service.Projects.Locations.Functions.Create(
		createCloudFunctionsPathString(CLOUDFUNCTIONS_REGION_ONLY, project, region, ""), function).Do()
	if err != nil {
		return err
	}
	if createOpt.Done == false {
		_, err := getCloudFunctionsOperationsResults(createOpt.Name, service)
		if err != nil {
			return err
		}
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

	d.Set("name", name)
	d.Set("entry_point", getOpt.EntryPoint)
	d.Set("memory", getOpt.AvailableMemoryMb)
	timeout, err := strconv.Atoi(strings.Replace(getOpt.Timeout, "s", "", -1))
	if err != nil {
		return err
	}
	d.Set("timeout", timeout)
	d.Set("region", funcRegion)
	d.Set("project", funcProject)
	if getOpt.SourceArchiveUrl != "" {
		d.Set("source", getOpt.SourceArchiveUrl) //TODO: Make other options of source here
	}
	if getOpt.HttpsTrigger != nil {
		d.Set("trigger_http", true)
	}
	if getOpt.EventTrigger != nil {
		switch getOpt.EventTrigger.EventType {
		//TODO: Get other event triggers: Read https://github.com/google/google-api-go-client/blob/master/cloudfunctions/v1/cloudfunctions-gen.go#L335
		}
	}

	return nil
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
		return fmt.Errorf("Error reading cloud function name %s.", name, err)
	}

	deleteOpt, err := service.Projects.Locations.Functions.Delete(
		createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
	if err != nil {
		return err
	}
	if deleteOpt.Done == false {
		_, err := getCloudFunctionsOperationsResults(deleteOpt.Name, service)
		if err != nil {
			return err
		}
	}

	return nil
}
