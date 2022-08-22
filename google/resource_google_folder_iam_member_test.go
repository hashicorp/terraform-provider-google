package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Test that an IAM binding can be applied to a folder
func TestAccFolderIamMember_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "tf-test-" + randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccFolderAssociateMemberBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that multiple IAM bindings can be applied to a folder
func TestAccFolderIamMember_multiple(t *testing.T) {
	skipIfVcr(t)
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "tf-test-" + randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccFolderAssociateMemberBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Apply another IAM binding
			{
				Config: testAccFolderAssociateMemberMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that an IAM binding can be removed from a folder
func TestAccFolderIamMember_remove(t *testing.T) {
	skipIfVcr(t)
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "tf-test-" + randString(t, 10)
	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply multiple IAM bindings
			{
				Config: testAccFolderAssociateMemberMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"},
					}, org, fname),
				),
			},
			// Remove the bindings
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
		},
	})
}

func testAccFolderAssociateMemberBasic(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_member" "acceptance" {
  folder = google_folder.acceptance.name
  member = "user:admin@hashicorptest.com"
  role   = "roles/compute.instanceAdmin"
}
`, org, fname)
}

func testAccFolderAssociateMemberMultiple(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_member" "acceptance" {
  folder = google_folder.acceptance.name
  member = "user:admin@hashicorptest.com"
  role   = "roles/compute.instanceAdmin"
}

resource "google_folder_iam_member" "multiple" {
  folder = google_folder.acceptance.name
  member = "user:gterraformtest1@gmail.com"
  role   = "roles/compute.instanceAdmin"
}

resource "google_folder_iam_member" "condition" {
  folder = google_folder.acceptance.name
  member = "user:gterraformtest1@gmail.com"
  role   = "roles/compute.instanceAdmin"
  condition {
    title       = "expires_after_2019_12_31"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
`, org, fname)
}
