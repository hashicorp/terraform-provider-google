package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var loggingBucketConfigSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `The resource name of the bucket`,
	},
	"location": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The location of the bucket. The supported locations are: "global" "us-central1"`,
	},
	"bucket_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The name of the logging bucket. Logging automatically creates two log buckets: _Required and _Default.`,
	},
	"description": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: `An optional description for this bucket.`,
	},
	"retention_days": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     30,
		Description: `Logs will be retained by default for this amount of time, after which they will automatically be deleted. The minimum retention period is 1 day. If this value is set to zero at bucket creation time, the default time of 30 days will be used.`,
	},
	"lifecycle_state": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `The bucket's lifecycle such as active or deleted.`,
	},
}

type loggingBucketConfigIDFunc func(d *schema.ResourceData, config *Config) (string, error)

// ResourceLoggingBucketConfig creates a resource definition by merging a unique field (eg: folder) to a generic logging bucket
// config resource. In practice the only difference between these resources is the url location.
func ResourceLoggingBucketConfig(parentType string, parentSpecificSchema map[string]*schema.Schema, iDFunc loggingBucketConfigIDFunc) *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingBucketConfigAcquire(iDFunc),
		Read:   resourceLoggingBucketConfigRead,
		Update: resourceLoggingBucketConfigUpdate,
		Delete: resourceLoggingBucketConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceLoggingBucketConfigImportState(parentType),
		},
		Schema: mergeSchemas(loggingBucketConfigSchema, parentSpecificSchema),
	}
}

var loggingBucketConfigIDRegex = regexp.MustCompile("(.+)/(.+)/locations/(.+)/buckets/(.+)")

func resourceLoggingBucketConfigImportState(parent string) schema.StateFunc {
	return func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		parts := loggingBucketConfigIDRegex.FindStringSubmatch(d.Id())
		if parts == nil {
			return nil, fmt.Errorf("unable to parse logging sink id %#v", d.Id())
		}

		if len(parts) != 5 {
			return nil, fmt.Errorf("Invalid id format. Format should be '{{parent}}/{{parent_id}}/locations/{{location}}/buckets/{{bucket_id}} with parent in %s", loggingSinkResourceTypes)
		}

		validLoggingType := false
		for _, v := range loggingSinkResourceTypes {
			if v == parts[1] {
				validLoggingType = true
				break
			}
		}
		if !validLoggingType {
			return nil, fmt.Errorf("Logging parent type %s is not valid. Valid resource types: %#v", parts[1],
				loggingSinkResourceTypes)
		}

		d.Set(parent, parts[1]+"/"+parts[2])

		d.Set("location", parts[3])

		d.Set("bucket_id", parts[4])

		return []*schema.ResourceData{d}, nil
	}
}

func resourceLoggingBucketConfigAcquire(iDFunc loggingBucketConfigIDFunc) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		id, err := iDFunc(d, config)
		if err != nil {
			return err
		}

		d.SetId(id)

		return resourceLoggingBucketConfigUpdate(d, meta)
	}
}

func resourceLoggingBucketConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[DEBUG] Fetching logging bucket config: %#v", d.Id())

	url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", "", url, nil)
	if err != nil {
		log.Printf("[WARN] Unable to acquire logging bucket config at %s", d.Id())

		d.SetId("")
		return err
	}

	d.Set("name", res["name"])
	d.Set("description", res["description"])
	d.Set("lifecycle_state", res["lifecycleState"])
	d.Set("retention_days", res["retentionDays"])

	return nil

}

func resourceLoggingBucketConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})

	url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
	if err != nil {
		return err
	}

	obj["retentionDays"] = d.Get("retention_days")
	obj["description"] = d.Get("description")

	updateMask := []string{}
	if d.HasChange("retention_days") {
		updateMask = append(updateMask, "retentionDays")
	}
	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	_, err = sendRequestWithTimeout(config, "PATCH", "", url, obj, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("Error updating Logging Bucket Config %q: %s", d.Id(), err)
	}

	return resourceLoggingBucketConfigRead(d, meta)

}

func resourceLoggingBucketConfigDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[WARN] Logging bucket configs cannot be deleted. Removing logging bucket config from state: %#v", d.Id())
	d.SetId("")

	return nil
}
