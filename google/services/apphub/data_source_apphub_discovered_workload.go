// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apphub

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceApphubDiscoveredWorkload() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApphubDiscoveredWorkloadRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workload_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workload_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"workload_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gcp_project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApphubDiscoveredWorkloadRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{ApphubBasePath}}projects/{{project}}/locations/{{location}}/discoveredWorkloads:lookup?uri={{workload_uri}}"))
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
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("ApphubDiscoveredWorkload %q", d.Id()), url)
	}

	if err := d.Set("name", flattenApphubDiscoveredWorkloadName(res["discoveredWorkload"].(map[string]interface{})["name"], d, config)); err != nil {
		return fmt.Errorf("Error setting workload name: %s", err)
	}

	if err := d.Set("workload_reference", flattenApphubDiscoveredWorkloadReference(res["discoveredWorkload"].(map[string]interface{})["workloadReference"], d, config)); err != nil {
		return fmt.Errorf("Error setting service reference: %s", err)
	}

	if err := d.Set("workload_properties", flattenApphubDiscoveredWorkloadProperties(res["discoveredWorkload"].(map[string]interface{})["workloadProperties"], d, config)); err != nil {
		return fmt.Errorf("Error setting workload properties: %s", err)
	}

	d.SetId(res["discoveredWorkload"].(map[string]interface{})["name"].(string))

	return nil

}

func flattenApphubDiscoveredWorkloadReference(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["uri"] = flattenApphubDiscoveredWorkloadDataUri(original["uri"], d, config)
	return []interface{}{transformed}
}

func flattenApphubDiscoveredWorkloadProperties(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["gcp_project"] = flattenApphubDiscoveredWorkloadDataGcpProject(original["gcpProject"], d, config)
	transformed["location"] = flattenApphubDiscoveredWorkloadDataLocation(original["location"], d, config)
	transformed["zone"] = flattenApphubDiscoveredWorkloadDataZone(original["zone"], d, config)
	return []interface{}{transformed}
}

func flattenApphubDiscoveredWorkloadName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredWorkloadDataUri(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredWorkloadDataGcpProject(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredWorkloadDataLocation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenApphubDiscoveredWorkloadDataZone(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
