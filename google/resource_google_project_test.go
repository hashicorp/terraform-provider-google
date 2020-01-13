package google

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

var (
	pname          = "Terraform Acceptance Tests"
	originalPolicy *cloudresourcemanager.Policy
	testPrefix     = "tf-test"
)

func init() {
	resource.AddTestSweepers("Project", &resource.Sweeper{
		Name: "Project",
		F:    testSweepProject,
	})
}

func testSweepProject(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	token := ""
	for paginate := true; paginate; {
		// Filter for projects with test prefix
		filter := "id:" + testPrefix + "*"
		found, err := config.clientResourceManager.Projects.List().Filter(filter).PageToken(token).Do()
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error listing projects: %s", err)
			return nil
		}
		for _, project := range found.Projects {
			if project.LifecycleState != "ACTIVE" {
				continue
			}
			log.Printf("[INFO][SWEEPER_LOG] Sweeping Project id: %s", project.ProjectId)

			_, err := config.clientResourceManager.Projects.Delete(project.ProjectId).Do()

			if err != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error, failed to delete project %s: %s", project.Name, err)
				continue
			}
		}
		token = found.NextPageToken
		paginate = token != ""
	}

	return nil
}

// Test that a Project resource can be created without an organization
func TestAccProject_createWithoutOrg(t *testing.T) {
	t.Parallel()

	creds := multiEnvSearch(credsEnvVars)
	if strings.Contains(creds, "iam.gserviceaccount.com") {
		t.Skip("Service accounts cannot create projects without a parent. Requires user credentials.")
	}

	pid := acctest.RandomWithPrefix(testPrefix)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// This step creates a new project
			{
				Config: testAccProject_createWithoutOrg(pid, pname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
				),
			},
		},
	})
}

// Test that a Project resource can be created and an IAM policy
// associated
func TestAccProject_create(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := acctest.RandomWithPrefix(testPrefix)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// This step creates a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
				),
			},
		},
	})
}

// Test that a Project resource can be created with an associated
// billing account
func TestAccProject_billing(t *testing.T) {
	t.Parallel()
	org := getTestOrgFromEnv(t)
	skipIfEnvNotSet(t, "GOOGLE_BILLING_ACCOUNT_2")
	billingId2 := os.Getenv("GOOGLE_BILLING_ACCOUNT_2")
	billingId := getTestBillingAccountFromEnv(t)
	pid := acctest.RandomWithPrefix(testPrefix)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// This step creates a new project with a billing account
			{
				Config: testAccProject_createBilling(pid, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasBillingAccount("google_project.acceptance", pid, billingId),
				),
			},
			// Make sure import supports billing account
			{
				ResourceName:            "google_project.acceptance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_delete"},
			},
			// Update to a different  billing account
			{
				Config: testAccProject_createBilling(pid, pname, org, billingId2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasBillingAccount("google_project.acceptance", pid, billingId2),
				),
			},
			// Unlink the billing account
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasBillingAccount("google_project.acceptance", pid, ""),
				),
			},
		},
	})
}

// Test that a Project resource can be created with labels
func TestAccProject_labels(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := acctest.RandomWithPrefix(testPrefix)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProject_labels(pid, pname, org, map[string]string{"test": "that"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectHasLabels("google_project.acceptance", pid, map[string]string{"test": "that"}),
				),
			},
			// Make sure import supports labels
			{
				ResourceName:            "google_project.acceptance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_delete"},
			},
			// update project with labels
			{
				Config: testAccProject_labels(pid, pname, org, map[string]string{"label": "label-value"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
					testAccCheckGoogleProjectHasLabels("google_project.acceptance", pid, map[string]string{"label": "label-value"}),
				),
			},
			// update project delete labels
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
					testAccCheckGoogleProjectHasNoLabels("google_project.acceptance", pid),
				),
			},
		},
	})
}

func TestAccProject_deleteDefaultNetwork(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := acctest.RandomWithPrefix(testPrefix)
	billingId := getTestBillingAccountFromEnv(t)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProject_deleteDefaultNetwork(pid, pname, org, billingId),
			},
		},
	})
}

func TestAccProject_parentFolder(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := acctest.RandomWithPrefix(testPrefix)
	folderDisplayName := testPrefix + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProject_parentFolder(pid, pname, folderDisplayName, org),
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

		projectId := fmt.Sprintf("projects/%s", pid)
		if rs.Primary.ID != projectId {
			return fmt.Errorf("Expected project %q to match ID %q in state", projectId, rs.Primary.ID)
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
		if v, ok := rs.Primary.Attributes["labels.%"]; ok && v != "0" {
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

func testAccProject_createWithoutOrg(pid, name string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
}
`, pid, name)
}

func testAccProject_createBilling(pid, name, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}
`, pid, name, org, billing)
}

func testAccProject_labels(pid, name, org string, labels map[string]string) string {
	r := fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
  labels = {`, pid, name, org)

	l := ""
	for key, value := range labels {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}

func testAccProject_deleteDefaultNetwork(pid, name, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id          = "%s"
  name                = "%s"
  org_id              = "%s"
  billing_account     = "%s" # requires billing to enable compute API
  auto_create_network = false
}
`, pid, name, org, billing)
}

func testAccProject_parentFolder(pid, projectName, folderName, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"

  # ensures we can set both org_id and folder_id as long as only one is not empty.
  org_id    = ""
  folder_id = google_folder.folder1.id
}

resource "google_folder" "folder1" {
  display_name = "%s"
  parent       = "organizations/%s"
}
`, pid, projectName, folderName, org)
}

func skipIfEnvNotSet(t *testing.T, envs ...string) {
	if t == nil {
		log.Printf("[DEBUG] Not running inside of test - skip skipping")
		return
	}

	for _, k := range envs {
		if os.Getenv(k) == "" {
			t.Skipf("Environment variable %s is not set", k)
		}
	}
}
