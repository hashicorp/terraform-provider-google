// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanagerregional

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceParameterManagerRegionalRegionalParameters() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceParameterManagerRegionalRegionalParameter().Schema)

	return &schema.Resource{
		Read: dataSourceParameterManagerRegionalRegionalParametersRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter": {
				Type: schema.TypeString,
				Description: `Filter string, adhering to the rules in List-operation filtering. List only parameters matching the filter. 
If filter is empty, all regional parameters are listed from specific location.`,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceParameterManagerRegionalRegionalParametersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ParameterManagerRegionalBasePath}}projects/{{project}}/locations/{{location}}/parameters")
	if err != nil {
		return err
	}

	filter, has_filter := d.GetOk("filter")

	if has_filter {
		url, err = transport_tpg.AddQueryParams(url, map[string]string{"filter": filter.(string)})
		if err != nil {
			return err
		}
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("error fetching project for Regional Parameters: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// To handle the pagination locally
	allParameters := make([]interface{}, 0)
	token := ""
	for paginate := true; paginate; {
		if token != "" {
			url, err = transport_tpg.AddQueryParams(url, map[string]string{"pageToken": token})
			if err != nil {
				return err
			}
		}
		parameters, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ParameterManagerRegionalParameters %q", d.Id()))
		}
		parametersInterface := parameters["parameters"]
		if parametersInterface != nil {
			allParameters = append(allParameters, parametersInterface.([]interface{})...)
		}
		tokenInterface := parameters["nextPageToken"]
		if tokenInterface == nil {
			paginate = false
		} else {
			paginate = true
			token = tokenInterface.(string)
		}
	}

	if err := d.Set("parameters", flattenParameterManagerRegionalRegionalParameterParameters(allParameters, d, config)); err != nil {
		return fmt.Errorf("error setting regional parameters: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}

	if err := d.Set("filter", filter); err != nil {
		return fmt.Errorf("error setting filter: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/parameters")
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	if has_filter {
		id += "/filter=" + filter.(string)
	}
	d.SetId(id)

	return nil
}

func flattenParameterManagerRegionalRegionalParameterParameters(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"format":           flattenParameterManagerRegionalRegionalParameterFormat(original["format"], d, config),
			"labels":           flattenParameterManagerRegionalRegionalParameterEffectiveLabels(original["labels"], d, config),
			"effective_labels": flattenParameterManagerRegionalRegionalParameterEffectiveLabels(original["labels"], d, config),
			"terraform_labels": flattenParameterManagerRegionalRegionalParameterEffectiveLabels(original["labels"], d, config),
			"create_time":      flattenParameterManagerRegionalRegionalParameterCreateTime(original["createTime"], d, config),
			"update_time":      flattenParameterManagerRegionalRegionalParameterUpdateTime(original["updateTime"], d, config),
			"policy_member":    flattenParameterManagerRegionalRegionalParameterPolicyMember(original["policyMember"], d, config),
			"name":             flattenParameterManagerRegionalRegionalParameterName(original["name"], d, config),
			"project":          getDataFromName(original["name"], 1),
			"location":         getDataFromName(original["name"], 3),
			"parameter_id":     getDataFromName(original["name"], 5),
		})
	}
	return transformed
}

func getDataFromName(v interface{}, part int) string {
	name := v.(string)
	split := strings.Split(name, "/")
	return split[part]
}
