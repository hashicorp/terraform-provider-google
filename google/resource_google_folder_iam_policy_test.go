package google

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

func TestAccFolderIamPolicy_basic(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org

	policy1 := &resourceManagerV2Beta1.Policy{
		Bindings: []*resourceManagerV2Beta1.Binding{
			{
				Role: "roles/viewer",
				Members: []string{
					"user:admin@hashicorptest.com",
				},
			},
		},
	}
	policy2 := &resourceManagerV2Beta1.Policy{
		Bindings: []*resourceManagerV2Beta1.Binding{
			{
				Role: "roles/editor",
				Members: []string{
					"user:admin@hashicorptest.com",
				},
			},
			{
				Role: "roles/viewer",
				Members: []string{
					"user:admin@hashicorptest.com",
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderIamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFolderIamPolicy_basic(folderDisplayName, parent, "roles/viewer", "user:admin@hashicorptest.com"),
				Check:  testAccCheckGoogleFolderIamPolicy("google_folder_iam_policy.test", policy1),
			},
			{
				Config: testAccFolderIamPolicy_basic2(folderDisplayName, parent, "roles/editor", "user:admin@hashicorptest.com", "roles/viewer", "user:admin@hashicorptest.com"),
				Check:  testAccCheckGoogleFolderIamPolicy("google_folder_iam_policy.test", policy2),
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

func testAccCheckGoogleFolderIamPolicy(n string, policy *resourceManagerV2Beta1.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		p, err := config.clientResourceManagerV2Beta1.Folders.GetIamPolicy(rs.Primary.ID, &resourceManagerV2Beta1.GetIamPolicyRequest{}).Do()
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(p.Bindings, policy.Bindings) {
			return fmt.Errorf("Incorrect iam policy bindings. Expected '%v', got '%v'", policy.Bindings, p.Bindings)
		}

		if _, ok = rs.Primary.Attributes["etag"]; !ok {
			return fmt.Errorf("Etag should be set.")
		}

		if rs.Primary.Attributes["etag"] != p.Etag {
			return fmt.Errorf("Incorrect etag value. Expected '%s', got '%s'", p.Etag, rs.Primary.Attributes["etag"])
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
