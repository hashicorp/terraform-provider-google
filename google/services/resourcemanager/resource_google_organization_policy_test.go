// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var DENIED_ORG_POLICIES = []string{
	"doubleclicksearch.googleapis.com",
	"replicapoolupdater.googleapis.com",
}

// Since each test here is acting on the same organization, run the tests serially to
// avoid race conditions and aborted operations.
func TestAccOrganizationPolicy(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"boolean":                testAccOrganizationPolicy_boolean,
		"list_allowAll":          testAccOrganizationPolicy_list_allowAll,
		"list_allowSome":         testAccOrganizationPolicy_list_allowSome,
		"list_denySome":          testAccOrganizationPolicy_list_denySome,
		"list_update":            testAccOrganizationPolicy_list_update,
		"list_inheritFromParent": testAccOrganizationPolicy_list_inheritFromParent,
		"restore_policy":         testAccOrganizationPolicy_restore_defaultTrue,
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

func testAccOrganizationPolicy_boolean(t *testing.T) {
	org := envvar.GetTestOrgTargetFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Test creation of an enforced boolean policy
				Config: testAccOrganizationPolicyConfig_boolean(org, true),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy(t, "bool", true),
			},
			{
				// Test update from enforced to not
				Config: testAccOrganizationPolicyConfig_boolean(org, false),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy(t, "bool", false),
			},
			{
				Config:  " ",
				Destroy: true,
			},
			{
				// Test creation of a not enforced boolean policy
				Config: testAccOrganizationPolicyConfig_boolean(org, false),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy(t, "bool", false),
			},
			{
				// Test update from not enforced to enforced
				Config: testAccOrganizationPolicyConfig_boolean(org, true),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy(t, "bool", true),
			},
			{
				ResourceName:      "google_organization_policy.bool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testAccOrganizationPolicy_list_allowAll(t *testing.T) {
	org := envvar.GetTestOrgTargetFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicyConfig_list_allowAll(org),
				Check:  testAccCheckGoogleOrganizationListPolicyAll(t, "list", "ALLOW"),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationPolicy_list_allowSome(t *testing.T) {
	org := envvar.GetTestOrgTargetFromEnv(t)
	project := envvar.GetTestProjectFromEnv()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicyConfig_list_allowSome(org, project),
				Check:  testAccCheckGoogleOrganizationListPolicyAllowedValues(t, "list", []string{"projects/" + project, "projects/debian-cloud"}),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationPolicy_list_denySome(t *testing.T) {
	org := envvar.GetTestOrgTargetFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicyConfig_list_denySome(org),
				Check:  testAccCheckGoogleOrganizationListPolicyDeniedValues(t, "list", DENIED_ORG_POLICIES),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationPolicy_list_update(t *testing.T) {
	org := envvar.GetTestOrgTargetFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicyConfig_list_allowAll(org),
				Check:  testAccCheckGoogleOrganizationListPolicyAll(t, "list", "ALLOW"),
			},
			{
				Config: testAccOrganizationPolicyConfig_list_denySome(org),
				Check:  testAccCheckGoogleOrganizationListPolicyDeniedValues(t, "list", DENIED_ORG_POLICIES),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationPolicy_list_inheritFromParent(t *testing.T) {
	org := envvar.GetTestOrgTargetFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicyConfig_list_inheritFromParent(org),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationPolicy_restore_defaultTrue(t *testing.T) {
	org := envvar.GetTestOrgTargetFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicyConfig_restore_defaultTrue(org),
				Check:  testAccCheckGoogleOrganizationRestoreDefaultTrue(t, "restore", &cloudresourcemanager.RestoreDefault{}),
			},
			{
				ResourceName:      "google_organization_policy.restore",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleOrganizationPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_organization_policy" {
				continue
			}

			org := "organizations/" + rs.Primary.Attributes["org_id"]
			constraint := resourcemanager.CanonicalOrgPolicyConstraint(rs.Primary.Attributes["constraint"])
			policy, err := config.NewResourceManagerClient(config.UserAgent).Organizations.GetOrgPolicy(org, &cloudresourcemanager.GetOrgPolicyRequest{
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

func testAccCheckGoogleOrganizationBooleanPolicy(t *testing.T, n string, enforced bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleOrganizationPolicyTestResource(t, s, n)
		if err != nil {
			return err
		}

		if policy.BooleanPolicy.Enforced != enforced {
			return fmt.Errorf("Expected boolean policy enforcement to be '%t', got '%t'", enforced, policy.BooleanPolicy.Enforced)
		}

		return nil
	}
}

func testAccCheckGoogleOrganizationListPolicyAll(t *testing.T, n, policyType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleOrganizationPolicyTestResource(t, s, n)
		if err != nil {
			return err
		}

		if len(policy.ListPolicy.AllowedValues) > 0 || len(policy.ListPolicy.DeniedValues) > 0 {
			return fmt.Errorf("The `values` field shouldn't be set")
		}

		if policy.ListPolicy.AllValues != policyType {
			return fmt.Errorf("Expected the list policy to '%s' all values, got '%s'", policyType, policy.ListPolicy.AllValues)
		}

		return nil
	}
}

func testAccCheckGoogleOrganizationListPolicyAllowedValues(t *testing.T, n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleOrganizationPolicyTestResource(t, s, n)
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

func testAccCheckGoogleOrganizationListPolicyDeniedValues(t *testing.T, n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleOrganizationPolicyTestResource(t, s, n)
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

func testAccCheckGoogleOrganizationRestoreDefaultTrue(t *testing.T, n string, policyDefault *cloudresourcemanager.RestoreDefault) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		policy, err := getGoogleOrganizationPolicyTestResource(t, s, n)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(policy.RestoreDefault, policyDefault) {
			return fmt.Errorf("Expected the restore default '%s', instead denied, %s", policyDefault, policy.RestoreDefault)
		}

		return nil
	}
}

func getGoogleOrganizationPolicyTestResource(t *testing.T, s *terraform.State, n string) (*cloudresourcemanager.OrgPolicy, error) {
	rn := "google_organization_policy." + n
	rs, ok := s.RootModule().Resources[rn]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", rn)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	config := acctest.GoogleProviderConfig(t)

	return config.NewResourceManagerClient(config.UserAgent).Organizations.GetOrgPolicy("organizations/"+rs.Primary.Attributes["org_id"], &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: rs.Primary.Attributes["constraint"],
	}).Do()
}

func testAccOrganizationPolicyConfig_boolean(org string, enforced bool) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "bool" {
  org_id     = "%s"
  constraint = "constraints/compute.disableSerialPortAccess"

  boolean_policy {
    enforced = %t
  }
}
`, org, enforced)
}

func testAccOrganizationPolicyConfig_list_allowAll(org string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "list" {
  org_id     = "%s"
  constraint = "constraints/serviceuser.services"

  list_policy {
    allow {
      all = true
    }
  }
}
`, org)
}

func testAccOrganizationPolicyConfig_list_allowSome(org, project string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "list" {
  org_id     = "%s"
  constraint = "constraints/compute.trustedImageProjects"

  list_policy {
    allow {
      values = [
        "projects/%s",
        "projects/debian-cloud",
      ]
    }
  }
}
`, org, project)
}

func testAccOrganizationPolicyConfig_list_denySome(org string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "list" {
  org_id     = "%s"
  constraint = "serviceuser.services"

  list_policy {
    deny {
      values = [
        "doubleclicksearch.googleapis.com",
        "replicapoolupdater.googleapis.com",
      ]
    }
  }
}
`, org)
}

func testAccOrganizationPolicyConfig_list_inheritFromParent(org string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "list" {
  org_id     = "%s"
  constraint = "serviceuser.services"

  list_policy {
    deny {
      values = [
        "doubleclicksearch.googleapis.com",
        "replicapoolupdater.googleapis.com",
      ]
    }
    inherit_from_parent = true
  }
}
`, org)
}

func testAccOrganizationPolicyConfig_restore_defaultTrue(org string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "restore" {
  org_id     = "%s"
  constraint = "serviceuser.services"

  restore_policy {
    default = true
  }
}
`, org)
}
