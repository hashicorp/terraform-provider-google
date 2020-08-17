package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  `The camel case role id to use for this role. Cannot contain - characters.`,
				ValidateFunc: validateIAMCustomRoleID,
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
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The project that the service account will be created in. Defaults to the provider project configuration.`,
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
				Description: `The name of the role in the format projects/{{project}}/roles/{{role_id}}. Like id, this field can be used as a reference in other resources such as IAM role bindings.`,
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

func extractProjectFromProjectIamCustomRoleID(id string) string {
	parts := strings.Split(id, "/")

	return parts[1]
}

func resourceGoogleProjectIamCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project := extractProjectFromProjectIamCustomRoleID(d.Id())

	role, err := config.clientIAM.Projects.Roles.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Id())
	}

	d.Set("role_id", GetResourceNameFromSelfLink(role.Name))
	d.Set("title", role.Title)
	d.Set("name", role.Name)
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

	// We want to update the role to some undeleted state.
	// Make sure the role with given ID exists and is un-deleted before patching.
	r, err := config.clientIAM.Projects.Roles.Get(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("unable to find custom project role %s to update: %v", d.Id(), err)
	}
	if r.Deleted {
		_, err := config.clientIAM.Projects.Roles.Undelete(d.Id(), &iam.UndeleteRoleRequest{}).Do()
		if err != nil {
			return fmt.Errorf("Error undeleting the custom project role %s: %s", d.Get("title").(string), err)
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
			return fmt.Errorf("Error updating the custom project role %s: %s", d.Get("title").(string), err)
		}
		d.SetPartial("title")
		d.SetPartial("description")
		d.SetPartial("stage")
		d.SetPartial("permissions")
	}

	d.Partial(false)
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
