package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

func dataSourceGoogleFolder() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFolderRead,
		Schema: map[string]*schema.Schema{
			"folder": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lookup_organization": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"organization": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFolderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	folderName := d.Get("folder").(string)

	folder, err := config.clientResourceManagerV2Beta1.Folders.Get(canonicalFolderName(folderName)).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Folder Not Found : %s", folderName))
	}

	d.SetId(GetResourceNameFromSelfLink(folder.Name))
	d.Set("name", folder.Name)
	d.Set("parent", folder.Parent)
	d.Set("display_name", folder.DisplayName)
	d.Set("lifecycle_state", folder.LifecycleState)
	d.Set("create_time", folder.CreateTime)

	if v, ok := d.GetOk("lookup_organization"); ok && v.(bool) {
		organization, err := lookupOrganizationName(folder, config)

		if err != nil {
			return err
		}

		d.Set("organization", organization)
	}

	return nil
}

func canonicalFolderName(ba string) string {
	if strings.HasPrefix(ba, "folders/") {
		return ba
	}

	return "folders/" + ba
}

func lookupOrganizationName(folder *resourceManagerV2Beta1.Folder, config *Config) (string, error) {
	parent := folder.Parent

	if parent == "" || strings.HasPrefix(parent, "organizations/") {
		return parent, nil
	} else if strings.HasPrefix(parent, "folders/") {
		parentFolder, err := config.clientResourceManagerV2Beta1.Folders.Get(parent).Do()

		if err != nil {
			return "", fmt.Errorf("Error getting parent folder '%s': %s", parent, err)
		}

		return lookupOrganizationName(parentFolder, config)
	} else {
		return "", fmt.Errorf("Unknown parent type '%s' on folder '%s'", parent, folder.Name)
	}
}
