package google

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func TestAccResourceGoogleProjectDefaultServiceAccountsBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_project_default_service_accounts.acceptance"
	org := getTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-project-%d", randInt(t))
	billingAccount := getTestBillingAccountFromEnv(t)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsBasic(org, project, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "projects/"+project),
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "action"),
					resource.TestCheckResourceAttrSet(resourceName, "project"),
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
}
`, project, project, org, billingAccount)
}

func TestAccResourceGoogleProjectDefaultServiceAccountsDisable(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-project-%d", randInt(t))
	billingAccount := getTestBillingAccountFromEnv(t)
	action := "DISABLE"
	restorePolicy := "REACTIVATE"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectDefaultServiceAccountsRevert(t, project, action),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsAdvanced(org, project, billingAccount, action, restorePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "id", "projects/"+project),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "action", action),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					preventRaceCondition(), // Looks like the API is eventually consistent by couple seconds
					testAccCheckGoogleProjectDefaultServiceAccountsChanges(t, project, action),
				),
			},
		},
	})
}

func TestAccResourceGoogleProjectDefaultServiceAccountsDelete(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-project-%d", randInt(t))
	billingAccount := getTestBillingAccountFromEnv(t)
	action := "DELETE"
	restorePolicy := "REACTIVATE"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectDefaultServiceAccountsRevert(t, project, action),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsAdvanced(org, project, billingAccount, action, restorePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "id", "projects/"+project),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "action", action),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					preventRaceCondition(), // Looks like the API is eventually consistent by couple seconds
					testAccCheckGoogleProjectDefaultServiceAccountsChanges(t, project, action),
				),
			},
		},
	})
}

func TestAccResourceGoogleProjectDefaultServiceAccountsDeprivilege(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	project := fmt.Sprintf("tf-project-%d", randInt(t))
	billingAccount := getTestBillingAccountFromEnv(t)
	action := "DEPRIVILEGE"
	restorePolicy := "REACTIVATE"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectDefaultServiceAccountsRevert(t, project, action),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleProjectDefaultServiceAccountsAdvanced(org, project, billingAccount, action, restorePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "id", "projects/"+project),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					resource.TestCheckResourceAttr("google_project_default_service_accounts.acceptance", "action", action),
					resource.TestCheckResourceAttrSet("google_project_default_service_accounts.acceptance", "project"),
					preventRaceCondition(), // Looks like the API is eventually consistent by couple seconds
					testAccCheckGoogleProjectDefaultServiceAccountsChanges(t, project, action),
				),
			},
		},
	})
}

func preventRaceCondition() func(s *terraform.State) error {
	return func(s *terraform.State) error {
		time.Sleep(5 * time.Second)
		return nil
	}
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
		config := googleProviderConfig(t)
		response, err := config.NewIamClient(config.userAgent).Projects.ServiceAccounts.List(prefixedProject(project)).Do()
		if err != nil {
			return fmt.Errorf("failed to list service accounts on project %q: %v", project, err)
		}
		for _, sa := range response.Accounts {
			switch strings.ToLower(sa.DisplayName) {
			case "compute engine default service account":
				switch action {
				case "DISABLE":
					if !sa.Disabled {
						return fmt.Errorf("compute engine default service account is not disabled, disable field is %t", sa.Disabled)
					}
				case "DELETE":
					return fmt.Errorf("compute engine default service account is not deleted")
				case "DEPRIVILEGE":
					iamPolicy, err := config.NewResourceManagerClient(config.userAgent).Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{
						Options:         &cloudresourcemanager.GetPolicyOptions{},
						ForceSendFields: []string{},
						NullFields:      []string{},
					}).Do()
					if err != nil {
						return fmt.Errorf("cannot get IAM policy on project %s: %v", project, err)
					}
					for _, bind := range iamPolicy.Bindings {
						if bind.Role == "roles/editor" {
							for _, member := range bind.Members {
								if member == fmt.Sprintf("serviceAccount:%s", sa.Email) {
									return fmt.Errorf("compute engine default service account is not deprivileged")
								}
							}
						}
					}
					return nil
				}
			case "app engine default service account":
				switch action {
				case "DISABLE":
					if !sa.Disabled {
						return fmt.Errorf("app engine default service account is not disabled, disable field is %t", sa.Disabled)
					}
				case "DELETE":
					return fmt.Errorf("app engine default service account is not deleted")
				case "DEPRIVILEGE":
					iamPolicy, err := config.NewResourceManagerClient(config.userAgent).Projects.GetIamPolicy(project, &cloudresourcemanager.GetIamPolicyRequest{
						Options:         &cloudresourcemanager.GetPolicyOptions{},
						ForceSendFields: []string{},
						NullFields:      []string{},
					}).Do()
					if err != nil {
						return fmt.Errorf("cannot get IAM policy on project %s: %v", project, err)
					}
					for _, bind := range iamPolicy.Bindings {
						if bind.Role == "roles/editor" {
							for _, member := range bind.Members {
								if member == fmt.Sprintf("serviceAccount:%s", sa.Email) {
									return fmt.Errorf("app engine default service account is not deprivileged")
								}
							}
						}
					}
				}
			default:
				continue
			}
		}
		return nil
	}
}

func testAccCheckGoogleProjectDefaultServiceAccountsRevert(t *testing.T, project, action string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		response, err := config.NewIamClient(config.userAgent).Projects.ServiceAccounts.List(prefixedProject(project)).Do()
		if err != nil {
			return fmt.Errorf("failed to list service accounts on project %q: %v", project, err)
		}
		for _, sa := range response.Accounts {
			switch strings.ToLower(sa.DisplayName) {
			case "compute engine default service account":
				switch action {
				case "DISABLE":
					if sa.Disabled {
						return fmt.Errorf("compute engine default service account is not enabled, disable field is %t", sa.Disabled)
					}
				case "DELETE":
					return nil
				case "DEPRIVILEGE":
					return nil
				}
			case "app engine default service account":
				switch action {
				case "DISABLE":
					if sa.Disabled {
						return fmt.Errorf("app engine default service account is not enabled, disable field is %t", sa.Disabled)
					}
				case "DELETE":
					return nil
				case "DEPRIVILEGE":
					return nil
				}
			default:
				continue
			}
		}
		if action == "DELETE" {
			// if it does not return from the case statement for DELETE action is because
			// it failed to undelete it
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
