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

// Test that an IAM binding can be applied to a project
func TestAccGoogleProjectIamBinding_basic(t *testing.T) {
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			resource.TestStep{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			resource.TestStep{
				Config: testAccGoogleProjectAssociateBindingBasic(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, pid),
				),
			},
		},
	})
}

// Test that multiple IAM bindings can be applied to a project
func TestAccGoogleProjectIamBinding_multiple(t *testing.T) {
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			resource.TestStep{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			resource.TestStep{
				Config: testAccGoogleProjectAssociateBindingBasic(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, pid),
				),
			},
			// Apply another IAM binding
			resource.TestStep{
				Config: testAccGoogleProjectAssociateBindingMultiple(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_binding.multiple", &cloudresourcemanager.Binding{
						Role:    "roles/viewer",
						Members: []string{"user:paddy@hashicorp.com"},
					}, pid),
				),
			},
		},
	})
}

// Test that an IAM binding can be updated once applied to a project
func TestAccGoogleProjectIamBinding_update(t *testing.T) {
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			resource.TestStep{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			resource.TestStep{
				Config: testAccGoogleProjectAssociateBindingBasic(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, pid),
				),
			},
			// Apply an updated IAM binding
			resource.TestStep{
				Config: testAccGoogleProjectAssociateBindingUpdated(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_binding.updated", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:paddy@hashicorp.com"},
					}, pid),
				),
			},
		},
	})
}

// Test that an IAM binding can be removed from a project
func TestAccGoogleProjectIamBinding_remove(t *testing.T) {
	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			resource.TestStep{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
			// Apply multiple IAM bindings
			resource.TestStep{
				Config: testAccGoogleProjectAssociateBindingMultiple(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_binding.multiple", &cloudresourcemanager.Binding{
						Role:    "roles/viewer",
						Members: []string{"user:paddy@hashicorp.com"},
					}, pid),
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_binding.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, pid),
				),
			},
			// Remove the bindings
			resource.TestStep{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
		},
	})
}

func testAccCheckGoogleProjectIamBindingExists(key string, expected *cloudresourcemanager.Binding, pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		projectPolicy, err := getProjectIamPolicy(pid, config)
		if err != nil {
			return fmt.Errorf("Failed to retrieve IAM policy for project %q: %s", pid, err)
		}

		var result *cloudresourcemanager.Binding
		for _, binding := range projectPolicy.Bindings {
			if binding.Role == expected.Role {
				result = binding
				break
			}
		}
		if result == nil {
			return fmt.Errorf("IAM policy for project %q had no role %q", pid, expected.Role)
		}
		if len(result.Members) != len(expected.Members) {
			return fmt.Errorf("Got %v as members for role %q of project %q, expected %v", result.Members, expected.Role, pid, expected.Members)
		}
		sort.Strings(result.Members)
		sort.Strings(expected.Members)
		for pos, exp := range expected.Members {
			if result.Members[pos] != exp {
				return fmt.Errorf("Expected members for role %q of project %q to be %v, got %v", expected.Role, pid, expected.Members, result.Members)
			}
		}
		return nil
	}
}

func testAccGoogleProjectAssociateBindingBasic(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
    project_id = "%s"
    name = "%s"
    org_id = "%s"
}
resource "google_project_iam_binding" "acceptance" {
    project = "${google_project.acceptance.id}"
    members = ["user:admin@hashicorptest.com"]
    role = "roles/compute.instanceAdmin"
}
`, pid, name, org)
}

func testAccGoogleProjectAssociateBindingMultiple(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
    project_id = "%s"
    name = "%s"
    org_id = "%s"
}
resource "google_project_iam_binding" "acceptance" {
    project = "${google_project.acceptance.id}"
    members = ["user:admin@hashicorptest.com"]
    role = "roles/compute.instanceAdmin"
}
resource "google_project_iam_binding" "multiple" {
    project = "${google_project.acceptance.id}"
    members = ["user:paddy@hashicorp.com"]
    role = "roles/viewer"
}
`, pid, name, org)
}

func testAccGoogleProjectAssociateBindingUpdated(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
    project_id = "%s"
    name = "%s"
    org_id = "%s"
}
resource "google_project_iam_binding" "acceptance" {
    project = "${google_project.acceptance.id}"
    members = ["user:admin@hashicorptest.com", "user:paddy@hashicorp.com"]
    role = "roles/compute.instanceAdmin"
}
`, pid, name, org)
}
