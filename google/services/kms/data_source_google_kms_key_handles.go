// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleKmsKeyHandles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsKeyHandlesRead,
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
			"resource_type_selector": {
				Type:     schema.TypeString,
				Required: true,
				Description: `
					The resource_type_selector argument is used to add a filter query parameter that limits which key handles are retrieved by the data source: ?filter=resource_type_selector="{{resource_type_selector}}".
					Example values:
					* resource_type_selector="{SERVICE}.googleapis.com/{TYPE}".
					[See the documentation about using filters](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyHandles/list)
				`,
			},
			"key_handles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all the retrieved key handles",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type_selector": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

}

func dataSourceGoogleKmsKeyHandlesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	resourceTypeSelector := ""
	if fl, ok := d.GetOk("resource_type_selector"); ok {
		resourceTypeSelector = strings.Replace(fl.(string), "\"", "%22", -1)
	}

	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}projects/{{project}}/locations/{{location}}/keyHandles")
	if err != nil {
		return err
	}
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	params := make(map[string]string)
	var keyHandles []interface{}
	for {
		newUrl, err := addQueryParams(url, resourceTypeSelector, params)
		if err != nil {
			return err
		}
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              billingProject,
			RawURL:               newUrl,
			UserAgent:            userAgent,
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429RetryableQuotaError},
		})
		if err != nil {
			return fmt.Errorf("Error retrieving keyhandles: %s", err)
		}

		if res["keyHandles"] == nil {
			break
		}
		pageKeyHandles, err := flattenKMSKeyHandlesList(config, res["keyHandles"])
		if err != nil {
			return fmt.Errorf("error flattening key handle list: %s", err)
		}
		keyHandles = append(keyHandles, pageKeyHandles...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}
	log.Printf("[DEBUG] Found %d key handles", len(keyHandles))
	if err := d.Set("key_handles", keyHandles); err != nil {
		return fmt.Errorf("error setting key handles: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/locations/%s/keyHandles?filter=resource_type_selector=%s", project, d.Get("location"), resourceTypeSelector))
	return nil
}

// transport_tpg.AddQueryParams() encodes the filter=resource_type_selector="value" into
// filter=resource_type_selector%3D%22value%22
// The encoding of '=' into %3D is currently causing issue with ListKeyHandle api.
// To to handle this case currently, as part of this function,
// we are manually adding filter as a query param to the url
func addQueryParams(url string, resourceTypeSelector string, params map[string]string) (string, error) {
	quoteEncoding := "%22"
	if len(params) == 0 {
		return fmt.Sprintf("%s?filter=resource_type_selector=%s%s%s", url, quoteEncoding, resourceTypeSelector, quoteEncoding), nil
	} else {
		url, err := transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return "", nil
		}
		return fmt.Sprintf("%s&filter=resource_type_selector=%s%s%s", url, quoteEncoding, resourceTypeSelector, quoteEncoding), nil
	}
}

// flattenKMSKeyHandlesList flattens a list of key handles
func flattenKMSKeyHandlesList(config *transport_tpg.Config, keyHandlesList interface{}) ([]interface{}, error) {
	var keyHandles []interface{}
	for _, k := range keyHandlesList.([]interface{}) {
		keyHandle := k.(map[string]interface{})

		data := map[string]interface{}{}
		data["name"] = keyHandle["name"]
		data["kms_key"] = keyHandle["kmsKey"]
		data["resource_type_selector"] = keyHandle["resourceTypeSelector"]

		keyHandles = append(keyHandles, data)
	}

	return keyHandles, nil
}
