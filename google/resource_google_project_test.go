package google

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var (
	org = multiEnvSearch([]string{
		"GOOGLE_ORG",
	})

	pname          = "Terraform Acceptance Tests"
	originalPolicy *cloudresourcemanager.Policy
)

// Test that a Project resource can be created without an organization
func TestAccGoogleProject_createWithoutOrg(t *testing.T) {
	creds := multiEnvSearch(credsEnvVars)
	if strings.Contains(creds, "iam.gserviceaccount.com") {
		t.Skip("Service accounts cannot create projects without a parent. Requires user credentials.")
	}

	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// This step creates a new project
			resource.TestStep{
				Config: testAccGoogleProject_createWithoutOrg(pid, pname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
				),
			},
		},
	})
}

// Test that a Project resource can be created and an IAM policy
// associated
func TestAccGoogleProject_create(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
		}...,
	)

	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// This step creates a new project
			resource.TestStep{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
				),
			},
		},
	})
}

// Test that a Project resource can be created with an associated
// billing account
func TestAccGoogleProject_createBilling(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
		}...,
	)

	billingId := os.Getenv("GOOGLE_BILLING_ACCOUNT")
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// This step creates a new project with a billing account
			resource.TestStep{
				Config: testAccGoogleProject_createBilling(pid, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasBillingAccount("google_project.acceptance", pid, billingId),
				),
			},
		},
	})
}

// Test that a Project resource can be created with labels
func TestAccGoogleProject_createLabels(t *testing.T) {
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleProject_createLabels(pid, pname, org, "test", "that"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasLabels("google_project.acceptance", pid, map[string]string{"test": "that"}),
				),
			},
		},
	})
}

// Test that a Project resource can be created and updated
// with billing account information
func TestAccGoogleProject_updateBilling(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
			"GOOGLE_BILLING_ACCOUNT_2",
		}...,
	)

	billingId := os.Getenv("GOOGLE_BILLING_ACCOUNT")
	billingId2 := os.Getenv("GOOGLE_BILLING_ACCOUNT_2")
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// This step creates a new project without a billing account
			resource.TestStep{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
				),
			},
			// Update to include a billing account
			resource.TestStep{
				Config: testAccGoogleProject_createBilling(pid, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasBillingAccount("google_project.acceptance", pid, billingId),
				),
			},
			// Update to a different  billing account
			resource.TestStep{
				Config: testAccGoogleProject_createBilling(pid, pname, org, billingId2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasBillingAccount("google_project.acceptance", pid, billingId2),
				),
			},
			// Unlink the billing account
			resource.TestStep{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasBillingAccount("google_project.acceptance", pid, ""),
				),
			},
		},
	})
}

// Test that a Project resource merges the IAM policies that already
// exist, and won't lock people out.
func TestAccGoogleProject_merge(t *testing.T) {
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// when policy_data is set, merge
			{
				Config: testAccGoogleProject_toMerge(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
					testAccCheckGoogleProjectHasMoreBindingsThan(pid, 1),
				),
			},
			// when policy_data is unset, restore to what it was
			{
				Config: testAccGoogleProject_mergeEmpty(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
					testAccCheckGoogleProjectHasMoreBindingsThan(pid, 0),
				),
			},
		},
	})
}

// Test that a Project resource can be updated with labels
func TestAccGoogleProject_updateLabels(t *testing.T) {
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// create project without labels
			{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
				),
			},
			// update project with labels
			{
				Config: testAccGoogleProject_updateLabels(pid, pname, org, map[string]string{"label": "label-value"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
					testAccCheckGoogleProjectHasLabels("google_project.acceptance", pid, map[string]string{"label": "label-value"}),
				),
			},
			// update project with other labels
			{
				Config: testAccGoogleProject_updateLabels(pid, pname, org, map[string]string{"empty-label": ""}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
					testAccCheckGoogleProjectHasLabels("google_project.acceptance", pid, map[string]string{"empty-label": ""}),
				),
			},
			// update project delete labels
			{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
					testAccCheckGoogleProjectHasNoLabels("google_project.acceptance", pid),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectExists(r, pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.ID != pid {
			return fmt.Errorf("Expected project %q to match ID %q in state", pid, rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckGoogleProjectHasBillingAccount(r, pid, billingId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		// State should match expected
		if rs.Primary.Attributes["billing_account"] != billingId {
			return fmt.Errorf("Billing ID in state (%s) does not match expected value (%s)", rs.Primary.Attributes["billing_account"], billingId)
		}

		// Actual value in API should match state and expected
		// Read the billing account
		config := testAccProvider.Meta().(*Config)
		ba, err := config.clientBilling.Projects.GetBillingInfo(prefixedProject(pid)).Do()
		if err != nil {
			return fmt.Errorf("Error reading billing account for project %q: %v", prefixedProject(pid), err)
		}
		if billingId != strings.TrimPrefix(ba.BillingAccountName, "billingAccounts/") {
			return fmt.Errorf("Billing ID returned by API (%s) did not match expected value (%s)", ba.BillingAccountName, billingId)
		}
		return nil
	}
}

func testAccCheckGoogleProjectHasMoreBindingsThan(pid string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, err := getProjectIamPolicy(pid, testAccProvider.Meta().(*Config))
		if err != nil {
			return err
		}
		if len(policy.Bindings) <= count {
			return fmt.Errorf("Expected more than %d bindings, got %d: %#v", count, len(policy.Bindings), policy.Bindings)
		}
		return nil
	}
}

func testAccCheckGoogleProjectHasLabels(r, pid string, expected map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		// State should have the same number of labels
		if rs.Primary.Attributes["labels.%"] != strconv.Itoa(len(expected)) {
			return fmt.Errorf("Expected %d labels, got %s", len(expected), rs.Primary.Attributes["labels.%"])
		}

		// Actual value in API should match state and expected
		config := testAccProvider.Meta().(*Config)

		found, err := config.clientResourceManager.Projects.Get(pid).Do()
		if err != nil {
			return err
		}

		actual := found.Labels
		if !reflect.DeepEqual(actual, expected) {
			// Determine only the different attributes
			for k, v := range expected {
				if av, ok := actual[k]; ok && v == av {
					delete(expected, k)
					delete(actual, k)
				}
			}

			spewConf := spew.NewDefaultConfig()
			spewConf.SortKeys = true
			return fmt.Errorf(
				"Labels not equivalent. Difference is shown below. Top is actual, bottom is expected."+
					"\n\n%s\n\n%s",
				spewConf.Sdump(actual), spewConf.Sdump(expected),
			)
		}
		return nil
	}
}

func testAccCheckGoogleProjectHasNoLabels(r, pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		// State should have zero labels
		if rs.Primary.Attributes["labels.%"] != "0" {
			return fmt.Errorf("Expected 0 labels, got %s", rs.Primary.Attributes["labels.%"])
		}

		// Actual value in API should match state and expected
		config := testAccProvider.Meta().(*Config)

		found, err := config.clientResourceManager.Projects.Get(pid).Do()
		if err != nil {
			return err
		}

		spewConf := spew.NewDefaultConfig()
		spewConf.SortKeys = true
		if found.Labels != nil {
			return fmt.Errorf("Labels should be empty. Actual \n%s", spewConf.Sdump(found.Labels))
		}
		return nil
	}
}

func testAccGoogleProject_toMerge(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
    project_id = "%s"
    name = "%s"
    org_id = "%s"
}

resource "google_project_iam_policy" "acceptance" {
    project = "${google_project.acceptance.project_id}"
    policy_data = "${data.google_iam_policy.acceptance.policy_data}"
}

data "google_iam_policy" "acceptance" {
    binding {
        role = "roles/storage.objectViewer"
	members = [
	  "user:evanbrown@google.com",
	]
    }
}`, pid, name, org)
}

func testAccGoogleProject_mergeEmpty(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
    project_id = "%s"
    name = "%s"
    org_id = "%s"
}`, pid, name, org)
}

func testAccGoogleProject_createLabels(pid, name, org, key, value string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
    project_id = "%s"
    name = "%s"
    org_id = "%s"
	labels {
		"%s" = "%s"
	}
}`, pid, name, org, key, value)
}

func testAccGoogleProject_updateLabels(pid, name, org string, labels map[string]string) string {
	r := fmt.Sprintf(`
resource "google_project" "acceptance" {
    project_id = "%s"
    name = "%s"
    org_id = "%s"
	labels {`, pid, name, org)

	l := ""
	for key, value := range labels {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}

func skipIfEnvNotSet(t *testing.T, envs ...string) {
	for _, k := range envs {
		if os.Getenv(k) == "" {
			t.Skipf("Environment variable %s is not set", k)
		}
	}
}
