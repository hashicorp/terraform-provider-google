// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package accesscontextmanager

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAccessContextManagerAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccessContextManagerAccessPolicyRead,
		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scopes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccessContextManagerAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{AccessContextManagerBasePath}}accessPolicies?parent={{parent}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})

	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("AccessContextManagerAccessPolicy %q", d.Id()), url)
	}

	if res == nil {
		return fmt.Errorf("Error fetching policies: %s", err)
	}

	policies, err := parse_policies_response(res)
	if err != nil {
		return fmt.Errorf("Error parsing list policies response: %s", err)
	}

	// Find the matching policy in the list of policies response. Both the parent and scopes
	// should match
	for _, fetched_policy := range policies {
		scopes_match := compare_scopes(d.Get("scopes").([]interface{}), fetched_policy.Scopes)
		if fetched_policy.Parent == d.Get("parent").(string) && scopes_match {
			name_without_prefix := strings.Split(fetched_policy.Name, "accessPolicies/")[1]
			d.SetId(name_without_prefix)
			if err := d.Set("name", name_without_prefix); err != nil {
				return fmt.Errorf("Error setting policy name: %s", err)
			}

			if err := d.Set("title", fetched_policy.Title); err != nil {
				return fmt.Errorf("Error setting policy title: %s", err)
			}

			return nil
		}
	}

	return nil
}

func parse_policies_response(res map[string]interface{}) ([]AccessPolicy, error) {
	var policies []AccessPolicy
	if _, ok := res["accessPolicies"].([]interface{}); !ok {
		// response did not include any policies
		return policies, nil
	}

	for _, res_policy := range res["accessPolicies"].([]interface{}) {
		parsed_policy := &AccessPolicy{}

		err := tpgresource.Convert(res_policy, parsed_policy)
		if err != nil {
			return nil, err
		}

		policies = append(policies, *parsed_policy)
	}
	return policies, nil
}

func compare_scopes(config_scopes []interface{}, policy_scopes []string) bool {
	// converts []interface{} to []string
	var config_scopes_slice []string
	for _, scope := range config_scopes {
		config_scopes_slice = append(config_scopes_slice, scope.(string))
	}

	return slices.Equal(config_scopes_slice, policy_scopes)
}

type AccessPolicy struct {
	Name   string   `json:"name"`
	Title  string   `json:"title"`
	Parent string   `json:"parent"`
	Scopes []string `json:"scopes"`
}
