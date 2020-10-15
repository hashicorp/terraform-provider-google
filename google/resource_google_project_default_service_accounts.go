package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
)

// resourceGoogleProjectDefaultServiceAccounts returns a *schema.Resource that allows a customer
// to manage all the default serviceAccounts.
// It does mean that terraform tried to perform the action in the SA at some point but does not ensure that
// all defaults serviceAccounts where managed. Eg.: API was activated after project creation.
func resourceGoogleProjectDefaultServiceAccounts() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 0,

		Create: resourceGoogleProjectDefaultServiceAccountsCreate,
		Read:   resourceGoogleProjectDefaultServiceAccountsReadAndUpdate,
		Update: resourceGoogleProjectDefaultServiceAccountsReadAndUpdate,
		Delete: resourceGoogleProjectDefaultServiceAccountsDelete,

		// This resource should not have import, right?
		// Importer: &schema.ResourceImporter{
		// 	State: resourceGoogleProjectDefaultServiceAccountsImportState,
		// },

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
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
				ValidateFunc: validateServiceAccountAction(),
				Description:  `The action to be performed in the default service accounts. Valid values are: deprivilege, delete, disable`,
			},
			"restore_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validateRestorePolicy(),
				Description: `The action to be performed in the default service accounts on the resource destroy.
				Valid values are NONE and REACTIVATE. If set to REACTIVATE it will attempt to restore all default SAs`,
			},
		},
	}
}

func resourceGoogleProjectDefaultServiceAccountsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	pid, ok := d.Get("project_id").(string)
	if !ok {
		return fmt.Errorf("Cannot get project_id variable")
	}
	action, ok := d.Get("action").(string)
	if !ok {
		return fmt.Errorf("Cannot get action variable")
	}

	serviceAccounts, err := resourceGoogleProjectDefaultServiceAccountsList(config, d, userAgent)
	if err != nil {
		return fmt.Errorf("Error listing service accounts on project %s: %v", pid, err)
	}
	for _, sa := range serviceAccounts {
		switch action {
		// TODO: Add all cases and code apiCalls
		case "delete":
			log.Printf("[INFO] - Deleting service account %s on project %s", sa.Email, pid)
			return nil
		}
	}

	d.SetId(prefixedProject(pid))
	err = resourceGoogleProjectDefaultServiceAccountsReadAndUpdate(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func resourceGoogleProjectDefaultServiceAccountsList(config *Config, d *schema.ResourceData, userAgent string) ([]*iam.ServiceAccount, error) {
	pid, ok := d.Get("project_id").(string)
	if !ok {
		return nil, fmt.Errorf("Cannot get project_id variable")
	}
	// TODO: Add filter based on SA name as per documentation https://cloud.google.com/iam/docs/service-accounts#default
	// to filter only default service accounts
	response, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.List(prefixedProject(pid)).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list service accounts on project %q: %v", pid, err)
	}
	return response.Accounts, nil
}

func resourceGoogleProjectDefaultServiceAccountsReadAndUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if err = d.Set("project_id", p.ProjectId); err != nil {
		return fmt.Errorf("Error setting project_id: %s", err)
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
