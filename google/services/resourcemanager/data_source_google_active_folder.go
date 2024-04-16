// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
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
			"api_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Provides the REST method through which to find the folder. LIST is recommended as it is strongly consistent.",
				Default:      "LIST",
				ValidateFunc: verify.ValidateEnum([]string{"LIST", "SEARCH"}),
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
	apiMethod := d.Get("api_method").(string)

	if apiMethod == "LIST" {
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
	} else {
		queryString := fmt.Sprintf("lifecycleState=ACTIVE AND parent=%s AND displayName=\"%s\"", parent, displayName)
		searchRequest := config.NewResourceManagerV3Client(userAgent).Folders.Search()
		searchRequest.Query(queryString)
		searchResponse, err := searchRequest.Do()
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Folder Not Found : %s", displayName))
		}

		for _, folder := range searchResponse.Folders {
			if folder.DisplayName == displayName {
				folderMatch = folder
				break
			}
		}
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
