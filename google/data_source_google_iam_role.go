package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	config := meta.(*Config)
	roleName := d.Get("name").(string)
	role, err := config.clientIAM.Roles.Get(roleName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Error reading IAM Role %s: %s", roleName, err))
	}

	d.SetId(role.Name)
	d.Set("title", role.Title)
	d.Set("stage", role.Stage)
	d.Set("included_permissions", role.IncludedPermissions)

	return nil
}
