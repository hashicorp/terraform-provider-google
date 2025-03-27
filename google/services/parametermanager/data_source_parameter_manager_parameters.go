// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanager

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceParameterManagerParameters() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceParameterManagerParameter().Schema)

	return &schema.Resource{
		Read: dataSourceParameterManagerParametersRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"filter": {
				Type: schema.TypeString,
				Description: `Filter string, adhering to the rules in List-operation filtering. List only parameters matching the filter. 
If filter is empty, all parameters are listed.`,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
		},
	}
}

func dataSourceParameterManagerParametersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ParameterManagerBasePath}}projects/{{project}}/locations/global/parameters")
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
		return fmt.Errorf("error fetching project for Parameters: %s", err)
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
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("ParameterManagerParameters %q", d.Id()))
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

	if err := d.Set("parameters", flattenParameterManagerParameterParameters(allParameters, d, config)); err != nil {
		return fmt.Errorf("error setting parameters: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}

	if err := d.Set("filter", filter); err != nil {
		return fmt.Errorf("error setting filter: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/global/parameters")
	if err != nil {
		return fmt.Errorf("error constructing id: %s", err)
	}
	if has_filter {
		id += "/filter=" + filter.(string)
	}
	d.SetId(id)

	return nil
}

func flattenParameterManagerParameterParameters(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"format":           flattenParameterManagerParameterFormat(original["format"], d, config),
			"labels":           flattenParameterManagerParameterEffectiveLabels(original["labels"], d, config),
			"effective_labels": flattenParameterManagerParameterEffectiveLabels(original["labels"], d, config),
			"terraform_labels": flattenParameterManagerParameterEffectiveLabels(original["labels"], d, config),
			"create_time":      flattenParameterManagerParameterCreateTime(original["createTime"], d, config),
			"update_time":      flattenParameterManagerParameterUpdateTime(original["updateTime"], d, config),
			"policy_member":    flattenParameterManagerParameterPolicyMember(original["policyMember"], d, config),
			"name":             flattenParameterManagerParameterName(original["name"], d, config),
			"kms_key":          flattenParameterManagerParameterKmsKey(original["kmsKey"], d, config),
			"project":          getDataFromName(original["name"], 1),
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
