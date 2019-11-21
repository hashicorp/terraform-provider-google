package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func projectIamBindingImportStep(resourceName, pid, role string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportStateId:     fmt.Sprintf("%s %s", pid, role),
		ImportState:       true,
		ImportStateVerify: true,
	}
}

// Test that an IAM binding can be applied to a project
func TestAccProjectIamBinding_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	role := "roles/compute.instanceAdmin"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingBasic(pid, pname, org, role),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
		},
	})
}

// Test that multiple IAM bindings can be applied to a project, one at a time
func TestAccProjectIamBinding_multiple(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	role := "roles/compute.instanceAdmin"
	role2 := "roles/viewer"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingBasic(pid, pname, org, role),
			},
			// Apply another IAM binding
			{
				Config: testAccProjectAssociateBindingMultiple(pid, pname, org, role, role2),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
			projectIamBindingImportStep("google_project_iam_binding.multiple", pid, role2),
		},
	})
}

// Test that multiple IAM bindings can be applied to a project all at once
func TestAccProjectIamBinding_multipleAtOnce(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	role := "roles/compute.instanceAdmin"
	role2 := "roles/viewer"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingMultiple(pid, pname, org, role, role2),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
			projectIamBindingImportStep("google_project_iam_binding.multiple", pid, role2),
		},
	})
}

// Test that an IAM binding can be updated once applied to a project
func TestAccProjectIamBinding_update(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	role := "roles/compute.instanceAdmin"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingBasic(pid, pname, org, role),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),

			// Apply an updated IAM binding
			{
				Config: testAccProjectAssociateBindingUpdated(pid, pname, org, role),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),

			// Drop the original member
			{
				Config: testAccProjectAssociateBindingDropMemberFromBasic(pid, pname, org, role),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
		},
	})
}

// Test that an IAM binding can be removed from a project
func TestAccProjectIamBinding_remove(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	role := "roles/compute.instanceAdmin"
	role2 := "roles/viewer"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply multiple IAM bindings
			{
				Config: testAccProjectAssociateBindingMultiple(pid, pname, org, role, role2),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
			projectIamBindingImportStep("google_project_iam_binding.multiple", pid, role2),

			// Remove the bindings
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
		},
	})
}

// Test that an IAM binding with no members can be applied to a project
func TestAccProjectIamBinding_noMembers(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	role := "roles/compute.instanceAdmin"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingNoMembers(pid, pname, org, role),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
		},
	})
}

func testAccProjectAssociateBindingBasic(pid, name, org, role string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["user:admin@hashicorptest.com"]
  role    = "%s"
}
`, pid, name, org, role)
}

func testAccProjectAssociateBindingMultiple(pid, name, org, role, role2 string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["user:admin@hashicorptest.com"]
  role    = "%s"
}

resource "google_project_iam_binding" "multiple" {
  project = google_project.acceptance.project_id
  members = ["user:paddy@hashicorp.com"]
  role    = "%s"
}
`, pid, name, org, role, role2)
}

func testAccProjectAssociateBindingUpdated(pid, name, org, role string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["user:admin@hashicorptest.com", "user:paddy@hashicorp.com"]
  role    = "%s"
}
`, pid, name, org, role)
}

func testAccProjectAssociateBindingDropMemberFromBasic(pid, name, org, role string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["user:paddy@hashicorp.com"]
  role    = "%s"
}
`, pid, name, org, role)
}

func testAccProjectAssociateBindingNoMembers(pid, name, org, role string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = []
  role    = "%s"
}
`, pid, name, org, role)
}
