package google

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

func TestAccFolderIamPolicy_basic(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderIamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFolderIamPolicy_basic(folderDisplayName, parent, "roles/viewer", "user:admin@hashicorptest.com"),
			},
			{
				ResourceName:      "google_folder_iam_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFolderIamPolicy_basic2(folderDisplayName, parent, "roles/editor", "user:admin@hashicorptest.com", "roles/viewer", "user:admin@hashicorptest.com"),
			},
			{
				ResourceName:      "google_folder_iam_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test that an IAM policy with an audit config can be applied to a folder
func TestAccFolderIamPolicy_basicAuditConfig(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply an IAM policy from a data source. The application
			// merges policies, so we validate the expected state.
			{
				Config: testAccFolderAssociatePolicyAuditConfigBasic(folderDisplayName, parent),
			},
			{
				ResourceName: "google_folder_iam_policy.acceptance",
				ImportState:  true,
			},
		},
	})
}

// Test that a non-collapsed IAM policy with AuditConfig doesn't perpetually diff
func TestAccFolderIamPolicy_expandedAuditConfig(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFolderAssociatePolicyAuditConfigExpanded(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamPolicyExists("google_folder_iam_policy.acceptance", "data.google_iam_policy.expanded"),
				),
			},
		},
	})
}

func testAccCheckGoogleFolderIamPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_folder_iam_policy" {
			continue
		}

		folder := rs.Primary.Attributes["folder"]
		policy, err := config.clientResourceManagerV2Beta1.Folders.GetIamPolicy(folder, &resourceManagerV2Beta1.GetIamPolicyRequest{}).Do()

		if err != nil && len(policy.Bindings) > 0 {
			return fmt.Errorf("Folder '%s' policy hasn't been deleted.", folder)
		}
	}
	return nil
}

func getGoogleFolderIamPolicyFromResource(resource *terraform.InstanceState) (cloudresourcemanager.Policy, error) {
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

func getGoogleFolderIamPolicyFromState(s *terraform.State, res string) (cloudresourcemanager.Policy, error) {
	folder, err := getStatePrimaryResource(s, res, "")
	if err != nil {
		return cloudresourcemanager.Policy{}, err
	}
	return getGoogleFolderIamPolicyFromResource(folder)
}

func testAccCheckGoogleFolderIamPolicyExists(folder, policyRes string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		folderPolicy, err := getGoogleFolderIamPolicyFromState(s, folder)
		if err != nil {
			return fmt.Errorf("Error retrieving IAM policy for folder from state: %s", err)
		}
		policyPolicy, err := getGoogleFolderIamPolicyFromState(s, policyRes)
		if err != nil {
			return fmt.Errorf("Error retrieving IAM policy for data_policy from state: %s", err)
		}

		// The bindings in both policies should be identical
		if !compareBindings(folderPolicy.Bindings, policyPolicy.Bindings) {
			return fmt.Errorf("Folder and data source policies do not match: folder policy is %+v, data resource policy is  %+v", debugPrintBindings(folderPolicy.Bindings), debugPrintBindings(policyPolicy.Bindings))
		}

		// The audit configs in both policies should be identical
		if !compareAuditConfigs(folderPolicy.AuditConfigs, policyPolicy.AuditConfigs) {
			return fmt.Errorf("Folder and data source policies do not match: folder policy is %+v, data resource policy is  %+v", debugPrintAuditConfigs(folderPolicy.AuditConfigs), debugPrintAuditConfigs(policyPolicy.AuditConfigs))
		}
		return nil
	}
}

// Confirm that a folder has an IAM policy with at least 1 binding
func testAccFolderExistingPolicy(org, fname string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Config)
		var err error
		originalPolicy, err = getFolderIamPolicyByParentAndDisplayName("organizations/"+org, fname, c)
		if err != nil {
			return fmt.Errorf("Failed to retrieve IAM Policy for folder %q: %s", fname, err)
		}
		if len(originalPolicy.Bindings) == 0 {
			return fmt.Errorf("Refuse to run test against folder with zero IAM Bindings. This is likely an error in the test code that is not properly identifying the IAM policy of a folder.")
		}
		return nil
	}
}

func testAccFolderAssociatePolicyAuditConfigBasic(folder, parent string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}
resource "google_folder_iam_policy" "acceptance" {
    folder = "${google_folder.acceptance.id}"
    policy_data = "${data.google_iam_policy.admin.policy_data}"
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
      log_type = "DATA_READ"
      exempted_members = ["user:paddy@hashicorp.com"]
    }
    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
  audit_config {
    service = "cloudsql.googleapis.com"
    audit_log_configs {
      log_type = "DATA_READ"
      exempted_members = ["user:paddy@hashicorp.com"]
    }
    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
}
`, folder, parent)
}

func testAccFolder_create(folder, parent string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}`, folder, parent)
}

func testAccFolderIamPolicy_basic(folder, parent, role, member string) string {
	return fmt.Sprintf(`
resource "google_folder" "permissiontest" {
  display_name = "%s"
  parent       = "%s"
}
data "google_iam_policy" "test" {
  binding {
    role    = "%s"
    members = ["%s"]
  }
}
resource "google_folder_iam_policy" "test" {
  folder      = google_folder.permissiontest.name
  policy_data = data.google_iam_policy.test.policy_data
}
`, folder, parent, role, member)
}

func testAccFolderIamPolicy_basic2(folder, parent, role, member, role2, member2 string) string {
	return fmt.Sprintf(`
resource "google_folder" "permissiontest" {
  display_name = "%s"
  parent       = "%s"
}
data "google_iam_policy" "test" {
  binding {
    role    = "%s"
    members = ["%s"]
  }
  binding {
    role    = "%s"
    members = ["%s"]
  }
}
resource "google_folder_iam_policy" "test" {
  folder      = google_folder.permissiontest.name
  policy_data = data.google_iam_policy.test.policy_data
}
`, folder, parent, role, member, role2, member2)
}

func testAccFolderAssociatePolicyAuditConfigExpanded(folder, parent string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}
resource "google_folder_iam_policy" "acceptance" {
    folder = "${google_folder.acceptance.id}"
    policy_data = "${data.google_iam_policy.expanded.policy_data}"
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
      log_type = "DATA_READ"
      exempted_members = ["user:paddy@hashicorp.com"]
    }
    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
  audit_config {
    service = "cloudkms.googleapis.com"
    audit_log_configs {
      log_type = "DATA_READ"
      exempted_members = ["user:paddy@hashicorp.com"]
    }
    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
}`, folder, parent)
}
