package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The role id to use for this role.`,
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The numeric ID of the organization in which you want to create a custom role.`,
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `A human-readable title for the role.`,
			},
			"permissions": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: `The names of the permissions this role grants when bound in an IAM policy. At least one permission must be specified.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"stage": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "GA",
				Description:      `The current launch stage of the role. Defaults to GA.`,
				ValidateFunc:     validation.StringInSlice([]string{"ALPHA", "BETA", "GA", "DEPRECATED", "DISABLED", "EAP"}, false),
				DiffSuppressFunc: emptyOrDefaultStringSuppress("ALPHA"),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A human-readable description for the role.`,
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The current deleted state of the role.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the role in the format organizations/{{org_id}}/roles/{{role_id}}. Like id, this field can be used as a reference in other resources such as IAM role bindings.`,
			},
		},
	}
}

func resourceGoogleOrganizationIamCustomRoleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	org := d.Get("org_id").(string)
	roleId := fmt.Sprintf("organizations/%s/roles/%s", org, d.Get("role_id").(string))
	orgId := fmt.Sprintf("organizations/%s", org)

	// Look for role with given ID.
	// If it exists in deleted state, update to match "created" role state
	// If it exists and and is enabled, return error - we should not try to recreate.
	r, err := config.clientIAM.Organizations.Roles.Get(roleId).Do()
	if err == nil {
		if r.Deleted {
			// This role was soft-deleted; update to match new state.
			d.SetId(r.Name)
			if err := resourceGoogleOrganizationIamCustomRoleUpdate(d, meta); err != nil {
				// If update failed, make sure it wasn't actually added to state.
				d.SetId("")
				return err
			}
		} else {
			// If a role with same name exists and is enabled, just return error
			return fmt.Errorf("Custom project role %s already exists and must be imported", roleId)
		}
	} else if err := handleNotFoundError(err, d, fmt.Sprintf("Custom Organization Role %q", roleId)); err == nil {
		// If no role was found, actually create a new role.
		role, err := config.clientIAM.Organizations.Roles.Create(orgId, &iam.CreateRoleRequest{
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
	} else {
		return fmt.Errorf("Unable to verify whether custom org role %s already exists and must be undeleted: %v", roleId, err)
	}

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
	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("permissions", role.IncludedPermissions)
	d.Set("stage", role.Stage)
	d.Set("deleted", role.Deleted)

	return nil
}

func resourceGoogleOrganizationIamCustomRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	// We want to update the role to some undeleted state.
	// Make sure the role with given ID exists and is un-deleted before patching.
	r, err := config.clientIAM.Organizations.Roles.Get(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("unable to find custom project role %s to update: %v", d.Id(), err)
	}

	if r.Deleted {
		_, err := config.clientIAM.Organizations.Roles.Undelete(d.Id(), &iam.UndeleteRoleRequest{}).Do()
		if err != nil {
			return fmt.Errorf("Error undeleting the custom organization role %s: %s", d.Get("title").(string), err)
		}

		d.SetPartial("deleted")
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

	r, err := config.clientIAM.Organizations.Roles.Get(d.Id()).Do()
	if err == nil && r != nil && r.Deleted && d.Get("deleted").(bool) {
		// This role has already been deleted, don't try again.
		return nil
	}

	_, err = config.clientIAM.Organizations.Roles.Delete(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting the custom organization role %s: %s", d.Get("title").(string), err)
	}

	return nil
}
