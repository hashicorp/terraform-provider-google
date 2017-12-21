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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	domain, domainOk := d.GetOk("domain")
	name, nameOk := d.GetOk("name")
	if domainOk == nameOk {
		return fmt.Errorf("One of ['domain', 'name'] must be set to read organizations")
	}

	var organization *cloudresourcemanager.Organization
	if domainOk {
		filter := fmt.Sprintf("domain=%s", domain.(string))
		resp, err := config.clientResourceManager.Organizations.Search(&cloudresourcemanager.SearchOrganizationsRequest{
			Filter: filter,
		}).Do()
		if err != nil {
			return fmt.Errorf("Error reading organization: %s", err)
		}

		if len(resp.Organizations) == 0 {
			return fmt.Errorf("Organization not found: %s", domain)
		}
		if len(resp.Organizations) > 1 {
			return fmt.Errorf("More than one matching organization found")
		}

		organization = resp.Organizations[0]
	} else {
		resp, err := config.clientResourceManager.Organizations.Get(name.(string)).Do()
		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
				return fmt.Errorf("Organization not found: %s", name)
			}

			return fmt.Errorf("Error reading organization: %s", err)
		}

		organization = resp
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
