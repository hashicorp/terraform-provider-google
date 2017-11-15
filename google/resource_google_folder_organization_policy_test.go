package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
	"reflect"
	"testing"
)

func TestAccGoogleFolderOrganizationPolicy_boolean(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test creation of an enforced boolean policy
				Config: testAccGoogleFolderOrganizationPolicy_boolean(org, true),
				Check:  testAccCheckGoogleFolderOrganizationBooleanPolicy("bool", true),
			},
			{
				// Test update from enforced to not
				Config: testAccGoogleFolderOrganizationPolicy_boolean(org, false),
				Check:  testAccCheckGoogleFolderOrganizationBooleanPolicy("bool", false),
			},
			{
				Config:  " ",
				Destroy: true,
			},
			{
				// Test creation of a not enforced boolean policy
				Config: testAccGoogleFolderOrganizationPolicy_boolean(org, false),
				Check:  testAccCheckGoogleFolderOrganizationBooleanPolicy("bool", false),
			},
			{
				// Test update from not enforced to enforced
				Config: testAccGoogleFolderOrganizationPolicy_boolean(org, true),
				Check:  testAccCheckGoogleFolderOrganizationBooleanPolicy("bool", true),
			},
		},
	})
}

func TestAccGoogleFolderOrganizationPolicy_list_allowAll(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleFolderOrganizationPolicy_list_allowAll(org),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyAll("list", "ALLOW"),
			},
		},
	})
}

func TestAccGoogleFolderOrganizationPolicy_list_allowSome(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_ORG")
	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleFolderOrganizationPolicy_list_allowSome(org, project),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyAllowedValues("list", []string{project}),
			},
		},
	})
}

func TestAccGoogleFolderOrganizationPolicy_list_denySome(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleFolderOrganizationPolicy_list_denySome(org),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyDeniedValues("list", DENIED_ORG_POLICIES),
			},
		},
	})
}

func TestAccGoogleFolderOrganizationPolicy_list_update(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleFolderOrganizationPolicy_list_allowAll(org),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyAll("list", "ALLOW"),
			},
			{
				Config: testAccGoogleFolderOrganizationPolicy_list_denySome(org),
				Check:  testAccCheckGoogleFolderOrganizationListPolicyDeniedValues("list", DENIED_ORG_POLICIES),
			},
		},
	})
}

func testAccCheckGoogleFolderOrganizationPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_folder_organization_policy" {
			continue
		}

		folder := rs.Primary.Attributes["folder"]
		constraint := canonicalOrgPolicyConstraint(rs.Primary.Attributes["constraint"])
		policy, err := config.clientResourceManager.Folders.GetOrgPolicy(folder, &cloudresourcemanager.GetOrgPolicyRequest{
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

func testAccCheckGoogleFolderOrganizationBooleanPolicy(n string, enforced bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleFolderOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if policy.BooleanPolicy.Enforced != enforced {
			return fmt.Errorf("Expected boolean policy enforcement to be '%t', got '%t'", enforced, policy.BooleanPolicy.Enforced)
		}

		return nil
	}
}

func testAccCheckGoogleFolderOrganizationListPolicyAll(n, policyType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleFolderOrganizationPolicyTestResource(s, n)
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

func testAccCheckGoogleFolderOrganizationListPolicyAllowedValues(n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleFolderOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(policy.ListPolicy.AllowedValues, values) {
			return fmt.Errorf("Expected the list policy to allow '%s', instead allowed '%s'", values, policy.ListPolicy.AllowedValues)
		}

		return nil
	}
}

func testAccCheckGoogleFolderOrganizationListPolicyDeniedValues(n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleFolderOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(policy.ListPolicy.DeniedValues, values) {
			return fmt.Errorf("Expected the list policy to deny '%s', instead denied '%s'", values, policy.ListPolicy.DeniedValues)
		}

		return nil
	}
}

func getGoogleFolderOrganizationPolicyTestResource(s *terraform.State, n string) (*cloudresourcemanager.OrgPolicy, error) {
	rn := "google_folder_organization_policy." + n
	rs, ok := s.RootModule().Resources[rn]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", rn)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	config := testAccProvider.Meta().(*Config)

	return config.clientResourceManager.Folders.GetOrgPolicy(rs.Primary.Attributes["folder"], &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: rs.Primary.Attributes["constraint"],
	}).Do()
}

func testAccGoogleFolderOrganizationPolicy_boolean(org string, enforced bool) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "bool" {
	folder     = "${google_folder.orgpolicy.name}"
	constraint = "constraints/compute.disableSerialPortAccess"

	boolean_policy {
		enforced = %t
	}
}
`, acctest.RandomWithPrefix("tf-test"), "organizations/"+org, enforced)
}

func testAccGoogleFolderOrganizationPolicy_list_allowAll(org string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "list" {
	folder     = "${google_folder.orgpolicy.name}"
	constraint = "constraints/serviceuser.services"

	list_policy {
		allow {
			all = true
		}
	}
}
`, acctest.RandomWithPrefix("tf-test"), "organizations/"+org)
}

func testAccGoogleFolderOrganizationPolicy_list_allowSome(org, project string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "list" {
	folder     = "${google_folder.orgpolicy.name}"
	constraint = "constraints/compute.trustedImageProjects"

	list_policy {
		allow {
			values = [
				"%s",
			]
		}
  }
}
`, acctest.RandomWithPrefix("tf-test"), "organizations/"+org, project)
}

func testAccGoogleFolderOrganizationPolicy_list_denySome(org string) string {
	return fmt.Sprintf(`
resource "google_folder" "orgpolicy" {
  display_name = "%s"
  parent       = "%s"
}

resource "google_folder_organization_policy" "list" {
	folder     = "${google_folder.orgpolicy.name}"
 	constraint = "serviceuser.services"

  	list_policy {
		deny {
			values = [
				"maps-ios-backend.googleapis.com",
				"placesios.googleapis.com",
			]
		}
	}
}
`, acctest.RandomWithPrefix("tf-test"), "organizations/"+org)
}
