package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	roleId := fmt.Sprintf("projects/%s/roles/%s", project, d.Get("role_id").(string))
	r, err := config.NewIamClient(userAgent).Projects.Roles.Get(roleId).Do()
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
		role, err := config.NewIamClient(userAgent).Projects.Roles.Create("projects/"+project, &iam.CreateRoleRequest{
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project := extractProjectFromProjectIamCustomRoleID(d.Id())

	role, err := config.NewIamClient(userAgent).Projects.Roles.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, d.Id())
	}

	if err := d.Set("role_id", GetResourceNameFromSelfLink(role.Name)); err != nil {
		return fmt.Errorf("Error setting role_id: %s", err)
	}
	if err := d.Set("title", role.Title); err != nil {
		return fmt.Errorf("Error setting title: %s", err)
	}
	if err := d.Set("name", role.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", role.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("permissions", role.IncludedPermissions); err != nil {
		return fmt.Errorf("Error setting permissions: %s", err)
	}
	if err := d.Set("stage", role.Stage); err != nil {
		return fmt.Errorf("Error setting stage: %s", err)
	}
	if err := d.Set("deleted", role.Deleted); err != nil {
		return fmt.Errorf("Error setting deleted: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}

func resourceGoogleProjectIamCustomRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	d.Partial(true)

	// We want to update the role to some undeleted state.
	// Make sure the role with given ID exists and is un-deleted before patching.
	r, err := config.NewIamClient(userAgent).Projects.Roles.Get(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("unable to find custom project role %s to update: %v", d.Id(), err)
	}
	if r.Deleted {
		_, err := config.NewIamClient(userAgent).Projects.Roles.Undelete(d.Id(), &iam.UndeleteRoleRequest{}).Do()
		if err != nil {
			return fmt.Errorf("Error undeleting the custom project role %s: %s", d.Get("title").(string), err)
		}
	}

	if d.HasChange("title") || d.HasChange("description") || d.HasChange("stage") || d.HasChange("permissions") {
		_, err := config.NewIamClient(userAgent).Projects.Roles.Patch(d.Id(), &iam.Role{
			Title:               d.Get("title").(string),
			Description:         d.Get("description").(string),
			Stage:               d.Get("stage").(string),
			IncludedPermissions: convertStringSet(d.Get("permissions").(*schema.Set)),
		}).Do()

		if err != nil {
			return fmt.Errorf("Error updating the custom project role %s: %s", d.Get("title").(string), err)
		}
	}

	d.Partial(false)
	return nil
}

func resourceGoogleProjectIamCustomRoleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	_, err = config.NewIamClient(userAgent).Projects.Roles.Delete(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting the custom project role %s: %s", d.Get("title").(string), err)
	}

	return nil
}
