// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	resourceManagerV3 "google.golang.org/api/cloudresourcemanager/v3"
)

func DataSourceGoogleActiveFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleActiveFolderRead,

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleActiveFolderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	var folderMatch *resourceManagerV3.Folder
	parent := d.Get("parent").(string)
	displayName := d.Get("display_name").(string)
	token := ""

	for paginate := true; paginate; {
		resp, err := config.NewResourceManagerV3Client(userAgent).Folders.List().Parent(parent).PageSize(300).PageToken(token).Do()
		if err != nil {
			return fmt.Errorf("error reading folder list: %s", err)
		}

		for _, folder := range resp.Folders {
			if folder.DisplayName == displayName && folder.State == "ACTIVE" {
				if folderMatch != nil {
					return fmt.Errorf("more than one matching folder found")
				}
				folderMatch = folder
			}
		}
		token = resp.NextPageToken
		paginate = token != ""
	}

	if folderMatch == nil {
		return fmt.Errorf("folder not found: %s", displayName)
	}

	d.SetId(folderMatch.Name)
	if err := d.Set("name", folderMatch.Name); err != nil {
		return fmt.Errorf("Error setting folder name: %s", err)
	}

	return nil
}
