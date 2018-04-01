package google

import (
	"bytes"
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
				Optional: true,
				ConflictsWith: []string{
					"name",
					"parent",
					"display_name",
					"lifecycle_state",
					"create_time",
				},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"folder"},
			},
			"parent": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"folder"},
			},
			"display_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"folder"},
			},
			"lifecycle_state": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"folder"},
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

	var folder *resourceManagerV2Beta1.Folder

	if queryString, ok := generateFolderReadQueryString(d); ok {
		searchRequest := &resourceManagerV2Beta1.SearchFoldersRequest{
			Query: queryString,
		}
		searchResponse, err := config.clientResourceManagerV2Beta1.Folders.Search(searchRequest).Do()
		if err != nil || (err == nil && len(searchResponse.Folders) == 0) {
			return handleNotFoundError(err, d, fmt.Sprintf("Folder Not Found With Query : %s", queryString))
		}

		folders := searchResponse.Folders
		if len(folders) > 1 {
			return fmt.Errorf("More than one folder found")
		}

		folder = folders[0]
	} else if v, ok := d.GetOk("folder"); ok {
		resp, err := config.clientResourceManagerV2Beta1.Folders.Get(canonicalFolderName(v.(string))).Do()

		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Folder Not Found : %s", v))
		}

		folder = resp
	} else {
		return fmt.Errorf("at least one of folder, name, parent, display_name, lifecycle_state must be set")
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

func generateFolderReadQueryString(d *schema.ResourceData) (string, bool) {
	var buffer bytes.Buffer

	conditionals := map[string]string{
		"name":            "name",
		"parent":          "parent",
		"display_name":    "displayName",
		"lifecycle_state": "lifecycleState",
	}

	firstConditional := true

	for conditionalInput, conditionalQueryColumn := range conditionals {
		if v, ok := d.GetOk(conditionalInput); ok {
			if !firstConditional {
				buffer.WriteString(" AND ")
			}

			conditional := fmt.Sprintf("%s=%s", conditionalQueryColumn, v.(string))

			buffer.WriteString(conditional)

			if firstConditional {
				firstConditional = false
			}
		}
	}

	queryString := buffer.String()

	ok := len(queryString) > 0

	return queryString, ok
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
