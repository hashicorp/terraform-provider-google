package google

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var DENIED_ORG_POLICIES = []string{
	"doubleclicksearch.googleapis.com",
	"replicapoolupdater.googleapis.com",
}

// Since each test here is acting on the same organization, run the tests serially to
// avoid race conditions and aborted operations.

func TestAccOrganizationPolicy_boolean(t *testing.T) {
	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				// Test creation of an enforced boolean policy
				Config: testAccOrganizationPolicy_boolean(org, true),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", true),
			},
			{
				// Test update from enforced to not
				Config: testAccOrganizationPolicy_boolean(org, false),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", false),
			},
			{
				Config:  " ",
				Destroy: true,
			},
			{
				// Test creation of a not enforced boolean policy
				Config: testAccOrganizationPolicy_boolean(org, false),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", false),
			},
			{
				// Test update from not enforced to enforced
				Config: testAccOrganizationPolicy_boolean(org, true),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", true),
			},
			{
				ResourceName:      "google_organization_policy.bool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func TestAccOrganizationPolicy_list_allowAll(t *testing.T) {
	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicy_list_allowAll(org),
				Check:  testAccCheckGoogleOrganizationListPolicyAll("list", "ALLOW"),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationPolicy_list_allowSome(t *testing.T) {
	org := getTestOrgFromEnv(t)
	project := getTestProjectFromEnv()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicy_list_allowSome(org, project),
				Check:  testAccCheckGoogleOrganizationListPolicyAllowedValues("list", []string{"projects/" + project}),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationPolicy_list_denySome(t *testing.T) {
	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicy_list_denySome(org),
				Check:  testAccCheckGoogleOrganizationListPolicyDeniedValues("list", DENIED_ORG_POLICIES),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationPolicy_list_update(t *testing.T) {
	org := getTestOrgFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationPolicy_list_allowAll(org),
				Check:  testAccCheckGoogleOrganizationListPolicyAll("list", "ALLOW"),
			},
			{
				Config: testAccOrganizationPolicy_list_denySome(org),
				Check:  testAccCheckGoogleOrganizationListPolicyDeniedValues("list", DENIED_ORG_POLICIES),
			},
			{
				ResourceName:      "google_organization_policy.list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGoogleOrganizationPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_organization_policy" {
			continue
		}

		org := "organizations/" + rs.Primary.Attributes["org_id"]
		constraint := canonicalOrgPolicyConstraint(rs.Primary.Attributes["constraint"])
		policy, err := config.clientResourceManager.Organizations.GetOrgPolicy(org, &cloudresourcemanager.GetOrgPolicyRequest{
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

func testAccCheckGoogleOrganizationBooleanPolicy(n string, enforced bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleOrganizationPolicyTestResource(s, n)
		if err != nil {
			return err
		}

		if policy.BooleanPolicy.Enforced != enforced {
			return fmt.Errorf("Expected boolean policy enforcement to be '%t', got '%t'", enforced, policy.BooleanPolicy.Enforced)
		}

		return nil
	}
}

func testAccCheckGoogleOrganizationListPolicyAll(n, policyType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleOrganizationPolicyTestResource(s, n)
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

func testAccCheckGoogleOrganizationListPolicyAllowedValues(n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleOrganizationPolicyTestResource(s, n)
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

func testAccCheckGoogleOrganizationListPolicyDeniedValues(n string, values []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getGoogleOrganizationPolicyTestResource(s, n)
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

func getGoogleOrganizationPolicyTestResource(s *terraform.State, n string) (*cloudresourcemanager.OrgPolicy, error) {
	rn := "google_organization_policy." + n
	rs, ok := s.RootModule().Resources[rn]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", rn)
	}

	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No ID is set")
	}

	config := testAccProvider.Meta().(*Config)

	return config.clientResourceManager.Organizations.GetOrgPolicy("organizations/"+rs.Primary.Attributes["org_id"], &cloudresourcemanager.GetOrgPolicyRequest{
		Constraint: rs.Primary.Attributes["constraint"],
	}).Do()
}

func testAccOrganizationPolicy_boolean(org string, enforced bool) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "bool" {
	org_id = "%s"
	constraint = "constraints/compute.disableSerialPortAccess"

	boolean_policy {
		enforced = %t
	}
}
`, org, enforced)
}

func testAccOrganizationPolicy_list_allowAll(org string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "list" {
	org_id = "%s"
	constraint = "constraints/serviceuser.services"

	list_policy {
		allow {
			all = true
		}
	}
}
`, org)
}

func testAccOrganizationPolicy_list_allowSome(org, project string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "list" {
	org_id = "%s"
	constraint = "constraints/compute.trustedImageProjects"

	list_policy {
		allow {
			values = [
				"projects/%s",
			]
		}
  }
}
`, org, project)
}

func testAccOrganizationPolicy_list_denySome(org string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "list" {
	org_id = "%s"
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
