// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpuv2

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceTpuV2AcceleratorTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTpuV2AcceleratorTypesRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTpuV2AcceleratorTypesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{TpuV2BasePath}}projects/{{project}}/locations/{{zone}}/acceleratorTypes")
	if err != nil {
		return err
	}

	typesRaw, err := tpgresource.PaginatedListRequest(project, url, userAgent, config, flattenTpuV2AcceleratorTypes)
	if err != nil {
		return fmt.Errorf("error listing TPU v2 accelerator types: %s", err)
	}

	types := make([]string, len(typesRaw))
	for i, typeRaw := range typesRaw {
		types[i] = typeRaw.(string)
	}
	sort.Strings(types)

	log.Printf("[DEBUG] Received Google TPU v2 accelerator types: %q", types)

	if err := d.Set("types", types); err != nil {
		return fmt.Errorf("error setting types: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("error setting zone: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/zones/%s", project, zone))

	return nil
}

func flattenTpuV2AcceleratorTypes(resp map[string]interface{}) []interface{} {
	typeObjList := resp["acceleratorTypes"].([]interface{})
	types := make([]interface{}, len(typeObjList))
	for i, typ := range typeObjList {
		typeObj := typ.(map[string]interface{})
		types[i] = typeObj["type"]
	}
	return types
}
