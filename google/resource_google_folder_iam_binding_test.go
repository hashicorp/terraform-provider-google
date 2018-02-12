package google

import (
	"fmt"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Test that an IAM binding can be applied to a folder
func TestAccGoogleFolderIamBinding_basic(t *testing.T) {
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
				Config: testAccGoogleFolderAssociateBindingBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that multiple IAM bindings can be applied to a folder, one at a time
func TestAccGoogleFolderIamBinding_multiple(t *testing.T) {
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
				Config: testAccGoogleFolderAssociateBindingBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Apply another IAM binding
			{
				Config: testAccGoogleFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.multiple", &cloudresourcemanager.Binding{
						Role:    "roles/viewer",
						Members: []string{"user:paddy@hashicorp.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.multiple", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that multiple IAM bindings can be applied to a folder all at once
func TestAccGoogleFolderIamBinding_multipleAtOnce(t *testing.T) {
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
				Config: testAccGoogleFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.multiple", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that an IAM binding can be updated once applied to a folder
func TestAccGoogleFolderIamBinding_update(t *testing.T) {
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
				Config: testAccGoogleFolderAssociateBindingBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Apply an updated IAM binding
			{
				Config: testAccGoogleFolderAssociateBindingUpdated(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.updated", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:paddy@hashicorp.com"},
					}, org, fname),
				),
			},
			// Drop the original member
			{
				Config: testAccGoogleFolderAssociateBindingDropMemberFromBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.dropped", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:paddy@hashicorp.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that an IAM binding can be removed from a folder
func TestAccGoogleFolderIamBinding_remove(t *testing.T) {
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
				Config: testAccGoogleFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.multiple", &cloudresourcemanager.Binding{
						Role:    "roles/viewer",
						Members: []string{"user:paddy@hashicorp.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists("google_folder_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
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

func testAccCheckGoogleFolderIamBindingExists(key string, expected *cloudresourcemanager.Binding, org, fname string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		folderPolicy, err := getFolderIamPolicyByParentAndDisplayName("organizations/"+org, fname, config)
		if err != nil {
			return fmt.Errorf("Failed to retrieve IAM policy for folder %q: %s", fname, err)
		}

		var result *cloudresourcemanager.Binding
		for _, binding := range folderPolicy.Bindings {
			if binding.Role == expected.Role {
				result = binding
				break
			}
		}
		if result == nil {
			return fmt.Errorf("IAM policy for folder %q had no role %q", fname, expected.Role)
		}
		if len(result.Members) != len(expected.Members) {
			return fmt.Errorf("Got %v as members for role %q of folder %q, expected %v", result.Members, expected.Role, fname, expected.Members)
		}
		sort.Strings(result.Members)
		sort.Strings(expected.Members)
		for pos, exp := range expected.Members {
			if result.Members[pos] != exp {
				return fmt.Errorf("Expected members for role %q of folder %q to be %v, got %v", expected.Role, fname, expected.Members, result.Members)
			}
		}
		return nil
	}
}

func testAccGoogleFolderIamBasic(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}
`, org, fname)
}

func testAccGoogleFolderAssociateBindingBasic(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_binding" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  members = ["user:admin@hashicorptest.com"]
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}

func testAccGoogleFolderAssociateBindingMultiple(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_binding" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  members = ["user:admin@hashicorptest.com"]
  role    = "roles/compute.instanceAdmin"
}

resource "google_folder_iam_binding" "multiple" {
  folder = "${google_folder.acceptance.name}"
  members = ["user:paddy@hashicorp.com"]
  role    = "roles/viewer"
}
`, org, fname)
}

func testAccGoogleFolderAssociateBindingUpdated(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_binding" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  members = ["user:admin@hashicorptest.com", "user:paddy@hashicorp.com"]
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}

func testAccGoogleFolderAssociateBindingDropMemberFromBasic(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_binding" "dropped" {
  folder = "${google_folder.acceptance.name}"
  members = ["user:paddy@hashicorp.com"]
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}
