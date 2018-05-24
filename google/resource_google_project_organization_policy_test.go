package google

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

/*
Tests for `google_project_organization_policy`

These are *not* run in parallel, as they all use the same project
and I end up with 409 Conflict errors from the API when they are
run in parallel.
*/

func TestAccProjectOrganizationPolicy_boolean(t *testing.T) {
	projectId := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test creation of an enforced boolean policy
				Config: testAccProjectOrganizationPolicy_boolean(projectId, true),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy("bool", true),
			},
			{
				// Test update from enforced to not
				Config: testAccProjectOrganizationPolicy_boolean(projectId, false),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy("bool", false),
			},
			{
				Config:  " ",
				Destroy: true,
			},
			{
				// Test creation of a not enforced boolean policy
				Config: testAccProjectOrganizationPolicy_boolean(projectId, false),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy("bool", false),
			},
			{
				// Test update from not enforced to enforced
				Config: testAccProjectOrganizationPolicy_boolean(projectId, true),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy("bool", true),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_list_allowAll(t *testing.T) {
	projectId := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_list_allowAll(projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAll("list", "ALLOW"),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_list_allowSome(t *testing.T) {
	project := getTestProjectFromEnv()
	canonicalProject := canonicalProjectId(project)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_list_allowSome(project),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAllowedValues("list", []string{canonicalProject}),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_list_denySome(t *testing.T) {
	projectId := getTestProjectFromEnv()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_list_denySome(projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyDeniedValues("list", DENIED_ORG_POLICIES),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_list_update(t *testing.T) {
	projectId := getTestProjectFromEnv()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_list_allowAll(projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAll("list", "ALLOW"),
			},
			{
				Config: testAccProjectOrganizationPolicy_list_denySome(projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyDeniedValues("list", DENIED_ORG_POLICIES),
			},
		},
	})
}

func TestAccProjectOrganizationPolicy_restore_defaultTrue(t *testing.T) {
	projectId := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleProjectOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicy_restore_defaultTrue(projectId),
				Check:  getGoogleProjectOrganizationRestoreDefaultTrue("restore", &cloudresourcemanager.RestoreDefault{}),
			},
		},
	})
}

func testAccCheckGoogleProjectOrganizationPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_project_organization_policy" {
			continue
		}

		projectId := canonicalProjectId(rs.Primary.Attributes["project"])
		constraint := canonicalOrgPolicyConstraint(rs.Primary.Attributes["constraint"])
		policy, err := config.clientResourceManager.Projects.GetOrgPolicy(projectId, &cloudresourcemanager.GetOrgPolicyRequest{
			Constraint: constraint,
		}).Do()

		if err != nil {
			return err
		}

		if policy.ListPolicy != nil || policy.BooleanPolicy != nil {
			return fmt.Errorf("Org policy with constraint '%s' hasn't been cleared", constraint)
		}
	}
	return nil
}

func testAccCheckGoogleProjectOrganizationBooleanPolicy(n string, enforced bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if policy.BooleanPolicy.Enforced != enforced {
			return fmt.Errorf("Expected boolean policy enforcement to be '%t', got '%t'", enforced, policy.BooleanPolicy.Enforced)
		}

		return nil
	}
}

func testAccCheckGoogleProjectOrganizationListPolicyAll(n, policyType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if policy.ListPolicy == nil {
			return nil
		}

		if len(policy.ListPolicy.AllowedValues) > 0 || len(policy.ListPolicy.DeniedValues) > 0 {
			return fmt.Errorf("The `values` field shouldn't be set")
		}

		if policy.ListPolicy.AllValues != policyType {
			return fmt.Errorf("The list policy should %s all values", policyType)
		}

		return nil
	}
}

func testAccCheckGoogleProjectOrganizationListPolicyAllowedValues(n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		sort.Strings(policy.ListPolicy.AllowedValues)
		sort.Strings(values)
		if !reflect.DeepEqual(policy.ListPolicy.AllowedValues, values) {
			return fmt.Errorf("Expected the list policy to allow '%s', instead allowed '%s'", values, policy.ListPolicy.AllowedValues)
		}

		return nil
	}
}

func testAccCheckGoogleProjectOrganizationListPolicyDeniedValues(n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		sort.Strings(policy.ListPolicy.DeniedValues)
		sort.Strings(values)
		if !reflect.DeepEqual(policy.ListPolicy.DeniedValues, values) {
			return fmt.Errorf("Expected the list policy to deny '%s', instead denied '%s'", values, policy.ListPolicy.DeniedValues)
		}

		return nil
	}
}

func getGoogleProjectOrganizationRestoreDefaultTrue(n string, policyDefault *cloudresourcemanager.RestoreDefault) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		policy, err := getGoogleProjectOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(policy.RestoreDefault, policyDefault) {
			return fmt.Errorf("Expected the restore default '%s', instead denied, %s", policyDefault, policy.RestoreDefault)
		}

		return nil
	}
}

func getGoogleProjectOrganizationPolicyTestResource(s *terraform.State, n string) (*cloudresourcemanager.OrgPolicy, error) {
	rn := "google_project_organization_policy." + n
	rs, ok := s.RootModule().Resources[rn]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", rn)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	config := testAccProvider.Meta().(*Config)
	projectId := canonicalProjectId(rs.Primary.Attributes["project"])

	return config.clientResourceManager.Projects.GetOrgPolicy(projectId, &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: rs.Primary.Attributes["constraint"],
	}).Do()
}

func testAccProjectOrganizationPolicy_boolean(pid string, enforced bool) string {
	return fmt.Sprintf(`
resource "google_project_organization_policy" "bool" {
  project    = "%s"
  constraint = "constraints/compute.disableSerialPortAccess"

  boolean_policy {
    enforced = %t
  }
}
`, pid, enforced)
}

func testAccProjectOrganizationPolicy_list_allowAll(pid string) string {
	return fmt.Sprintf(`
resource "google_project_organization_policy" "list" {
  project    = "%s"
  constraint = "constraints/serviceuser.services"

  list_policy {
    allow {
      all = true
    }
  }
}
`, pid)
}

func testAccProjectOrganizationPolicy_list_allowSome(pid string) string {
	return fmt.Sprintf(`

resource "google_project_organization_policy" "list" {
  project    = "%s"
  constraint = "constraints/compute.trustedImageProjects"

  list_policy {
    allow {
      values = ["projects/%s"]
    }
  }
}
`, pid, pid)
}

func testAccProjectOrganizationPolicy_list_denySome(pid string) string {
	return fmt.Sprintf(`

resource "google_project_organization_policy" "list" {
  project    = "%s"
  constraint = "constraints/serviceuser.services"

  list_policy {
    deny {
      values = [
        "doubleclicksearch.googleapis.com",
        "replicapoolupdater.googleapis.com",
      ]
    }
  }
}
`, pid)
}

func testAccProjectOrganizationPolicy_restore_defaultTrue(pid string) string {
	return fmt.Sprintf(`
resource "google_project_organization_policy" "restore" {
  project    = "%s"
  constraint = "constraints/serviceuser.services"

    restore_policy {
        default = true
    }
}
`, pid)
}

func canonicalProjectId(project string) string {
	if strings.HasPrefix(project, "projects/") {
		return project
	}
	return fmt.Sprintf("projects/%s", project)
}
