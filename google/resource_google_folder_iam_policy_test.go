package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func TestAccFolderIamPolicy_auditConfigs(t *testing.T) {
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
				Config: testAccFolderIamPolicy_auditConfigs(folderDisplayName, parent, "roles/viewer", "user:admin@hashicorptest.com"),
			},
			{
				ResourceName:      "google_folder_iam_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccFolderIamPolicy_basic(folder, parent, role, member string) string {
	return fmt.Sprintf(`
resource "google_folder" "permissiontest" {
  display_name = "%s"
  parent = "%s"
}

data "google_iam_policy" "test" {
  binding {
    role = "%s"
    members = ["%s"]
  }
}

resource "google_folder_iam_policy" "test" {
  folder = "${google_folder.permissiontest.name}"
  policy_data = "${data.google_iam_policy.test.policy_data}"
}
`, folder, parent, role, member)
}

func testAccFolderIamPolicy_basic2(folder, parent, role, member, role2, member2 string) string {
	return fmt.Sprintf(`
resource "google_folder" "permissiontest" {
  display_name = "%s"
  parent = "%s"
}

data "google_iam_policy" "test" {
  binding {
    role = "%s"
    members = ["%s"]
  }

  binding {
    role = "%s"
    members = ["%s"]
  }
}

resource "google_folder_iam_policy" "test" {
  folder = "${google_folder.permissiontest.name}"
  policy_data = "${data.google_iam_policy.test.policy_data}"
}
`, folder, parent, role, member, role2, member2)
}

func testAccFolderIamPolicy_auditConfigs(folder, parent, role, member string) string {
	return fmt.Sprintf(`
resource "google_folder" "permissiontest" {
  display_name = "%s"
  parent = "%s"
}

data "google_iam_policy" "test" {
  binding {
    role = "%s"
    members = ["%s"]
  }
  audit_config {
    service = "cloudkms.googleapis.com"
    audit_log_configs {
      log_type = "DATA_READ"
      exempted_members = ["%s"]
    }

    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
  audit_config {
    service = "cloudsql.googleapis.com"
    audit_log_configs {
      log_type = "DATA_READ"
      exempted_members = ["%s"]
    }

    audit_log_configs {
      log_type = "DATA_WRITE"
    }
  }
}

resource "google_folder_iam_policy" "test" {
  folder = "${google_folder.permissiontest.name}"
  policy_data = "${data.google_iam_policy.test.policy_data}"
}
`, folder, parent, role, member, member, member)
}
