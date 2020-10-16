package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
)

// resourceGoogleProjectDefaultServiceAccounts returns a *schema.Resource that allows a customer
// to manage all the default serviceAccounts.
// It does mean that terraform tried to perform the action in the SA at some point but does not ensure that
// all defaults serviceAccounts where managed. Eg.: API was activated after project creation.
func resourceGoogleProjectDefaultServiceAccounts() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectDefaultServiceAccountsCreate,
		Read:   resourceGoogleProjectDefaultServiceAccountsRead,
		Update: resourceGoogleProjectDefaultServiceAccountsUpdate,
		Delete: resourceGoogleProjectDefaultServiceAccountsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateProjectID(),
				Description:  `The project ID where service accounts are created.`,
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "deprivilege",
				ValidateFunc: validation.StringInSlice([]string{"deprivilege", "delete", "disable"}, false),
				Description:  `The action to be performed in the default service accounts. Valid values are: deprivilege, delete, disable.`,
			},
			"restore_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "REACTIVATE"}, false),
				Description: `The action to be performed in the default service accounts on the resource destroy.
				Valid values are NONE and REACTIVATE. If set to REACTIVATE it will attempt to restore all default SAs.`,
			},
		},
	}
}

func resourceGoogleProjectDefaultServiceAccountsDeleteAction(d *schema.ResourceData, meta interface{}, selflink string) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	_, err = config.NewIamClient(userAgent).Projects.ServiceAccounts.Delete(selflink).Do()
	if err != nil {
		return fmt.Errorf("Cannot delete service account: %v", err)
	}
	return nil
}

func resourceGoogleProjectDefaultServiceAccountsDoAction(d *schema.ResourceData, meta interface{}, action, email, project string) error {
	switch action {
	case "delete":
		var serviceAccountSelfLink = fmt.Sprintf("projects/%s/serviceAccounts/%s", project, email)
		err := resourceGoogleProjectDefaultServiceAccountsDeleteAction(d, meta, serviceAccountSelfLink)
		if err != nil {
			return fmt.Errorf("Cannot delete %s: %v", serviceAccountSelfLink, err)
		}
		return nil
	case "disable":
		return fmt.Errorf("not implemented yet")
	case "deprivilege":
		return fmt.Errorf("not implemented yet")
	default:
		return fmt.Errorf("Action %s is not a valid action", action)
	}
}

func resourceGoogleProjectDefaultServiceAccountsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	pid, ok := d.Get("project").(string)
	if !ok {
		return fmt.Errorf("Cannot get project")
	}
	action, ok := d.Get("action").(string)
	if !ok {
		return fmt.Errorf("Cannot get action")
	}

	serviceAccounts, err := resourceGoogleProjectDefaultServiceAccountsList(config, d, userAgent)
	if err != nil {
		return fmt.Errorf("Error listing service accounts on project %s: %v", pid, err)
	}

	for _, sa := range serviceAccounts {
		// As per documentation https://cloud.google.com/iam/docs/service-accounts#default
		// we have just two default SAs and the e-mail may change. So, it is been filtered
		// by the Display Name
		switch sa.DisplayName {
		case "Compute Engine default service account":
			err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, sa.Email, pid)
			if err != nil {
				return fmt.Errorf("Error doing action %s on Service Account %s", action, sa.Email)
			}
		case "App Engine default service account":
			err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, sa.Email, pid)
			if err != nil {
				return fmt.Errorf("Error doing action %s on Service Account %s", action, sa.Email)
			}
		default:
			continue
		}
	}

	d.SetId(prefixedProject(pid))

	return resourceGoogleProjectDefaultServiceAccountsRead(d, meta)
}

func resourceGoogleProjectDefaultServiceAccountsList(config *Config, d *schema.ResourceData, userAgent string) ([]*iam.ServiceAccount, error) {
	pid, ok := d.Get("project").(string)
	if !ok {
		return nil, fmt.Errorf("Cannot get project")
	}
	// TODO: Add filter based on SA name as per documentation https://cloud.google.com/iam/docs/service-accounts#default
	// to filter only default service accounts
	response, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.List(prefixedProject(pid)).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list service accounts on project %q: %v", pid, err)
	}
	return response.Accounts, nil
}

func resourceGoogleProjectDefaultServiceAccountsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	// Test if the project exists and permissions are set
	p, err := readGoogleProject(d, config, userAgent)
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 403 && strings.Contains(gerr.Message, "caller does not have permission") {
			return fmt.Errorf("the user does not have permission to access Project %q or it may not exist", p.ProjectId)
		}
		return handleNotFoundError(err, d, fmt.Sprintf("Project %q", p.ProjectId))
	}

	if err = d.Set("project", p.ProjectId); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err = d.Set("action", d.Get("action")); err != nil {
		return fmt.Errorf("Error setting action: %s", err)
	}
	if err = d.Set("restore_policy", d.Get("restore_policy")); err != nil {
		return fmt.Errorf("Error setting restore_policy: %s", err)
	}

	d.SetId(d.Id())

	return nil
}

func resourceGoogleProjectDefaultServiceAccountsDelete(d *schema.ResourceData, meta interface{}) error {
	// TODO: create func to handle actions on destroy depending on the current action and restore_policy
	d.SetId("")
	return nil
}

func resourceGoogleProjectDefaultServiceAccountsUpdate(d *schema.ResourceData, meta interface{}) error {

	// Restore policy has changed
	if ok := d.HasChange("restore_policy"); ok {
		if err := d.Set("restore_policy", d.Get("restore_policy")); err != nil {
			return fmt.Errorf("Error setting restore_policy: %s", err)
		}
	}

	return resourceGoogleProjectDefaultServiceAccountsRead(d, meta)
}
