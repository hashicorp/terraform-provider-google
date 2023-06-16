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

func DataSourceGoogleContainerAttachedInstallManifest() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleContainerAttachedInstallManifestRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"manifest": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleContainerAttachedInstallManifestRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	clusterId := d.Get("cluster_id").(string)
	platformVersion := d.Get("platform_version").(string)

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

	url, err := tpgresource.ReplaceVars(d, config, "{{ContainerAttachedBasePath}}projects/{{project}}/locations/{{location}}:generateAttachedClusterInstallManifest")
	if err != nil {
		return err
	}
	params := map[string]string{
		"attached_cluster_id": clusterId,
		"platform_version":    platformVersion,
	}
	url, err = transport_tpg.AddQueryParams(url, params)
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

	if err := d.Set("manifest", res["manifest"]); err != nil {
		return fmt.Errorf("Error setting manifest: %s", err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
