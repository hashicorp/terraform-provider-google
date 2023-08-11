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

func TestAccFolderOrganizationPolicy_boolean(t *testing.T) {
	t.Parallel()

	folder := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	org := envvar.GetTestOrgFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleFolderOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Test creation of an enforced boolean policy
				Config: testAccFolderOrganizationPolicy_boolean(org, folder, true),
				Check:  testAccCheckGoogleFolderOrganizationBooleanPolicy(t, "bool", true),
			},
			{
				// Test update from enforced to not
				Config: testAccFolderOrganizationPolicy_boolean(org, folder, false),
				Check:  testAccCheckGoogleFolderOrganizationBooleanPolicy(t, "bool", false),
			},
			{
				Config:  " ",
				Destroy: true,
			},
			{
				// Test creation of a not enforced boolean policy
				Config: testAccFolderOrganizationPolicy_boolean(org, folder, false),
				Check:  testAccCheckGoogleFolderOrganizationBooleanPolicy(t, "bool", false),
			},
			{
				// Test update from not enforced to enforced
				Config: testAccFolderOrganizationPolicy_boolean(org, folder, true),
				Check:  testAccCheckGoogleFolderOrganizationBooleanPolicy(t, "bool", true),
			},
		},
	})
}

func TestAccFolderOrganizationPolicy_list_allowAll(t *testing.T) {
	t.Parallel()

	folder := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	org := envvar.GetTestOrgFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleFolderOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFolderOrganizationPolicy_list_allowAll(org, folder),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyAll(t, "list", "ALLOW"),
			},
			{
				ResourceName:      "google_folder_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFolderOrganizationPolicy_list_allowSome(t *testing.T) {
	t.Parallel()

	folder := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	org := envvar.GetTestOrgFromEnv(t)
	project := envvar.GetTestProjectFromEnv()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleFolderOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFolderOrganizationPolicy_list_allowSome(org, folder, project),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyAllowedValues(t, "list", []string{"projects/" + project}),
			},
			{
				ResourceName:      "google_folder_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFolderOrganizationPolicy_list_denySome(t *testing.T) {
	t.Parallel()

	folder := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	org := envvar.GetTestOrgFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleFolderOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFolderOrganizationPolicy_list_denySome(org, folder),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyDeniedValues(t, "list", DENIED_ORG_POLICIES),
			},
			{
				ResourceName:      "google_folder_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFolderOrganizationPolicy_list_update(t *testing.T) {
	t.Parallel()

	folder := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	org := envvar.GetTestOrgFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleFolderOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFolderOrganizationPolicy_list_allowAll(org, folder),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyAll(t, "list", "ALLOW"),
			},
			{
				Config: testAccFolderOrganizationPolicy_list_denySome(org, folder),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyDeniedValues(t, "list", DENIED_ORG_POLICIES),
			},
			{
				ResourceName:      "google_folder_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFolderOrganizationPolicy_restore_defaultTrue(t *testing.T) {
	t.Parallel()

	folder := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	org := envvar.GetTestOrgFromEnv(t)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleOrganizationPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFolderOrganizationPolicy_restore_defaultTrue(org, folder),
				Check:  getGoogleFolderOrganizationRestoreDefaultTrue(t, "restore", &cloudresourcemanager.RestoreDefault{}),
			},
			{
				ResourceName:      "google_folder_organization_policy.restore",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleFolderOrganizationPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_folder_organization_policy" {
				continue
			}

			folder := resourcemanager.CanonicalFolderId(rs.Primary.Attributes["folder"])
			constraint := resourcemanager.CanonicalOrgPolicyConstraint(rs.Primary.Attributes["constraint"])
			policy, err := config.NewResourceManagerClient(config.UserAgent).Folders.GetOrgPolicy(folder, &cloudresourcemanager.GetOrgPolicyRequest{
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

func testAccCheckGoogleFolderOrganizationBooleanPolicy(t *testing.T, n string, enforced bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleFolderOrganizationPolicyTestResource(t, s, n)
		if err != nil {
			return err
		}

		if policy.BooleanPolicy.Enforced != enforced {
			return fmt.Errorf("Expected boolean policy enforcement to be '%t', got '%t'", enforced, policy.BooleanPolicy.Enforced)
		}

		return nil
	}
}

func testAccCheckGoogleFolderOrganizationListPolicyAll(t *testing.T, n, policyType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleFolderOrganizationPolicyTestResource(t, s, n)
		if err != nil {
			return err
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

func testAccCheckGoogleFolderOrganizationListPolicyAllowedValues(t *testing.T, n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleFolderOrganizationPolicyTestResource(t, s, n)
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

func testAccCheckGoogleFolderOrganizationListPolicyDeniedValues(t *testing.T, n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleFolderOrganizationPolicyTestResource(t, s, n)
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

func getGoogleFolderOrganizationRestoreDefaultTrue(t *testing.T, n string, policyDefault *cloudresourcemanager.RestoreDefault) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		policy, err := getGoogleFolderOrganizationPolicyTestResource(t, s, n)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(policy.RestoreDefault, policyDefault) {
			return fmt.Errorf("Expected the restore default '%s', instead denied, %s", policyDefault, policy.RestoreDefault)
		}

		return nil
	}
}

func getGoogleFolderOrganizationPolicyTestResource(t *testing.T, s *terraform.State, n string) (*cloudresourcemanager.OrgPolicy, error) {
	rn := "google_folder_organization_policy." + n
	rs, ok := s.RootModule().Resources[rn]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", rn)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	config := acctest.GoogleProviderConfig(t)
	folder := resourcemanager.CanonicalFolderId(rs.Primary.Attributes["folder"])

	return config.NewResourceManagerClient(config.UserAgent).Folders.GetOrgPolicy(folder, &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: rs.Primary.Attributes["constraint"],
	}).Do()
}

func testAccFolderOrganizationPolicy_boolean(org, folder string, enforced bool) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "bool" {
  # Test numeric folder ID.
  folder     = replace(google_folder.orgpolicy.name, "folders/", "")
  constraint = "constraints/compute.disableSerialPortAccess"

  boolean_policy {
    enforced = %t
  }
}
`, folder, "organizations/"+org, enforced)
}

func testAccFolderOrganizationPolicy_list_allowAll(org, folder string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "list" {
  folder     = google_folder.orgpolicy.name
  constraint = "constraints/serviceuser.services"

  list_policy {
    allow {
      all = true
    }
  }
}
`, folder, "organizations/"+org)
}

func testAccFolderOrganizationPolicy_list_allowSome(org, folder, project string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "list" {
  folder     = google_folder.orgpolicy.name
  constraint = "constraints/compute.trustedImageProjects"

  list_policy {
    allow {
      values = ["projects/%s"]
    }
  }
}
`, folder, "organizations/"+org, project)
}

func testAccFolderOrganizationPolicy_list_denySome(org, folder string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "list" {
  folder     = google_folder.orgpolicy.name
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
`, folder, "organizations/"+org)
}

func testAccFolderOrganizationPolicy_restore_defaultTrue(org, folder string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "restore" {
  folder     = google_folder.orgpolicy.name
  constraint = "serviceuser.services"

  restore_policy {
    default = true
  }
}
`, folder, "organizations/"+org)
}
