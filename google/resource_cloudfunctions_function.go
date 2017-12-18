package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/cloudfunctions/v1"

	"fmt"
)

const DEFAULT_FUNCTION_TIMEOUT_IN_SEC = 60

//Min is 1 second, max is 9 minutes 540 sec
const FUNCTION_TIMEOUT_MAX = 540
const FUNCTION_TIMEOUT_MIN = 1

func resourceCloudFunctionsFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFunctionsCreate,
		Read:   resourceCloudFunctionsRead,
		Delete: resourceCloudFunctionsDestroy,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"entry_point": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"memory": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	trigger := &cloudfunctions.HttpsTrigger{
		Url: "",
	}

	function := &cloudfunctions.CloudFunction{
		AvailableMemoryMb: 128,
		Description:       "Testing Function",
		EntryPoint:        "helloGET",
		HttpsTrigger:      trigger,
		Name:              createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, funcName),
		SourceArchiveUrl:  "gs://test-cloudfunctions-sk/index.zip",
	}

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
	d.Set("region", funcRegion)
	d.Set("project", funcProject)
	if getOpt.SourceArchiveUrl != "" {
		d.Set("source", getOpt.SourceArchiveUrl) //TODO: Make other options of source here
	}
	if getOpt.HttpsTrigger != nil {
		d.Set("trigger_http", getOpt.HttpsTrigger.Url)
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
