// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package containerattached

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleContainerAttachedVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleContainerAttachedVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"valid_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleContainerAttachedVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if len(location) == 0 {
		return fmt.Errorf("Cannot determine location: set location in this data source or at provider-level")
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{ContainerAttachedBasePath}}projects/{{project}}/locations/{{location}}/attachedServerConfig")
	if err != nil {
		return err
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return err
	}
	var validVersions []string
	for _, v := range res["validVersions"].([]interface{}) {
		vm := v.(map[string]interface{})
		validVersions = append(validVersions, vm["version"].(string))
	}
	if err := d.Set("valid_versions", validVersions); err != nil {
		return err
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
