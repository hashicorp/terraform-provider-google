package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
	"google.golang.org/api/googleapi"
)

func dataSourceGoogleActiveFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleActiveFolderRead,

		Schema: map[string]*schema.Schema{
			"parent": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
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

	queryString := fmt.Sprintf("lifecycleState=ACTIVE AND parent=%s AND displayName=%s", parent, displayName)
	searchRequest := &resourceManagerV2Beta1.SearchFoldersRequest{
		Query: queryString,
	}
	searchResponse, err := config.clientResourceManagerV2Beta1.Folders.Search(searchRequest).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return fmt.Errorf("Folder Not Found : %s", d.Get("name"))
		}

		return fmt.Errorf("Error reading folders: %s", err)
	}

	folders := searchResponse.Folders
	if len(folders) != 1 {
		return fmt.Errorf("More than one folder found")
	}

	d.SetId(folders[0].Name)
	d.Set("name", folders[0].Name)
	return nil
}
