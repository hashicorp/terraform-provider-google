// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	compute "google.golang.org/api/compute/v1"
)

func DataSourceGoogleComputeInstanceGuestAttributes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeInstanceGuestAttributesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"query_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"variable_key"},
			},

			"variable_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"query_path"},
			},

			"variable_value": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"query_value": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeInstanceGuestAttributesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	project, zone, name, err := tpgresource.GetZonalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, name)
	instanceGuestAttributes := &compute.GuestAttributes{}

	// You can either query based on variable_key, query_path or just get the first value
	if d.Get("query_path").(string) != "" {
		instanceGuestAttributes, err = config.NewComputeClient(userAgent).Instances.GetGuestAttributes(project, zone, name).QueryPath(d.Get("query_path").(string)).Do()
	} else if d.Get("variable_key").(string) != "" {
		instanceGuestAttributes, err = config.NewComputeClient(userAgent).Instances.GetGuestAttributes(project, zone, name).VariableKey(d.Get("variable_key").(string)).Do()
	} else {
		instanceGuestAttributes, err = config.NewComputeClient(userAgent).Instances.GetGuestAttributes(project, zone, name).Do()
	}
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("Instance's Guest Attributes %s", name), id)
	}

	// Set query results
	if err := d.Set("variable_value", instanceGuestAttributes.VariableValue); err != nil {
		return fmt.Errorf("Error variable_value: %s", err)
	}
	if err := d.Set("query_value", flattenQueryValues(instanceGuestAttributes.QueryValue)); err != nil {
		return fmt.Errorf("Error query_value: %s", err)
	}

	d.SetId(fmt.Sprintf(instanceGuestAttributes.SelfLink))
	return nil
}

func flattenQueryValues(queryValue *compute.GuestAttributesValue) []map[string]interface{} {
	if queryValue == nil {
		return nil
	}
	queryValueItems := make([]map[string]interface{}, 0)
	for _, item := range queryValue.Items {
		queryValueItems = append(queryValueItems, map[string]interface{}{
			"key":       item.Key,
			"namespace": item.Namespace,
			"value":     item.Value,
		})
	}
	return queryValueItems
}
