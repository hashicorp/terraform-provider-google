package google

import (
	"fmt"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV3 "google.golang.org/api/cloudresourcemanager/v3"
)

// Test that an IAM binding can be applied to a folder
func TestAccFolderIamBinding_basic(t *testing.T) {
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
				Config: testAccFolderAssociateBindingBasic(org, fname),
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

// Test that multiple IAM bindings can be applied to a folder, one at a time
func TestAccFolderIamBinding_multiple(t *testing.T) {
	// Multiple fine-grained resources
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
				Config: testAccFolderAssociateBindingBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Apply another IAM binding
			{
				Config: testAccFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/viewer",
						Members: []string{"user:gterraformtest1@gmail.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
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
	// Multiple fine-grained resources
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
				Config: testAccFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
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
				Config: testAccFolderAssociateBindingBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, org, fname),
				),
			},
			// Apply an updated IAM binding
			{
				Config: testAccFolderAssociateBindingUpdated(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"},
					}, org, fname),
				),
			},
			// Drop the original member
			{
				Config: testAccFolderAssociateBindingDropMemberFromBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:gterraformtest1@gmail.com"},
					}, org, fname),
				),
			},
		},
	})
}

// Test that an IAM binding can be removed from a folder
func TestAccFolderIamBinding_remove(t *testing.T) {
	// Multiple fine-grained resources
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
				Config: testAccFolderAssociateBindingMultiple(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/viewer",
						Members: []string{"user:gterraformtest1@gmail.com"},
					}, org, fname),
					testAccCheckGoogleFolderIamBindingExists(t, &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
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

func testAccCheckGoogleFolderIamBindingExists(t *testing.T, expected *cloudresourcemanager.Binding, org, fname string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
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
	var folderMatch *resourceManagerV3.Folder
	token := ""

	for paginate := true; paginate; {
		resp, err := config.NewResourceManagerV3Client(config.userAgent).Folders.List().Parent(parent).PageSize(300).PageToken(token).Do()
		if err != nil {
			return nil, fmt.Errorf("Error reading folder list: %s", err)
		}

		for _, folder := range resp.Folders {
			if folder.DisplayName == displayName {
				if folderMatch != nil {
					return nil, fmt.Errorf("More than one matching folder found")
				}
				folderMatch = folder
			}
		}

		token = resp.NextPageToken
		paginate = token != ""
	}

	if folderMatch == nil {
		return nil, fmt.Errorf("Folder not found: %s", displayName)
	}

	return getFolderIamPolicyByFolderName(folderMatch.Name, config.userAgent, config)
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
  members = ["user:gterraformtest1@gmail.com"]
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
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
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
  members = ["user:gterraformtest1@gmail.com"]
  role    = "roles/compute.instanceAdmin"
}
`, org, fname)
}
