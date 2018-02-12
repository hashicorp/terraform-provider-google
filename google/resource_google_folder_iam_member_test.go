package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Test that an IAM binding can be applied to a folder
func TestAccGoogleFolderIamMember_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccGoogleFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleFolderExistingPolicy(org, fname),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccGoogleFolderAssociateMemberBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_member.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that multiple IAM bindings can be applied to a folder
func TestAccGoogleFolderIamMember_multiple(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccGoogleFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleFolderExistingPolicy(org, fname),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccGoogleFolderAssociateMemberBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_member.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Apply another IAM binding
			{
				Config: testAccGoogleFolderAssociateMemberMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_member.multiple", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:paddy@hashicorp.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that an IAM binding can be removed from a folder
func TestAccGoogleFolderIamMember_remove(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccGoogleFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleFolderExistingPolicy(org, fname),
				),
			},
			// Apply multiple IAM bindings
			{
				Config: testAccGoogleFolderAssociateMemberMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_member.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:paddy@hashicorp.com"},
					}, org, fname),
				),
			},
			// Remove the bindings
			{
				Config: testAccGoogleFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleFolderExistingPolicy(org, fname),
				),
			},
		},
	})
}

func testAccGoogleFolderAssociateMemberBasic(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_member" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  member  = "user:admin@hashicorptest.com"
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}

func testAccGoogleFolderAssociateMemberMultiple(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_member" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  member  = "user:admin@hashicorptest.com"
  role    = "roles/compute.instanceAdmin"
}

resource "google_folder_iam_member" "multiple" {
  folder = "${google_folder.acceptance.name}"
  member  = "user:paddy@hashicorp.com"
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}
