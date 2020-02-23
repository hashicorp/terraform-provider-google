package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func dataSourceGoogleOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"organization"},
			},
			"organization": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"domain"},
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
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

	var organization *cloudresourcemanager.Organization
	if v, ok := d.GetOk("domain"); ok {
		filter := fmt.Sprintf("domain=%s", v.(string))
		var resp *cloudresourcemanager.SearchOrganizationsResponse
		err := retryTimeDuration(func() (err error) {
			resp, err = config.clientResourceManager.Organizations.Search(&cloudresourcemanager.SearchOrganizationsRequest{
				Filter: filter,
			}).Do()
			return err
		}, d.Timeout(schema.TimeoutRead))
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
	} else if v, ok := d.GetOk("organization"); ok {
		var resp *cloudresourcemanager.Organization
		err := retryTimeDuration(func() (err error) {
			resp, err = config.clientResourceManager.Organizations.Get(canonicalOrganizationName(v.(string))).Do()
			return err
		}, d.Timeout(schema.TimeoutRead))
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Organization Not Found : %s", v))
		}

		organization = resp
	} else {
		return fmt.Errorf("one of domain or organization must be set")
	}

	d.SetId(organization.Name)
	d.Set("name", organization.Name)
	d.Set("org_id", GetResourceNameFromSelfLink(organization.Name))
	d.Set("domain", organization.DisplayName)
	d.Set("create_time", organization.CreationTime)
	d.Set("lifecycle_state", organization.LifecycleState)
	if organization.Owner != nil {
		d.Set("directory_customer_id", organization.Owner.DirectoryCustomerId)
	}

	return nil
}

func canonicalOrganizationName(ba string) string {
	if strings.HasPrefix(ba, "organizations/") {
		return ba
	}

	return "organizations/" + ba
}
