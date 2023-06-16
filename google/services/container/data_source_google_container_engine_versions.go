// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleContainerEngineVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleContainerEngineVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_master_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_node_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_master_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"valid_node_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"release_channel_default_version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"release_channel_latest_version": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGoogleContainerEngineVersionsRead(d *schema.ResourceData, meta interface{}) error {
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

	location = fmt.Sprintf("projects/%s/locations/%s", project, location)
	resp, err := config.NewContainerClient(userAgent).Projects.Locations.GetServerConfig(location).Do()
	if err != nil {
		return fmt.Errorf("Error retrieving available container cluster versions: %s", err.Error())
	}

	validMasterVersions := make([]string, 0)
	for _, v := range resp.ValidMasterVersions {
		if strings.HasPrefix(v, d.Get("version_prefix").(string)) {
			validMasterVersions = append(validMasterVersions, v)
		}
	}

	validNodeVersions := make([]string, 0)
	for _, v := range resp.ValidNodeVersions {
		if strings.HasPrefix(v, d.Get("version_prefix").(string)) {
			validNodeVersions = append(validNodeVersions, v)
		}
	}

	if err := d.Set("valid_master_versions", validMasterVersions); err != nil {
		return fmt.Errorf("Error setting valid_master_versions: %s", err)
	}
	if len(validMasterVersions) > 0 {
		if err := d.Set("latest_master_version", validMasterVersions[0]); err != nil {
			return fmt.Errorf("Error setting latest_master_version: %s", err)
		}
	}

	if err := d.Set("valid_node_versions", validNodeVersions); err != nil {
		return fmt.Errorf("Error setting valid_node_versions: %s", err)
	}
	if len(validNodeVersions) > 0 {
		if err := d.Set("latest_node_version", validNodeVersions[0]); err != nil {
			return fmt.Errorf("Error setting latest_node_version: %s", err)
		}
	}

	if err := d.Set("default_cluster_version", resp.DefaultClusterVersion); err != nil {
		return fmt.Errorf("Error setting default_cluster_version: %s", err)
	}

	releaseChannelDefaultVersion := map[string]string{}
	releaseChannelLatestVersion := map[string]string{}
	for _, channelResp := range resp.Channels {
		releaseChannelDefaultVersion[channelResp.Channel] = channelResp.DefaultVersion
		for _, v := range channelResp.ValidVersions {
			if strings.HasPrefix(v, d.Get("version_prefix").(string)) {
				releaseChannelLatestVersion[channelResp.Channel] = v
				break
			}
		}
	}

	if err := d.Set("release_channel_default_version", releaseChannelDefaultVersion); err != nil {
		return fmt.Errorf("Error setting release_channel_default_version: %s", err)
	}
	if err := d.Set("release_channel_latest_version", releaseChannelLatestVersion); err != nil {
		return fmt.Errorf("Error setting release_channel_latest_version: %s", err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
