package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/iam/v1"
)

func resourceGoogleProjectIamCustomRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectIamCustomRoleCreate,
		Read:   resourceGoogleProjectIamCustomRoleRead,
		Update: resourceGoogleProjectIamCustomRoleUpdate,
		Delete: resourceGoogleProjectIamCustomRoleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_id": {
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
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stage": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "GA",
				ValidateFunc:     validation.StringInSlice([]string{"ALPHA", "BETA", "GA", "DEPRECATED", "DISABLED", "EAP"}, false),
				DiffSuppressFunc: emptyOrDefaultStringSuppress("ALPHA"),
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

func resourceGoogleProjectIamCustomRoleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	if d.Get("deleted").(bool) {
		return fmt.Errorf("Cannot create a custom project role with a deleted state. `deleted` field should be false.")
	}

	roleId := fmt.Sprintf("projects/%s/roles/%s", project, d.Get("role_id").(string))
	r, err := config.clientIAM.Projects.Roles.Get(roleId).Do()
	if err == nil {
		if r.Deleted {
			// This role was soft-deleted; update to match new state.
			d.SetId(r.Name)
			if err := resourceGoogleProjectIamCustomRoleUpdate(d, meta); err != nil {
				// If update failed, make sure it wasn't actually added to state.
				d.SetId("")
				return err
			}
		} else {
			// If a role with same name exists and is enabled, just return error
			return fmt.Errorf("Custom project role %s already exists and must be imported", roleId)
		}
	} else if err := handleNotFoundError(err, d, fmt.Sprintf("Custom Project Role %q", roleId)); err == nil {
		// If no role is found, actually create a new role.
		role, err := config.clientIAM.Projects.Roles.Create("projects/"+project, &iam.CreateRoleRequest{
			RoleId: d.Get("role_id").(string),
			Role: &iam.Role{
				Title:               d.Get("title").(string),
				Description:         d.Get("description").(string),
				Stage:               d.Get("stage").(string),
				IncludedPermissions: convertStringSet(d.Get("permissions").(*schema.Set)),
			},
		}).Do()
		if err != nil {
			return fmt.Errorf("Error creating the custom project role %s: %v", roleId, err)
		}

		d.SetId(role.Name)
	} else {
		return fmt.Errorf("Unable to verify whether custom project role %s already exists and must be undeleted: %v", roleId, err)
	}

	return resourceGoogleProjectIamCustomRoleRead(d, meta)
}

func resourceGoogleProjectIamCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	role, err := config.clientIAM.Projects.Roles.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Id())
	}

	d.Set("role_id", GetResourceNameFromSelfLink(role.Name))
	d.Set("title", role.Title)
	d.Set("description", role.Description)
	d.Set("permissions", role.IncludedPermissions)
	d.Set("stage", role.Stage)
	d.Set("deleted", role.Deleted)
	d.Set("project", project)

	return nil
}

func resourceGoogleProjectIamCustomRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	if d.Get("deleted").(bool) {
		if d.HasChange("deleted") {
			// If other fields were changed, we need to update those first and then delete.
			// If we don't update, we will get diffs from re-apply
			// If we delete and then try to update, we will get an error.
			if err := resourceGoogleProjectIamCustomRoleUpdateNonDeletedFields(d, meta); err != nil {
				return err
			}
			if err := resourceGoogleProjectIamCustomRoleDelete(d, meta); err != nil {
				return err
			}

			d.SetPartial("deleted")
			d.Partial(false)
			return nil
		} else {
			return fmt.Errorf("cannot make changes to deleted custom project role %s", d.Id())
		}
	}

	// We want to update the role to some undeleted state.
	// Make sure the role with given ID exists and is un-deleted before patching.
	r, err := config.clientIAM.Projects.Roles.Get(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("unable to find custom project role %s to update: %v", d.Id(), err)
	}
	if r.Deleted {
		// Undelete if deleted previously
		if err := resourceGoogleProjectIamCustomRoleUndelete(d, meta); err != nil {
			return err
		}
		d.SetPartial("deleted")
	}

	if err := resourceGoogleProjectIamCustomRoleUpdateNonDeletedFields(d, meta); err != nil {
		return err
	}
	d.Partial(false)

	return nil
}

func resourceGoogleProjectIamCustomRoleUpdateNonDeletedFields(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("title") || d.HasChange("description") || d.HasChange("stage") || d.HasChange("permissions") {
		_, err := config.clientIAM.Projects.Roles.Patch(d.Id(), &iam.Role{
			Title:               d.Get("title").(string),
			Description:         d.Get("description").(string),
			Stage:               d.Get("stage").(string),
			IncludedPermissions: convertStringSet(d.Get("permissions").(*schema.Set)),
		}).Do()

		if err != nil {
			return fmt.Errorf("Error updating the custom project role %s: %s", d.Get("title").(string), err)
		}
		d.SetPartial("title")
		d.SetPartial("description")
		d.SetPartial("stage")
		d.SetPartial("permissions")
	}
	return nil
}

func resourceGoogleProjectIamCustomRoleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientIAM.Projects.Roles.Delete(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting the custom project role %s: %s", d.Get("title").(string), err)
	}

	return nil
}

func resourceGoogleProjectIamCustomRoleUndelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientIAM.Projects.Roles.Undelete(d.Id(), &iam.UndeleteRoleRequest{}).Do()
	if err != nil {
		return fmt.Errorf("Error undeleting the custom project role %s: %s", d.Get("title").(string), err)
	}

	return nil
}
