package google

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Test that an IAM policy can be applied to a project
func TestAccProjectIamPolicy_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM policy from a data source. The application
			// merges policies, so we validate the expected state.
			{
				Config: testAccProjectAssociatePolicyBasic(pid, pname, org),
			},
			{
				ResourceName: "google_project_iam_policy.acceptance",
				ImportState:  true,
			},
		},
	})
}

// Test that an IAM policy with empty members does not cause a permadiff.
func TestAccProjectIamPolicy_emptyMembers(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIamPolicyEmptyMembers(pid, pname, org),
			},
		},
	})
}

// Test that a non-collapsed IAM policy doesn't perpetually diff
func TestAccProjectIamPolicy_expanded(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectAssociatePolicyExpanded(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamPolicyExists("google_project_iam_policy.acceptance", "data.google_iam_policy.expanded", pid),
				),
			},
		},
	})
}

// Test that an IAM policy with an audit config can be applied to a project
func TestAccProjectIamPolicy_basicAuditConfig(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "tf-acctest-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM policy from a data source. The application
			// merges policies, so we validate the expected state.
			{
				Config: testAccProjectAssociatePolicyAuditConfigBasic(pid, pname, org),
			},
			{
				ResourceName: "google_project_iam_policy.acceptance",
				ImportState:  true,
			},
		},
	})
}

// Test that a non-collapsed IAM policy with AuditConfig doesn't perpetually diff
func TestAccProjectIamPolicy_expandedAuditConfig(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "tf-acctest-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectAssociatePolicyAuditConfigExpanded(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamPolicyExists("google_project_iam_policy.acceptance", "data.google_iam_policy.expanded", pid),
				),
			},
		},
	})
}

func getStatePrimaryResource(s *terraform.State, res, expectedID string) (*terraform.InstanceState, error) {
	// Get the project resource
	resource, ok := s.RootModule().Resources[res]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", res)
	}
	if resource.Primary.Attributes["id"] != expectedID && expectedID != "" {
		return nil, fmt.Errorf("Expected project %q to match ID %q in state", resource.Primary.ID, expectedID)
	}
	return resource.Primary, nil
}

func getGoogleProjectIamPolicyFromResource(resource *terraform.InstanceState) (cloudresourcemanager.Policy, error) {
	var p cloudresourcemanager.Policy
	ps, ok := resource.Attributes["policy_data"]
	if !ok {
		return p, fmt.Errorf("Resource %q did not have a 'policy_data' attribute. Attributes were %#v", resource.ID, resource.Attributes)
	}
	if err := json.Unmarshal([]byte(ps), &p); err != nil {
		return p, fmt.Errorf("Could not unmarshal %s:\n: %v", ps, err)
	}
	return p, nil
}

func getGoogleProjectIamPolicyFromState(s *terraform.State, res, expectedID string) (cloudresourcemanager.Policy, error) {
	project, err := getStatePrimaryResource(s, res, expectedID)
	if err != nil {
		return cloudresourcemanager.Policy{}, err
	}
	return getGoogleProjectIamPolicyFromResource(project)
}

func testAccCheckGoogleProjectIamPolicyExists(projectRes, policyRes, pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		projectPolicy, err := getGoogleProjectIamPolicyFromState(s, projectRes, pid)
		if err != nil {
			return fmt.Errorf("Error retrieving IAM policy for project from state: %s", err)
		}
		policyPolicy, err := getGoogleProjectIamPolicyFromState(s, policyRes, "")
		if err != nil {
			return fmt.Errorf("Error retrieving IAM policy for data_policy from state: %s", err)
		}

		// The bindings in both policies should be identical
		if !compareBindings(projectPolicy.Bindings, policyPolicy.Bindings) {
			return fmt.Errorf("Project and data source policies do not match: project policy is %+v, data resource policy is  %+v", debugPrintBindings(projectPolicy.Bindings), debugPrintBindings(policyPolicy.Bindings))
		}

		// The audit configs in both policies should be identical
		if !compareAuditConfigs(projectPolicy.AuditConfigs, policyPolicy.AuditConfigs) {
			return fmt.Errorf("Project and data source policies do not match: project policy is %+v, data resource policy is  %+v", debugPrintAuditConfigs(projectPolicy.AuditConfigs), debugPrintAuditConfigs(policyPolicy.AuditConfigs))
		}
		return nil
	}
}

// Confirm that a project has an IAM policy with at least 1 binding
func testAccProjectExistingPolicy(pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Config)
		var err error
		originalPolicy, err = getProjectIamPolicy(pid, c)
		if err != nil {
			return fmt.Errorf("Failed to retrieve IAM Policy for project %q: %s", pid, err)
		}
		if len(originalPolicy.Bindings) == 0 {
			return fmt.Errorf("Refuse to run test against project with zero IAM Bindings. This is likely an error in the test code that is not properly identifying the IAM policy of a project.")
		}
		return nil
	}
}

func testAccProjectAssociatePolicyBasic(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_policy" "acceptance" {
  project     = google_project.acceptance.id
  policy_data = data.google_iam_policy.admin.policy_data
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "user:evanbrown@google.com",
    ]
  }
  binding {
    role = "roles/compute.instanceAdmin"
    members = [
      "user:evanbrown@google.com",
      "user:evandbrown@gmail.com",
    ]
  }
}
`, pid, name, org)
}

func testAccProjectAssociatePolicyAuditConfigBasic(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_policy" "acceptance" {
  project     = google_project.acceptance.id
  policy_data = data.google_iam_policy.admin.policy_data
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "user:evanbrown@google.com",
    ]
  }
  binding {
    role = "roles/compute.instanceAdmin"
    members = [
      "user:evanbrown@google.com",
      "user:evandbrown@gmail.com",
    ]
  }
  audit_config {
    service = "cloudkms.googleapis.com"
    audit_log_configs {
      log_type         = "DATA_READ"
      exempted_members = ["user:paddy@hashicorp.com"]
    }

    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
  audit_config {
    service = "cloudsql.googleapis.com"
    audit_log_configs {
      log_type         = "DATA_READ"
      exempted_members = ["user:paddy@hashicorp.com"]
    }

    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
}
`, pid, name, org)
}

func testAccProject_create(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}
`, pid, name, org)
}

func testAccProjectIamPolicyEmptyMembers(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_policy" "acceptance" {
  project     = google_project.acceptance.id
  policy_data = data.google_iam_policy.expanded.policy_data
}

data "google_iam_policy" "expanded" {
  binding {
    role    = "roles/viewer"
    members = []
  }
}
`, pid, name, org)
}

func testAccProjectAssociatePolicyExpanded(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_policy" "acceptance" {
  project     = google_project.acceptance.id
  policy_data = data.google_iam_policy.expanded.policy_data
}

data "google_iam_policy" "expanded" {
  binding {
    role = "roles/viewer"
    members = [
      "user:paddy@carvers.co",
    ]
  }

  binding {
    role = "roles/viewer"
    members = [
      "user:paddy@hashicorp.com",
    ]
  }
}
`, pid, name, org)
}

func testAccProjectAssociatePolicyAuditConfigExpanded(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_policy" "acceptance" {
  project     = google_project.acceptance.id
  policy_data = data.google_iam_policy.expanded.policy_data
}

data "google_iam_policy" "expanded" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "user:evanbrown@google.com",
    ]
  }
  binding {
    role = "roles/compute.instanceAdmin"
    members = [
      "user:evanbrown@google.com",
      "user:evandbrown@gmail.com",
    ]
  }
  audit_config {
    service = "cloudkms.googleapis.com"
    audit_log_configs {
      log_type         = "DATA_READ"
      exempted_members = ["user:paddy@hashicorp.com"]
    }

    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
  audit_config {
    service = "cloudkms.googleapis.com"
    audit_log_configs {
      log_type         = "DATA_READ"
      exempted_members = ["user:paddy@hashicorp.com"]
    }

    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
}
`, pid, name, org)
}
