package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
	"os"
	"reflect"
	"testing"
)

var DENIED_ORG_POLICIES = []string{
	"maps-ios-backend.googleapis.com",
	"placesios.googleapis.com",
}

func TestAccGoogleOrganizationPolicy_boolean_enforced(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_boolean(org, true),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", true),
			},
		},
	})

}

func TestAccGoogleOrganizationPolicy_boolean_notEnforced(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_boolean(org, false),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", false),
			},
		},
	})
}

func TestAccGoogleOrganizationPolicy_boolean_update(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_boolean(org, true),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", true),
			},
			{
				Config: testAccGoogleOrganizationPolicy_boolean(org, false),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", false),
			},
			{
				Config: testAccGoogleOrganizationPolicy_boolean(org, true),
				Check:  testAccCheckGoogleOrganizationBooleanPolicy("bool", true),
			},
		},
	})
}

func TestAccGoogleOrganizationPolicy_list_allowAll(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_list_allowAll(org),
				Check:  testAccCheckGoogleOrganizationListPolicyAll("listAll", "ALLOW"),
			},
		},
	})
}

func TestAccGoogleOrganizationPolicy_list_allowSome(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")
	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_list_allowSome(org, project),
				Check:  testAccCheckGoogleOrganizationListPolicyAllowedValues("listSome", []string{project}),
			},
		},
	})
}

func TestAccGoogleOrganizationPolicy_list_denySome(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_list_denySome(org),
				Check:  testAccCheckGoogleOrganizationListPolicyDeniedValues("listSome", DENIED_ORG_POLICIES),
			},
		},
	})
}

func TestAccGoogleOrganizationPolicy_list_update(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_list_allowAll(org),
				Check:  testAccCheckGoogleOrganizationListPolicyAll("listAll", "ALLOW"),
			},
			{
				Config: testAccGoogleOrganizationPolicy_list_denySome(org),
				Check:  testAccCheckGoogleOrganizationListPolicyDeniedValues("listSome", DENIED_ORG_POLICIES),
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

func testAccGoogleOrganizationPolicy_boolean(org string, enforced bool) string {
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

func testAccGoogleOrganizationPolicy_list_allowAll(org string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "listAll" {
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

func testAccGoogleOrganizationPolicy_list_allowSome(org, project string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "listSome" {
	org_id = "%s"
	constraint = "constraints/compute.trustedImageProjects"

	list_policy {
		allow {
			values = [
				"%s",
			]
		}
  }
}
`, org, project)
}

func testAccGoogleOrganizationPolicy_list_denySome(org string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "listSome" {
	org_id = "%s"
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
`, org)
}
