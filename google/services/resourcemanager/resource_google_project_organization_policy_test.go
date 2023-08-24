// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Since each test here is acting on the same project, run the tests serially to
// avoid race conditions and aborted operations.
func TestAccProjectOrganizationPolicy(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"boolean":        testAccProjectOrganizationPolicy_boolean,
		"list_allowAll":  testAccProjectOrganizationPolicy_list_allowAll,
		"list_allowSome": testAccProjectOrganizationPolicy_list_allowSome,
		"list_denySome":  testAccProjectOrganizationPolicy_list_denySome,
		"list_update":    testAccProjectOrganizationPolicy_list_update,
		"restore_policy": testAccProjectOrganizationPolicy_restore_defaultTrue,
		"empty_policy":   testAccProjectOrganizationPolicy_none,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccProjectOrganizationPolicy_boolean(t *testing.T) {
	projectId := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Test creation of an enforced boolean policy
				Config: testAccProjectOrganizationPolicyConfig_boolean(projectId, true),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy(t, "bool", true),
			},
			{
				// Test update from enforced to not
				Config: testAccProjectOrganizationPolicyConfig_boolean(projectId, false),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy(t, "bool", false),
			},
			{
				Config:  " ",
				Destroy: true,
			},
			{
				// Test creation of a not enforced boolean policy
				Config: testAccProjectOrganizationPolicyConfig_boolean(projectId, false),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy(t, "bool", false),
			},
			{
				// Test update from not enforced to enforced
				Config: testAccProjectOrganizationPolicyConfig_boolean(projectId, true),
				Check:  testAccCheckGoogleProjectOrganizationBooleanPolicy(t, "bool", true),
			},
		},
	})
}

func testAccProjectOrganizationPolicy_list_allowAll(t *testing.T) {
	projectId := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicyConfig_list_allowAll(projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAll(t, "list", "ALLOW"),
			},
			{
				ResourceName:      "google_project_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectOrganizationPolicy_list_allowSome(t *testing.T) {
	project := envvar.GetTestProjectFromEnv()
	canonicalProject := canonicalProjectId(project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicyConfig_list_allowSome(project),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAllowedValues(t, "list", []string{canonicalProject}),
			},
			{
				ResourceName:      "google_project_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectOrganizationPolicy_list_denySome(t *testing.T) {
	projectId := envvar.GetTestProjectFromEnv()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicyConfig_list_denySome(projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyDeniedValues(t, "list", DENIED_ORG_POLICIES),
			},
			{
				ResourceName:      "google_project_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectOrganizationPolicy_list_update(t *testing.T) {
	projectId := envvar.GetTestProjectFromEnv()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicyConfig_list_allowAll(projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyAll(t, "list", "ALLOW"),
			},
			{
				Config: testAccProjectOrganizationPolicyConfig_list_denySome(projectId),
				Check:  testAccCheckGoogleProjectOrganizationListPolicyDeniedValues(t, "list", DENIED_ORG_POLICIES),
			},
			{
				ResourceName:      "google_project_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectOrganizationPolicy_restore_defaultTrue(t *testing.T) {
	projectId := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicyConfig_restore_defaultTrue(projectId),
				Check:  getGoogleProjectOrganizationRestoreDefaultTrue(t, "restore", &cloudresourcemanager.RestoreDefault{}),
			},
			{
				ResourceName:      "google_project_organization_policy.restore",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectOrganizationPolicy_none(t *testing.T) {
	projectId := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectOrganizationPolicyConfig_none(projectId),
				Check:  testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t),
			},
			{
				ResourceName:      "google_project_organization_policy.none",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleProjectOrganizationPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_project_organization_policy" {
				continue
			}

			projectId := canonicalProjectId(rs.Primary.Attributes["project"])
			constraint := resourcemanager.CanonicalOrgPolicyConstraint(rs.Primary.Attributes["constraint"])
			policy, err := config.NewResourceManagerClient(config.UserAgent).Projects.GetOrgPolicy(projectId, &cloudresourcemanager.GetOrgPolicyRequest{
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
}

func testAccCheckGoogleProjectOrganizationBooleanPolicy(t *testing.T, n string, enforced bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(t, s, n)
		if err != nil {
			return err
		}

		if policy.BooleanPolicy.Enforced != enforced {
			return fmt.Errorf("Expected boolean policy enforcement to be '%t', got '%t'", enforced, policy.BooleanPolicy.Enforced)
		}

		return nil
	}
}

func testAccCheckGoogleProjectOrganizationListPolicyAll(t *testing.T, n, policyType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(t, s, n)
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

func testAccCheckGoogleProjectOrganizationListPolicyAllowedValues(t *testing.T, n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(t, s, n)
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

func testAccCheckGoogleProjectOrganizationListPolicyDeniedValues(t *testing.T, n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleProjectOrganizationPolicyTestResource(t, s, n)
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

func getGoogleProjectOrganizationRestoreDefaultTrue(t *testing.T, n string, policyDefault *cloudresourcemanager.RestoreDefault) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		policy, err := getGoogleProjectOrganizationPolicyTestResource(t, s, n)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(policy.RestoreDefault, policyDefault) {
			return fmt.Errorf("Expected the restore default '%s', instead denied, %s", policyDefault, policy.RestoreDefault)
		}

		return nil
	}
}

func getGoogleProjectOrganizationPolicyTestResource(t *testing.T, s *terraform.State, n string) (*cloudresourcemanager.OrgPolicy, error) {
	rn := "google_project_organization_policy." + n
	rs, ok := s.RootModule().Resources[rn]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", rn)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	config := acctest.GoogleProviderConfig(t)
	projectId := canonicalProjectId(rs.Primary.Attributes["project"])

	return config.NewResourceManagerClient(config.UserAgent).Projects.GetOrgPolicy(projectId, &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: rs.Primary.Attributes["constraint"],
	}).Do()
}

func testAccProjectOrganizationPolicyConfig_boolean(pid string, enforced bool) string {
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

func testAccProjectOrganizationPolicyConfig_list_allowAll(pid string) string {
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

func testAccProjectOrganizationPolicyConfig_list_allowSome(pid string) string {
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

func testAccProjectOrganizationPolicyConfig_list_denySome(pid string) string {
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

func testAccProjectOrganizationPolicyConfig_restore_defaultTrue(pid string) string {
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

func testAccProjectOrganizationPolicyConfig_none(pid string) string {
	return fmt.Sprintf(`
resource "google_project_organization_policy" "none" {
  project    = "%s"
  constraint = "constraints/serviceuser.services"
}
`, pid)
}

func canonicalProjectId(project string) string {
	if strings.HasPrefix(project, "projects/") {
		return project
	}
	return fmt.Sprintf("projects/%s", project)
}
