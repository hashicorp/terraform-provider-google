// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func TestAccResourceGoogleProjectDefaultServiceAccountsBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_project_default_service_accounts.acceptance"
	org := envvar.GetTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsBasic(org, project, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "projects/"+project),
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "action"),
					resource.TestCheckResourceAttrSet(resourceName, "restore_policy"),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectDefaultServiceAccountsBasic(org, project, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	project_id = "%s"
	name       = "%s"
	org_id  = "%s"
	billing_account = "%s"
}

resource "google_project_default_service_accounts" "acceptance" {
	project = google_project.acceptance.project_id
	action = "DISABLE"
}
`, project, project, org, billingAccount)
}

func TestAccResourceGoogleProjectDefaultServiceAccountsDisable(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	action := "DISABLE"
	restorePolicy := "REVERT"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectDefaultServiceAccountsRevert(t, project, action),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsAdvanced(org, project, billingAccount, action, restorePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "id", "projects/"+project),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "action", action),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					acctest.SleepInSecondsForTest(5),
					testAccCheckGoogleProjectDefaultServiceAccountsChanges(t, project, action),
				),
			},
		},
	})
}

func TestAccResourceGoogleProjectDefaultServiceAccountsDelete(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	action := "DELETE"
	restorePolicy := "REVERT"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectDefaultServiceAccountsRevert(t, project, action),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsAdvanced(org, project, billingAccount, action, restorePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "id", "projects/"+project),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "action", action),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					acctest.SleepInSecondsForTest(10),
					testAccCheckGoogleProjectDefaultServiceAccountsChanges(t, project, action),
				),
			},
		},
	})
}

func TestAccResourceGoogleProjectDefaultServiceAccountsDeleteRevertIgnoreFailure(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	action := "DELETE"
	restorePolicy := "REVERT_AND_IGNORE_FAILURE"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsAdvanced(org, project, billingAccount, action, restorePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "id", "projects/"+project),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "action", action),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					acctest.SleepInSecondsForTest(10),
					testAccCheckGoogleProjectDefaultServiceAccountsChanges(t, project, action),
				),
			},
		},
	})
}

func TestAccResourceGoogleProjectDefaultServiceAccountsDeprivilege(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	action := "DEPRIVILEGE"
	restorePolicy := "REVERT"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectDefaultServiceAccountsRevert(t, project, action),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsAdvanced(org, project, billingAccount, action, restorePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "id", "projects/"+project),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "action", action),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					acctest.SleepInSecondsForTest(5),
					testAccCheckGoogleProjectDefaultServiceAccountsChanges(t, project, action),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectDefaultServiceAccountsAdvanced(org, project, billingAccount, action, restorePolicy string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	project_id = "%s"
	name       = "%s"
	org_id  = "%s"
	billing_account = "%s"
}

resource "google_project_service" "acceptance" {
	project = google_project.acceptance.project_id
	service = "compute.googleapis.com"

	disable_dependent_services = true
}

resource "google_project_default_service_accounts" "acceptance" {
	depends_on = [google_project_service.acceptance]
	project = google_project.acceptance.project_id
	action = "%s"
	restore_policy = "%s"
}
`, project, project, org, billingAccount, action, restorePolicy)
}

func testAccCheckGoogleProjectDefaultServiceAccountsChanges(t *testing.T, project, action string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		response, err := config.NewIamClient(config.UserAgent).Projects.ServiceAccounts.List(resourcemanager.PrefixedProject(project)).Do()
		if err != nil {
			return fmt.Errorf("failed to list service accounts on project %q: %v", project, err)
		}
		for _, sa := range response.Accounts {
			if testAccIsDefaultServiceAccount(sa.DisplayName) {
				switch action {
				case "DISABLE":
					if !sa.Disabled {
						return fmt.Errorf("compute engine default service account is not disabled, disable field is %t", sa.Disabled)
					}
				case "DELETE":
					return fmt.Errorf("compute engine default service account is not deleted")
				case "DEPRIVILEGE":
					iamPolicy, err := config.NewResourceManagerClient(config.UserAgent).Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{}).Do()
					if err != nil {
						return fmt.Errorf("cannot get IAM policy on project %s: %v", project, err)
					}
					for _, bind := range iamPolicy.Bindings {
						for _, member := range bind.Members {
							if member == fmt.Sprintf("serviceAccount:%s", sa.Email) {
								return fmt.Errorf("compute engine default service account is not deprivileged")
							}
						}
					}
					return nil
				}
			}
		}
		return nil
	}
}

// Test if actions were reverted properly
func testAccCheckGoogleProjectDefaultServiceAccountsRevert(t *testing.T, project, action string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		response, err := config.NewIamClient(config.UserAgent).Projects.ServiceAccounts.List(resourcemanager.PrefixedProject(project)).Do()
		if err != nil {
			return fmt.Errorf("failed to list service accounts on project %q: %v", project, err)
		}
		for _, sa := range response.Accounts {
			if testAccIsDefaultServiceAccount(sa.DisplayName) {
				// We agreed to not revert the DEPRIVILEGE action because will be hard to track the roles over the time
				if action == "DISABLE" {
					if sa.Disabled {
						return fmt.Errorf("compute engine default service account is not enabled, disable field is %t", sa.Disabled)
					}
				} else if action == "DELETE" {
					// A deleted service account was found meaning the undelete action triggered
					// on destroy worked
					return nil
				}
			}
		}
		// if action is DELETE, the service account should be found in the previous loop
		// due to undelete action
		if action == "DELETE" {
			return fmt.Errorf("service account changes were not reverted after destroy")
		}

		return nil
	}
}

// testAccIsDefaultServiceAccount is a helper function to facilitate TDD when there is a need
// to update how we determine whether it's a default SA or not.
// If you follow TDD, it is going to be different from isDefaultServiceAccount func while coding
// but they must be identical before commit/push
func testAccIsDefaultServiceAccount(displayName string) bool {
	gceDefaultSA := "compute engine default service account"
	appEngineDefaultSA := "app engine default service account"
	saDisplayName := strings.ToLower(displayName)
	if saDisplayName == gceDefaultSA || saDisplayName == appEngineDefaultSA {
		return true
	}

	return false
}
