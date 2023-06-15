// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleIamTestablePermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleIamTestablePermissionsRead,
		Schema: map[string]*schema.Schema{
			"full_resource_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stages": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ALPHA", "BETA", "GA", "DEPRECATED"}, true),
				},
			},
			"custom_support_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "SUPPORTED",
				ValidateFunc: validation.StringInSlice([]string{"NOT_SUPPORTED", "SUPPORTED", "TESTING"}, true),
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_support_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleIamTestablePermissionsRead(d *schema.ResourceData, meta interface{}) (err error) {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	body := make(map[string]interface{})
	body["pageSize"] = 500
	permissions := make([]map[string]interface{}, 0)

	custom_support_level := strings.ToUpper(d.Get("custom_support_level").(string))
	stages := []string{}
	for _, e := range d.Get("stages").([]interface{}) {
		stages = append(stages, strings.ToUpper(e.(string)))
	}
	if len(stages) == 0 {
		// Since schema.TypeLists cannot specify defaults, we'll specify it here
		stages = append(stages, "GA")
	}
	for {
		url := "https://iam.googleapis.com/v1/permissions:queryTestablePermissions"
		body["fullResourceName"] = d.Get("full_resource_name").(string)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			RawURL:    url,
			UserAgent: userAgent,
			Body:      body,
		})
		if err != nil {
			return fmt.Errorf("Error retrieving permissions: %s", err)
		}

		pagePermissions := flattenTestablePermissionsList(res["permissions"], custom_support_level, stages)
		permissions = append(permissions, pagePermissions...)
		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			body["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err = d.Set("permissions", permissions); err != nil {
		return fmt.Errorf("Error retrieving permissions: %s", err)
	}

	d.SetId(d.Get("full_resource_name").(string))
	return nil
}

func flattenTestablePermissionsList(v interface{}, custom_support_level string, stages []string) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})
	permissions := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		p := raw.(map[string]interface{})

		if _, ok := p["name"]; ok {
			var csl bool
			if custom_support_level == "SUPPORTED" {
				csl = p["customRolesSupportLevel"] == nil || p["customRolesSupportLevel"] == "SUPPORTED"
			} else {
				csl = p["customRolesSupportLevel"] == custom_support_level
			}
			if csl && p["stage"] != nil && tpgresource.StringInSlice(stages, p["stage"].(string)) {
				permissions = append(permissions, map[string]interface{}{
					"name":                 p["name"],
					"title":                p["title"],
					"stage":                p["stage"],
					"api_disabled":         p["apiDisabled"],
					"custom_support_level": p["customRolesSupportLevel"],
				})
			}
		}
	}

	return permissions
}
