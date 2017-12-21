package google

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
)

func dataSourceGoogleOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"domain"},
			},
			"directory_customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var organization *cloudresourcemanager.Organization
	if v, ok := d.GetOk("domain"); ok {
		filter := fmt.Sprintf("domain=%s", v.(string))
		resp, err := config.clientResourceManager.Organizations.Search(&cloudresourcemanager.SearchOrganizationsRequest{
			Filter: filter,
		}).Do()
		if err != nil {
			return fmt.Errorf("Error reading organization: %s", err)
		}

		if len(resp.Organizations) == 0 {
			return fmt.Errorf("Organization not found: %s", v)
		}
		if len(resp.Organizations) > 1 {
			return fmt.Errorf("More than one matching organization found")
		}

		organization = resp.Organizations[0]
	} else if v, ok := d.GetOk("name"); ok {
		resp, err := config.clientResourceManager.Organizations.Get(v.(string)).Do()
		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
				return fmt.Errorf("Organization not found: %s", v)
			}

			return fmt.Errorf("Error reading organization: %s", err)
		}

		organization = resp
	} else {
		return fmt.Errorf("one of domain or name must be set")
	}

	d.SetId(GetResourceNameFromSelfLink(organization.Name))
	d.Set("name", organization.Name)
	d.Set("domain", organization.DisplayName)
	d.Set("create_time", organization.CreationTime)
	d.Set("lifecycle_state", organization.LifecycleState)
	if organization.Owner != nil {
		d.Set("directory_customer_id", organization.Owner.DirectoryCustomerId)
	}

	return nil
}
