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

func DataSourceTpuV2RuntimeVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTpuV2RuntimeVersionsRead,
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
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTpuV2RuntimeVersionsRead(d *schema.ResourceData, meta interface{}) error {
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

	url, err := tpgresource.ReplaceVars(d, config, "{{TpuV2BasePath}}projects/{{project}}/locations/{{zone}}/runtimeVersions")
	if err != nil {
		return err
	}

	versionsRaw, err := tpgresource.PaginatedListRequest(project, url, userAgent, config, flattenTpuV2RuntimeVersions)
	if err != nil {
		return fmt.Errorf("error listing TPU v2 runtime versions: %s", err)
	}

	versions := make([]string, len(versionsRaw))
	for i, ver := range versionsRaw {
		versions[i] = ver.(string)
	}
	sort.Strings(versions)

	log.Printf("[DEBUG] Received Google TPU v2 runtime versions: %q", versions)

	if err := d.Set("versions", versions); err != nil {
		return fmt.Errorf("error setting versions: %s", err)
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

func flattenTpuV2RuntimeVersions(resp map[string]interface{}) []interface{} {
	verObjList := resp["runtimeVersions"].([]interface{})
	versions := make([]interface{}, len(verObjList))
	for i, v := range verObjList {
		verObj := v.(map[string]interface{})
		versions[i] = verObj["version"]
	}
	return versions
}
