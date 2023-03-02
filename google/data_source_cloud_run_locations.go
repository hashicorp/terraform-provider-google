package google

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleCloudRunLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCloudRunLocationsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleCloudRunLocationsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "https://run.googleapis.com/v1/projects/{{project}}/locations")
	if err != nil {
		return err
	}

	res, err := SendRequest(config, "GET", project, url, userAgent, nil)
	if err != nil {
		return fmt.Errorf("Error listing Cloud Run Locations : %s", err)
	}

	locationsRaw := flattenCloudRunLocations(res)

	locations := make([]string, len(locationsRaw))
	for i, loc := range locationsRaw {
		locations[i] = loc.(string)
	}
	sort.Strings(locations)

	log.Printf("[DEBUG] Received Google Cloud Run Locations: %q", locations)

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("locations", locations); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s", project))

	return nil
}

func flattenCloudRunLocations(resp map[string]interface{}) []interface{} {
	regionList := resp["locations"].([]interface{})
	regions := make([]interface{}, len(regionList))
	for i, v := range regionList {
		regionObj := v.(map[string]interface{})
		regions[i] = regionObj["locationId"]
	}
	return regions
}
