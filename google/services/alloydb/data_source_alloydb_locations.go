// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbLocations() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceAlloydbLocationsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project.`,
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `Resource name for the location, which may vary between implementations. For example: "projects/example-project/locations/us-east1`,
						},
						"location_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The canonical id for this location. For example: "us-east1".`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The friendly name for this location, typically a nearby city name. For example, "Tokyo".`,
						},
						"labels": {
							Type:        schema.TypeMap,
							Computed:    true,
							Optional:    true,
							Description: `Cross-service attributes for the location. For example {"cloud.googleapis.com/region": "us-east1"}`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"metadata": {
							Type:        schema.TypeMap,
							Computed:    true,
							Optional:    true,
							Description: `Service-specific metadata. For example the available capacity at the given location.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlloydbLocationsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{AlloydbBasePath}}projects/{{project}}/locations")
	if err != nil {
		return fmt.Errorf("Error setting api endpoint")
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Locations %q", d.Id()))
	}
	var locations []map[string]interface{}
	for {
		fetchedLocations := res["locations"].([]interface{})
		for _, loc := range fetchedLocations {
			locationDetails := make(map[string]interface{})
			l := loc.(map[string]interface{})
			if l["name"] != nil {
				locationDetails["name"] = l["name"].(string)
			}
			if l["locationId"] != nil {
				locationDetails["location_id"] = l["locationId"].(string)
			}
			if l["displayName"] != nil {
				locationDetails["display_id"] = l["displayName"].(string)
			}
			if l["labels"] != nil {
				labels := make(map[string]string)
				for k, v := range l["labels"].(map[string]interface{}) {
					labels[k] = v.(string)
				}
				locationDetails["labels"] = labels
			}
			if l["metadata"] != nil {
				metadata := make(map[string]string)
				for k, v := range l["metadata"].(map[interface{}]interface{}) {
					metadata[k.(string)] = v.(string)
				}
				locationDetails["metadata"] = metadata
			}
			locations = append(locations, locationDetails)
		}
		if res["nextPageToken"] == nil || res["nextPageToken"].(string) == "" {
			break
		}
		url, err = tpgresource.ReplaceVars(d, config, "{{AlloydbBasePath}}projects/{{project}}/locations?pageToken="+res["nextPageToken"].(string))
		if err != nil {
			return fmt.Errorf("Error setting api endpoint")
		}
		res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Locations %q", d.Id()))
		}
	}

	if err := d.Set("locations", locations); err != nil {
		return fmt.Errorf("Error setting locations: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/locations", project))
	return nil
}
