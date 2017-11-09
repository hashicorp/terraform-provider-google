package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/iam/v1"
)

func resourceGoogleProjectIamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectIamRoleCreate,
		Read:   resourceGoogleProjectIamRoleRead,
		Update: resourceGoogleProjectIamRoleUpdate,
		Delete: resourceGoogleProjectIamRoleDelete,

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
				ForceNew: true,
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

func resourceGoogleProjectIamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	if d.Get("deleted").(bool) {
		return fmt.Errorf("Cannot create a project with a deleted state. `deleted` field should be false.")
	}

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
		return fmt.Errorf("Error creating the project role %s: %s", d.Get("title").(string), err)
	}

	d.SetId(role.Name)

	return resourceGoogleProjectIamRoleRead(d, meta)
}

func resourceGoogleProjectIamRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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

	return nil
}

func resourceGoogleProjectIamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	if d.HasChange("deleted") {
		if d.Get("deleted").(bool) {
			if err := resourceGoogleProjectIamRoleDelete(d, meta); err != nil {
				return err
			}
		} else {
			if err := resourceGoogleProjectIamRoleUndelete(d, meta); err != nil {
				return err
			}
		}
		d.SetPartial("deleted")
	}

	if d.HasChange("title") || d.HasChange("description") || d.HasChange("stage") || d.HasChange("permissions") {
		_, err := config.clientIAM.Projects.Roles.Patch(d.Id(), &iam.Role{
			Title:               d.Get("title").(string),
			Description:         d.Get("description").(string),
			Stage:               d.Get("stage").(string),
			IncludedPermissions: convertStringSet(d.Get("permissions").(*schema.Set)),
		}).Do()

		if err != nil {
			return fmt.Errorf("Error updating the project role %s: %s", d.Get("title").(string), err)
		}
		d.SetPartial("title")
		d.SetPartial("description")
		d.SetPartial("stage")
		d.SetPartial("permissions")
	}

	d.Partial(false)

	return nil
}

func resourceGoogleProjectIamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientIAM.Projects.Roles.Delete(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting the project role %s: %s", d.Get("title").(string), err)
	}

	return nil
}

func resourceGoogleProjectIamRoleUndelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientIAM.Projects.Roles.Undelete(d.Id(), &iam.UndeleteRoleRequest{}).Do()
	if err != nil {
		return fmt.Errorf("Error undeleting the project role %s: %s", d.Get("title").(string), err)
	}

	return nil
}
