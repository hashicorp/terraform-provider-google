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

var SOME_ORG_POLICIES = []string{"compute.googleapis.com", "cloudresourcemanager.googleapis.com"}

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
				Config: testAccGoogleOrganizationPolicy_listAll(org, "allow"),
				Check:  testAccCheckGoogleOrganizationListPolicyAll("listAll", "ALLOW"),
			},
		},
	})
}

func TestAccGoogleOrganizationPolicy_list_allowSome(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_listSome(org, "allow"),
				Check:  testAccCheckGoogleOrganizationListPolicyAllowedValues("listSome", SOME_ORG_POLICIES),
			},
		},
	})
}

func TestAccGoogleOrganizationPolicy_list_denyAll(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")
	org := os.Getenv("GOOGLE_ORG")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleOrganizationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleOrganizationPolicy_listAll(org, "deny"),
				Check:  testAccCheckGoogleOrganizationListPolicyAll("listAll", "DENY"),
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
				Config: testAccGoogleOrganizationPolicy_listSome(org, "deny"),
				Check:  testAccCheckGoogleOrganizationListPolicyDeniedValues("listSome", SOME_ORG_POLICIES),
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
				Config: testAccGoogleOrganizationPolicy_listAll(org, "allow"),
				Check:  testAccCheckGoogleOrganizationListPolicyAll("listAll", "ALLOW"),
			},
			{
				Config: testAccGoogleOrganizationPolicy_listSome(org, "deny"),
				Check:  testAccCheckGoogleOrganizationListPolicyDeniedValues("listSome", SOME_ORG_POLICIES),
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
		constraint := rs.Primary.Attributes["constraint"]
		_, err := config.clientResourceManager.Organizations.GetOrgPolicy(org, &cloudresourcemanager.GetOrgPolicyRequest{
			Constraint: constraint,
		}).Do()

		if err != nil {
			return fmt.Errorf("Org policy with constraint '%s' hasn't been deleted", constraint)
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

func testAccGoogleOrganizationPolicy_listAll(org, policyType string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "listAll" {
	org_id = "%s"
	constraint = "constraints/serviceuser.services"

	list_policy {
		%s {
			all = true
		}
	}
}
`, org, policyType)
}

func testAccGoogleOrganizationPolicy_listSome(org, policyType string) string {
	return fmt.Sprintf(`
resource "google_organization_policy" "listSome" {
	org_id = "%s"
 	constraint = "serviceuser.services"

  	list_policy {
		%s {
			values = [
        		"cloudresourcemanager.googleapis.com",
        		"compute.googleapis.com",
      		]
		}
  }
}
`, org, policyType)
}
