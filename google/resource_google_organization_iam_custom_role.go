package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/iam/v1"
)

func resourceGoogleOrganizationIamCustomRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleOrganizationIamCustomRoleCreate,
		Read:   resourceGoogleOrganizationIamCustomRoleRead,
		Update: resourceGoogleOrganizationIamCustomRoleUpdate,
		Delete: resourceGoogleOrganizationIamCustomRoleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"stage": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "GA",
				ValidateFunc: validation.StringInSlice([]string{"ALPHA", "BETA", "GA", "DEPRECATED", "DISABLED", "EAP"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deleted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceGoogleOrganizationIamCustomRoleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.Get("deleted").(bool) {
		return fmt.Errorf("Cannot create a custom organization role with a deleted state. `deleted` field should be false.")
	}

	org := d.Get("org_id").(string)
	roleId := fmt.Sprintf("organizations/%s/roles/%s", org, d.Get("role_id").(string))
	r, err := config.clientIAM.Organizations.Roles.Get(roleId).Do()
	if err == nil {
		if r.Deleted {
			// Roles have soft deletes - creating a role with the same name
			// as a recently deleted role must instead be undelete/update.
			d.SetId(r.Name)
			return resourceGoogleOrganizationIamCustomRoleUpdate(d, meta)
		}
		// If old role with same name exists, just return error
		return fmt.Errorf("Custom project role %s already exists and must be imported", roleId)
	}

	// If no role is found, actually create a new role.
	role, err := config.clientIAM.Organizations.Roles.Create("organizations/"+d.Get("org_id").(string), &iam.CreateRoleRequest{
		RoleId: d.Get("role_id").(string),
		Role: &iam.Role{
			Title:               d.Get("title").(string),
			Description:         d.Get("description").(string),
			Stage:               d.Get("stage").(string),
			IncludedPermissions: convertStringSet(d.Get("permissions").(*schema.Set)),
		},
	}).Do()

	if err != nil {
		return fmt.Errorf("Error creating the custom organization role %s: %s", d.Get("title").(string), err)
	}

	d.SetId(role.Name)

	return resourceGoogleOrganizationIamCustomRoleRead(d, meta)
}

func resourceGoogleOrganizationIamCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	role, err := config.clientIAM.Organizations.Roles.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Id())
	}

	parsedRoleName, err := ParseOrganizationCustomRoleName(role.Name)
	if err != nil {
		return err
	}

	d.Set("role_id", parsedRoleName.Name)
	d.Set("org_id", parsedRoleName.OrgId)
	d.Set("title", role.Title)
	d.Set("description", role.Description)
	d.Set("permissions", role.IncludedPermissions)
	d.Set("stage", role.Stage)
	d.Set("deleted", role.Deleted)

	return nil
}

func resourceGoogleOrganizationIamCustomRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	if d.HasChange("deleted") {
		if d.Get("deleted").(bool) {
			if err := resourceGoogleOrganizationIamCustomRoleDelete(d, meta); err != nil {
				return err
			}
			d.SetPartial("deleted")
		}
	}

	// If role is not deleted, make sure it exists and undelete if needed.
	// TODO(emilymye): Change deleted to computed so it can't be updated.
	if !d.Get("deleted").(bool) {
		r, err := config.clientIAM.Organizations.Roles.Get(d.Id()).Do()
		if err != nil {
			return fmt.Errorf("unable to find custom project role %s to update: %v", d.Id(), err)
		}
		if r.Deleted {
			// Undelete if deleted previously
			if err := resourceGoogleOrganizationIamCustomRoleUndelete(d, meta); err != nil {
				return err
			}
			d.SetPartial("deleted")
		}
	}

	if d.HasChange("title") || d.HasChange("description") || d.HasChange("stage") || d.HasChange("permissions") {
		_, err := config.clientIAM.Organizations.Roles.Patch(d.Id(), &iam.Role{
			Title:               d.Get("title").(string),
			Description:         d.Get("description").(string),
			Stage:               d.Get("stage").(string),
			IncludedPermissions: convertStringSet(d.Get("permissions").(*schema.Set)),
		}).Do()

		if err != nil {
			return fmt.Errorf("Error updating the custom organization role %s: %s", d.Get("title").(string), err)
		}
		d.SetPartial("title")
		d.SetPartial("description")
		d.SetPartial("stage")
		d.SetPartial("permissions")
	}

	d.Partial(false)

	return nil
}

func resourceGoogleOrganizationIamCustomRoleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientIAM.Organizations.Roles.Delete(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting the custom organization role %s: %s", d.Get("title").(string), err)
	}

	return nil
}

func resourceGoogleOrganizationIamCustomRoleUndelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientIAM.Organizations.Roles.Undelete(d.Id(), &iam.UndeleteRoleRequest{}).Do()
	if err != nil {
		return fmt.Errorf("Error undeleting the custom organization role %s: %s", d.Get("title").(string), err)
	}

	return nil
}
