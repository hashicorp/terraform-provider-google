package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func projectIamMemberImportStep(resourceName, pid, role, member string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportStateId:     fmt.Sprintf("%s %s %s", pid, role, member),
		ImportState:       true,
		ImportStateVerify: true,
	}
}

// Test that an IAM binding can be applied to a project
func TestAccProjectIamMember_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	resourceName := "google_project_iam_member.acceptance"
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"
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
				Config: testAccProjectAssociateMemberBasic(pid, pname, org, role, member),
			},
			projectIamMemberImportStep(resourceName, pid, role, member),
		},
	})
}

// Test that multiple IAM bindings can be applied to a project
func TestAccProjectIamMember_multiple(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	skipIfEnvNotSet(t, "GOOGLE_ORG")

	pid := "terraform-" + acctest.RandString(10)
	resourceName := "google_project_iam_member.acceptance"
	resourceName2 := "google_project_iam_member.multiple"
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"
	member2 := "user:paddy@hashicorp.com"

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
				Config: testAccProjectAssociateMemberBasic(pid, pname, org, role, member),
			},
			projectIamMemberImportStep(resourceName, pid, role, member),

			// Apply another IAM binding
			{
				Config: testAccProjectAssociateMemberMultiple(pid, pname, org, role, member, role, member2),
			},
			projectIamMemberImportStep(resourceName, pid, role, member),
			projectIamMemberImportStep(resourceName2, pid, role, member2),
		},
	})
}

// Test that an IAM binding can be removed from a project
func TestAccProjectIamMember_remove(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	skipIfEnvNotSet(t, "GOOGLE_ORG")

	pid := "terraform-" + acctest.RandString(10)
	resourceName := "google_project_iam_member.acceptance"
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"
	member2 := "user:paddy@hashicorp.com"

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
				Config: testAccProjectAssociateMemberMultiple(pid, pname, org, role, member, role, member2),
			},
			projectIamMemberImportStep(resourceName, pid, role, member),
			projectIamMemberImportStep(resourceName, pid, role, member2),

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

func testAccProjectAssociateMemberBasic(pid, name, org, role, member string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_member" "acceptance" {
  project = google_project.acceptance.project_id
  role    = "%s"
  member  = "%s"
}
`, pid, name, org, role, member)
}

func testAccProjectAssociateMemberMultiple(pid, name, org, role, member, role2, member2 string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_member" "acceptance" {
  project = google_project.acceptance.project_id
  role    = "%s"
  member  = "%s"
}

resource "google_project_iam_member" "multiple" {
  project = google_project.acceptance.project_id
  role    = "%s"
  member  = "%s"
}
`, pid, name, org, role, member, role2, member2)
}
