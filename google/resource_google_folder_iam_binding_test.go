package google

import (
	"fmt"
	"sort"
	"testing"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

// Test that an IAM binding can be applied to a folder
func TestAccFolderIamBinding_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, fname),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccFolderAssociateBindingBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that multiple IAM bindings can be applied to a folder, one at a time
func TestAccFolderIamBinding_multiple(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, fname),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccFolderAssociateBindingBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Apply another IAM binding
			{
				Config: testAccFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/viewer",
						Members: []string{"user:paddy@hashicorp.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that multiple IAM bindings can be applied to a folder all at once
func TestAccFolderIamBinding_multipleAtOnce(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, fname),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that an IAM binding can be updated once applied to a folder
func TestAccFolderIamBinding_update(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, fname),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccFolderAssociateBindingBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Apply an updated IAM binding
			{
				Config: testAccFolderAssociateBindingUpdated(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:paddy@hashicorp.com"},
					}, org, fname),
				),
			},
			// Drop the original member
			{
				Config: testAccFolderAssociateBindingDropMemberFromBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:paddy@hashicorp.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that an IAM binding can be removed from a folder
func TestAccFolderIamBinding_remove(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	fname := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, fname),
				),
			},
			// Apply multiple IAM bindings
			{
				Config: testAccFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/viewer",
						Members: []string{"user:paddy@hashicorp.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists(&cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Remove the bindings
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, fname),
				),
			},
		},
	})
}

func testAccCheckGoogleFolderIamBindingExists(expected *cloudresourcemanager.Binding, org, fname string) resource.TestCheckFunc {
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
			return fmt.Errorf("IAM policy for folder %q had no role %q, got %#v", fname, expected.Role, folderPolicy.Bindings)
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

func getFolderIamPolicyByParentAndDisplayName(parent, displayName string, config *Config) (*cloudresourcemanager.Policy, error) {
	queryString := fmt.Sprintf("lifecycleState=ACTIVE AND parent=%s AND displayName=%s", parent, displayName)
	searchRequest := &resourceManagerV2Beta1.SearchFoldersRequest{
		Query: queryString,
	}
	searchResponse, err := config.clientResourceManagerV2Beta1.Folders.Search(searchRequest).Do()
	if err != nil {
		if isGoogleApiErrorWithCode(err, 404) {
			return nil, fmt.Errorf("Folder not found: %s,%s", parent, displayName)
		}

		return nil, errwrap.Wrapf("Error reading folders: {{err}}", err)
	}

	folders := searchResponse.Folders
	if len(folders) != 1 {
		return nil, fmt.Errorf("expected exactly 1 folder, found %d", len(folders))
	}

	return getFolderIamPolicyByFolderName(folders[0].Name, config)
}

func testAccFolderIamBasic(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}
`, org, fname)
}

func testAccFolderAssociateBindingBasic(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_binding" "acceptance" {
  folder  = google_folder.acceptance.name
  members = ["user:admin@hashicorptest.com"]
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}

func testAccFolderAssociateBindingMultiple(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_binding" "acceptance" {
  folder  = google_folder.acceptance.name
  members = ["user:admin@hashicorptest.com"]
  role    = "roles/compute.instanceAdmin"
}

resource "google_folder_iam_binding" "multiple" {
  folder  = google_folder.acceptance.name
  members = ["user:paddy@hashicorp.com"]
  role    = "roles/viewer"
}
`, org, fname)
}

func testAccFolderAssociateBindingUpdated(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_binding" "acceptance" {
  folder  = google_folder.acceptance.name
  members = ["user:admin@hashicorptest.com", "user:paddy@hashicorp.com"]
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}

func testAccFolderAssociateBindingDropMemberFromBasic(org, fname string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_binding" "acceptance" {
  folder  = google_folder.acceptance.name
  members = ["user:paddy@hashicorp.com"]
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}
