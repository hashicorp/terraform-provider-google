// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	resourceName := "google_project_iam_member.acceptance"
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateMemberBasic(pid, org, role, member),
			},
			projectIamMemberImportStep(resourceName, pid, role, member),
		},
	})
}

// Test that multiple IAM bindings can be applied to a project
func TestAccProjectIamMember_multiple(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	acctest.SkipIfEnvNotSet(t, "GOOGLE_ORG")

	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	resourceName := "google_project_iam_member.acceptance"
	resourceName2 := "google_project_iam_member.multiple"
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"
	member2 := "user:gterraformtest1@gmail.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateMemberBasic(pid, org, role, member),
			},
			projectIamMemberImportStep(resourceName, pid, role, member),

			// Apply another IAM binding
			{
				Config: testAccProjectAssociateMemberMultiple(pid, org, role, member, role, member2),
			},
			projectIamMemberImportStep(resourceName, pid, role, member),
			projectIamMemberImportStep(resourceName2, pid, role, member2),
		},
	})
}

// Test that an IAM binding can be removed from a project
func TestAccProjectIamMember_remove(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	acctest.SkipIfEnvNotSet(t, "GOOGLE_ORG")

	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	resourceName := "google_project_iam_member.acceptance"
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"
	member2 := "user:gterraformtest1@gmail.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},

			// Apply multiple IAM bindings
			{
				Config: testAccProjectAssociateMemberMultiple(pid, org, role, member, role, member2),
			},
			projectIamMemberImportStep(resourceName, pid, role, member),
			projectIamMemberImportStep(resourceName, pid, role, member2),

			// Remove the bindings
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
		},
	})
}

func TestAccProjectIamMember_withCondition(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	resourceName := "google_project_iam_member.acceptance"
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"
	conditionTitle := "expires_after_2019_12_31"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateMember_withCondition(pid, org, role, member, conditionTitle),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     fmt.Sprintf("%s %s %s %s", pid, role, member, conditionTitle),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProjectIamMember_invalidMembers(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectAssociateMemberBasic(pid, org, role, "admin@hashicorptest.com"),
				ExpectError: regexp.MustCompile("invalid value for member \\(IAM members must have one of the values outlined here: https://cloud.google.com/billing/docs/reference/rest/v1/Policy#Binding\\)"),
			},
			{
				Config: testAccProjectAssociateMemberBasic(pid, org, role, "user:admin@hashicorptest.com"),
			},
		},
	})
}

func testAccProjectAssociateMemberBasic(pid, org, role, member string) string {
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
`, pid, pid, org, role, member)
}

func testAccProjectAssociateMemberMultiple(pid, org, role, member, role2, member2 string) string {
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
`, pid, pid, org, role, member, role2, member2)
}

func testAccProjectAssociateMember_withCondition(pid, org, role, member, conditionTitle string) string {
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
  condition {
    title       = "%s"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
`, pid, pid, org, role, member, conditionTitle)
}
