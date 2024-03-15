// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudquotas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudQuotasQuotaInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCloudQuotasQuotaInfoRead,

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"quota_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_precise": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"refresh_interval": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dimensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"metric_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_increase_eligibility": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_eligible": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ineligibility_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"is_fixed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dimensions_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dimensions": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"applicable_locations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"is_concurrent": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"service_request_quota_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		UseJSONNumber: true,
	}
}

func dataSourceGoogleCloudQuotasQuotaInfoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{CloudQuotasBasePath}}{{parent}}/locations/global/services/{{service}}/quotaInfos/{{quota_id}}")
	if err != nil {
		return fmt.Errorf("error setting api endpoint")
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: userAgent,
	})

	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("CloudQuotasQuotaInfo %q", d.Id()))
	}

	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo name: %s", err)
	}
	if err := d.Set("quota_id", res["quotaId"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo quota_id: %s", err)
	}
	if err := d.Set("metric", res["metric"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo metric: %s", err)
	}
	if err := d.Set("service", res["service"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo service: %s", err)
	}
	if err := d.Set("is_precise", res["isPrecise"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo is_precise: %s", err)
	}
	if err := d.Set("refresh_interval", res["refreshInterval"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo refresh_interval: %s", err)
	}
	if err := d.Set("container_type", res["containerType"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo container_type: %s", err)
	}
	if err := d.Set("dimensions", res["dimensions"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo dimensions: %s", err)
	}
	if err := d.Set("metric_display_name", res["metricDisplayName"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo metric_display_name: %s", err)
	}
	if err := d.Set("quota_display_name", res["quotaDisplayName"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo quota_display_name: %s", err)
	}
	if err := d.Set("metric_unit", res["metricUnit"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo metric_unit: %s", err)
	}
	if err := d.Set("quota_increase_eligibility", flattenCloudQuotasQuotaInfoQuotaIncreaseEligibility(res["quotaIncreaseEligibility"], d, config)); err != nil {
		return fmt.Errorf("error reading QuotaInfo quota_increase_eligibility: %s", err)
	}
	if err := d.Set("is_fixed", res["isFixed"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo is_fixed: %s", err)
	}
	if err := d.Set("dimensions_infos", flattenCloudQuotasQuotaInfoDimensionsInfos(res["dimensionsInfos"], d, config)); err != nil {
		return fmt.Errorf("error reading QuotaInfo dimensions_infos: %s", err)
	}
	if err := d.Set("is_concurrent", res["isConcurrent"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo is_concurrent: %s", err)
	}
	if err := d.Set("service_request_quota_uri", res["serviceRequestQuotaUri"]); err != nil {
		return fmt.Errorf("error reading QuotaInfo service_request_quota_uri: %s", err)
	}

	d.SetId(res["name"].(string))
	return nil
}

func flattenCloudQuotasQuotaInfoQuotaIncreaseEligibility(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["is_eligible"] = original["is_eligible"]
	transformed["ineligibility_reason"] = original["ineligibility_reason"]
	return []interface{}{transformed}
}

func flattenCloudQuotasQuotaInfoDimensionsInfos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) []interface{} {
	if v == nil {
		return make([]interface{}, 0)
	}

	original := v.([]interface{})
	dimensionsInfos := make([]interface{}, 0, len(original))

	for _, raw := range original {
		data := make(map[string]interface{})
		data["details"] = flattenCloudQuotasQuotaInfoDetails(raw.(map[string]interface{})["details"], d, config)
		data["applicable_locations"] = raw.(map[string]interface{})["applicableLocations"]
		data["dimensions"] = raw.(map[string]interface{})["dimensions"]

		dimensionsInfos = append(dimensionsInfos, data)
	}
	return dimensionsInfos
}

func flattenCloudQuotasQuotaInfoDetails(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	original, ok := v.(map[string]interface{})
	if !ok || len(original) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{"value": original["value"]},
	}
}
