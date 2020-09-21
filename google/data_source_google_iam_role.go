package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleIamRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleIamRoleRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"included_permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"stage": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleIamRoleRead(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}
	config := meta.(*Config)
	config.clientIAM.UserAgent = fmt.Sprintf("%s %s", config.clientIAM.UserAgent, m.ModuleName)

	roleName := d.Get("name").(string)
	role, err := config.clientIAM.Roles.Get(roleName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Error reading IAM Role %s: %s", roleName, err))
	}

	d.SetId(role.Name)
	if err := d.Set("title", role.Title); err != nil {
		return fmt.Errorf("Error setting title: %s", err)
	}
	if err := d.Set("stage", role.Stage); err != nil {
		return fmt.Errorf("Error setting stage: %s", err)
	}
	if err := d.Set("included_permissions", role.IncludedPermissions); err != nil {
		return fmt.Errorf("Error setting included_permissions: %s", err)
	}

	return nil
}
