package google

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

func dataSourceGoogleActiveFolder() *schema.Resource {
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
	config := meta.(*Config)

	parent := d.Get("parent").(string)
	displayName := d.Get("display_name").(string)

	queryString := fmt.Sprintf("lifecycleState=ACTIVE AND parent=%s AND displayName=%s", parent, url.QueryEscape(displayName))
	searchRequest := &resourceManagerV2Beta1.SearchFoldersRequest{
		Query: queryString,
	}
	searchResponse, err := config.clientResourceManagerV2Beta1.Folders.Search(searchRequest).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Folder Not Found : %s", displayName))
	}

	for _, folder := range searchResponse.Folders {
		if folder.DisplayName == displayName {
			d.SetId(folder.Name)
			d.Set("name", folder.Name)
			return nil
		}
	}
	return fmt.Errorf("Folder not found")
}
