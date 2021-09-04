package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var loggingViewResourceTypes = []string{
	"billingAccounts",
	"folders",
	"organizations",
	"projects",
}

var loggingViewSchema = map[string]*schema.Schema{
	"view_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The ID of the view",
	},
	"name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `The resource name of the view`,
	},
	"bucket": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The resource name of the bucket.`,
	},
	"description": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: `An optional description for this bucket.`,
	},
	"filter": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: `Filter that restricts which log entries in a bucket are visible in this view.`,
	},
}

func ResourceLoggingView() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingViewCreate,
		Read:   resourceLoggingViewRead,
		Update: resourceLoggingViewUpdate,
		Delete: resourceLoggingViewDelete,
		Importer: &schema.ResourceImporter{
			State: resourceLoggingViewImportState,
		},
		Schema:        loggingViewSchema,
		UseJSONNumber: true,
	}
}

var loggingViewIDRegex = regexp.MustCompile("((.+)/.+/locations/.+/buckets/.+)/views/(.+)")

func resourceLoggingViewImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := loggingViewIDRegex.FindStringSubmatch(d.Id())
	if parts == nil {
		return nil, fmt.Errorf("Unable to parse Log View id %#v", d.Id())
	}

	if len(parts) != 4 {
		return nil, fmt.Errorf("Invalid id format. Format should be '{{parent}}/{{parent_id}}/locations/{{location}}/buckets/{{bucket_id}}/views/{{view_id}} with parent in %s", loggingSinkResourceTypes)
	}

	validLoggingType := false
	for _, v := range loggingViewResourceTypes {
		if v == parts[2] {
			validLoggingType = true
			break
		}
	}
	if !validLoggingType {
		return nil, fmt.Errorf("Logging parent type %s is not valid. Valid resource types: %#v", parts[1],
			loggingViewResourceTypes)
	}

	if err := d.Set("bucket", parts[1]); err != nil {
		return nil, fmt.Errorf("Error setting bucket: %s", err)
	}
	if err := d.Set("view_id", parts[3]); err != nil {
		return nil, fmt.Errorf("Error setting view_id: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceLoggingViewCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	obj["filter"] = d.Get("filter")
	obj["description"] = d.Get("description")

	url, err := replaceVars(d, config, "{{LoggingBasePath}}{{bucket}}/views?viewId={{view_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Log View: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Log View: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/views/%s", d.Get("bucket"), d.Get("view_id")))
	log.Printf("[DEBUG] Finished creating Log View %q: %#v", d.Id(), res)

	return resourceLoggingViewRead(d, meta)
}

func resourceLoggingViewRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Fetching Log View: %#v", d.Id())

	url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", "", url, userAgent, nil)
	if err != nil {
		log.Printf("[WARN] Unable to acquire Log View at %s", d.Id())

		d.SetId("")
		return err
	}

	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("filter", res["filter"]); err != nil {
		return fmt.Errorf("Error setting filter: %s", err)
	}
	if err := d.Set("description", res["description"]); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}

	return nil
}

func resourceLoggingViewUpdate(d *schema.ResourceData, meta interface{}) error {
	var updateMask []string
	if d.HasChange("filter") {
		updateMask = append(updateMask, "filter")
	}
	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}
	return resourceLoggingViewUpdateWithUpdateMask(d, meta, updateMask)
}

func resourceLoggingViewUpdateWithUpdateMask(d *schema.ResourceData, meta interface{}, updateMask []string) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	for _, field := range updateMask {
		obj[field] = d.Get(field)
	}

	url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
	if err != nil {
		return err
	}

	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}
	_, err = sendRequestWithTimeout(config, "PATCH", "", url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("Error updating Log View %q: %s", d.Id(), err)
	}

	return resourceLoggingViewRead(d, meta)
}

func resourceLoggingViewDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)

	url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
	if err != nil {
		return err
	}

	_, err = sendRequest(config, "DELETE", "", url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, "Log View")
	}

	d.SetId("")
	return nil
}
