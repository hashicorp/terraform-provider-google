// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleKmsKeyRings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsKeyRingsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The canonical id for the location. For example: "us-east1".`,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `
					The filter argument is used to add a filter query parameter that limits which keys are retrieved by the data source: ?filter={{filter}}.
					Example values:
					
					* "name:my-key-" will retrieve key rings that contain "my-key-" anywhere in their name. Note: names take the form projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}.
					* "name=projects/my-project/locations/global/keyRings/my-key-ring" will only retrieve a key ring with that exact name.
					
					[See the documentation about using filters](https://cloud.google.com/kms/docs/sorting-and-filtering)
				`,
			},
			"key_rings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all the retrieved key rings",
				Elem: &schema.Resource{
					// schema isn't used from resource_kms_key_ring due to having project and location fields which are empty when grabbed in a list.
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleKmsKeyRingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/keyRings")
	if err != nil {
		return err
	}
	if filter, ok := d.GetOk("filter"); ok {
		id += "/filter=" + filter.(string)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Searching for keyrings")
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for keyRings: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	var keyRings []interface{}

	params := make(map[string]string)
	if filter, ok := d.GetOk("filter"); ok {
		log.Printf("[DEBUG] Search for key rings using filter ?filter=%s", filter.(string))
		params["filter"] = filter.(string)
		if err != nil {
			return err
		}
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}projects/{{project}}/locations/{{location}}/keyRings")
	if err != nil {
		return err
	}

	for {
		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              billingProject,
			RawURL:               url,
			UserAgent:            userAgent,
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429RetryableQuotaError},
		})
		if err != nil {
			return fmt.Errorf("Error retrieving buckets: %s", err)
		}

		if res["keyRings"] == nil {
			break
		}
		pageKeyRings, err := flattenKMSKeyRingsList(config, res["keyRings"])
		if err != nil {
			return fmt.Errorf("error flattening key rings list: %s", err)
		}
		keyRings = append(keyRings, pageKeyRings...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	log.Printf("[DEBUG] Found %d key rings", len(keyRings))
	if err := d.Set("key_rings", keyRings); err != nil {
		return fmt.Errorf("error setting key rings: %s", err)
	}

	return nil
}

// flattenKMSKeyRingsList flattens a list of key rings
func flattenKMSKeyRingsList(config *transport_tpg.Config, keyRingsList interface{}) ([]interface{}, error) {
	var keyRings []interface{}
	for _, k := range keyRingsList.([]interface{}) {
		keyRing := k.(map[string]interface{})

		parsedId, err := parseKmsKeyRingId(keyRing["name"].(string), config)
		if err != nil {
			return nil, err
		}

		data := map[string]interface{}{}
		// The google_kms_key_rings resource and dataset set
		// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{name}})
		// and set name is set as just {{name}}.
		data["id"] = keyRing["name"]
		data["name"] = parsedId.Name

		keyRings = append(keyRings, data)
	}

	return keyRings, nil
}
