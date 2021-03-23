package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	resourceManagerV2 "google.golang.org/api/cloudresourcemanager/v2"
)

func dataSourceGoogleActiveFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleActiveFolderRead,

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Optional: true,
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	queryString := ""

	// parent is optional
	if parent, ok := d.GetOk("parent"); ok {
		queryString = fmt.Sprintf("lifecycleState=ACTIVE AND parent=%s AND displayName=\"%s\"", parent.(string), displayName)
	} else {
		queryString = fmt.Sprintf("lifecycleState=ACTIVE AND displayName=\"%s\"", displayName)
	}

	searchRequest := &resourceManagerV2.SearchFoldersRequest{
		Query: queryString,
	}
	searchResponse, err := config.NewResourceManagerV2Client(userAgent).Folders.Search(searchRequest).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Folder Not Found : %s", displayName))
	}

	for _, folder := range searchResponse.Folders {
		if folder.DisplayName == displayName {
			d.SetId(folder.Name)
			if err := d.Set("name", folder.Name); err != nil {
				return fmt.Errorf("Error setting folder name: %s", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Folder not found")
}
