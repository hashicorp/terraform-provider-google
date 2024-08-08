// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package artifactregistry

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleArtifactRegistryLocations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleArtifactRegistryLocationsRead,
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

func dataSourceGoogleArtifactRegistryLocationsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "https://artifactregistry.googleapis.com/v1/projects/{{project}}/locations")
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error listing Artifact Registry Locations : %s", err)
	}

	locationsRaw := flattenArtifactRegistryLocations(res)

	locations := make([]string, len(locationsRaw))
	for i, loc := range locationsRaw {
		locations[i] = loc.(string)
	}
	sort.Strings(locations)

	log.Printf("[DEBUG] Received Artifact Registry Locations: %q", locations)

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("locations", locations); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s", project))

	return nil
}

func flattenArtifactRegistryLocations(resp map[string]interface{}) []interface{} {
	regionList := resp["locations"].([]interface{})
	regions := make([]interface{}, len(regionList))
	for i, v := range regionList {
		regionObj := v.(map[string]interface{})
		regions[i] = regionObj["locationId"]
	}
	return regions
}
