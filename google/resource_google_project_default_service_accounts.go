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
			"service_accounts": {
				Type:        schema.TypeMap,
				Optional:    true,
				Default:     "",
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
	case "delete":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Delete(serviceAccountSelfLink).Do()
		if err != nil {
			return fmt.Errorf("Cannot delete service account %s: %v", serviceAccountSelfLink, err)
		}
	case "undelete":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Undelete(serviceAccountSelfLink, &iam.UndeleteServiceAccountRequest{}).Do()
		if err != nil {
			return fmt.Errorf("Cannot delete service account %s: %v", serviceAccountSelfLink, err)
		}
	case "disable":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Disable(serviceAccountSelfLink, &iam.DisableServiceAccountRequest{}).Do()
		if err != nil {
			return fmt.Errorf("Cannot disable service account %s: %v", serviceAccountSelfLink, err)
		}
	case "enable":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Enable(serviceAccountSelfLink, &iam.EnableServiceAccountRequest{}).Do()
		if err != nil {
			return fmt.Errorf("Cannot disable service account %s: %v", serviceAccountSelfLink, err)
		}
	case "deprivilege":
		iamPolicy, err := config.NewResourceManagerClient(userAgent).Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{
			Options:         &cloudresourcemanager.GetPolicyOptions{},
			ForceSendFields: []string{},
			NullFields:      []string{},
		}).Do()
		if err != nil {
			return fmt.Errorf("Cannot get IAM policy on project %s: %v", project, err)
		}

		for _, bind := range iamPolicy.Bindings {
			newMembers := []string{}
			if bind.Role == "roles/editor" {
				for _, member := range bind.Members {
					if member != fmt.Sprintf("serviceAccount:%s", email) {
						newMembers = append(newMembers, member)
					}
				}
			}
			bind.Members = newMembers
		}
		_, err = config.NewResourceManagerClient(userAgent).Projects.SetIamPolicy(project, &cloudresourcemanager.SetIamPolicyRequest{
			Policy:          iamPolicy,
			UpdateMask:      "",
			ForceSendFields: []string{},
			NullFields:      []string{},
		}).Do()
		if err != nil {
			return fmt.Errorf("Cannot update IAM policy on project %s: %v", project, err)
		}
	default:
		return fmt.Errorf("Action %s is not a valid action", action)
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
	changedServiceAccounts := make(map[string]interface{})
	for _, sa := range serviceAccounts {
		// As per documentation https://cloud.google.com/iam/docs/service-accounts#default
		// we have just two default SAs and the e-mail may change. So, it is been filtered
		// by the Display Name
		switch sa.DisplayName {
		case "Compute Engine default service account":
			changedServiceAccounts[sa.UniqueId] = fmt.Sprintf("%s:%s", sa.Email, action)
			err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, sa.UniqueId, sa.Email, pid)
			if err != nil {
				return fmt.Errorf("Error doing action %s on Service Account %s: %v", action, sa.Email, err)
			}
		case "App Engine default service account":
			changedServiceAccounts[sa.UniqueId] = fmt.Sprintf("%s:%s", sa.Email, action)
			err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, sa.UniqueId, sa.Email, pid)
			if err != nil {
				return fmt.Errorf("Error doing action %s on Service Account %s: %v", action, sa.Email, err)
			}
		default:
			continue
		}
		if changedServiceAccounts != nil {
			if err := d.Set("service_accounts", changedServiceAccounts); err != nil {
				return fmt.Errorf("Error setting action: %s", err)
			}
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
	response, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.List(prefixedProject(pid)).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list service accounts on project %q: %v", pid, err)
	}
	return response.Accounts, nil
}

func resourceGoogleProjectDefaultServiceAccountsRead(d *schema.ResourceData, meta interface{}) error {
	if err := d.Set("project", d.Get("project").(string)); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("action", d.Get("action").(string)); err != nil {
		return fmt.Errorf("Error setting action: %s", err)
	}
	if err := d.Set("restore_policy", d.Get("restore_policy").(string)); err != nil {
		return fmt.Errorf("Error setting restore_policy: %s", err)
	}
	if err := d.Set("service_accounts", d.Get("service_accounts").(map[string]interface{})); err != nil {
		return fmt.Errorf("Error setting service_accounts: %s", err)
	}
	d.SetId(d.Id())

	return nil
}

func resourceGoogleProjectDefaultServiceAccountsDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("restore_policy").(string) != "NONE" {
		pid, ok := d.Get("project").(string)
		if !ok {
			return fmt.Errorf("Cannot get project")
		}
		for saUniqueID, a := range d.Get("service_accounts").(map[string]interface{}) {
			data := strings.Split(a.(string), ":")
			saEmail := data[0]
			action := data[1]
			switch action {
			case "disable":
				action := "enable"
				err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, saUniqueID, saEmail, pid)
				if err != nil {
					return fmt.Errorf("Error doing action %s on Service Account %s: %v", action, saUniqueID, err)
				}
			case "deprivilege":
				action := "grantRole"
				err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, saUniqueID, saEmail, pid)
				if err != nil {
					return fmt.Errorf("Error doing action %s on Service Account %s: %v", action, saUniqueID, err)
				}
			case "delete":
				action := "undelete"
				err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, action, saUniqueID, saEmail, pid)
				if err != nil {
					return fmt.Errorf("Error doing action %s on Service Account %s: %v", action, saUniqueID, err)
				}
			}
		}
	}

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
