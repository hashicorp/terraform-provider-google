package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Test that an IAM binding can be applied to a project
func TestAccGoogleProjectIamMember_basic(t *testing.T) {
	t.Parallel()

	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccGoogleProjectAssociateMemberBasic(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_member.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, pid),
				),
			},
		},
	})
}

// Test that multiple IAM bindings can be applied to a project
func TestAccGoogleProjectIamMember_multiple(t *testing.T) {
	t.Parallel()

	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccGoogleProjectAssociateMemberBasic(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_member.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com"},
					}, pid),
				),
			},
			// Apply another IAM binding
			{
				Config: testAccGoogleProjectAssociateMemberMultiple(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_member.multiple", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:paddy@hashicorp.com"},
					}, pid),
				),
			},
		},
	})
}

// Test that an IAM binding can be removed from a project
func TestAccGoogleProjectIamMember_remove(t *testing.T) {
	t.Parallel()

	pid := "terraform-" + acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
			// Apply multiple IAM bindings
			{
				Config: testAccGoogleProjectAssociateMemberMultiple(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectIamBindingExists("google_project_iam_member.acceptance", &cloudresourcemanager.Binding{
						Role:    "roles/compute.instanceAdmin",
						Members: []string{"user:admin@hashicorptest.com", "user:paddy@hashicorp.com"},
					}, pid),
				),
			},
			// Remove the bindings
			{
				Config: testAccGoogleProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccGoogleProjectExistingPolicy(pid),
				),
			},
		},
	})
}

func testAccGoogleProjectAssociateMemberBasic(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_member" "acceptance" {
  project = "${google_project.acceptance.project_id}"
  member  = "user:admin@hashicorptest.com"
  role    = "roles/compute.instanceAdmin"
}
`, pid, name, org)
}

func testAccGoogleProjectAssociateMemberMultiple(pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_member" "acceptance" {
  project = "${google_project.acceptance.project_id}"
  member  = "user:admin@hashicorptest.com"
  role    = "roles/compute.instanceAdmin"
}

resource "google_project_iam_member" "multiple" {
  project = "${google_project.acceptance.project_id}"
  member  = "user:paddy@hashicorp.com"
  role    = "roles/compute.instanceAdmin"
}
`, pid, name, org)
}
