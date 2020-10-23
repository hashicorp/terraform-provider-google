package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
)

// resourceGoogleProjectDefaultServiceAccounts returns a *schema.Resource that allows a customer
// to manage all the default serviceAccounts.
// It does mean that terraform tried to perform the action in the SA at some point but does not ensure that
// all defaults serviceAccounts where managed. Eg.: API was activated after project creation.
func resourceGoogleProjectDefaultServiceAccounts() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectDefaultServiceAccountsCreate,
		Read:   schema.Noop,
		Update: schema.Noop,
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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEPRIVILEGE", "DELETE", "DISABLE"}, false),
				Description:  `The action to be performed in the default service accounts. Valid values are: DEPRIVILEGE, DELETE, DISABLE.`,
			},
			"restore_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "REACTIVATE",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "REACTIVATE"}, false),
				Description: `The action to be performed in the default service accounts on the resource destroy.
				Valid values are NONE and REACTIVATE. If set to REACTIVATE it will attempt to restore all default SAs.`,
			},
			"service_accounts": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `The Service Accounts changed by this resource. It is used for revert the action on the destroy`,
			},
		},
	}
}

func resourceGoogleProjectDefaultServiceAccountsDoAction(d *schema.ResourceData, meta interface{}, action, uniqueID, email, project string) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	serviceAccountSelfLink := fmt.Sprintf("projects/%s/serviceAccounts/%s", project, uniqueID)
	switch action {
	case "DELETE":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Delete(serviceAccountSelfLink).Do()
		if err != nil {
			return fmt.Errorf("cannot delete service account %s: %v", serviceAccountSelfLink, err)
		}
	case "UNDELETE":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Undelete(serviceAccountSelfLink, &iam.UndeleteServiceAccountRequest{}).Do()
		if err != nil {
			return fmt.Errorf("cannot undelete service account %s: %v", serviceAccountSelfLink, err)
		}
	case "DISABLE":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Disable(serviceAccountSelfLink, &iam.DisableServiceAccountRequest{}).Do()
		if err != nil {
			return fmt.Errorf("cannot disable service account %s: %v", serviceAccountSelfLink, err)
		}
	case "ENABLE":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Enable(serviceAccountSelfLink, &iam.EnableServiceAccountRequest{}).Do()
		if err != nil {
			return fmt.Errorf("cannot enable service account %s: %v", serviceAccountSelfLink, err)
		}
	case "DEPRIVILEGE":
		iamPolicy, err := config.NewResourceManagerClient(userAgent).Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{}).Do()
		if err != nil {
			return fmt.Errorf("cannot get IAM policy on project %s: %v", project, err)
		}

		for _, bind := range iamPolicy.Bindings {
			newMembers := []string{}
			// Google only adds editor role when creating default service accounts
			// That said, we just want to check when the role binded is editor.
			if bind.Role == "roles/editor" {
				for _, member := range bind.Members {
					if member != fmt.Sprintf("serviceAccount:%s", email) {
						newMembers = append(newMembers, member)
					}
				}
			}
			bind.Members = newMembers
		}
		_, err = config.NewResourceManagerClient(userAgent).Projects.SetIamPolicy(project, &cloudresourcemanager.SetIamPolicyRequest{}).Do()
		if err != nil {
			return fmt.Errorf("cannot update IAM policy on project %s: %v", project, err)
		}
	default:
		return fmt.Errorf("action %s is not a valid action", action)
	}

	return nil
}

func resourceGoogleProjectDefaultServiceAccountsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	pid, ok := d.Get("project").(string)
	if !ok {
		return fmt.Errorf("cannot get project")
	}
	action, ok := d.Get("action").(string)
	if !ok {
		return fmt.Errorf("cannot get action")
	}

	serviceAccounts, err := resourceGoogleProjectDefaultServiceAccountsList(config, d, userAgent)
	if err != nil {
		return fmt.Errorf("error listing service accounts on project %s: %v", pid, err)
	}
	changedServiceAccounts := make(map[string]interface{})
	for _, sa := range serviceAccounts {
		// As per documentation https://cloud.google.com/iam/docs/service-accounts#default
		// we have just two default SAs and the e-mail may change. So, it is been filtered
		// by the Display Name
		if isDefaultServiceAccount(sa.DisplayName) {
			err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, sa.UniqueId, sa.Email, pid)
			if err != nil {
				return fmt.Errorf("error doing action %s on Service Account %s: %v", action, sa.Email, err)
			}
			changedServiceAccounts[sa.UniqueId] = sa.Email
		}
	}
	if err := d.Set("service_accounts", changedServiceAccounts); err != nil {
		return fmt.Errorf("error setting service_accounts: %s", err)
	}
	d.SetId(prefixedProject(pid))

	return nil
}

func resourceGoogleProjectDefaultServiceAccountsList(config *Config, d *schema.ResourceData, userAgent string) ([]*iam.ServiceAccount, error) {
	pid, ok := d.Get("project").(string)
	if !ok {
		return nil, fmt.Errorf("cannot get project")
	}
	response, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.List(prefixedProject(pid)).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list service accounts on project %q: %v", pid, err)
	}
	return response.Accounts, nil
}

func resourceGoogleProjectDefaultServiceAccountsDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("restore_policy").(string) == "NONE" {
		d.SetId("")
		return nil
	}

	pid, ok := d.Get("project").(string)
	if !ok {
		return fmt.Errorf("cannot get project")
	}
	for saUniqueID, saEmail := range d.Get("service_accounts").(map[string]interface{}) {
		action := d.Get("action").(string)
		switch action {
		case "DISABLE":
			action := "ENABLE"
			err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, saUniqueID, saEmail.(string), pid)
			if err != nil {
				return fmt.Errorf("error doing action %s on Service Account %s: %v", action, saUniqueID, err)
			}
		case "DELETE":
			action := "UNDELETE"
			err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, saUniqueID, saEmail.(string), pid)
			if err != nil {
				return fmt.Errorf("error doing action %s on Service Account %s: %v", action, saUniqueID, err)
			}
		}
	}

	d.SetId("")

	return nil
}

func isDefaultServiceAccount(displayName string) bool {
	gceDefaultSA := "compute engine default service account"
	appEngineDefaultSA := "app engine default service account"
	saDisplayName := strings.ToLower(displayName)
	if saDisplayName == gceDefaultSA || saDisplayName == appEngineDefaultSA {
		return true
	}

	return false
}
