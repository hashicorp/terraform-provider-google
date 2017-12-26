package google

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleCloudFunctionsFunction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCloudFunctionsFunctionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"entry_point": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"storage_bucket": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"storage_object": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"trigger_bucket": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"trigger_http": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"trigger_topic": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleCloudFunctionsFunctionRead(d *schema.ResourceData, meta interface{}) error {
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

	//Name of function should be unique
	d.SetId(d.Get("name").(string))

	return nil
}
