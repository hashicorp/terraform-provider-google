// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"
)

// ResourceGoogleProjectDefaultServiceAccounts returns a *schema.Resource that allows a customer
// to manage all the default serviceAccounts.
// It does mean that terraform tried to perform the action in the SA at some point but does not ensure that
// all defaults serviceAccounts where managed. Eg.: API was activated after project creation.
func ResourceGoogleProjectDefaultServiceAccounts() *schema.Resource {
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
				ValidateFunc: verify.ValidateProjectID(),
				Description:  `The project ID where service accounts are created.`,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEPRIVILEGE", "DELETE", "DISABLE"}, false),
				Description: `The action to be performed in the default service accounts. Valid values are: DEPRIVILEGE, DELETE, DISABLE.
				Note that DEPRIVILEGE action will ignore the REVERT configuration in the restore_policy.`,
			},
			"restore_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "REVERT",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "REVERT", "REVERT_AND_IGNORE_FAILURE"}, false),
				Description: `The action to be performed in the default service accounts on the resource destroy.
				Valid values are NONE, REVERT and REVERT_AND_IGNORE_FAILURE. It is applied for any action but in the DEPRIVILEGE.`,
			},
			"service_accounts": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `The Service Accounts changed by this resource. It is used for revert the action on the destroy.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceGoogleProjectDefaultServiceAccountsDoAction(d *schema.ResourceData, meta interface{}, action, uniqueID, email, project string) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	restorePolicy := d.Get("restore_policy").(string)
	serviceAccountSelfLink := fmt.Sprintf("projects/%s/serviceAccounts/%s", project, uniqueID)
	switch action {
	case "DELETE":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Delete(serviceAccountSelfLink).Do()
		if err != nil {
			return fmt.Errorf("cannot delete service account %s: %v", serviceAccountSelfLink, err)
		}
	case "UNDELETE":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Undelete(serviceAccountSelfLink, &iam.UndeleteServiceAccountRequest{}).Do()
		errExpected := restorePolicy == "REVERT_AND_IGNORE_FAILURE"
		errReceived := err != nil
		if errReceived {
			if !errExpected {
				return fmt.Errorf("cannot undelete service account %s: %v", serviceAccountSelfLink, err)
			}
			log.Printf("cannot undelete service account %s: %v", serviceAccountSelfLink, err)
			log.Printf("restore policy is %s... ignoring error", restorePolicy)
		}
	case "DISABLE":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Disable(serviceAccountSelfLink, &iam.DisableServiceAccountRequest{}).Do()
		if err != nil {
			return fmt.Errorf("cannot disable service account %s: %v", serviceAccountSelfLink, err)
		}
	case "ENABLE":
		_, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.Enable(serviceAccountSelfLink, &iam.EnableServiceAccountRequest{}).Do()
		errReceived := err != nil
		errExpected := restorePolicy == "REVERT_AND_IGNORE_FAILURE"
		if errReceived {
			if !errExpected {
				return fmt.Errorf("cannot enable service account %s: %v", serviceAccountSelfLink, err)
			}
			log.Printf("cannot enable service account %s: %v", serviceAccountSelfLink, err)
			log.Printf("restore policy is %s... ignoring error", restorePolicy)
		}
	case "DEPRIVILEGE":
		iamPolicy, err := config.NewResourceManagerClient(userAgent).Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{}).Do()
		if err != nil {
			return fmt.Errorf("cannot get IAM policy on project %s: %v", project, err)
		}

		// Creates a new slice with all members but the service account
		for _, bind := range iamPolicy.Bindings {
			newMembers := []string{}
			for _, member := range bind.Members {
				if member != fmt.Sprintf("serviceAccount:%s", email) {
					newMembers = append(newMembers, member)
				}
			}
			bind.Members = newMembers
		}
		updateRequest := &cloudresourcemanager.SetIamPolicyRequest{
			Policy:     iamPolicy,
			UpdateMask: "bindings,etag,auditConfigs",
		}
		_, err = config.NewResourceManagerClient(userAgent).Projects.SetIamPolicy(project, updateRequest).Do()
		if err != nil {
			return fmt.Errorf("cannot update IAM policy on project %s: %v", project, err)
		}
	default:
		return fmt.Errorf("action %s is not a valid action", action)
	}

	return nil
}

func resourceGoogleProjectDefaultServiceAccountsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	pid := d.Get("project").(string)
	action := d.Get("action").(string)

	serviceAccounts, err := listServiceAccounts(config, d, userAgent)
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
	d.SetId(PrefixedProject(pid))

	return nil
}

func listServiceAccounts(config *transport_tpg.Config, d *schema.ResourceData, userAgent string) ([]*iam.ServiceAccount, error) {
	pid := d.Get("project").(string)
	response, err := config.NewIamClient(userAgent).Projects.ServiceAccounts.List(PrefixedProject(pid)).Do()
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

	pid := d.Get("project").(string)
	for saUniqueID, saEmail := range d.Get("service_accounts").(map[string]interface{}) {
		origAction := d.Get("action").(string)
		newAction := ""
		// We agreed to not revert the DEPRIVILEGE because Morgante said it is not required.
		// It may be an enhancement. https://github.com/hashicorp/terraform-provider-google/issues/4135#issuecomment-709480278
		if origAction == "DISABLE" {
			newAction = "ENABLE"
		} else if origAction == "DELETE" {
			newAction = "UNDELETE"
		}
		if newAction != "" {
			err := resourceGoogleProjectDefaultServiceAccountsDoAction(d, meta, newAction, saUniqueID, saEmail.(string), pid)
			if err != nil {
				return fmt.Errorf("error doing action %s on Service Account %s: %v", newAction, saUniqueID, err)
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
